package database

import (
	"blog/util"

	"gorm.io/gorm"
)

type User struct { //默认情况下，驼峰对应的蛇形复数(users)即为数据库里的表名
	Id     int    `gorm:"column:id;primaryKey"` //默认情况下，id为主键
	Name   string `gorm:"column:name"`          //默认情况下，驼峰对应的蛇形即为数据库里的列名。name
	PassWd string `gorm:"column:password"`      //pass_wd
}

func (User) TableName() string { //显示指定表名
	return "user"
}

var (
	_all_user_field = util.GetGormFields(User{})
)

// 根据用户名检索用户
func GetUserByName(name string) *User {
	db := GetBlogDBConnection()
	var user User
	if err := db.Select([]string{"id", "name", "password"}).Where("name=?", name).First(&user).Error; err != nil { //name不存在重复，所以用First即可
		if err != gorm.ErrRecordNotFound { // 如果是用户名不存在，不需要打错误日志
			util.LogRus.Errorf("get password of user %s failed: %s", name, err) // 系统性异常，才打错误日志
		}
		return nil
	}
	return &user
}

// 创建一个用户
func CreateUser(name, pass string) {
	db := GetBlogDBConnection()
	pass = util.Md5(pass)                          //密码经过md5
	user := User{Name: name, PassWd: pass}         //ORM
	if err := db.Create(&user).Error; err != nil { //必须传指针，因为要给user的主键赋值
		util.LogRus.Errorf("create user %s failed: %s", name, err)
	} else {
		util.LogRus.Infof("create user id %d", user.Id)
	}
}

// 删除一个用户
func DeleteUser(name string) {
	db := GetBlogDBConnection()
	if err := db.Where("name=?", name).Delete(User{}).Error; err != nil { //Delete操作必须有where条件
		util.LogRus.Errorf("delete user %s failed: %s", name, err)
	}
}
