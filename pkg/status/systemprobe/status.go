// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package systemprobe fetch information about the system probe
package systemprobe

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/DataDog/datadog-agent/comp/core/status"
	"github.com/DataDog/datadog-agent/comp/core/sysprobeconfig"
	sysprobeclient "github.com/DataDog/datadog-agent/pkg/system-probe/api/client"
)

// GetStatus returns the expvar stats of the system probe
func GetStatus(stats map[string]interface{}, socketPath string) {
	client := sysprobeclient.Get(socketPath)
	systemProbeDetails, err := getStats(client)
	if err != nil {
		stats["systemProbeStats"] = map[string]interface{}{
			"Errors": fmt.Sprintf("issue querying stats from system probe: %v", err),
		}
		return
	}
	stats["systemProbeStats"] = systemProbeDetails
}

func getStats(client *http.Client) (map[string]interface{}, error) {
	url := sysprobeclient.DebugURL("/stats")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("conn request failed: url: %s, status code: %d", req.URL, resp.StatusCode)
	}

	body, err := sysprobeclient.ReadAllResponseBody(resp)
	if err != nil {
		return nil, err
	}

	stats := make(map[string]interface{})
	err = json.Unmarshal(body, &stats)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

// Provider provides the functionality to populate the status output
type Provider struct {
	SocketPath string
}

// GetProvider if system probe is enabled returns status.Provider otherwise returns nil
func GetProvider(config sysprobeconfig.Component) status.Provider {
	systemProbeConfig := config.SysProbeObject()

	if systemProbeConfig.Enabled {
		return Provider{
			SocketPath: systemProbeConfig.SocketAddress,
		}
	}

	return nil
}

//go:embed status_templates
var templatesFS embed.FS

// Name returns the name
func (Provider) Name() string {
	return "System Probe"
}

// Section return the section
func (Provider) Section() string {
	return "System Probe"
}

// JSON populates the status map
func (p Provider) JSON(_ bool, stats map[string]interface{}) error {
	GetStatus(stats, p.SocketPath)

	return nil
}

// Text renders the text output
func (p Provider) Text(_ bool, buffer io.Writer) error {
	return status.RenderText(templatesFS, "systemprobe.tmpl", buffer, p.getStatusInfo())
}

// HTML renders the html output
func (p Provider) HTML(_ bool, _ io.Writer) error {
	return nil
}

func (p Provider) getStatusInfo() map[string]interface{} {
	stats := make(map[string]interface{})

	GetStatus(stats, p.SocketPath)

	return stats
}
