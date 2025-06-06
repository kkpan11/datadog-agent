// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2025-present Datadog, Inc.

//go:build linux && !linux_bpf

package module

type nopNetworkCollector struct{}

func newNetworkCollector(_ *discoveryConfig) (networkCollector, error) {
	return &nopNetworkCollector{}, nil
}

func (c *nopNetworkCollector) close() {
}

func (c *nopNetworkCollector) getStats(_ pidSet) (map[uint32]NetworkStats, error) {
	return nil, nil
}
