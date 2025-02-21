package helpers

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"orders-service/logger"
	"orders-service/models"

	_ "github.com/go-sql-driver/mysql"
)

const migrationsDir = "./migrations"

type DataBaseConnection struct {
	host     string
	db       string
	user     string
	password string
}

var DB *sql.DB

// InitDataBase подключение к БД сервиса
func InitDataBase(config models.Database) (*sql.DB, error) {
	dataSource := GetDatabaseConnectionString(config)
	db, err := sql.Open("mysql", dataSource)

	if err != nil {
		logger.Fatalf("Failed to connect to service database. Error: %v", err)
		return nil, err
	}

	DB = db

	return db, nil
}

func GetDatabaseConnectionString(config models.Database) string {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		config.User, config.Password, config.Address, config.Port, config.DB)
	logger.Info(connectionString)
	return connectionString
}

// CreateTables Выполнить запросы на создание таблиц
func CreateTables() {
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		logger.Errorf("Failed to get list of migration files. Error: %v", err)
	} else {
		logger.Debug("List of migration files retrieved successfully.")
	}

	dataFirstMigration, err := os.ReadFile(migrationsDir + "/" + models.FirstVersion)
	sqlQuery := strings.ReplaceAll(string(dataFirstMigration), "\r\n", "")
	_, err = DB.Exec(sqlQuery)

	if err != nil {
		logger.Errorf("Failed to execute first migration. Error: %v", err)
	} else {
		logger.Debug("The First migration was successfully applied")
	}

	for _, file := range files {
		if !file.IsDir() {
			migration := models.Migration{
				Version: file.Name(),
			}

			if !migration.HasExistsRow(DB) {
				data, err := os.ReadFile(migrationsDir + "/" + file.Name())

				if err != nil {
					logger.Errorf("Failed to read migration file: %v. Error: %v", file.Name(), err)
				} else {
					logger.Debugf("The migration file was successfully read: %v.", file.Name())
				}

				sqlQuery := strings.ReplaceAll(string(data), "\r\n", "")
				result, err := DB.Exec(sqlQuery)
				migration.InsertRow(DB)

				if err == nil && result != nil {
					logger.Infof("Migration has been applied successfully: %v.", file.Name())
				}
			}
		}
	}
}
