package helpers

import (
	"database/sql"

	"orders-service/models"

	"github.com/IvanSkripnikov/go-gormdb"
	logger "github.com/IvanSkripnikov/go-logger"
	"gorm.io/gorm"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB
var GormDB *gorm.DB

func InitDatabase(config gormdb.Database) {
	gormDatabase, err := gormdb.AddMysql(models.ServiceDatabase, config)
	if err != nil {
		logger.Fatalf("Cant initialize DB: %v", err)
	}
	db, err := gormDatabase.DB()
	if err != nil {
		logger.Fatalf("Cant get DB: %v", err)
	}

	DB = db
	GormDB = gormdb.GetClient(models.ServiceDatabase)
}
