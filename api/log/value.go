package log

import (
	"context"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	defaultDepth     = 3
	DefaultCaller    = Caller(defaultDepth)
	DefaultTinestamp = Timestamp(time.RFC3339)
)

type Valuer func(ctx context.Context) interface{}

func Value(ctx context.Context, v interface{}) interface{} {
	if v, ok := v.(Valuer); ok {
		return v(ctx)
	}
	return v
}

func Caller(depth int) Valuer {
	return func(context.Context) interface{} {
		_, file, line, _ := runtime.Caller(depth)
		if strings.LastIndex(file, "/log/filer.go") > 0 {
			depth++
			_, file, line, _ = runtime.Caller(depth)
		}
		if strings.LastIndex(file, "/log/helper.go") > 0 {
			depth++
			_, file, line, _ = runtime.Caller(depth)

		}
		idx := strings.LastIndexByte(file, '/')
		return file[idx+1:] + ":" + strconv.Itoa(line)
	}
}

func Timestamp(layout string) Valuer {
	return func(context.Context) interface{} {
		return time.Now().Format(layout)
	}
}

func bindValues(ctx context.Context, keyvals []interface{}) {
	for i := 1; i < len(keyvals); i += 2 {
		if v, ok := keyvals[i].(Valuer); ok {
			keyvals[i] = v(ctx)
		}
	}
}

func containsValuer(keyvals []interface{}) bool {
	for i := 1; i < len(keyvals); i += 2 {
		if _, ok := keyvals[i].(Valuer); ok {
			return true
		}
	}
	return false
}
