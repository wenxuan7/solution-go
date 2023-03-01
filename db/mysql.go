package db

import (
	"database/sql"
	"fmt"
	"github.com/json-iterator/go"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"os"
	"time"
)

var (
	MysqlDB *gorm.DB
)

type mysqlConfig struct {
	Host         string `json:"host"`
	User         string `json:"user"`
	Password     string `json:"password"`
	Database     string `json:"database"`
	MaxIdleConns int    `json:"maxIdleConns"`
	MaxOpenConns int    `json:"maxOpenConns"`
}

func ConnectMysql() {
	var (
		err       error
		bs        []byte
		dia       gorm.Dialector
		mysqlPool *sql.DB
		logger    *zap.Logger
	)

	if bs, err = os.ReadFile("./db/config.json"); err != nil {
		panic(err)
	}

	mysqlConf := &mysqlConfig{}
	if err = jsoniter.Unmarshal(bs, mysqlConf); err != nil {
		panic(err)
	}
	url := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		mysqlConf.User, mysqlConf.Password, mysqlConf.Host, mysqlConf.Database)

	dia = mysql.New(mysql.Config{
		DSN:                       url,   // DSN data source name
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	})

	gConf := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 表名单数
		},
		CreateBatchSize: 100, // 批次创建最大值
	}

	if MysqlDB, err = gorm.Open(dia, gConf); err != nil {
		panic(err)
	}

	if mysqlPool, err = MysqlDB.DB(); err != nil {
		panic(err)
	}
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	mysqlPool.SetMaxIdleConns(mysqlConf.MaxIdleConns)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	mysqlPool.SetMaxOpenConns(mysqlConf.MaxOpenConns)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	mysqlPool.SetConnMaxLifetime(time.Hour)

	if logger, err = zap.NewProduction(); err != nil {
		panic(err)
	}
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()
	sugar.Infof("Connect mysql URL: %s", url)
}
