package migrate_data

import (
	"fmt"
	"uims/command"
	"uims/command/commands/version"
	"uims/internal/model"
	"uims/internal/service/uuid"
	"uims/pkg/db"
)

var CMDInitUserCode = &command.Command{
	UsageLine: "init:user_code",
	Short:     "初始化空的用户编码",
	Long:      `初始化空的用户编码`,
	PreRun:    func(cmd *command.Command, args []string) { version.ShowShortVersionBanner() },
	Run:       initUserCode,
}

func init() {
	command.CMD.Register(CMDInitUserCode)
}

func initUserCode(_ *command.Command, args []string) int {
	var userData []model.User
	rows := db.Def().
		Table("uims_user_auth")
	if len(args) > 0 {
		for i := 0; i < len(args); i++ {
			rows = rows.Where(args[i])
		}
	}
	rows = rows.Scan(&userData)
	if rows == nil {
		fmt.Println("获取结算系统用户数据失败")
		return 0
	}

	for _, v := range userData {
		if v.UserCode == "" {
			fmt.Println(v.ID, v.UserCode)
			var err error
			var userAuth model.User
			userCode := uuid.GenerateForUIMS().String()
			err = db.Def().Model(&userAuth).
				Where("id = ?", v.ID).
				Update("user_code", userCode).Error

			if err != nil {
				fmt.Println("用户鉴权修改失败：", err)
				return 0
			}

			var userInfo model.UserInfo
			err = db.Def().Model(&userInfo).
				Where("user_id = ?", v.ID).
				Update("user_code", userCode).Error
			if err != nil {
				fmt.Println("用户信息修改失败：", err)
				return 0
			}
		}
	}
	fmt.Println("填充不存在的用户编码成功")
	return 0
}
