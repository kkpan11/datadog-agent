// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2024-present Datadog, Inc.

//go:build test

package ciscosdwan

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/DataDog/datadog-agent/comp/aggregator/demultiplexer"
	"github.com/DataDog/datadog-agent/comp/aggregator/demultiplexer/demultiplexerimpl"
	"github.com/DataDog/datadog-agent/comp/core"
	"github.com/DataDog/datadog-agent/comp/core/autodiscovery/integration"
	"github.com/DataDog/datadog-agent/comp/forwarder/defaultforwarder"
	"github.com/DataDog/datadog-agent/pkg/aggregator/mocksender"
	checkid "github.com/DataDog/datadog-agent/pkg/collector/check/id"
	"github.com/DataDog/datadog-agent/pkg/collector/corechecks/network-devices/cisco-sdwan/client"
	"github.com/DataDog/datadog-agent/pkg/collector/corechecks/network-devices/cisco-sdwan/payload"
	"github.com/DataDog/datadog-agent/pkg/collector/corechecks/network-devices/cisco-sdwan/report"
	"github.com/DataDog/datadog-agent/pkg/util/fxutil"
)

type deps struct {
	fx.In
	Demultiplexer demultiplexer.Mock
}

func createDeps(t *testing.T) deps {
	return fxutil.Test[deps](t, demultiplexerimpl.MockModule(), defaultforwarder.MockModule(), core.MockBundle())
}

// mockTimeNow mocks time.Now
var mockTimeNow = func() time.Time {
	layout := "2006-01-02 15:04:05"
	str := "2024-02-26 10:22:00" // Set to have a one-hour uptime from API fixtures
	t, _ := time.Parse(layout, str)
	return t
}

