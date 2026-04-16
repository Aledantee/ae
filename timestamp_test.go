package ae_test

import (
	"errors"
	"testing"
	"time"

	"go.aledante.io/ae"
)

func TestTimestamp_ErrorWithoutInterface(t *testing.T) {
	t.Parallel()

	got := ae.Timestamp(errors.New("plain"))
	if !got.IsZero() {
		t.Errorf("Timestamp(plainErr) = %v, want zero time", got)
	}
}

func TestTimestamp_ErrorImplementingInterface(t *testing.T) {
	t.Parallel()

	when := time.Date(2025, 1, 2, 3, 4, 5, 0, time.UTC)
	err := stubErr{msg: "x", timestamp: when}
	if got := ae.Timestamp(err); !got.Equal(when) {
		t.Errorf("Timestamp(stubErr) = %v, want %v", got, when)
	}
}

func TestTimestamp_AeBuilderSetsTimestamp(t *testing.T) {
	t.Parallel()

	when := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
	err := ae.New().Timestamp(when).Msg("fail")
	if got := ae.Timestamp(err); !got.Equal(when) {
		t.Errorf("Timestamp on builder = %v, want %v", got, when)
	}
}

func TestTimestamp_AeBuilderNowSetsCurrentTime(t *testing.T) {
	t.Parallel()

	before := time.Now()
	err := ae.New().Now().Msg("fail")
	after := time.Now()

	got := ae.Timestamp(err)
	if got.Before(before) || got.After(after) {
		t.Errorf("Timestamp(Now()) = %v, want between %v and %v", got, before, after)
	}
}
