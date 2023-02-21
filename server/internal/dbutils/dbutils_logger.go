package dbutils

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"time"
)

type Logger struct {
	SlowThreshold time.Duration
	Debug         bool
}

func NewLogger(isDebug bool) *Logger {
	return &Logger{
		SlowThreshold: time.Second,
		Debug:         isDebug,
	}
}

func (l *Logger) LogMode(gormLogger.LogLevel) gormLogger.Interface {
	// ignored
	return l
}

func (l *Logger) Info(ctx context.Context, s string, args ...interface{}) {
	log.WithContext(ctx).Infof(s, args)
}

func (l *Logger) Warn(ctx context.Context, s string, args ...interface{}) {
	log.WithContext(ctx).Warnf(s, args)
}

func (l *Logger) Error(ctx context.Context, s string, args ...interface{}) {
	log.WithContext(ctx).Errorf(s, args)
}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)

	sql, rowsAffected := fc()

	fields := log.Fields{
		"module":  "db",
		"elapsed": elapsed,
		"rows":    rowsAffected,
	}

	if l.Debug {
		fields["src"] = utils.FileWithLineNum()
	}

	// ignore "record not found" if not running in debug mode
	if err != nil && (l.Debug || (!l.Debug && !errors.Is(err, gorm.ErrRecordNotFound))) {
		log.WithContext(ctx).
			WithFields(fields).
			WithError(err).
			Errorf("%s", sql)

		return
	}

	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		log.
			WithContext(ctx).
			WithFields(fields).
			Warnf("%s", sql)

		return
	}

	if l.Debug {
		log.
			WithContext(ctx).
			WithFields(fields).
			Debugf("%s", sql)
	}
}
