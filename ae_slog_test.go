package ae_test

import (
	"log/slog"
	"strings"
	"testing"
	"time"

	"go.aledante.io/ae"
)

// flattenAttrs walks a slog.Value (expected to be a Group) and collects every
// leaf attribute into a flat "group.group.key" -> value map. Nested
// slog.LogValuer implementations are resolved so that, e.g., cause errors
// that themselves implement LogValue expand into their own nested groups.
func flattenAttrs(v slog.Value) map[string]any {
	out := make(map[string]any)
	var walk func(prefix string, attrs []slog.Attr)
	walk = func(prefix string, attrs []slog.Attr) {
		for _, a := range attrs {
			key := a.Key
			if prefix != "" {
				key = prefix + "." + a.Key
			}
			val := a.Value.Resolve()
			if val.Kind() == slog.KindAny {
				if lv, ok := val.Any().(slog.LogValuer); ok {
					val = lv.LogValue().Resolve()
				}
			}
			if val.Kind() == slog.KindGroup {
				walk(key, val.Group())
				continue
			}
			out[key] = val.Any()
		}
	}
	if v.Kind() == slog.KindGroup {
		walk("", v.Group())
	}
	return out
}

// logValue extracts the LogValue from an Ae error. Builder terminal operations
// return the error interface; Ae itself is the concrete value-type receiver of
// LogValue, so we assert the error is an *ae.Ae pointer via the slog.LogValuer
// interface to sidestep that indirection.
func logValue(t *testing.T, err error) slog.Value {
	t.Helper()
	lv, ok := err.(slog.LogValuer)
	if !ok {
		t.Fatalf("%T does not implement slog.LogValuer", err)
	}
	return lv.LogValue()
}

func TestAe_LogValue_RootFieldsPresent(t *testing.T) {
	t.Parallel()

	ts := time.Date(2026, 4, 17, 12, 0, 0, 0, time.UTC)

	// UserMsg is a terminal operation — build all non-terminal fields first,
	// then finish with UserMsg to seed both msg and userMsg.
	err := ae.New().
		Hint("try again").
		Timestamp(ts).
		Code("E_X").
		ExitCode(5).
		Tag("a").
		Attr("k", "v").
		UserMsg("dev msg", "user msg")

	attrs := flattenAttrs(logValue(t, err))

	if attrs["msg"] != "dev msg" {
		t.Errorf("msg = %v, want 'dev msg'", attrs["msg"])
	}
	if attrs["user_msg"] != "user msg" {
		t.Errorf("user_msg = %v, want 'user msg'", attrs["user_msg"])
	}
	if attrs["hint"] != "try again" {
		t.Errorf("hint = %v, want 'try again'", attrs["hint"])
	}
	if attrs["code"] != "E_X" {
		t.Errorf("code = %v, want 'E_X'", attrs["code"])
	}
	if attrs["exit_code"] != int64(5) {
		t.Errorf("exit_code = %v, want 5 (int64)", attrs["exit_code"])
	}
	if attrs["recoverable"] != true {
		t.Errorf("recoverable = %v, want true", attrs["recoverable"])
	}
	if _, ok := attrs["timestamp"]; !ok {
		t.Errorf("timestamp missing from LogValue output")
	}
	if attrs["tags"] != "a" {
		t.Errorf("tags = %v, want 'a'", attrs["tags"])
	}
	if attrs["attributes.k"] != "v" {
		t.Errorf("attributes.k = %v, want 'v'", attrs["attributes.k"])
	}
}

func TestAe_LogValue_OmitsEmptyFields(t *testing.T) {
	t.Parallel()

	attrs := flattenAttrs(logValue(t, ae.New().Msg("plain")))

	// msg + recoverable are always present; everything else is omitempty.
	for _, k := range []string{"user_msg", "hint", "timestamp", "code", "exit_code", "tags"} {
		if _, present := attrs[k]; present {
			t.Errorf("LogValue emitted empty %q = %v", k, attrs[k])
		}
	}
	for k := range attrs {
		if strings.HasPrefix(k, "attributes.") ||
			strings.HasPrefix(k, "causes.") ||
			strings.HasPrefix(k, "related.") {
			t.Errorf("LogValue emitted unexpected group key %q", k)
		}
	}
}

func TestAe_LogValue_CausesAndRelatedGrouped(t *testing.T) {
	t.Parallel()

	err := ae.New().
		Cause(ae.New().Msg("root-cause")).
		Related(ae.New().Msg("side-effect")).
		Msg("outer")

	attrs := flattenAttrs(logValue(t, err))

	// Each cause / related nests another Ae, which re-renders its own group.
	// Asserting the top-most leaf key of each nested message is enough to
	// show the groups wire through correctly.
	if attrs["causes.0.msg"] != "root-cause" {
		t.Errorf("causes.0.msg = %v, want 'root-cause'", attrs["causes.0.msg"])
	}
	if attrs["related.0.msg"] != "side-effect" {
		t.Errorf("related.0.msg = %v, want 'side-effect'", attrs["related.0.msg"])
	}
}
