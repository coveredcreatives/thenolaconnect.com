package main

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
	"github.com/coveredcreatives/thenolaconnect.com/pkg/devtools"
	"github.com/coveredcreatives/thenolaconnect.com/pkg/model"
	"github.com/docker/go-connections/nat"
	"github.com/spf13/viper"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	cli "github.com/urfave/cli/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ContainerDatabase(ctx *cli.Context, v *viper.Viper) (func(), error) {
	terminate, _, err := SetupPostgresContainer(log.Default(), v)
	if err != nil {
		alog.WithError(err).Error("failed to setup postgres container")
		return terminate, err
	}

	log.Println("DB_USERNAME: ", os.Getenv("DB_USERNAME"))
	log.Println("DB_PASSWORD: ", os.Getenv("DB_PASSWORD"))
	log.Println("DB_PORT: ", os.Getenv("DB_PORT"))
	log.Println("DB_HOSTNAME: ", os.Getenv("DB_HOSTNAME"))
	log.Println("DB_NAME: ", os.Getenv("DB_NAME"))

	log.Println("\nto migrate database: ")
	log.Println("psql -U", os.Getenv("DB_USERNAME"), "-h", os.Getenv("DB_HOSTNAME"), "-p", os.Getenv("DB_PORT"), "-d", os.Getenv("DB_NAME"), "-f", os.Getenv("DB_MIGRATION_FILE"))

	log.Println("\nto insert default rows: ")
	log.Println("psql -U", os.Getenv("DB_USERNAME"), "-h", os.Getenv("DB_HOSTNAME"), "-p", os.Getenv("DB_PORT"), "-d", os.Getenv("DB_NAME"), "-f", os.Getenv("DB_SEEDER_FILE"))

	log.Println("\nto connect an active session: ")
	log.Println("psql -U ", os.Getenv("DB_USERNAME"), " -h ", os.Getenv("DB_HOSTNAME"), " -p ", os.Getenv("DB_PORT"), "-d", os.Getenv("DB_NAME"))

	return terminate, nil
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
func SetupPostgresContainer(logger *log.Logger, v *viper.Viper) (func(), *sql.DB, error) {
	ctx := context.Background()

	config, err := devtools.DatabaseConnectionLoadConfig(v)
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
