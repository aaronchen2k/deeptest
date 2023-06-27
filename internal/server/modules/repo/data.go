package repo

import (
	"database/sql"
	"fmt"
	"github.com/aaronchen2k/deeptest/internal/pkg/config"
	"gorm.io/gorm"
)

type DataRepo struct {
	DB *gorm.DB `inject:""`
}

// CreateMySqlDb 创建数据库(mysql)
func (s *DataRepo) CreateMySqlDb() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/",
		config.CONFIG.Mysql.Username, config.CONFIG.Mysql.Password,
		config.CONFIG.Mysql.Url)
	createSql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_general_ci;",
		config.CONFIG.Mysql.Dbname)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)
	if err = db.Ping(); err != nil {
		return err
	}
	_, err = db.Exec(createSql)
	return err
}

func (s *DataRepo) SetSqlMode() (err error) {
	sql := "SET sql_mode = '';"
	err = s.DB.Raw(sql).Error

	return
}
