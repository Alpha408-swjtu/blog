package test

import (
	"blog/util"
	"fmt"
	"path"
	"runtime"
	"testing"
	"time"
)

func TestProjectRootPath(t *testing.T) {
	fmt.Println(util.ProjectRootPath)
}

func TestConfig(t *testing.T) {
	dbViper := util.CreateConfig("mysql")
	dbViper.WatchConfig() //确保在调用WatchConfig()之前添加了所有的配置路径(AddConfigPath)
	//读取配置的第一种方式
	if !dbViper.IsSet("blog.port") { //检查有没有此项配置
		t.Fail()
	}
	port := dbViper.GetInt("blog.port")
	fmt.Println("port", port)
	time.Sleep(10 * time.Second) //10秒之内修改一下配置文件，看看viper能不能读取最新值
	port = dbViper.GetInt("blog.port")
	fmt.Println("port", port)

	//读取配置的第二种方式
	logViper := util.CreateConfig("log")
	logViper.WatchConfig()
	type LogConfig struct {
		Level string `mapstructure:"level"` //Tag
		File  string `mapstructure:"file"`
	}
	var config LogConfig
	if err := logViper.Unmarshal(&config); err != nil {
		fmt.Println(err)
		t.Fail()
	} else {
		fmt.Println(config.Level)
		fmt.Println(config.File)
	}
}

func af() (string, int) {
	_, filename, line, _ := runtime.Caller(2) //参数换成0、1、2试试
	return filename, line
}

func bf() (string, int) {
	return af()
}

// TestCaller() -> bf()  ->  af()     56 -> 51 -> 46
func TestCaller(t *testing.T) {
	filename, line := bf()
	fmt.Println(filename, line)
	fmt.Println(path.Dir(filename) + "/../../")
	fmt.Println(path.Dir(path.Dir(filename) + "/../../"))
}

// go test -v .\util\test\ -run=^TestCaller$ -count=1
// go test -v .\util\test\ -run=^TestProjectRootPath$ -count=1
// go test -v .\util\test\ -run=^TestConfig$ -count=1
