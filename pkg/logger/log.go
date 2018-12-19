package logger

import (
	"os"

	logging "github.com/op/go-logging"
)

// DefaultLogger return a default logger
func DefaultLogger() *logging.Logger {
	var log = logging.MustGetLogger("service-copy")
	var format = logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level}%{color:reset} %{message}`,
	)
	backend2 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2Formatter := logging.NewBackendFormatter(backend2, format)
	logging.SetBackend(backend2Formatter)

	return log
}
