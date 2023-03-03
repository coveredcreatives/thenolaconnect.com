package main

import (
	"log"
	"os"

	alog "github.com/apex/log"
	cli "github.com/urfave/cli/v2"
	"gitlab.com/the-new-orleans-connection/qr-code/devtools"
)

func ContainerDatabase(ctx *cli.Context) (func(), error) {
	terminate, _, err := devtools.SetupPostgresContainer(log.Default())
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

	// gormconfig := devtools.GormConfig()

	// gormdb, err := gorm.Open(postgres.New(postgres.Config{
	// 	Conn: db,
	// }), gormconfig)
	// if err != nil {
	// 	alog.WithError(err).Error("unable to initialize gorm")
	// 	return terminate, err
	// }

	// Migrate database based on project models
	// for key, table := range model.Tables([]string{}) {
	// 	alog.WithField("name", key).Info("migrating table")
	// 	err = gormdb.AutoMigrate(table)
	// 	if err != nil {
	// 		alog.WithField("name", key).WithError(err).Error("unable to migrate table")
	// 		return terminate, err
	// 	}
	// }
	return terminate, nil
}
