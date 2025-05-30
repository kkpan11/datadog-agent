// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

//go:build !windows

package config

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/DataDog/datadog-agent/comp/core/config"
)

func TestValidateShouldSucceedWithValidConfigs(t *testing.T) {
	validConfigs := []*LogsConfig{
		{Type: FileType, Path: "/var/log/foo.log"},
		{Type: TCPType, Port: 1234},
		{Type: UDPType, Port: 5678},
		{Type: DockerType},
		{Type: JournaldType, ProcessingRules: []*ProcessingRule{{Name: "foo", Type: ExcludeAtMatch, Pattern: ".*"}}},
	}

	for _, config := range validConfigs {
		err := config.Validate()
		assert.Nil(t, err)
	}
}

func TestValidateShouldFailWithInvalidConfigs(t *testing.T) {
	invalidConfigs := []*LogsConfig{
		{},
		{Type: FileType},
		{Type: TCPType},
		{Type: UDPType},
		{Type: DockerType, ProcessingRules: []*ProcessingRule{{Name: "foo"}}},
		{Type: DockerType, ProcessingRules: []*ProcessingRule{{Name: "foo", Type: "bar"}}},
		{Type: DockerType, ProcessingRules: []*ProcessingRule{{Name: "foo", Type: ExcludeAtMatch}}},
		{Type: DockerType, ProcessingRules: []*ProcessingRule{{Name: "foo", Pattern: ".*"}}},
		{Type: DockerType, ProcessingRules: []*ProcessingRule{{Type: ExcludeAtMatch, Pattern: ".*"}}},
		{Type: DockerType, ProcessingRules: []*ProcessingRule{{Type: ExcludeAtMatch}}},
		{Type: DockerType, ProcessingRules: []*ProcessingRule{{Pattern: ".*"}}},
	}

	for _, config := range invalidConfigs {
		err := config.Validate()
		assert.NotNil(t, err)
	}
}

func TestAutoMultilineEnabled(t *testing.T) {
	decode := func(cfg string) *LogsConfig {
		lc := LogsConfig{}
		json.Unmarshal([]byte(cfg), &lc)
		return &lc
	}

	mockConfig := config.NewMock(t)
	mockConfig.SetWithoutSource("logs_config.auto_multi_line_detection", false)
	assert.False(t, decode(`{"auto_multi_line_detection":false}`).AutoMultiLineEnabled(mockConfig))

	mockConfig.SetWithoutSource("logs_config.auto_multi_line_detection", true)
	assert.False(t, decode(`{"auto_multi_line_detection":false}`).AutoMultiLineEnabled(mockConfig))

	mockConfig.SetWithoutSource("logs_config.auto_multi_line_detection", true)
	assert.True(t, decode(`{}`).AutoMultiLineEnabled(mockConfig))

	mockConfig.SetWithoutSource("logs_config.auto_multi_line_detection", false)
	assert.True(t, decode(`{"auto_multi_line_detection":true}`).AutoMultiLineEnabled(mockConfig))

	mockConfig.SetWithoutSource("logs_config.auto_multi_line_detection", true)
	assert.True(t, decode(`{"auto_multi_line_detection":true}`).AutoMultiLineEnabled(mockConfig))

	mockConfig.SetWithoutSource("logs_config.auto_multi_line_detection", false)
	assert.False(t, decode(`{}`).AutoMultiLineEnabled(mockConfig))

}

