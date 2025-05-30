// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package utilizationtracker provides a utility to track the utilization of a component.
package utilizationtracker

import (
	"math/rand"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helpers

func newTracker(_ *testing.T) (*UtilizationTracker, *clock.Mock) {
	clk := clock.NewMock()
	ut := newUtilizationTrackerWithClock(
		100*time.Millisecond,
		clk,
		0.25,
	)

	return ut, clk
}

// Tests

func TestUtilizationTracker(t *testing.T) {
	ut, clk := newTracker(t)
	defer ut.Stop()

	old := 0.0
	newValue := 0.0

	// After some time without any checks running, the utilization
	// should be a constant zero value
	clk.Add(300 * time.Millisecond)
	ut.Tick()
	old, newValue = newValue, <-ut.Output
	require.Equal(t, old, newValue)

	clk.Add(300 * time.Millisecond)
	// Ramp up the expected utilization
	ut.Started()
	old, newValue = newValue, <-ut.Output
	require.Equal(t, old, newValue)

	clk.Add(250 * time.Millisecond)
	ut.Tick()
	old, newValue = newValue, <-ut.Output
	require.Greater(t, newValue, old)

	clk.Add(550 * time.Millisecond)
	ut.Tick()
	old, newValue = newValue, <-ut.Output
	require.Greater(t, newValue, old)

	// Ramp down the expected utilization
	ut.Finished()
	old, newValue = newValue, <-ut.Output
	require.Equal(t, old, newValue) //no time have passed

	clk.Add(250 * time.Millisecond)
	ut.Tick()
	old, newValue = newValue, <-ut.Output
	require.Less(t, newValue, old)

	clk.Add(550 * time.Millisecond)
	ut.Tick()
	require.Less(t, newValue, old)
}

func TestUtilizationTrackerCheckLifecycle(t *testing.T) {
	ut, clk := newTracker(t)
	defer ut.Stop()

	var old, newValue float64

	// No tasks should equal no utilization
	clk.Add(250 * time.Millisecond)
	ut.Tick()
	old, newValue = newValue, <-ut.Output
	assert.Equal(t, old, newValue)

	for idx := 0; idx < 3; idx++ {
		// Ramp up utilization
		ut.Started()
		old, newValue = newValue, <-ut.Output
		assert.Equal(t, old, newValue)

		clk.Add(250 * time.Millisecond)
		ut.Tick()
		old, newValue = newValue, <-ut.Output
		assert.Greater(t, newValue, old)

		clk.Add(250 * time.Millisecond)
		ut.Tick()
		old, newValue = newValue, <-ut.Output
		assert.Greater(t, newValue, old)

		// Ramp down utilization
		ut.Finished()
		old, newValue = newValue, <-ut.Output
		assert.Equal(t, newValue, old)

		clk.Add(250 * time.Millisecond)
		ut.Tick()
		old, newValue = newValue, <-ut.Output
		assert.Less(t, newValue, old)

		clk.Add(250 * time.Millisecond)
		ut.Tick()
		old, newValue = newValue, <-ut.Output
		assert.Less(t, newValue, old)
	}
}

func TestUtilizationTrackerAccuracy(t *testing.T) {
	ut, clk := newTracker(t)

	val := 0.0

	// It would be nice to figure out a way to compute bounds for the
	// smoothed value that would work for any random sequence.
	r := rand.New(rand.NewSource(1))

	for checkIdx := 1; checkIdx <= 100; checkIdx++ {
		// This should provide about 30% utilization
		// Range for the full loop would be between 100-200ms
		totalMs := r.Int31n(100) + 100
		runtimeMs := (totalMs * 30) / 100

		ut.Started()
		<-ut.Output

		runtimeDuration := time.Duration(runtimeMs) * time.Millisecond
		clk.Add(runtimeDuration)

		ut.Finished()
		val = <-ut.Output

		idleDuration := time.Duration(totalMs-runtimeMs) * time.Millisecond
		clk.Add(idleDuration)

		if checkIdx > 30 {
			require.InDelta(t, 0.3, val, 0.07)
		}
	}

	require.InDelta(t, 0.3, val, 0.07)
}
