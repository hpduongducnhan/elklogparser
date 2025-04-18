package elkcollector

import (
	"runtime"

	"github.com/rs/zerolog/log"
)

func runWithErrLog(fn func() error) error {
	res := fn()
	if res != nil {
		// Capture caller information (file and line number)
		// Get caller information (function name, file, and line number)
		pc, file, line, ok := runtime.Caller(1)
		if !ok {
			file = "unknown"
			line = 0
		}

		var funcName string
		if pc != 0 {
			fnPtr := runtime.FuncForPC(pc)
			if fnPtr != nil {
				funcName = fnPtr.Name()
			} else {
				funcName = "unknown"
			}
		} else {
			funcName = "unknown"
		}
		log.Error().Err(res).Str("function", funcName).
			Str("file", file).
			Int("line", line).Msg("get error")
		return res
	}
	return nil
}
