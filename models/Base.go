package models

import (
	"fmt"
	"leedoc/conf"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	var (
		err error
	)
	host := conf.GetConfig("db_host", "127.0.0.1").(string)
	port := conf.GetConfig("db_port", "3306").(string)
	user := conf.GetConfig("db_username", "root").(string)
	password := conf.GetConfig("db_password", "123456").(string)
	db_adapter := conf.GetConfig("db_adapter", "mysql").(string)

	if db_adapter == "mysql" {
		db_database_mysql := conf.GetConfig("db_database_mysql", "leedoc").(string)
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", user, password, host, port, db_database_mysql)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			panic(err)
		}
	} else {
		db_database_sqlite := conf.GetConfig("db_database_sqlite", "./database/leedoc.db").(string)
		if strings.HasPrefix(db_database_sqlite, "./") {
			db_database_sqlite = filepath.Join(conf.WorkingDirectory, string(db_database_sqlite[1:]))
		}
		if p, err := filepath.Abs(db_database_sqlite); err == nil {
			db_database_sqlite = p
		}

		dbPath := filepath.Dir(db_database_sqlite)

		if _, err := os.Stat(dbPath); err != nil && os.IsNotExist(err) {
			_ = os.MkdirAll(dbPath, 0777)
		}
		db, err = gorm.Open(sqlite.Open(db_database_sqlite), &gorm.Config{})
		if err != nil {
			panic(err)
		}
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)

	//check db connection
	if err := sqlDB.Ping(); err != nil {
		sqlDB.Close()
		log.Println("database connection failed")
		panic(err)
	}
	dm := db.Migrator()
	if !dm.HasTable(&Member{}) {
		dm.CreateTable(&Member{}, &Attachment{}, &Book{}, &Document{})
		RegisterAdmin("admin", "123456")
		RegisterDefaultBook(1)
	} else {
		db.AutoMigrate(&Member{}, &Attachment{}, &Book{}, &Document{})

	}
}
