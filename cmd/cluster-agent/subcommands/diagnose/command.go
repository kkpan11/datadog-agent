// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

//go:build !windows && kubeapiserver

// Package diagnose implements 'cluster-agent diagnose'.
package diagnose

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"go.uber.org/fx"

	"github.com/DataDog/datadog-agent/cmd/cluster-agent/command"
	"github.com/DataDog/datadog-agent/comp/core"
	"github.com/DataDog/datadog-agent/comp/core/config"
	diagnose "github.com/DataDog/datadog-agent/comp/core/diagnose/def"
	"github.com/DataDog/datadog-agent/comp/core/diagnose/format"
	diagnosefx "github.com/DataDog/datadog-agent/comp/core/diagnose/fx"
	log "github.com/DataDog/datadog-agent/comp/core/log/def"
	"github.com/DataDog/datadog-agent/comp/core/secrets"
	"github.com/DataDog/datadog-agent/pkg/diagnose/connectivity"
	"github.com/DataDog/datadog-agent/pkg/util/fxutil"
)

// Commands returns a slice of subcommands for the 'cluster-agent' command.
func Commands(globalParams *command.GlobalParams) []*cobra.Command {
	cmd := &cobra.Command{
		Use:   "diagnose",
		Short: "Execute some connectivity diagnosis on your system",
		Long:  ``,
		RunE: func(_ *cobra.Command, _ []string) error {
			return fxutil.OneShot(run,
				fx.Supply(core.BundleParams{
					ConfigParams: config.NewClusterAgentParams(globalParams.ConfFilePath),
					SecretParams: secrets.NewEnabledParams(),
					LogParams:    log.ForOneShot(command.LoggerName, "off", true), // no need to show regular logs
				}),
				core.Bundle(),
				diagnosefx.Module(),
			)
		},
	}

	return []*cobra.Command{cmd}
}

func run(_ config.Component, diagnoseComponent diagnose.Component) error {
	suite := diagnose.Suites{
		diagnose.AutodiscoveryConnectivity: func(_ diagnose.Config) []diagnose.Diagnosis {
			return connectivity.DiagnoseMetadataAutodiscoveryConnectivity()
		},
	}

	config := diagnose.Config{Verbose: true}

	result, err := diagnoseComponent.RunLocalSuite(suite, config)

	if err != nil {
		return err
	}

	return format.Text(color.Output, config, result)
}
