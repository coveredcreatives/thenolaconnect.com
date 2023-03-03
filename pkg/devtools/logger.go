package devtools

import (
	"context"
	"errors"
	"strconv"
	"time"

	alog "github.com/apex/log"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

const SlowThresholdConfig = 10

func NewGormLogger(log *alog.Logger) gormLogger.Interface {
	return &GormLogger{
		SlowThreshold: time.Duration(SlowThresholdConfig * time.Second),
		LogLevel:      gormLogger.Info,
		Logger:        log,
	}
}

// GormLogger struct
type GormLogger struct {
	LogLevel                  gormLogger.LogLevel
	Logger                    *alog.Logger
	IgnoreRecordNotFoundError bool
	SlowThreshold             time.Duration
}

func (g *GormLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	g.LogLevel = level
	return g
}

// Info print info
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormLogger.Info {
		entry := l.Logger.WithField("file_with_line_number", utils.FileWithLineNum())
		for i, d := range data {
			entry = entry.WithField(strconv.Itoa(i), d)
		}
		entry.Info(msg)
	}
}

// Warn print warn messages
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormLogger.Warn {
		entry := l.Logger.
			WithField("file_with_line_number", utils.FileWithLineNum())
		for i, d := range data {
			entry = entry.WithField(strconv.Itoa(i), d)
		}
		entry.Warn(msg)
	}
}

// Error print error messages
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormLogger.Error {
		entry := l.Logger.
			WithField("file_with_line_number", utils.FileWithLineNum())
		for i, d := range data {
			entry = entry.WithField(strconv.Itoa(i), d)
		}
		entry.Error(msg)
	}
}

// Trace print sql message
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= gormLogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= gormLogger.Error && (!errors.Is(err, gormLogger.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rowcount := fc()
		rows := ""
		if rowcount == -1 {
			rows = "-"
		} else {
			rows = strconv.FormatInt(rowcount, 10)
		}
		l.Logger.
			WithError(err).
			WithField("rows", rows).
			WithField("file_with_line_num", utils.FileWithLineNum()).
			WithField("elapsed", time.Duration(float64(elapsed.Nanoseconds())/1e6)).
			WithField("sql", sql).
			Error("Done")

	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= gormLogger.Warn:
		sql, rowcount := fc()
		rows := ""
		if rowcount == -1 {
			rows = "-"
		} else {
			rows = strconv.FormatInt(rowcount, 10)
		}
		l.Logger.
			WithField("rows", rows).
			WithField("file_with_line_num", utils.FileWithLineNum()).
			WithField("slow_sql", l.SlowThreshold).
			WithField("elapsed", time.Duration(float64(elapsed.Nanoseconds()))).
			WithField("sql", sql).
			Warn("Done")
	case l.LogLevel == gormLogger.Info:
		sql, rowcount := fc()
		rows := ""
		if rowcount == -1 {
			rows = "-"
		} else {
			rows = strconv.FormatInt(rowcount, 10)
		}
		l.Logger.
			WithField("rows", rows).
			WithField("file_with_line_num", utils.FileWithLineNum()).
			WithField("elapsed", time.Duration(float64(elapsed.Nanoseconds())/1e6)).
			WithField("sql", sql).
			Info("Done")
	}
}
