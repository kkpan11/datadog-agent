// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

//go:build linux_bpf

package usm

import (
	"errors"
	"fmt"
	"io"
	"syscall"
	"time"

	"github.com/DataDog/datadog-go/v5/statsd"
	"github.com/cilium/ebpf"
	"go.uber.org/atomic"

	manager "github.com/DataDog/ebpf-manager"

	ddebpf "github.com/DataDog/datadog-agent/pkg/ebpf"
	"github.com/DataDog/datadog-agent/pkg/network/config"
	filterpkg "github.com/DataDog/datadog-agent/pkg/network/filter"
	"github.com/DataDog/datadog-agent/pkg/network/protocols"
	"github.com/DataDog/datadog-agent/pkg/network/protocols/telemetry"
	usmconfig "github.com/DataDog/datadog-agent/pkg/network/usm/config"
	"github.com/DataDog/datadog-agent/pkg/network/usm/consts"
	usmstate "github.com/DataDog/datadog-agent/pkg/network/usm/state"
	"github.com/DataDog/datadog-agent/pkg/network/usm/utils"
	"github.com/DataDog/datadog-agent/pkg/process/monitor"
	"github.com/DataDog/datadog-agent/pkg/util/log"
)

var (
	startupError error
)

// Monitor is responsible for:
// * Creating a raw socket and attaching an eBPF filter to it;
// * Consuming HTTP transaction "events" that are sent from Kernel space;
// * Aggregating and emitting metrics based on the received HTTP transactions;
type Monitor struct {
	cfg *config.Config

	ebpfProgram *ebpfProgram

	processMonitor *monitor.ProcessMonitor

	// termination
	closeFilterFn func()

	lastUpdateTime *atomic.Int64

	telemetryStopChannel chan struct{}

	statsd statsd.ClientInterface
}

// NewMonitor returns a new Monitor instance
func NewMonitor(c *config.Config, connectionProtocolMap *ebpf.Map, statsd statsd.ClientInterface) (m *Monitor, err error) {
	defer func() {
		// capture error and wrap it
		if err != nil {
			usmstate.Set(usmstate.NotRunning)
			err = fmt.Errorf("could not initialize USM: %w", err)
			startupError = err
		}
	}()

	mgr, err := newEBPFProgram(c, connectionProtocolMap)
	if err != nil {
		return nil, fmt.Errorf("error setting up ebpf program: %w", err)
	}

	if len(mgr.enabledProtocols) == 0 {
		usmstate.Set(usmstate.Disabled)
		log.Debug("not enabling USM as no protocols monitoring were enabled.")
		return nil, nil
	}

	if err := mgr.Init(); err != nil {
		return nil, fmt.Errorf("error initializing ebpf program: %w", err)
	}

	filter, _ := mgr.GetProbe(manager.ProbeIdentificationPair{EBPFFuncName: protocolDispatcherSocketFilterFunction, UID: probeUID})
	if filter == nil {
		return nil, fmt.Errorf("error retrieving socket filter")
	}
	ddebpf.AddNameMappings(mgr.Manager.Manager, "usm_monitor")

	closeFilterFn, err := filterpkg.HeadlessSocketFilter(c, filter)
	if err != nil {
		return nil, fmt.Errorf("error enabling traffic inspection: %s", err)
	}

	processMonitor := monitor.GetProcessMonitor()

	usmstate.Set(usmstate.Running)

	usmMonitor := &Monitor{
		cfg:                  c,
		ebpfProgram:          mgr,
		closeFilterFn:        closeFilterFn,
		processMonitor:       processMonitor,
		telemetryStopChannel: make(chan struct{}),
		statsd:               statsd,
	}

	usmMonitor.lastUpdateTime = atomic.NewInt64(time.Now().Unix())

	return usmMonitor, nil
}

// Start USM monitor.
func (m *Monitor) Start() error {
	if m == nil {
		return nil
	}

	var err error

	defer func() {
		if err != nil {
			if errors.Is(err, syscall.ENOMEM) {
				err = fmt.Errorf("could not enable usm monitoring: not enough memory to attach http ebpf socket filter. please consider raising the limit via sysctl -w net.core.optmem_max=<LIMIT>")
			} else {
				err = fmt.Errorf("could not enable USM: %s", err)
			}

			m.Stop()

			startupError = err
		}
	}()

	err = m.ebpfProgram.Start()
	if err != nil {
		return err
	}

	ddebpf.AddProbeFDMappings(m.ebpfProgram.Manager.Manager)

	// Need to explicitly save the error in `err` so the defer function could save the startup error.
	if usmconfig.NeedProcessMonitor(m.cfg) {
		err = m.processMonitor.Initialize(m.cfg.EnableUSMEventStream)
	}

	if err != nil {
		return err
	}
	m.startTelemetryReporter()
	return nil
}

// Pause bypasses the eBPF programs in the monitor
func (m *Monitor) Pause() error {
	if m == nil {
		return nil
	}

	return m.ebpfProgram.Pause()
}

// Resume enables the previously bypassed eBPF programs in the monitor
func (m *Monitor) Resume() error {
	if m == nil {
		return nil
	}

	return m.ebpfProgram.Resume()
}

// GetUSMStats returns the current state of the USM monitor
func (m *Monitor) GetUSMStats() map[string]interface{} {
	response := map[string]interface{}{
		"state": usmstate.Get(),
	}

	if startupError != nil {
		response["error"] = startupError.Error()
	}

	response["blocked_processes"] = utils.GetBlockedPathIDsList(consts.USMModuleName)

	tracedPrograms := utils.GetTracedProgramList(consts.USMModuleName)
	response["traced_programs"] = tracedPrograms

	if m != nil {
		response["last_check"] = m.lastUpdateTime
	}
	return response
}

// GetProtocolStats returns the current stats for all protocols and a cleanup function to free resources.
func (m *Monitor) GetProtocolStats() (map[protocols.ProtocolType]interface{}, func()) {
	if m == nil {
		return nil, func() {}
	}

	defer func() {
		// Update update time
		now := time.Now().Unix()
		m.lastUpdateTime.Swap(now)
		telemetry.ReportPrometheus()
	}()

	return m.ebpfProgram.getProtocolStats()
}

// Stop HTTP monitoring
func (m *Monitor) Stop() {
	if m == nil {
		return
	}

	if usmstate.Get() == usmstate.Stopped {
		return
	}

	if m.telemetryStopChannel != nil {
		close(m.telemetryStopChannel)
	}

	m.processMonitor.Stop()

	ddebpf.RemoveNameMappings(m.ebpfProgram.Manager.Manager)

	m.ebpfProgram.Close()
	m.closeFilterFn()
	usmstate.Set(usmstate.Stopped)
}

// DumpMaps dumps the maps associated with the monitor
func (m *Monitor) DumpMaps(w io.Writer, maps ...string) error {
	return m.ebpfProgram.DumpMaps(w, maps...)
}

func (m *Monitor) startTelemetryReporter() {
	telemetry.SetStatsdClient(m.statsd)
	ticker := time.NewTicker(30 * time.Second)
	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				telemetry.ReportStatsd()
			case <-m.telemetryStopChannel:
				return
			}
		}
	}()
}
