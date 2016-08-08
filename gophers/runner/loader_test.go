package runner

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStepLoader(t *testing.T) {
	l, err := NewStepLoader(5, 5, 1, time.Second)
	require.NoError(t, err)
	assert.Equal(t, "StepLoader(5, 5, 1, 1s: 1 steps in 1s)", l.String())
	for d, c := range map[time.Duration]int{
		0 * time.Second: 5,
		1 * time.Second: -1,
		2 * time.Second: -1,
	} {
		assert.Equal(t, c, l.Count(d), "d = %s", d)
	}

	l, err = NewStepLoader(5, 10, 1, time.Second)
	require.NoError(t, err)
	assert.Equal(t, "StepLoader(5, 10, 1, 1s: 6 steps in 6s)", l.String())
	for d, c := range map[time.Duration]int{
		0 * time.Second: 5,
		1 * time.Second: 6,
		2 * time.Second: 7,
		3 * time.Second: 8,
		4 * time.Second: 9,
		5 * time.Second: 10,
		6 * time.Second: -1,
		7 * time.Second: -1,
	} {
		assert.Equal(t, c, l.Count(d), "d = %s", d)
	}

	l, err = NewStepLoader(5, 5, 2, time.Second)
	require.NoError(t, err)
	assert.Equal(t, "StepLoader(5, 5, 2, 1s: 1 steps in 1s)", l.String())
	for d, c := range map[time.Duration]int{
		0 * time.Second: 5,
		1 * time.Second: -1,
		2 * time.Second: -1,
	} {
		assert.Equal(t, c, l.Count(d), "d = %s", d)
	}

	l, err = NewStepLoader(5, 10, 2, time.Second)
	require.NoError(t, err)
	assert.Equal(t, "StepLoader(5, 10, 2, 1s: 4 steps in 4s)", l.String())
	for d, c := range map[time.Duration]int{
		0 * time.Second: 5,
		1 * time.Second: 7,
		2 * time.Second: 9,
		3 * time.Second: 10,
		4 * time.Second: -1,
		5 * time.Second: -1,
	} {
		assert.Equal(t, c, l.Count(d), "d = %s", d)
	}
}

func TestStepLoaderInfinite(t *testing.T) {
	l, err := NewStepLoader(5, 0, 1, time.Second)
	require.NoError(t, err)
	assert.Equal(t, "StepLoader(5, ∞, 1, 1s)", l.String())
	for d, c := range map[time.Duration]int{
		0 * time.Second: 5,
		1 * time.Second: 6,
		2 * time.Second: 7,
	} {
		assert.Equal(t, c, l.Count(d), "d = %s", d)
	}

	l, err = NewStepLoader(5, 0, 2, time.Second)
	require.NoError(t, err)
	assert.Equal(t, "StepLoader(5, ∞, 2, 1s)", l.String())
	for d, c := range map[time.Duration]int{
		0 * time.Second: 5,
		1 * time.Second: 7,
		2 * time.Second: 9,
	} {
		assert.Equal(t, c, l.Count(d), "d = %s", d)
	}
}
