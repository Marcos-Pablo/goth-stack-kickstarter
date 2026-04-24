package logging

import (
	"fmt"
	"log/slog"
	"net/url"
	"slices"
)

var sensitiveKeys = []string{
	"password", "key", "apikey", "secret", "pin", "token",
	"authorization",
}

func redactAttr(groups []string, a slog.Attr) slog.Attr {
	if slices.Contains(sensitiveKeys, a.Key) {
		return slog.String(a.Key, "[REDACTED]")
	}

	if a.Value.Kind() == slog.KindString {
		if u, err := url.Parse(a.Value.String()); err == nil {
			if _, ok := u.User.Password(); ok {
				u.User = url.UserPassword(u.User.Username(),
					"[REDACTED]")
				return slog.String(a.Key, u.String())
			}
		}
	}

	if a.Key == "err" || a.Key == "error" {
		if err, ok := a.Value.Any().(error); ok {
			if multi, ok := err.(interface{ Unwrap() []error }); ok {
				groups := make([]any, 0, len(multi.Unwrap()))
				for i, e := range multi.Unwrap() {
					groups = append(groups, slog.String(fmt.Sprintf("error_%d", i+1),
						e.Error()))
				}
				return slog.Group("errors", groups...)
			}
			return slog.String(a.Key, err.Error())
		}
	}

	return a
}
