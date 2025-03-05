package models

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	logger "github.com/IvanSkripnikov/go-logger"
)

const MigrationTableName = "migration"
const FirstVersion = "00_create_migration_table.sql"

type Migration struct {
	Version   string
	ApplyTime int64
}

// HasExistsRow Проверка наличия текущей миграции в таблице
func (model *Migration) HasExistsRow(dbConn *sql.DB) bool {
	sqlQuery := fmt.Sprintf("select count(*) as countRow from %s where `version` = '%s'",
		MigrationTableName, model.Version)
	rows, err := dbConn.Query(sqlQuery)

	if err != nil {
		if strings.Contains(model.Version, FirstVersion) {
			return false
		}

		logger.Errorf("Failed to get migration data for version %v. Error: %v", model.Version, err)
	} else {
		logger.Debugf("Migration data for version %v was successfully received.", model.Version)
	}

	var countRow int

	if rows.Next() {
		err := rows.Scan(&countRow)

		if err != nil {
			logger.Errorf("Failed to get current migration string for version %v. Error: %v", model.Version, err)
		} else {
			logger.Debugf("The current migration string for version %v was successfully retrieved.", model.Version)
		}
	}

	defer rows.Close()
	return countRow > 0
}

// InsertRow Вставка строки для текущей миграции
func (model *Migration) InsertRow(dbConn *sql.DB) {
	sqlQuery := fmt.Sprintf("insert into %s values ('%s', %d);",
		MigrationTableName, model.Version, time.Now().Unix())

	result, err := dbConn.Exec(sqlQuery)
	if err != nil || result == nil {
		logger.Errorf("Failed to insert migration for version %v. Error: %v", model.Version, err)
	} else {
		logger.Debugf("Migration for version %v inserted successfully.", model.Version)
	}
}
