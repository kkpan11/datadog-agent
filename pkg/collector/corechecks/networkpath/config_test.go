// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023-present Datadog, Inc.

package networkpath

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/DataDog/datadog-agent/comp/core/autodiscovery/integration"
	configmock "github.com/DataDog/datadog-agent/pkg/config/mock"
	"github.com/DataDog/datadog-agent/pkg/config/setup"
	"github.com/DataDog/datadog-agent/pkg/networkpath/payload"
)

func TestNewCheckConfig(t *testing.T) {
	mockConfig := configmock.New(t)
	mockConfig.SetWithoutSource("network_devices.namespace", "my-namespace")
	tests := []struct {
		name           string
		rawInstance    integration.Data
		rawInitConfig  integration.Data
		expectedConfig *CheckConfig
		expectedError  string
	}{
		{
			name: "valid minimal config",
			rawInstance: []byte(`
hostname: 1.2.3.4
`),
			rawInitConfig: []byte(``),
			expectedConfig: &CheckConfig{
				DestHostname:          "1.2.3.4",
				MinCollectionInterval: time.Duration(60) * time.Second,
				Namespace:             "my-namespace",
				Timeout:               setup.DefaultNetworkPathTimeout * time.Millisecond,
				MaxTTL:                setup.DefaultNetworkPathMaxTTL,
			},
		},
		{
			name:          "invalid raw instance json",
			rawInstance:   []byte(`{{{`),
			expectedError: "invalid instance config",
		},
		{
			name:          "invalid raw instance json",
			rawInstance:   []byte(`hostname: 1.2.3.4`),
			rawInitConfig: []byte(`{{{`),
			expectedError: "invalid init_config",
		},
		{
			name: "invalid min_collection_interval",
			rawInstance: []byte(`
hostname: 1.2.3.4
min_collection_interval: -1
`),
			expectedError: "min collection interval must be > 0",
		},
		{
			name: "min_collection_interval from instance",
			rawInstance: []byte(`
hostname: 1.2.3.4
min_collection_interval: 42
`),
			rawInitConfig: []byte(`
min_collection_interval: 10
`),
			expectedConfig: &CheckConfig{
				DestHostname:          "1.2.3.4",
				MinCollectionInterval: time.Duration(42) * time.Second,
				Namespace:             "my-namespace",
				Timeout:               setup.DefaultNetworkPathTimeout * time.Millisecond,
				MaxTTL:                setup.DefaultNetworkPathMaxTTL,
			},
		},
		{
			name: "min_collection_interval from init_config",
			rawInstance: []byte(`
hostname: 1.2.3.4
`),
			rawInitConfig: []byte(`
min_collection_interval: 10
`),
			expectedConfig: &CheckConfig{
				DestHostname:          "1.2.3.4",
				MinCollectionInterval: time.Duration(10) * time.Second,
				Namespace:             "my-namespace",
				Timeout:               setup.DefaultNetworkPathTimeout * time.Millisecond,
				MaxTTL:                setup.DefaultNetworkPathMaxTTL,
			},
		},
		{
			name: "min_collection_interval default value",
			rawInstance: []byte(`
hostname: 1.2.3.4
`),
			expectedConfig: &CheckConfig{
				DestHostname:          "1.2.3.4",
				MinCollectionInterval: time.Duration(1) * time.Minute,
				Namespace:             "my-namespace",
				Timeout:               setup.DefaultNetworkPathTimeout * time.Millisecond,
				MaxTTL:                setup.DefaultNetworkPathMaxTTL,
			},
		},
		{
			name: "source and destination service config",
			rawInstance: []byte(`
hostname: 1.2.3.4
source_service: service-a
destination_service: service-b
`),
			rawInitConfig: []byte(``),
			expectedConfig: &CheckConfig{
				DestHostname:          "1.2.3.4",
				SourceService:         "service-a",
				DestinationService:    "service-b",
				MinCollectionInterval: time.Duration(60) * time.Second,
				Namespace:             "my-namespace",
				Timeout:               setup.DefaultNetworkPathTimeout * time.Millisecond,
				MaxTTL:                setup.DefaultNetworkPathMaxTTL,
			},
		},
		{
			name: "lower case protocol",
			rawInstance: []byte(`
hostname: 1.2.3.4
protocol: udp
`),
			rawInitConfig: []byte(``),
			expectedConfig: &CheckConfig{
				DestHostname:          "1.2.3.4",
				MinCollectionInterval: time.Duration(60) * time.Second,
				Namespace:             "my-namespace",
				Protocol:              payload.ProtocolUDP,
				Timeout:               setup.DefaultNetworkPathTimeout * time.Millisecond,
				MaxTTL:                setup.DefaultNetworkPathMaxTTL,
			},
		},
		{
			name: "lower case protocol",
			rawInstance: []byte(`
hostname: 1.2.3.4
protocol: UDP
`),
			rawInitConfig: []byte(``),
			expectedConfig: &CheckConfig{
				DestHostname:          "1.2.3.4",
				MinCollectionInterval: time.Duration(60) * time.Second,
				Namespace:             "my-namespace",
				Protocol:              payload.ProtocolUDP,
				Timeout:               setup.DefaultNetworkPathTimeout * time.Millisecond,
				MaxTTL:                setup.DefaultNetworkPathMaxTTL,
			},
		},
		{
			name: "lower case protocol",
			rawInstance: []byte(`
hostname: 1.2.3.4
protocol: TCP
`),
			rawInitConfig: []byte(``),
			expectedConfig: &CheckConfig{
				DestHostname:          "1.2.3.4",
				MinCollectionInterval: time.Duration(60) * time.Second,
				Namespace:             "my-namespace",
				Protocol:              payload.ProtocolTCP,
				Timeout:               setup.DefaultNetworkPathTimeout * time.Millisecond,
				MaxTTL:                setup.DefaultNetworkPathMaxTTL,
			},
		},
		{
			name: "timeout from instance config",
			rawInstance: []byte(`
hostname: 1.2.3.4
timeout: 50000
min_collection_interval: 42
`),
			rawInitConfig: []byte(`
min_collection_interval: 10
`),
			expectedConfig: &CheckConfig{
				DestHostname:          "1.2.3.4",
				MinCollectionInterval: time.Duration(42) * time.Second,
				Namespace:             "my-namespace",
				Timeout:               50000 * time.Millisecond,
				MaxTTL:                setup.DefaultNetworkPathMaxTTL,
			},
		},
		{
			name: "timeout from instance config preferred over init config",
			rawInstance: []byte(`
hostname: 1.2.3.4
timeout: 50000
min_collection_interval: 42
`),
			rawInitConfig: []byte(`
min_collection_interval: 10
timeout: 70000
`),
			expectedConfig: &CheckConfig{
				DestHostname:          "1.2.3.4",
				MinCollectionInterval: time.Duration(42) * time.Second,
				Namespace:             "my-namespace",
				Timeout:               50000 * time.Millisecond,
				MaxTTL:                setup.DefaultNetworkPathMaxTTL,
			},
		},
		{
			name: "timeout from init config",
			rawInstance: []byte(`
hostname: 1.2.3.4
min_collection_interval: 42
`),
			rawInitConfig: []byte(`
min_collection_interval: 10
timeout: 70000
`),
			expectedConfig: &CheckConfig{
				DestHostname:          "1.2.3.4",
				MinCollectionInterval: time.Duration(42) * time.Second,
				Namespace:             "my-namespace",
				Timeout:               70000 * time.Millisecond,
				MaxTTL:                setup.DefaultNetworkPathMaxTTL,
			},
		},
		{
			name: "default timeout",
			rawInstance: []byte(`
hostname: 1.2.3.4
min_collection_interval: 42
`),
			rawInitConfig: []byte(`
min_collection_interval: 10
`),
			expectedConfig: &CheckConfig{
				DestHostname:          "1.2.3.4",
				MinCollectionInterval: time.Duration(42) * time.Second,
				Namespace:             "my-namespace",
				Timeout:               setup.DefaultNetworkPathTimeout * time.Millisecond,
				MaxTTL:                setup.DefaultNetworkPathMaxTTL,
			},
		},
		{
			name: "negative timeout returns an error",
			rawInstance: []byte(`
hostname: 1.2.3.4
min_collection_interval: 42
`),
			rawInitConfig: []byte(`
min_collection_interval: 10
timeout: -1
`),
			expectedError: "timeout must be > 0",
		},
		{
			name: "maxTTL from instance config",
			rawInstance: []byte(`
hostname: 1.2.3.4
max_ttl: 50
min_collection_interval: 42
`),
			rawInitConfig: []byte(`
min_collection_interval: 10
`),
			expectedConfig: &CheckConfig{
				DestHostname:          "1.2.3.4",
				MinCollectionInterval: time.Duration(42) * time.Second,
				Namespace:             "my-namespace",
				Timeout:               setup.DefaultNetworkPathTimeout * time.Millisecond,
				MaxTTL:                50,
			},
		},
		{
			name: "maxTTL from instance config preferred over init config",
			rawInstance: []byte(`
hostname: 1.2.3.4
max_ttl: 50
min_collection_interval: 42
`),
			rawInitConfig: []byte(`
min_collection_interval: 10
max_ttl: 64
`),
			expectedConfig: &CheckConfig{
				DestHostname:          "1.2.3.4",
				MinCollectionInterval: time.Duration(42) * time.Second,
				Namespace:             "my-namespace",
				Timeout:               setup.DefaultNetworkPathTimeout * time.Millisecond,
				MaxTTL:                50,
			},
		},
		{
			name: "maxTTL from init config",
			rawInstance: []byte(`
hostname: 1.2.3.4
min_collection_interval: 42
`),
			rawInitConfig: []byte(`
min_collection_interval: 10
max_ttl: 64
`),
			expectedConfig: &CheckConfig{
				DestHostname:          "1.2.3.4",
				MinCollectionInterval: time.Duration(42) * time.Second,
				Namespace:             "my-namespace",
				Timeout:               setup.DefaultNetworkPathTimeout * time.Millisecond,
				MaxTTL:                64,
			},
		},
		{
			name: "overriding the TCP method",
			rawInstance: []byte(`
hostname: 1.2.3.4
protocol: tcp
tcp_method: sack
`),
			rawInitConfig: []byte(``),
			expectedConfig: &CheckConfig{
				DestHostname:          "1.2.3.4",
				MinCollectionInterval: time.Duration(60) * time.Second,
				Namespace:             "my-namespace",
				Protocol:              payload.ProtocolTCP,
				Timeout:               setup.DefaultNetworkPathTimeout * time.Millisecond,
				MaxTTL:                setup.DefaultNetworkPathMaxTTL,
				TCPMethod:             payload.TCPConfigSACK,
			},
		},
		{
			name: "TCP method converts to lower case",
			rawInstance: []byte(`
hostname: 1.2.3.4
protocol: tcp
tcp_method: prefer_SACK
`),
			rawInitConfig: []byte(``),
			expectedConfig: &CheckConfig{
				DestHostname:          "1.2.3.4",
				MinCollectionInterval: time.Duration(60) * time.Second,
				Namespace:             "my-namespace",
				Protocol:              payload.ProtocolTCP,
				Timeout:               setup.DefaultNetworkPathTimeout * time.Millisecond,
				MaxTTL:                setup.DefaultNetworkPathMaxTTL,
				TCPMethod:             payload.TCPConfigPreferSACK,
			},
		},
		{
			name: "Enabling TCP SYN compatibility mode",
			rawInstance: []byte(`
hostname: 1.2.3.4
protocol: tcp
tcp_syn_paris_traceroute_mode: true
`),
			rawInitConfig: []byte(``),
			expectedConfig: &CheckConfig{
				DestHostname:              "1.2.3.4",
				MinCollectionInterval:     time.Duration(60) * time.Second,
				Namespace:                 "my-namespace",
				Protocol:                  payload.ProtocolTCP,
				Timeout:                   setup.DefaultNetworkPathTimeout * time.Millisecond,
				MaxTTL:                    setup.DefaultNetworkPathMaxTTL,
				TCPSynParisTracerouteMode: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := NewCheckConfig(tt.rawInstance, tt.rawInitConfig)
			assert.Equal(t, tt.expectedConfig, config)
			if tt.expectedError != "" {
				assert.ErrorContains(t, err, tt.expectedError)
			}
		})
	}
}

func TestFirstNonZero(t *testing.T) {
	tests := []struct {
		name          string
		values        []time.Duration
		expectedValue time.Duration
	}{
		{
			name:          "no value",
			expectedValue: 0,
		},
		{
			name: "one value",
			values: []time.Duration{
				time.Duration(10) * time.Second,
			},
			expectedValue: time.Duration(10) * time.Second,
		},
		{
			name: "multiple values - select first",
			values: []time.Duration{
				time.Duration(10) * time.Second,
				time.Duration(20) * time.Second,
				time.Duration(30) * time.Second,
			},
			expectedValue: time.Duration(10) * time.Second,
		},
		{
			name: "multiple values - select second",
			values: []time.Duration{
				time.Duration(0) * time.Second,
				time.Duration(20) * time.Second,
				time.Duration(30) * time.Second,
			},
			expectedValue: time.Duration(20) * time.Second,
		},
		{
			name: "multiple values - select third",
			values: []time.Duration{
				time.Duration(0) * time.Second,
				time.Duration(0) * time.Second,
				time.Duration(30) * time.Second,
			},
			expectedValue: time.Duration(30) * time.Second,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedValue, firstNonZero(tt.values...))
		})
	}
}
