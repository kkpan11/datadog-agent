// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package installer defines a Pulumi component for installing the Datadog Installer on a remote host in the
// provisioning step.
package installer

import (
	"fmt"
	"github.com/DataDog/datadog-agent/test/new-e2e/tests/installer/windows/consts"
	"github.com/DataDog/datadog-agent/test/new-e2e/tests/windows/common/pipeline"
	"github.com/DataDog/test-infra-definitions/common"
	"github.com/DataDog/test-infra-definitions/common/config"
	"github.com/DataDog/test-infra-definitions/common/namer"
	"github.com/DataDog/test-infra-definitions/components"
	"github.com/DataDog/test-infra-definitions/components/command"
	remoteComp "github.com/DataDog/test-infra-definitions/components/remote"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"strings"
)

// Output is an object that models the output of the resource creation
// from the Component.
// See https://www.pulumi.com/docs/concepts/resources/components/#registering-component-outputs
type Output struct {
	components.JSONImporter
}

// Component is a Datadog Installer component.
// See https://www.pulumi.com/docs/concepts/resources/components/
type Component struct {
	pulumi.ResourceState
	components.Component

	namer namer.Namer
	Host  *remoteComp.Host `pulumi:"host"`
}

// Export exports the output of this component
func (h *Component) Export(ctx *pulumi.Context, out *Output) error {
	return components.Export(ctx, h, out)
}

// Configuration represents the Windows NewDefender configuration
type Configuration struct {
	URL                  string
	AgentUser            string
	CreateInstallerPaths bool
}

// Option is an optional function parameter type for Configuration options
type Option = func(*Configuration) error

// WithInstallURL specifies the URL to use to retrieve the Datadog Installer
func WithInstallURL(url string) func(*Configuration) error {
	return func(p *Configuration) error {
		p.URL = url
		return nil
	}
}

// WithAgentUser specifies the ddagentuser for the installation
func WithAgentUser(user string) func(*Configuration) error {
	return func(p *Configuration) error {
		p.AgentUser = user
		return nil
	}
}

// CreateInstallerPaths creates some directories that are normally created when using the
// PowerShell install script but are not when using the installer directly.
func CreateInstallerPaths() func(*Configuration) error {
	return func(p *Configuration) error {
		p.CreateInstallerPaths = true
		return nil
	}
}

// NewConfig creates a default config
func NewConfig(env config.Env, options ...Option) (*Configuration, error) {
	options = append(options, CreateInstallerPaths())
	if env.PipelineID() != "" {
		artifactURL, err := pipeline.GetPipelineArtifact(env.PipelineID(), pipeline.AgentS3BucketTesting, pipeline.DefaultMajorVersion, func(artifact string) bool {
			return strings.Contains(artifact, "datadog-installer") && strings.HasSuffix(artifact, ".msi")
		})
		if err != nil {
			return nil, err
		}
		options = append(options, WithInstallURL(artifactURL))
	}
	return common.ApplyOption(&Configuration{}, options)
}

// NewInstaller creates a new instance of an on-host Agent Installer
func NewInstaller(e config.Env, host *remoteComp.Host, options ...Option) (*Component, error) {

	params, err := NewConfig(e, options...)
	if err != nil {
		return nil, err
	}

	agentUserArg := ""
	if params.AgentUser != "" {
		agentUserArg = "DDAGENTUSER_NAME=" + params.AgentUser
	}

	hostInstaller, err := components.NewComponent(e, e.CommonNamer().ResourceName("datadog-installer"), func(comp *Component) error {
		comp.namer = e.CommonNamer().WithPrefix(consts.InstallerPackage)
		comp.Host = host

		createCmd := fmt.Sprintf("Exit (Start-Process -Wait msiexec -PassThru -ArgumentList '/qn /i %s %s').ExitCode", params.URL, agentUserArg)
		if params.CreateInstallerPaths {
			for _, p := range consts.InstallerConfigPaths {
				createCmd = fmt.Sprintf("New-Item -Path \"%s\" -ItemType Directory -Force; %s", p, createCmd)
			}
		}
		_, err = host.OS.Runner().Command(comp.namer.ResourceName("install"), &command.Args{
			Create: pulumi.String(createCmd).ToStringOutput(),
			Delete: pulumi.Sprintf(`
$installerList = Get-ItemProperty "HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\*" | Where-Object {$_.DisplayName -like 'Datadog Installer'}
if (($installerList | measure).Count -ne 1) {
    Write-Error "Could not find the Datadog Installer"
} else {
    cmd /c $installerList.UninstallString
}
`),
		}, pulumi.Parent(comp))
		if err != nil {
			return err
		}

		return nil
	}, pulumi.Parent(host), pulumi.DeletedWith(host))
	if err != nil {
		return nil, err
	}

	return hostInstaller, nil
}
