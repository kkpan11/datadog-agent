// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

//go:build windows

package parser

import (
	"github.com/DataDog/datadog-agent/pkg/util/winutil"
)

// A basic mock for `winutil.SCMMonitor`.
type mockableSCM interface {
	GetServiceInfo(pid uint64) (*winutil.ServiceInfo, error)
}

// scmReader is a cross-platform compatibility wrapper around `winutil.SCMMonitor`.
// The non-windows version does nothing, and instead only exists so that we don't get compile errors.
type scmReader struct {
	scmMonitor mockableSCM
}

func newSCMReader() *scmReader {
	return &scmReader{
		scmMonitor: winutil.GetServiceMonitor(),
	}
}

func (s *scmReader) getServiceInfo(pid uint64) (*WindowsServiceInfo, error) {
	monitorServiceInfo, err := s.scmMonitor.GetServiceInfo(pid)
	if err != nil {
		return nil, err
	}

	if monitorServiceInfo == nil {
		return nil, nil
	}

	return &WindowsServiceInfo{
		ServiceName: monitorServiceInfo.ServiceName,
		DisplayName: monitorServiceInfo.DisplayName,
	}, nil
}
