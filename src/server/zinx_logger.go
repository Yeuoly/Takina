package server

import (
	"context"
	"fmt"

	"github.com/Yeuoly/Takina/src/helper"
)

type zinxLogger struct{}

func (l *zinxLogger) InfoF(format string, v ...interface{}) {
	helper.Info(format, v...)
}

func (l *zinxLogger) ErrorF(format string, v ...interface{}) {
	helper.Error(format, v...)
}

func (l *zinxLogger) DebugF(format string, v ...interface{}) {
	helper.Debug(format, v...)
}

func (l *zinxLogger) InfoFX(ctx context.Context, format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

func (l *zinxLogger) ErrorFX(ctx context.Context, format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

func (l *zinxLogger) DebugFX(ctx context.Context, format string, v ...interface{}) {
	fmt.Printf(format, v...)
}
