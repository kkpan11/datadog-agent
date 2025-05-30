// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package languagedetection determines the language that a process is written or compiled in.
package languagedetection

import (
	"bytes"
	"io"
	"net/http"
	"regexp"
	"runtime"
	"strings"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/DataDog/datadog-agent/pkg/config/model"
	"github.com/DataDog/datadog-agent/pkg/languagedetection/internal/detectors"
	"github.com/DataDog/datadog-agent/pkg/languagedetection/languagemodels"
	languagepb "github.com/DataDog/datadog-agent/pkg/proto/pbgo/languagedetection"
	sysprobeclient "github.com/DataDog/datadog-agent/pkg/system-probe/api/client"
	sysconfig "github.com/DataDog/datadog-agent/pkg/system-probe/config"
	"github.com/DataDog/datadog-agent/pkg/telemetry"
	"github.com/DataDog/datadog-agent/pkg/util/log"
)

var cliDetectors = []languagemodels.Detector{
	detectors.JRubyDetector{},
}

type languageFromCLI struct {
	name      languagemodels.LanguageName
	validator func(exe string) bool
}

var (
	rubyPattern = regexp.MustCompile(`^ruby\d+\.\d+$`)
	phpPattern  = regexp.MustCompile(`^php(?:-fpm)?\d(?:\.\d)?$`)
)

func matchesRubyPrefix(exe string) bool {
	return rubyPattern.MatchString(exe)
}

func matchesJavaPrefix(exe string) bool {
	return exe != "javac"
}

func matchesPHPPrefix(exe string) bool {
	return phpPattern.MatchString(exe)
}

// knownPrefixes maps languages names to their prefix
var knownPrefixes = map[string]languageFromCLI{
	"python": {name: languagemodels.Python},
	"java":   {name: languagemodels.Java, validator: matchesJavaPrefix},
	"ruby":   {name: languagemodels.Ruby, validator: matchesRubyPrefix},
	"php":    {name: languagemodels.PHP, validator: matchesPHPPrefix},
}

// exactMatches maps an exact exe name match to a prefix
var exactMatches = map[string]languageFromCLI{
	"py":      {name: languagemodels.Python},
	"python":  {name: languagemodels.Python},
	"java":    {name: languagemodels.Java},
	"npm":     {name: languagemodels.Node},
	"node":    {name: languagemodels.Node},
	"dotnet":  {name: languagemodels.Dotnet},
	"ruby":    {name: languagemodels.Ruby},
	"rubyw":   {name: languagemodels.Ruby},
	"php":     {name: languagemodels.PHP},
	"php-fpm": {name: languagemodels.PHP},
}

// languageNameFromCmdline returns a process's language from its command.
// If the language is not detected, languagemodels.Unknown is returned.
func languageNameFromCommand(command string) languagemodels.LanguageName {
	// First check to see if there is an exact match
	if lang, ok := exactMatches[command]; ok {
		return lang.name
	}

	for prefix, language := range knownPrefixes {
		if strings.HasPrefix(command, prefix) {
			if language.validator != nil {
				isValidResult := language.validator(command)
				if !isValidResult {
					continue
				}
			}
			return language.name
		}
	}

	return languagemodels.Unknown
}

const subsystem = "language_detection"

// prometheus.DefBuckets, converted to milliseconds.
var buckets = []float64{5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000}

var (
	detectLanguageRuntimeMs = telemetry.NewHistogram(subsystem, "detect_language_ms", nil,
		"The amount of time it took for the call to DetectLanguage to complete.", buckets)
	systemProbeLanguageDetectionMs = telemetry.NewHistogram(subsystem, "system_probe_rpc_ms", nil,
		"The amount of time it took for the process agent to message the system probe.", buckets)
)

// DetectLanguage uses a combination of commandline parsing and binary analysis to detect a process' language
func DetectLanguage(procs []languagemodels.Process, sysprobeConfig model.Reader) []*languagemodels.Language {
	detectLanguageStart := time.Now()
	defer func() {
		detectLanguageRuntimeMs.Observe(float64(time.Since(detectLanguageStart).Milliseconds()))
	}()

	langs := make([]*languagemodels.Language, len(procs))
	unknownPids := make([]int32, 0, len(procs))
	langsToModify := make(map[int32]*languagemodels.Language, len(procs))
	for i, proc := range procs {
		// Language-specific detectors should precede matches on the command/exe
		for _, detector := range cliDetectors {
			lang, err := detector.DetectLanguage(proc)
			if err != nil {
				log.Warnf("unable to detect language for process %d: %s", proc.GetPid(), err)
				continue
			}

			if lang.Name != languagemodels.Unknown {
				langs[i] = &lang
				break
			}
		}

		if langs[i] != nil {
			continue
		}

		exe := getExe(proc.GetCmdline())
		languageName := languageNameFromCommand(exe)
		if languageName == languagemodels.Unknown {
			languageName = languageNameFromCommand(proc.GetCommand())
		}
		lang := &languagemodels.Language{Name: languageName}
		langs[i] = lang
		if lang.Name == languagemodels.Unknown {
			unknownPids = append(unknownPids, proc.GetPid())
			langsToModify[proc.GetPid()] = lang
		}
	}

	if privilegedLanguageDetectionEnabled(sysprobeConfig) {
		rpcStart := time.Now()
		defer func() {
			systemProbeLanguageDetectionMs.Observe(float64(time.Since(rpcStart).Milliseconds()))
		}()

		log.Trace("[language detection] Requesting language from system probe")
		sysprobeClient := sysprobeclient.Get(sysprobeConfig.GetString("system_probe_config.sysprobe_socket"))
		privilegedLangs, err := detectLanguage(sysprobeClient, unknownPids)
		if err != nil {
			log.Warn("[language detection] Failed to request language:", err)
			return langs
		}

		for i, pid := range unknownPids {
			*langsToModify[pid] = privilegedLangs[i]
		}
	}
	return langs
}

func detectLanguage(client *http.Client, pids []int32) ([]languagemodels.Language, error) {
	procs := make([]*languagepb.Process, len(pids))
	for i, pid := range pids {
		procs[i] = &languagepb.Process{Pid: pid}
	}
	reqBytes, err := proto.Marshal(&languagepb.DetectLanguageRequest{Processes: procs})
	if err != nil {
		return nil, err
	}

	url := sysprobeclient.ModuleURL(sysconfig.LanguageDetectionModule, "/detect")
	req, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var resProto languagepb.DetectLanguageResponse
	err = proto.Unmarshal(resBody, &resProto)
	if err != nil {
		return nil, err
	}

	langs := make([]languagemodels.Language, len(pids))
	for i, lang := range resProto.Languages {
		langs[i] = languagemodels.Language{
			Name:    languagemodels.LanguageName(lang.Name),
			Version: lang.Version,
		}
	}
	return langs, nil
}

func privilegedLanguageDetectionEnabled(sysProbeConfig model.Reader) bool {
	if sysProbeConfig == nil {
		return false
	}

	// System probe language detection only works on linux operating systems for the moment.
	if runtime.GOOS != "linux" {
		return false
	}

	return sysProbeConfig.GetBool("system_probe_config.language_detection.enabled")
}
