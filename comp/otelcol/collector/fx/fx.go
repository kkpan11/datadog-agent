// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

//go:build otlp

// Package fx creates the modules for fx
package fx

import (
	collector "github.com/DataDog/datadog-agent/comp/otelcol/collector/def"
	collectorimpl "github.com/DataDog/datadog-agent/comp/otelcol/collector/impl"
	"github.com/DataDog/datadog-agent/pkg/util/fxutil"
	"go.uber.org/fx"
)

// team: opentelemetry-agent

// Module for OTel Agent
func Module(params collectorimpl.Params) fxutil.Module {
	return fxutil.Component(
		fx.Supply(params),
		fxutil.ProvideComponentConstructor(
			collectorimpl.NewComponent,
		),
		fxutil.ProvideOptional[collector.Component](),
	)
}

// ModuleNoAgent for OTel Agent with no Agent functionalities
func ModuleNoAgent() fxutil.Module {
	return fxutil.Component(
		fxutil.ProvideComponentConstructor(
			collectorimpl.NewComponentNoAgent,
		),
		fxutil.ProvideOptional[collector.Component](),
	)
}
