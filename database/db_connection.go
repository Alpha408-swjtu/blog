package database

import (
	"blog/util"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	ormlog "gorm.io/gorm/logger"
)

var (
	blog_mysql      *gorm.DB
	blog_mysql_once sync.Once
	dblog           ormlog.Interface

	blog_redis      *redis.Client
	blog_redis_once sync.Once
)

func init() {
	dblog = ormlog.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		ormlog.Config{
			SlowThreshold: 100 * time.Millisecond, // 慢 SQL 阈值
			LogLevel:      ormlog.Silent,          // Log level，Silent表示不输出日志
			Colorful:      false,                  // 禁用彩色打印
		},
	)
}

func createMysqlDB(dbname, host, user, pass string, port int) *gorm.DB {
	// data source name 是 tester:123456@tcp(localhost:3306)/blog?charset=utf8mb4&parseTime=True&loc=Local
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, dbname) //mb4兼容emoji表情符号
	var err error
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: dblog, PrepareStmt: true}) //启用PrepareStmt，SQL预编译，提高查询效率
	if err != nil {
		util.LogRus.Panicf("connect to mysql use dsn %s failed: %s", dsn, err) //panic() os.Exit(2)
	}
	//设置数据库连接池参数，提高并发性能
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(100) //设置数据库连接池最大连接数
	sqlDB.SetMaxIdleConns(20)  //连接池最大允许的空闲连接数，如果没有sql任务需要执行的连接数大于20，超过的连接会被连接池关闭。
	util.LogRus.Infof("connect to mysql db %s", dbname)
	return db
}

func GetBlogDBConnection() *gorm.DB { //单例
	blog_mysql_once.Do(func() {
		if blog_mysql == nil {
			dbName := "blog"
			viper := util.CreateConfig("mysql")
			host := viper.GetString(dbName + ".host")
			port := viper.GetInt(dbName + ".port")
			user := viper.GetString(dbName + ".user")
			pass := viper.GetString(dbName + ".pass")
			blog_mysql = createMysqlDB(dbName, host, user, pass, port)
		}
	})

	return blog_mysql
}

func createRedisClient(address, passwd string, db int) *redis.Client {
	cli := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: passwd,
		DB:       db,
	})
	if err := cli.Ping().Err(); err != nil {
		util.LogRus.Panicf("connect to redis %d failed %v", db, err)
	} else {
		util.LogRus.Infof("connect to redis %d", db) //能ping成功才说明连接成功
	}
	return cli
}

func GetRedisClient() *redis.Client {
	blog_redis_once.Do(func() {
		if blog_redis == nil {
			viper := util.CreateConfig("redis")
			addr := viper.GetString("addr")
			pass := viper.GetString("pass") //没对该配置项时，viper会赋默认值(即零值)，不会报错
			db := viper.GetInt("db")
			blog_redis = createRedisClient(addr, pass, db)
		}
	})
	return blog_redis
}
