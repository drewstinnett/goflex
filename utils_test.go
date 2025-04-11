package goflex

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDaysToDuration(t *testing.T) {
	d := daysToDuration(1)
	require.Equal(t, time.Hour*24, d)

	d = daysToDuration(2)
	require.Equal(t, time.Hour*48, d)

	d = daysToDuration(0)
	require.Equal(t, time.Duration(0), d)

	d = daysToDuration(-1)
	require.Equal(t, time.Duration(-24)*time.Hour, d)
}
