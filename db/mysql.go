package db

import (
	"database/sql"
	"fmt"
	"github.com/json-iterator/go"
	"github.com/solution-go/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"os"
	"path"
	"runtime"
	"time"
)

var (
	MysqlCli *gorm.DB
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
		err error
		dia gorm.Dialector
	)

	mysqlConf := getMysqlConfig()
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
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	}

	if MysqlCli, err = gorm.Open(dia, gConf); err != nil {
		log.Sugar.Panic(err)
	}

	// 初始化连接池参数
	initConnectPool(MysqlCli, mysqlConf)
	log.Sugar.Infof("Connect mysql URL: %s", url)
}

// getMysqlConfig 获取配置文件
// 密码 mysql连接池大小 数据库名
func getMysqlConfig() *mysqlConfig {
	var (
		bs      []byte
		currDir string
		err     error
	)

	_, currFile, _, _ := runtime.Caller(0)
	currDir = path.Dir(currFile)
	if bs, err = os.ReadFile(currDir + "/config.json"); err != nil {
		log.Sugar.Panic(err)
	}

	mysqlConf := &mysqlConfig{}
	if err = jsoniter.Unmarshal(bs, mysqlConf); err != nil {
		log.Sugar.Panic(err)
	}

	return mysqlConf
}

// initConnectPool 初始化连接池参数
func initConnectPool(db *gorm.DB, mysqlConf *mysqlConfig) {
	var (
		mysqlPool *sql.DB
		err       error
	)

	if mysqlPool, err = db.DB(); err != nil {
		log.Sugar.Panic(err)
	}
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	mysqlPool.SetMaxIdleConns(mysqlConf.MaxIdleConns)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	mysqlPool.SetMaxOpenConns(mysqlConf.MaxOpenConns)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	mysqlPool.SetConnMaxLifetime(time.Hour)
}