func TestCiscoSDWANCheck(t *testing.T) {
	payload.TimeNow = mockTimeNow
	report.TimeNow = mockTimeNow

	apiMockServer := client.SetupMockAPIServer()
	defer apiMockServer.Close()

	deps := createDeps(t)
	chk := newCheck()
	senderManager := deps.Demultiplexer

	url := strings.TrimPrefix(apiMockServer.URL, "http://")

	// language=yaml
	rawInstanceConfig := []byte(`
vmanage_endpoint: ` + url + `
username: admin
password: 'test-password'
use_http: true
namespace: test
min_collection_interval: 180
collect_bfd_session_status: true
collect_hardware_status: true
collect_cloud_applications_metrics: true
collect_bgp_neighbor_states: true
`)

	// Use ID to ensure the mock sender gets registered
	id := checkid.BuildID(CheckName, integration.FakeConfigHash, rawInstanceConfig, []byte(``))
	sender := mocksender.NewMockSenderWithSenderManager(id, senderManager)
	sender.On("Gauge", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
	sender.On("MonotonicCount", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
	sender.On("GaugeWithTimestamp", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
	sender.On("CountWithTimestamp", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
	sender.On("EventPlatformEvent", mock.Anything, mock.Anything).Return()

	sender.On("Commit").Return()

	err := chk.Configure(senderManager, integration.FakeConfigHash, rawInstanceConfig, []byte(``), "test")
	require.NoError(t, err)

	assert.Equal(t, 3*time.Minute, chk.Interval())

	err = chk.Run()
	require.NoError(t, err)

	// Assert hardware metrics
	ts := float64(1709050342874) / 1000

	sender.AssertMetricWithTimestamp(t, "GaugeWithTimestamp", "cisco_sdwan.cpu.usage", 0.7, "", []string{"system_ip:10.10.1.5"}, ts)
	sender.AssertMetricWithTimestamp(t, "GaugeWithTimestamp", "cisco_sdwan.memory.usage", 15, "", []string{"system_ip:10.10.1.5"}, ts)
	sender.AssertMetricWithTimestamp(t, "GaugeWithTimestamp", "cisco_sdwan.disk.usage", 3.888901264100724, "", []string{"system_ip:10.10.1.5"}, ts)

	// Assert interface metrics
	ts = float64(1709049697985) / 1000
	tags := []string{"system_ip:10.10.1.22", "interface:GigabitEthernet3", "vpn_id:0"}

	sender.AssertMetricWithTimestamp(t, "CountWithTimestamp", "cisco_sdwan.interface.tx_bits", 32, "", tags, ts)
	sender.AssertMetricWithTimestamp(t, "CountWithTimestamp", "cisco_sdwan.interface.rx_bits", 184, "", tags, ts)
	sender.AssertMetricWithTimestamp(t, "GaugeWithTimestamp", "cisco_sdwan.interface.rx_bps", 10400, "", tags, ts)
	sender.AssertMetricWithTimestamp(t, "GaugeWithTimestamp", "cisco_sdwan.interface.tx_bps", 9800, "", tags, ts)
	sender.AssertMetricWithTimestamp(t, "GaugeWithTimestamp", "cisco_sdwan.interface.rx_bandwidth_usage", 0, "", tags, ts)
	sender.AssertMetricWithTimestamp(t, "GaugeWithTimestamp", "cisco_sdwan.interface.tx_bandwidth_usage", 0.8, "", tags, ts)
	sender.AssertMetricWithTimestamp(t, "CountWithTimestamp", "cisco_sdwan.interface.rx_errors", 2, "", tags, ts)
	sender.AssertMetricWithTimestamp(t, "CountWithTimestamp", "cisco_sdwan.interface.tx_errors", 506, "", tags, ts)
	sender.AssertMetricWithTimestamp(t, "CountWithTimestamp", "cisco_sdwan.interface.rx_drops", 6, "", tags, ts)
	sender.AssertMetricWithTimestamp(t, "CountWithTimestamp", "cisco_sdwan.interface.tx_drops", 3, "", tags, ts)

	// Assert uptime metrics
	sender.AssertMetric(t, "Gauge", "cisco_sdwan.device.uptime", 360000, "", []string{"device_vendor:cisco", "hostname:Manager", "system_ip:10.10.1.1", "site_id:101", "type:vmanage"})

	// Assert application aware routing metrics
	ts = float64(1709050725125) / 1000
	tags = []string{
		"system_ip:10.10.1.13",
		"remote_system_ip:10.10.1.11",
		"local_color:mpls",
		"remote_color:public-internet",
		"state:Up",
	}

	sender.AssertMetricWithTimestamp(t, "GaugeWithTimestamp", "cisco_sdwan.tunnel.status", 1, "", tags, ts)
	sender.AssertMetricWithTimestamp(t, "GaugeWithTimestamp", "cisco_sdwan.tunnel.latency", 202, "", tags, ts)
	sender.AssertMetricWithTimestamp(t, "GaugeWithTimestamp", "cisco_sdwan.tunnel.jitter", 0, "", tags, ts)
	sender.AssertMetricWithTimestamp(t, "GaugeWithTimestamp", "cisco_sdwan.tunnel.loss", 0.301, "", tags, ts)
	sender.AssertMetricWithTimestamp(t, "GaugeWithTimestamp", "cisco_sdwan.tunnel.qoe", 2, "", tags, ts)
	sender.AssertMetricWithTimestamp(t, "CountWithTimestamp", "cisco_sdwan.tunnel.rx_bits", 0, "", tags, ts)
	sender.AssertMetricWithTimestamp(t, "CountWithTimestamp", "cisco_sdwan.tunnel.tx_bits", 0, "", tags, ts)
	sender.AssertMetricWithTimestamp(t, "CountWithTimestamp", "cisco_sdwan.tunnel.rx_packets", 0, "", tags, ts)
	sender.AssertMetricWithTimestamp(t, "CountWithTimestamp", "cisco_sdwan.tunnel.tx_packets", 0, "", tags, ts)

	// Assert control-connection metrics
	sender.AssertMetric(t, "Gauge", "cisco_sdwan.control_connection.status", 1, "", []string{
		"device_vendor:cisco",
		"device_namespace:test",
		"hostname:Manager",
		"system_ip:10.10.1.1",
		"device_id:test:10.10.1.1",
		"device_hostname:Manager",
		"device_ip:10.10.1.1",
		"site_id:101",
		"type:vmanage",
		"remote_system_ip:10.10.1.3",
		"private_ip:10.10.20.80",
		"local_color:default",
		"remote_color:default",
		"peer_type:vbond",
		"state:up",
	})

	// Assert OMP Peer metrics
	sender.AssertMetric(t, "Gauge", "cisco_sdwan.omp_peer.status", 1, "", []string{"system_ip:10.10.1.5", "remote_system_ip:10.10.1.13", "legit:yes", "refresh:supported", "state:up"})

	// Assert BFD Session metrics
	sender.AssertMetric(t, "Gauge", "cisco_sdwan.bfd_session.status", 1, "", []string{"system_ip:10.10.1.11", "remote_system_ip:10.10.1.13", "local_color:public-internet", "remote_color:public-internet", "proto:ipsec", "state:up"})

	// Assert device counters metrics
	sender.AssertMetric(t, "MonotonicCount", "cisco_sdwan.crash.count", 0, "", []string{"system_ip:10.10.1.12"})
	sender.AssertMetric(t, "MonotonicCount", "cisco_sdwan.reboot.count", 3, "", []string{"system_ip:10.10.1.12"})

	// Assert device status metrics
	sender.AssertMetric(t, "Gauge", "cisco_sdwan.device.reachable", 1, "", []string{"device_vendor:cisco", "device_namespace:test", "hostname:Manager", "system_ip:10.10.1.1", "site_id:101", "type:vmanage"})

	// Assert hardware status metrics
	sender.AssertMetric(t, "Gauge", "cisco_sdwan.hardware.status_ok", 1, "", []string{"system_ip:10.10.1.11", "status:OK", "class:Fans", "item:Tray 0 fan", "dev_index:1"})

	// Assert cloud applications metrics
	ts = float64(1721819993833) / 1000
	tags = []string{
		"system_ip:10.10.1.13",
		"gateway_system_ip:10.10.1.11",
		"local_color:mpls",
		"remote_color:mpls",
		"interface:-",
		"exit_type:Remote",
		"application_group:amazon-group",
		"application:amazon_aws",
		"best_path:FALSE",
		"vpn_id:1",
	}

	sender.AssertMetricWithTimestamp(t, "GaugeWithTimestamp", "cisco_sdwan.application.latency", 404, "", tags, ts)
	sender.AssertMetricWithTimestamp(t, "GaugeWithTimestamp", "cisco_sdwan.application.loss", 0, "", tags, ts)
	sender.AssertMetricWithTimestamp(t, "GaugeWithTimestamp", "cisco_sdwan.application.qoe", 5, "", tags, ts)

	// Assert BGP neighbor metrics
	sender.AssertMetric(t, "Gauge", "cisco_sdwan.bgp.neighbor", 1, "", []string{"system_ip:10.10.1.11", "peer_state:established", "remote_as:2024", "neighbor:10.60.1.1", "vpn_id:1", "afi:ipv4-unicast"})

	// Assert metadata
	// language=json
	event := []byte(`
{
  "namespace": "test",
  "integration": "cisco-sdwan",
  "devices": [
    {
      "id": "test:10.10.1.1",
      "id_tags": [
        "device_namespace:test",
        "system_ip:10.10.1.1"
      ],
      "tags": [
        "device_vendor:cisco",
        "device_namespace:test",
        "hostname:Manager",
        "system_ip:10.10.1.1",
        "site_id:101",
        "type:vmanage",
        "device_ip:10.10.1.1",
        "device_hostname:Manager",
        "device_id:test:10.10.1.1",
        "source:cisco-sdwan"
      ],
      "ip_address": "10.10.1.1",
      "status": 1,
      "name": "Manager",
      "location": "SITE_101",
      "vendor": "cisco",
      "serial_number": "61FA4073B0169C46F4F498B8CA2C5C7A4A5510F9",
      "version": "20.12.1",
      "product_name": "vmanage",
      "model": "vmanage",
      "os_name": "next",
      "integration": "cisco-sdwan",
      "device_type": "sd-wan"
    }
  ],
  "interfaces": [
    {
      "device_id": "test:10.10.1.5",
      "id_tags": [
        "interface:system"
      ],
      "index": 3,
      "name": "system",
      "admin_status": 1,
      "oper_status": 1
    },
    {
      "device_id": "test:10.10.1.17",
      "id_tags": [
        "interface:GigabitEthernet4"
      ],
      "index": 4,
      "name": "GigabitEthernet4",
      "mac_address": "52:54:00:0b:6e:90",
      "admin_status": 1,
      "oper_status": 1
    }
  ],
  "ip_addresses": [
    {
      "interface_id": "test:10.10.1.17:4",
      "ip_address": "10.0.0.1",
	  "prefixlen": 24
    }
  ],
  "collect_timestamp": 1708942920
}
`)
	compactEvent := new(bytes.Buffer)
	err = json.Compact(compactEvent, event)
	assert.NoError(t, err)

	sender.AssertEventPlatformEvent(t, compactEvent.Bytes(), "network-devices-metadata")
}

func TestDefaultNotSentConfig(t *testing.T) {
	payload.TimeNow = mockTimeNow
	report.TimeNow = mockTimeNow

	apiMockServer := client.SetupMockAPIServer()
	defer apiMockServer.Close()

	deps := createDeps(t)
	chk := newCheck()
	senderManager := deps.Demultiplexer

	url := strings.TrimPrefix(apiMockServer.URL, "http://")

	// language=yaml
	rawInstanceConfig := []byte(`
vmanage_endpoint: ` + url + `
username: admin
password: 'test-password'
use_http: true
namespace: test
`)

	// Use ID to ensure the mock sender gets registered
	id := checkid.BuildID(CheckName, integration.FakeConfigHash, rawInstanceConfig, []byte(``))
	sender := mocksender.NewMockSenderWithSenderManager(id, senderManager)
	sender.On("Gauge", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
	sender.On("MonotonicCount", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
	sender.On("GaugeWithTimestamp", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
	sender.On("CountWithTimestamp", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
	sender.On("EventPlatformEvent", mock.Anything, mock.Anything).Return()

	sender.On("Commit").Return()

	err := chk.Configure(senderManager, integration.FakeConfigHash, rawInstanceConfig, []byte(``), "test")
	require.NoError(t, err)

	err = chk.Run()
	require.NoError(t, err)

	// BFD sessions not enabled by default
	sender.AssertNotCalled(t, "Gauge", "cisco_sdwan.bfd_session.status", mock.Anything, mock.Anything, mock.Anything, mock.Anything)

	// Hardware status not enabled by default
	sender.AssertNotCalled(t, "Gauge", "cisco_sdwan.hardware.status_ok", mock.Anything, mock.Anything, mock.Anything, mock.Anything)

	// Cloud applications not enabled by default
	sender.AssertNotCalled(t, "GaugeWithTimestamp", "cisco_sdwan.application.latency", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	sender.AssertNotCalled(t, "GaugeWithTimestamp", "cisco_sdwan.application.loss", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	sender.AssertNotCalled(t, "GaugeWithTimestamp", "cisco_sdwan.application.qoe", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestNoMetricsAndMetadataCollection(t *testing.T) {
	payload.TimeNow = mockTimeNow
	report.TimeNow = mockTimeNow

	apiMockServer := client.SetupMockAPIServer()
	defer apiMockServer.Close()

	deps := createDeps(t)
	chk := newCheck()
	senderManager := deps.Demultiplexer

	url := strings.TrimPrefix(apiMockServer.URL, "http://")

	// language=yaml
	rawInstanceConfig := []byte(`
vmanage_endpoint: ` + url + `
username: admin
password: 'test-password'
use_http: true
namespace: test
send_ndm_metadata: false
collect_hardware_metrics: false
collect_interface_metrics: false
collect_tunnel_metrics: false
collect_control_connection_metrics: false
collect_omp_peer_metrics: false
collect_device_counters_metrics: false
collect_bfd_session_status: false
collect_hardware_status: false
collect_cloud_applications_metrics: false
`)

	// Use ID to ensure the mock sender gets registered
	id := checkid.BuildID(CheckName, integration.FakeConfigHash, rawInstanceConfig, []byte(``))
	sender := mocksender.NewMockSenderWithSenderManager(id, senderManager)
	sender.On("Gauge", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
	sender.On("MonotonicCount", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
	sender.On("GaugeWithTimestamp", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
	sender.On("CountWithTimestamp", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
	sender.On("EventPlatformEvent", mock.Anything, mock.Anything).Return()

	sender.On("Commit").Return()

	err := chk.Configure(senderManager, integration.FakeConfigHash, rawInstanceConfig, []byte(``), "test")
	require.NoError(t, err)

	err = chk.Run()
	require.NoError(t, err)

	// Assert hardware metrics not sent
	sender.AssertNotCalled(t, "GaugeWithTimestamp", "cisco_sdwan.cpu.usage", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	sender.AssertNotCalled(t, "GaugeWithTimestamp", "cisco_sdwan.memory.usage", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	sender.AssertNotCalled(t, "GaugeWithTimestamp", "cisco_sdwan.disk.usage", mock.Anything, mock.Anything, mock.Anything, mock.Anything)

	// Assert interface metrics not sent
	sender.AssertNotCalled(t, "CountWithTimestamp", "cisco_sdwan.interface.tx_bits", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	sender.AssertNotCalled(t, "CountWithTimestamp", "cisco_sdwan.interface.rx_bits", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	sender.AssertNotCalled(t, "GaugeWithTimestamp", "cisco_sdwan.interface.rx_bps", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	sender.AssertNotCalled(t, "GaugeWithTimestamp", "cisco_sdwan.interface.tx_bps", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	sender.AssertNotCalled(t, "GaugeWithTimestamp", "cisco_sdwan.interface.rx_bandwidth_usage", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	sender.AssertNotCalled(t, "GaugeWithTimestamp", "cisco_sdwan.interface.tx_bandwidth_usage", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	sender.AssertNotCalled(t, "CountWithTimestamp", "cisco_sdwan.interface.rx_errors", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	sender.AssertNotCalled(t, "CountWithTimestamp", "cisco_sdwan.interface.tx_errors", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	sender.AssertNotCalled(t, "CountWithTimestamp", "cisco_sdwan.interface.rx_drops", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	sender.AssertNotCalled(t, "CountWithTimestamp", "cisco_sdwan.interface.tx_drops", mock.Anything, mock.Anything, mock.Anything, mock.Anything)

	// Assert uptime metrics not sent
	sender.AssertNotCalled(t, "Gauge", "cisco_sdwan.device.uptime", mock.Anything, mock.Anything, mock.Anything)

	// Assert application aware routing metrics not sent
	sender.AssertNotCalled(t, "GaugeWithTimestamp", "cisco_sdwan.tunnel.status", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	sender.AssertNotCalled(t, "GaugeWithTimestamp", "cisco_sdwan.tunnel.latency", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	sender.AssertNotCalled(t, "GaugeWithTimestamp", "cisco_sdwan.tunnel.jitter", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	sender.AssertNotCalled(t, "GaugeWithTimestamp", "cisco_sdwan.tunnel.loss", 0, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	sender.AssertNotCalled(t, "GaugeWithTimestamp", "cisco_sdwan.tunnel.qoe", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	sender.AssertNotCalled(t, "CountWithTimestamp", "cisco_sdwan.tunnel.rx_bits", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	sender.AssertNotCalled(t, "CountWithTimestamp", "cisco_sdwan.tunnel.tx_bits", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	sender.AssertNotCalled(t, "CountWithTimestamp", "cisco_sdwan.tunnel.rx_packets", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	sender.AssertNotCalled(t, "CountWithTimestamp", "cisco_sdwan.tunnel.tx_packets", mock.Anything, mock.Anything, mock.Anything, mock.Anything)

	// Assert control-connection metrics not sent
	sender.AssertNotCalled(t, "Gauge", "cisco_sdwan.control_connection.status", mock.Anything, mock.Anything, mock.Anything)

	// Assert OMP Peer metrics not sent
	sender.AssertNotCalled(t, "Gauge", "cisco_sdwan.omp_peer.status", mock.Anything, mock.Anything, mock.Anything)

	// Assert BFD Session metrics not sent
	sender.AssertNotCalled(t, "Gauge", "cisco_sdwan.bfd_session.status", mock.Anything, mock.Anything, mock.Anything)

	// Assert device counters metrics not sent
	sender.AssertNotCalled(t, "MonotonicCount", "cisco_sdwan.crash.count", mock.Anything, mock.Anything, mock.Anything)
	sender.AssertNotCalled(t, "MonotonicCount", "cisco_sdwan.reboot.count", mock.Anything, mock.Anything, mock.Anything)

	// Assert device status metrics
	sender.AssertNotCalled(t, "Gauge", "cisco_sdwan.device.reachable", mock.Anything, mock.Anything, mock.Anything)

	// Assert hardware status metrics
	sender.AssertNotCalled(t, "Gauge", "cisco_sdwan.hardware.status_ok", mock.Anything, mock.Anything, mock.Anything)

	// Assert cloud application metrics
	sender.AssertNotCalled(t, "GaugeWithTimestamp", "cisco_sdwan.application.latency", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	sender.AssertNotCalled(t, "GaugeWithTimestamp", "cisco_sdwan.application.loss", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	sender.AssertNotCalled(t, "GaugeWithTimestamp", "cisco_sdwan.application.qoe", mock.Anything, mock.Anything, mock.Anything, mock.Anything)

	// Assert BGP peer metrics
	sender.AssertNotCalled(t, "Gauge", "cisco_sdwan.bgp.neighbor", mock.Anything, mock.Anything, mock.Anything)

	sender.AssertNotCalled(t, "EventPlatformEvent", mock.Anything, mock.Anything)
}
