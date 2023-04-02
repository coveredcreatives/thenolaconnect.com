package devtools

import (
	"context"
	"fmt"

	alog "github.com/apex/log"
	alogHandler "github.com/apex/log/handlers/cli"
	"github.com/jackc/pgx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Logger struct{}

func (l *Logger) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	entry := alog.WithField("loglevel", level)
	for key, value := range data {
		entry = entry.WithField(key, value)
	}
	switch level {
	case pgx.LogLevelError:
		entry.Error(msg)
	case pgx.LogLevelWarn:
		entry.Warn(msg)
	case pgx.LogLevelTrace:
		entry.Trace(msg)
	case pgx.LogLevelDebug:
		entry.Debug(msg)
	default:
		entry.Info(msg)
	}
}

func DatabaseConnection(context context.Context, config DatabaseConnectionConfig) (*gorm.DB, error) {
	dsn := ""
	if config.EnvDBUsername != "" {
		dsn = dsn + fmt.Sprint("user=", config.EnvDBUsername)
	}
	if config.EnvDBPassword != "" {
		dsn = dsn + fmt.Sprint(" password=", config.EnvDBPassword)
	}
	if config.EnvDBHostname != "" {
		dsn = dsn + fmt.Sprint(" host=", config.EnvDBHostname)
	}
	if config.EnvDBName != "" {
		dsn = dsn + fmt.Sprint(" dbname=", config.EnvDBName)
	}
	if config.EnvDBPort != 0 {
		dsn = dsn + fmt.Sprint(" port=", config.EnvDBPort)
	}
	alog.WithField("dsn", dsn).Info("concatenated dsn")
	return gorm.Open(postgres.Open(dsn), GormConfig())
}

func GormConfig() *gorm.Config {
	return &gorm.Config{
		QueryFields: true, // enumerate fields in query instead of using *
		Logger: NewGormLogger(&alog.Logger{
			Level:   alog.WarnLevel,
			Handler: alogHandler.Default,
		}),
	}
}
