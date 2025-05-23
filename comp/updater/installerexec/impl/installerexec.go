// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2025-present Datadog, Inc.

// Package installerexecimpl implements the installerexec component interface
package installerexecimpl

import (
	installerexec "github.com/DataDog/datadog-agent/comp/updater/installerexec/def"
	"github.com/DataDog/datadog-agent/pkg/fleet/installer/env"
	iexec "github.com/DataDog/datadog-agent/pkg/fleet/installer/exec"
	"github.com/DataDog/datadog-agent/pkg/util/defaultpaths"
)

// Requires defines the dependencies for the installerexec component
type Requires struct{}

// Provides defines the output of the installerexec component
type Provides struct {
	Comp installerexec.Component
}

// NewComponent creates a new installerexec component
func NewComponent(_ Requires) (Provides, error) {
	return Provides{
		Comp: iexec.NewInstallerExec(&env.Env{}, defaultpaths.InstallerPath),
	}, nil
}
