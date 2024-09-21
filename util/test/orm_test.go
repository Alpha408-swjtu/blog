package test

import (
	"blog/util"
	"fmt"
	"reflect"
	"testing"
	"time"
)

type MyStruct struct {
	Id         int    `gorm:"column:id;primaryKey"`          //Tag
	PassWd     string `json:"passwd" gorm:"column:password"` //需要做映射的Field必须是可导出的，否则从DB里查询出来不能给结构体赋值
	Name       string
	FamilyName string    `gorm:"-"` //family_name
	CreateTime time.Time `form:"create_time" binding:"required,before_today" time_format:"2006-01-02" time_utc:"8"`
	int
	gender int
}

// 打印结构体的成员变量信息
func PrintFieldInfo(object any) {
	tp := reflect.TypeOf(object) //通过reflect.Type获取类型相关的信息
	fieldNum := tp.NumField()    //成员变量的个数，包括未导出成员
	for i := 0; i < fieldNum; i++ {
		field := tp.Field(i)
		fmt.Printf("%d %s offset %d anonymous %t type %s exported %t gorm tag=%s json tag=%s\n", i,
			field.Name,            //变量名称
			field.Offset,          //相对于结构体首地址的内存偏移量，string类型会占据16个字节
			field.Anonymous,       //是否为匿名成员
			field.Type,            //数据类型，reflect.Type类型
			field.IsExported(),    //包外是否可见（即是否以大写字母开头）
			field.Tag.Get("gorm"), //获取成员变量后面``里面定义的gorm
			field.Tag.Get("json"), //获取成员变量后面``里面定义的tag
		)
	}
	fmt.Println()
}

func TestPrintFieldInfo(t *testing.T) {
	PrintFieldInfo(MyStruct{})
}

func TestGetGormFields(t *testing.T) {
	var p MyStruct //不管p是不是指针，都没问题
	// var p *MyStruct
	cols := util.GetGormFields(p)
	fmt.Println(cols)
}

// go test -v ./util/test/ -run=^TestPrintFieldInfo$ -count=1
// go test -v ./util/test/ -run=^TestGetGormFields$ -count=1
