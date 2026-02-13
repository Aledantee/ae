package ae

import (
	"fmt"
	"log/slog"
	"maps"
	"slices"
	"strings"
)

func (a Ae) LogValue() slog.Value {
	rootAttrs := []slog.Attr{
		slog.String("msg", a.msg),
		slog.Bool("recoverable", a.recoverable),
	}

	if a.userMsg != "" {
		rootAttrs = append(rootAttrs, slog.String("user_msg", a.userMsg))
	}
	if a.hint != "" {
		rootAttrs = append(rootAttrs, slog.String("hint", a.hint))
	}
	if !a.timestamp.IsZero() {
		rootAttrs = append(rootAttrs, slog.Time("timestamp", a.timestamp))
	}
	if a.code != "" {
		rootAttrs = append(rootAttrs, slog.String("code", a.code))
	}
	if a.exitCode > 0 {
		rootAttrs = append(rootAttrs, slog.Int("exit_code", a.exitCode))
	}

	if len(a.tags) > 0 {
		rootAttrs = append(rootAttrs, slog.String("tags", strings.Join(
			slices.Collect(maps.Keys(a.tags)), ", ")),
		)
	}

	if len(a.attributes) > 0 {
		var attrs []slog.Attr
		for k, v := range a.attributes {
			attrs = append(attrs, slog.Any(k, v))
		}
		rootAttrs = append(rootAttrs, slog.GroupAttrs("attributes", attrs...))
	}

	if len(a.causes) > 0 {
		var causeAttrs []slog.Attr
		for i, cause := range a.causes {
			causeAttrs = append(causeAttrs, slog.Any(fmt.Sprintf("%d", i), cause))
		}
		rootAttrs = append(rootAttrs, slog.GroupAttrs("causes", causeAttrs...))
	}

	if len(a.related) > 0 {
		var relatedAttrs []slog.Attr
		for i, rel := range a.related {
			relatedAttrs = append(relatedAttrs, slog.Any(fmt.Sprintf("%d", i), rel))
		}
		rootAttrs = append(rootAttrs, slog.GroupAttrs("related", relatedAttrs...))
	}

	return slog.GroupValue(
		rootAttrs...,
	)
}
