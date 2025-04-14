package elkcollector

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/rs/zerolog"
)

type CallerHook struct{}

func (h CallerHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	pc, _, line, ok := runtime.Caller(3)
	if !ok {
		e.Str("caller", "unknown").Int("line", 0)
		return
	}
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		e.Str("caller", "unknown").Int("line", 0)
		return
	}
	// Lấy tên hàm đầy đủ
	fullName := fn.Name()
	// Chỉ lấy phần tên hàm (bỏ package path)
	shortName := strings.Split(fullName, ".")
	// Thêm caller và line vào log
	e.Str("caller", fmt.Sprintf("%s.%d", shortName[len(shortName)-1], line))
}
