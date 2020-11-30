package command

//import (
//	"fmt"
//	"time"
//	"uims/boot"
//	"uims/db"
//	"uims/db/redis"
//	"uims/internal/model"
//)
//
//func init() {
//	boot.SetInCommand()
//	boot.Bootstrap()
//}
//
//// 命令的编写需要调用上边的 init() 函数加载框架
//// 使用自定义的mysql连接可以使用 db.New() 来创建连接使用
//func Main1() {
//	// 调用 mysql
//	u := model.User{}
//	db.Def().Order("id desc").First(&u)
//	fmt.Printf("users: %+v", u)
//
//	// 调用 redis
//	err := redis.Def().Set("test", "test", 0*time.Second).Err()
//	if err != nil {
//		panic(err.Error())
//	}
//}
