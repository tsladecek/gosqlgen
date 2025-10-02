package gosqlgen

import (
	"fmt"
	"runtime"
	"strings"
)

func callingFuncInfo() (int, string) {
	pc, _, n, ok := runtime.Caller(2)
	if !ok {
		return 0, ""
	}
	f := runtime.FuncForPC(pc)
	if f != nil {
		return n, strings.TrimPrefix(f.Name(), "github.com/tsladecek/")
	}
	return n, ""
}

func Errorf(format string, a ...any) error {
	lineNumber, fn := callingFuncInfo()
	format = "\n[%s:%d]:\t" + format
	args := []any{fn, lineNumber}
	args = append(args, a...)
	return fmt.Errorf(format, args...)
}