func TestLegacyAutoMultilineEnabled(t *testing.T) {
	decode := func(cfg string) *LogsConfig {
		lc := LogsConfig{}
		json.Unmarshal([]byte(cfg), &lc)
		return &lc
	}
	mockConfig := config.NewMock(t)
	mockConfig.SetWithoutSource("logs_config.auto_multi_line_detection", false)
	assert.False(t, decode(`{"auto_multi_line_detection":false}`).LegacyAutoMultiLineEnabled(mockConfig))
	assert.False(t, decode(`{"auto_multi_line_detection":true}`).LegacyAutoMultiLineEnabled(mockConfig))

	mockConfig = config.NewMock(t)
	mockConfig.SetWithoutSource("logs_config.auto_multi_line_detection", true)
	assert.False(t, decode(`{"auto_multi_line_detection":false}`).LegacyAutoMultiLineEnabled(mockConfig))
	assert.False(t, decode(`{"auto_multi_line_detection":true}`).LegacyAutoMultiLineEnabled(mockConfig))

	assert.True(t, decode(`{"auto_multi_line_sample_size": 2}`).LegacyAutoMultiLineEnabled(mockConfig))
	assert.True(t, decode(`{"auto_multi_line_match_threshold": 0.4}`).LegacyAutoMultiLineEnabled(mockConfig))

	mockConfig = config.NewMock(t)
	mockConfig.SetWithoutSource("logs_config.force_auto_multi_line_detection_v1", true)
	assert.True(t, decode(`{}`).LegacyAutoMultiLineEnabled(mockConfig))

	mockConfig = config.NewMock(t)
	mockConfig.SetWithoutSource("logs_config.auto_multi_line_detection", true)
	mockConfig.SetWithoutSource("logs_config.auto_multi_line_default_sample_size", 10)
	assert.True(t, decode(`{}`).LegacyAutoMultiLineEnabled(mockConfig))

	mockConfig = config.NewMock(t)
	mockConfig.SetWithoutSource("logs_config.auto_multi_line_detection", true)
	mockConfig.SetWithoutSource("logs_config.auto_multi_line_default_match_timeout", 100)
	assert.True(t, decode(`{}`).LegacyAutoMultiLineEnabled(mockConfig))

	mockConfig = config.NewMock(t)
	mockConfig.SetWithoutSource("logs_config.auto_multi_line_detection", true)
	mockConfig.SetWithoutSource("logs_config.auto_multi_line_default_match_threshold", 501)
	assert.True(t, decode(`{}`).LegacyAutoMultiLineEnabled(mockConfig))

	mockConfig = config.NewMock(t)
	mockConfig.SetWithoutSource("logs_config.force_auto_multi_line_detection_v1", true)
	assert.False(t, decode(`{"auto_multi_line_detection":false}`).LegacyAutoMultiLineEnabled(mockConfig))

	mockConfig = config.NewMock(t)
	mockConfig.SetWithoutSource("logs_config.auto_multi_line_default_sample_size", 10)
	assert.True(t, decode(`{"auto_multi_line_detection":true}`).LegacyAutoMultiLineEnabled(mockConfig))

	mockConfig = config.NewMock(t)
	mockConfig.SetWithoutSource("logs_config.auto_multi_line_default_match_timeout", 100)
	assert.True(t, decode(`{"auto_multi_line_detection":true}`).LegacyAutoMultiLineEnabled(mockConfig))

	mockConfig = config.NewMock(t)
	mockConfig.SetWithoutSource("logs_config.auto_multi_line_default_match_threshold", 501)
	assert.True(t, decode(`{"auto_multi_line_detection":true}`).LegacyAutoMultiLineEnabled(mockConfig))
}

func TestConfigDump(t *testing.T) {
	config := LogsConfig{Type: FileType, Path: "/var/log/foo.log"}
	dump := config.Dump(true)
	assert.Contains(t, dump, `Path: "/var/log/foo.log",`)
}

func TestPublicJSON(t *testing.T) {
	config := LogsConfig{
		Type:     FileType,
		Path:     "/var/log/foo.log",
		Encoding: "utf-8",
		Service:  "foo",
		Tags:     []string{"foo:bar"},
		Source:   "bar",
	}
	ret, err := config.PublicJSON()
	assert.NoError(t, err)

	expectedJSON := `{"type":"file","path":"/var/log/foo.log","encoding":"utf-8","service":"foo","source":"bar","tags":["foo:bar"]}`
	assert.Equal(t, expectedJSON, string(ret))
}
