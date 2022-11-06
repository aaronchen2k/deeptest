package dao

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/config"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	logUtils "github.com/aaronchen2k/deeptest/pkg/lib/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
	"os"
	"path/filepath"
	"sync"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	once sync.Once
	db   *gorm.DB
)

// GetDB 数据库单例
func GetDB() *gorm.DB {
	once.Do(func() {
		switch config.CONFIG.System.DbType {
		case "sqlite":
			db = GormSQLLite()
		case "mysql":
			db = GormMySQL()
		default:
			db = GormSQLLite()
		}
	})

	return db
}

// MysqlTables 注册数据库表专用
func MysqlTables(db *gorm.DB) {
	err := db.AutoMigrate()
	if err != nil {
		logUtils.Errorf("注册数据表错误", zap.Any("err", err))
		os.Exit(0)
	}
	logUtils.Infof("注册数据表成功")
}

func GormSQLLite() *gorm.DB {
	conn := DBFile()
	dialector := sqlite.Open(conn)

	db, err := gorm.Open(dialector, &gorm.Config{
		SkipDefaultTransaction: false,
		Logger:                 logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",
			SingularTable: false,
		},
	})

	if err != nil {
		logUtils.Info(err.Error())
	}

	_ = db.Use(
		dbresolver.Register(
			dbresolver.Config{ /* xxx */ }).
			SetConnMaxIdleTime(time.Hour).
			SetConnMaxLifetime(24 * time.Hour).
			SetMaxIdleConns(100).
			SetMaxOpenConns(200),
	)

	db.Session(&gorm.Session{FullSaveAssociations: true, AllowGlobalUpdate: false})

	//err = db.AutoMigrate(
	//	model.Models...,
	//)
	//if err != nil {
	//	logUtils.Info(err.SendErrorMsg())
	//}

	return db
}

func DBFile() string {
	path := filepath.Join(consts.WorkDir, consts.App+".db")
	return path
}

// GormMySQL 初始化Mysql数据库
func GormMySQL() *gorm.DB {
	m := config.CONFIG.Mysql
	if m.Dbname == "" {
		return nil
	}
	mysqlConfig := mysql.Config{
		DSN:                       m.Dsn(), // DSN data source name
		DefaultStringSize:         191,     // string 类型字段的默认长度
		DisableDatetimePrecision:  true,    // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,    // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,    // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,   // 根据版本自动配置
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig), gormConfig(m.LogMode)); err != nil {
		return nil
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
		return db
	}
}

// gormConfig 根据配置决定是否开启日志
func gormConfig(mod bool) *gorm.Config {
	var config = &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true}
	switch config.CONFIG.Mysql.LogZap {
	case "silent", "Silent":
		config.Logger = Default.LogMode(logger.Silent)
	case "error", "SendErrorMsg":
		config.Logger = Default.LogMode(logger.Error)
	case "warn", "Warn":
		config.Logger = Default.LogMode(logger.Warn)
	case "info", "Info":
		config.Logger = Default.LogMode(logger.Info)
	case "zap", "Zap":
		config.Logger = Default.LogMode(logger.Info)
	default:
		if mod {
			config.Logger = Default.LogMode(logger.Info)
			break
		}
		config.Logger = Default.LogMode(logger.Silent)
	}
	return config
}
