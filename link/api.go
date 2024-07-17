package link

import (
	"database/sql"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

var MysqlDb *sql.DB

func Mysql() {
	gormCon := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
	}
	conn, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn(), // DSN data source name
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), gormCon)
	if err != nil {
		panic(err)
	}

	MysqlDb, err := conn.DB()
	if err != nil {
		panic(err)
	}

	MysqlDb.SetMaxIdleConns(10)
	MysqlDb.SetMaxOpenConns(100)
	MysqlDb.SetConnMaxLifetime(time.Hour)
}

var RedisDb *redis.Client

func Redis() {
	redisConf := newRedisConf()
	RedisDb = redis.NewClient(&redis.Options{
		Addr:     redisConf.Addr,
		Password: redisConf.Password, // no password set
		DB:       0,                  // use default DB
	})
}
