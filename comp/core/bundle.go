// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package core implements the "core" bundle, providing services common to all
// agent flavors and binaries.
//
// The constituent components serve as utilities and are mostly independent of
// one another.  Other components should depend on any components they need.
//
// This bundle does not depend on any other bundles.
package core

import (
	"go.uber.org/fx"

	"github.com/DataDog/datadog-agent/comp/core/config"
	"github.com/DataDog/datadog-agent/comp/core/hostname/hostnameimpl"
	log "github.com/DataDog/datadog-agent/comp/core/log/def"
	logfx "github.com/DataDog/datadog-agent/comp/core/log/fx"
	"github.com/DataDog/datadog-agent/comp/core/pid/pidimpl"
	"github.com/DataDog/datadog-agent/comp/core/secrets"
	"github.com/DataDog/datadog-agent/comp/core/secrets/secretsimpl"
	"github.com/DataDog/datadog-agent/comp/core/sysprobeconfig/sysprobeconfigimpl"
	"github.com/DataDog/datadog-agent/comp/core/telemetry/telemetryimpl"
	"github.com/DataDog/datadog-agent/pkg/util/fxutil"
	"github.com/DataDog/datadog-agent/pkg/util/option"
)

// team: agent-runtimes

// Bundle defines the fx options for this bundle.
func Bundle() fxutil.BundleOptions {
	return fxutil.Bundle(
		// As `config.Module` expects `config.Params` as a parameter, it is require to define how to get `config.Params` from `BundleParams`.
		fx.Provide(func(params BundleParams) config.Params { return params.ConfigParams }),
		config.Module(),
		fx.Provide(func(params BundleParams) log.Params { return params.LogParams }),
		logfx.Module(),
		fx.Provide(func(params BundleParams) sysprobeconfigimpl.Params { return params.SysprobeConfigParams }),
		secretsimpl.Module(),
		fx.Provide(func(params BundleParams) secrets.Params { return params.SecretParams }),
		fx.Provide(func(secrets secrets.Component) option.Option[secrets.Component] { return option.New(secrets) }),
		sysprobeconfigimpl.Module(),
		telemetryimpl.Module(),
		hostnameimpl.Module(),
		pidimpl.Module(), // You must supply pidimpl.NewParams in order to use it
	)
}
