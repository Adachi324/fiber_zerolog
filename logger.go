package fiber_zerolog

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/rs/zerolog"
	"io"
	"os"
	"sync"
)

type Logger struct {
	logger zerolog.Logger
}

func NewLogger() *Logger {
	l := zerolog.New(os.Stderr).Level(zerolog.DebugLevel).With().Timestamp().Logger()
	return &Logger{logger: l}
}

func NewLoggerByZerolog(l zerolog.Logger) *Logger {
	return &Logger{logger: l}
}

var _ log.AllLogger = (*Logger)(nil)

func (lg *Logger) Trace(val ...interface{}) {
	event := lg.logger.Trace()
	for _, v := range val {
		event.Msgf("%v", v)
	}
}

func (lg *Logger) Debug(v ...interface{}) {
	lg.logger.Debug().Msgf("%v", v)
}

func (lg *Logger) Info(v ...interface{}) {
	lg.logger.Info().Msgf("%v", v)
}

func (lg *Logger) Warn(v ...interface{}) {
	lg.logger.Warn().Msgf("%v", v)
}

func (lg *Logger) Error(v ...interface{}) {
	lg.logger.Error().Msgf("%v", v)
}

func (lg *Logger) Fatal(v ...interface{}) {
	lg.logger.Fatal().Msgf("%v", v)
}

func (lg *Logger) Panic(v ...interface{}) {
	lg.logger.Panic().Msgf("%v", v)
}

func (lg *Logger) Tracef(format string, v ...interface{}) {
	lg.logger.Trace().Msgf(format, v)
}

func (lg *Logger) Debugf(format string, v ...interface{}) {
	lg.logger.Debug().Msgf(format, v)
}

func (lg *Logger) Infof(format string, v ...interface{}) {
	lg.logger.Info().Msgf(format, v)
}

func (lg *Logger) Warnf(format string, v ...interface{}) {
	lg.logger.Warn().Msgf(format, v)
}

func (lg *Logger) Errorf(format string, v ...interface{}) {
	lg.logger.Error().Msgf(format, v)
}

func (lg *Logger) Fatalf(format string, v ...interface{}) {
	lg.logger.Fatal().Msgf(format, v)
}

func (lg *Logger) Panicf(format string, v ...interface{}) {
	lg.logger.Panic().Msgf(format, v)
}

func (lg *Logger) Tracew(msg string, keysAndValues ...interface{}) {
	lg.handleKV(lg.logger.Trace(), msg, keysAndValues...).Msg(msg)
}

func (lg *Logger) Debugw(msg string, keysAndValues ...interface{}) {
	lg.handleKV(lg.logger.Debug(), msg, keysAndValues...).Msg(msg)
}

func (lg *Logger) Infow(msg string, keysAndValues ...interface{}) {
	lg.handleKV(lg.logger.Info(), msg, keysAndValues...).Msg(msg)
}

func (lg *Logger) Warnw(msg string, keysAndValues ...interface{}) {
	lg.handleKV(lg.logger.Warn(), msg, keysAndValues...).Msg(msg)
}

func (lg *Logger) Errorw(msg string, keysAndValues ...interface{}) {
	lg.handleKV(lg.logger.Error(), msg, keysAndValues...).Msg(msg)
}

func (lg *Logger) Fatalw(msg string, keysAndValues ...interface{}) {
	lg.handleKV(lg.logger.Fatal(), msg, keysAndValues...).Msg(msg)
}

func (lg *Logger) Panicw(msg string, keysAndValues ...interface{}) {
	lg.handleKV(lg.logger.Panic(), msg, keysAndValues...).Msg(msg)
}
func (lg *Logger) handleKV(event *zerolog.Event, msg string, kvs ...interface{}) *zerolog.Event {
	var once sync.Once
	isFirst := true
	for i := 0; i < len(kvs); i += 2 {
		if msg == "" && isFirst {
			once.Do(func() {
				event.Any(fmt.Sprintf("%s", kvs[i]), kvs[i+1])
				isFirst = false
			})
			continue
		}
		event.Any(fmt.Sprintf("%s", kvs[i]), kvs[i+1])
	}
	return event
}

func (lg *Logger) SetLevel(level log.Level) {
	lvl, err := zerolog.ParseLevel(zerolog.Level(int8(level - 1)).String())
	if err != nil {
		lg.logger.Error().Err(err).Msg("invalid log level")
		return
	}
	zerolog.SetGlobalLevel(lvl)
}

func (lg *Logger) SetOutput(writer io.Writer) {
	l := lg.logger.Output(writer)
	lg.logger = l
}

func (lg *Logger) WithContext(ctx context.Context) log.CommonLogger {
	l := zerolog.Ctx(ctx)
	if l == zerolog.DefaultContextLogger || l.GetLevel() == zerolog.Disabled {
		return lg
	}
	return NewLoggerByZerolog(*l)
}
