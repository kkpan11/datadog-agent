// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2024-present Datadog, Inc.

// Package agent contains a process-agent component
package agent

// team: container-experiences

// Component is the process agent component type
type Component interface {
	// Enabled returns whether the process agent is enabled
	Enabled() bool
}
