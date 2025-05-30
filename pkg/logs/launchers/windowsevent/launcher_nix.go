// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

//go:build !windows

// Package windowsevent is not supported on non-windows platforms
package windowsevent

import (
	auditor "github.com/DataDog/datadog-agent/comp/logs/auditor/def"
	"github.com/DataDog/datadog-agent/pkg/logs/launchers"
	"github.com/DataDog/datadog-agent/pkg/logs/pipeline"
	"github.com/DataDog/datadog-agent/pkg/logs/tailers"
)

// Launcher is a stub for non-windows platforms
type Launcher struct{}

// Start is a stub for non-windows platforms
func (l *Launcher) Start(_ launchers.SourceProvider, _ pipeline.Provider, _ auditor.Registry, _ *tailers.TailerTracker) {
}

// Stop is a stub for non-windows platforms
func (l *Launcher) Stop() {}

// NewLauncher is a stub for non-windows platforms
func NewLauncher() *Launcher {
	return &Launcher{}
}
