// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

//go:build linux_bpf

// Package util providers utility file functions to dynamic instrumentation
package util

import (
	"os"
	"time"

	"github.com/DataDog/datadog-agent/pkg/util/log"
)

// FileWatcher is used to track updates to a particular filepath
type FileWatcher struct {
	filePath string
	stop     chan bool
}

// NewFileWatcher creates a FileWatcher to track updates to a specified file
func NewFileWatcher(filePath string) *FileWatcher {
	return &FileWatcher{
		filePath: filePath,
		stop:     make(chan bool),
	}
}

func (fw *FileWatcher) readFile() ([]byte, error) {
	content, err := os.ReadFile(fw.filePath)
	if err != nil {
		return nil, err
	}
	return content, nil
}

// Watch watches the target file for changes and returns a channel that will receive
// the file's content whenever it changes.
// The initial implementation used fsnotify, but this was losing update events when running
// e2e tests - this simpler implementation behaves as expected, even if it's less efficient.
// Since this is meant to be used only for testing and development, it's fine to keep this
// implementation.
func (fw *FileWatcher) Watch() (<-chan []byte, error) {
	updateChan := make(chan []byte)
	prevContent := []byte{}
	ticker := time.NewTicker(100 * time.Millisecond)
	go func() {
		defer ticker.Stop()
		defer close(updateChan)
		for {
			select {
			case <-ticker.C:
				content, err := fw.readFile()
				if err != nil {
					log.Errorf("Error reading file %s: %s", fw.filePath, err)
					return
				}
				if len(content) > 0 && string(content) != string(prevContent) {
					prevContent = content
					updateChan <- content
				}
			case <-fw.stop:
				return
			}

		}
	}()
	return updateChan, nil
}

// Stop causes the FileWatcher to stop watching the file
func (fw *FileWatcher) Stop() {
	fw.stop <- true
}
