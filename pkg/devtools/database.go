package devtools

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	alog "github.com/apex/log"
	alogHandler "github.com/apex/log/handlers/cli"
	"github.com/jackc/pgx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gitlab.com/the-new-orleans-connection/qr-code/model"
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
	if config.EnvDBUsername == "" {
		alog.Error("one of DB_USER or DB_IAM_USER must be defined")
	}
	alog.WithField("config", config).Info("attempted to fetch db config from env")
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

type TestLogConsumer struct {
	logger *log.Logger
}

func (g *TestLogConsumer) Accept(l testcontainers.Log) {
	g.logger.Print(string(l.Content))
}

// SetupPostgresContainer will create a docker container for a postgres image
// for our subsequent tests to connect to, will set appropriate env variables
// to allow dal on the same process to access connection details
func SetupPostgresContainer(logger *log.Logger) (func(), *sql.DB, error) {
	ctx := context.Background()

	config, err := DatabaseConnectionConfigFromEnv()
	if err != nil {
		return nil, nil, err
	}

	nat_port, err := nat.NewPort("tcp", "5432")
	if err != nil {
		return nil, nil, err
	}

	req := testcontainers.ContainerRequest{
		Image:        "postgres:14.1-alpine",
		ExposedPorts: []string{"5432/tcp", fmt.Sprint(config.EnvDBPort, "/tcp")},
		Env: map[string]string{
			"POSTGRES_USER":     config.EnvDBUsername,
			"POSTGRES_PASSWORD": config.EnvDBPassword,
			"POSTGRES_DB":       config.EnvDBName,
		},
		WaitingFor: wait.ForListeningPort(nat_port),
	}

	container_postgres, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
		Logger:           logger,
	})

	if err != nil {
		return nil, nil, err
	}

	closeContainer := func() {
		logger.Print("BEGIN SetupPostgresContainer.(terminate)")
		container_postgres.Terminate(ctx)
		os.Exit(1)
	}

	err = container_postgres.StartLogProducer(ctx)
	if err != nil {
		logger.Print(err)
		return nil, nil, err
	}
	defer func() {
		err := container_postgres.StopLogProducer()
		logger.Print(err)
	}()

	container_postgres.FollowOutput(&TestLogConsumer{
		logger: logger,
	})

	host, _ := container_postgres.Host(ctx)
	nat_port, _ = container_postgres.MappedPort(ctx, "5432/tcp")
	port, _ := strconv.Atoi(nat_port.Port())

	connstring := fmt.Sprint("postgresql://", config.EnvDBUsername, ":", config.EnvDBPassword, "@", host, ":", port, "/", config.EnvDBName)

	var db *sql.DB
	db, err = sql.Open("pgx", connstring)
	if err != nil {
		alog.WithError(err).Error("failed to connect to database")
		if strings.Contains(err.Error(), "(FATAL: the database system is starting up (SQLSTATE 57P03)") {
			time.Sleep(1 * time.Second)
		} else {
			return nil, nil, err
		}
	}

	os.Setenv("DB_USERNAME", config.EnvDBUsername)
	os.Setenv("DB_PASSWORD", config.EnvDBPassword)
	os.Setenv("DB_PORT", nat_port.Port())
	os.Setenv("DB_HOSTNAME", host)
	os.Setenv("DB_NAME", config.EnvDBName)

	errConnecting := db.Ping()
	return closeContainer, db, errConnecting
}

func DatabaseMock(db *sql.DB) error {
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			Colorful: true,
			LogLevel: logger.Warn,
		}),
	})
	if err != nil {
		return err
	}

	for _, entity := range []interface{}{
		&model.FileStorageRecord{},
		&model.QRMapping{},
	} {
		err = gormDB.AutoMigrate(entity)
		if err != nil {
			return err
		}
	}

	err = DatabaseUnseed(gormDB)
	if err != nil {
		return err
	}
	return nil
}

func DatabaseUnseed(db *gorm.DB) error {
	session := &gorm.Session{
		AllowGlobalUpdate: true,
	}
	var result *gorm.DB
	result = db.Session(session).Unscoped().Delete(&model.QRMapping{})
	if result.Error != nil {
		return result.Error
	}

	result = db.Session(session).Unscoped().Delete(&model.FileStorageRecord{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

type DatabaseSeeder struct {
	Database *gorm.DB
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
