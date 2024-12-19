package utils

import (
	"log/slog"
)

func FailOnError(err error, message string) {
	if err != nil {
		slog.Error(message, "error", err)
	}
}