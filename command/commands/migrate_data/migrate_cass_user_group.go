package migrate_data

import (
	"fmt"
	"time"
	"uims/command"
	"uims/command/commands/version"
	"uims/internal/model"
	"uims/pkg/color"
	"uims/pkg/db"
	"uims/pkg/tool"
)

var CMDMigrateCassUserGroup = &command.Command{
	UsageLine: "migrate:cass_user_group [command]",
	Short:     "迁移结算系统用户关联组数据",
	Long:      `迁移结算系统已添加的用户关联组数据`,
	PreRun:    func(cmd *command.Command, args []string) { version.ShowShortVersionBanner() },
	Run:       migrateCassUserGroup,
}

func init() {
	command.CMD.Register(CMDMigrateCassUserGroup)
}

func migrateCassUserGroup(_ *command.Command, args []string) int {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(color.Red("连接结算系统数据库失败"))
		}
	}()
	var err error
	cassDBConn := db.Conn("cass")

	cassUser := []CassUser{}
	rows := cassDBConn.
		Table("vz_user").
		Model(CassUser{})
	if len(args) > 0 {
		for i := 0; i < len(args); i++ {
			rows = rows.Where(args[i])
		}
	}
	rows = rows.Scan(&cassUser)

	if rows == nil {
		fmt.Println("获取结算系统用户数据失败")
		return 0
	}
	//go func() {
	//	fmt.Println("协程迁移数据")
	//	time.Sleep(5 * time.Minute)
	//}()
	for _, item := range cassUser {

		var createdTime, updatedTime time.Time
		if item.CreatedAt == 0 {
			createdTime = time.Time{}
		} else {
			createdTime = tool.BigIntConvertTime(int(item.CreatedAt))
		}
		if item.UpdatedAt == 0 {
			updatedTime = time.Time{}
		} else {
			updatedTime = tool.BigIntConvertTime(int(item.UpdatedAt))
		}

		var uimsUserInfo model.User
		db.Def().Table("uims_user_auth").Where("account = ?", item.Account).First(&uimsUserInfo)
		if uimsUserInfo.ID == 0 {
			continue
		}

		var orgInfo CassGroup
		var uimsOrgInfo model.Org
		cassDBConn.Table("vz_group").
			Model(CassUser{}).
			Where("id = ?", item.GroupID).
			Scan(&orgInfo)
		db.Def().Table("uims_organization").Where("org_name_cn = ?", orgInfo.Name).First(&uimsOrgInfo)
		if uimsOrgInfo.ID == 0 {
			continue
		}

		var uimsUserRoleInfo model.UserRole
		db.Def().Table("uims_user_role").Where("user_id = ?", uimsUserInfo.ID).First(&uimsUserRoleInfo)
		if uimsUserRoleInfo.ID == 0 {
			continue
		}

		var userOrg model.UserOrg
		db.Def().Table("uims_user_org").
			Where("user_id = ? AND client_id = ? AND org_id = ?", uimsUserInfo.ID, uimsUserRoleInfo.ClientID, uimsOrgInfo.ID).
			First(&userOrg)
		if userOrg.ID != 0 {
			continue
		}

		var saveUserOrg = model.UserOrg{
			UserID:   uimsUserInfo.ID,
			ClientID: uimsUserRoleInfo.ClientID,
			OrgID:    uimsOrgInfo.ID,
			CommonModel: &model.CommonModel{
				CreatedAt: createdTime,
				UpdatedAt: updatedTime,
			},
		}
		err = db.Def().Create(&saveUserOrg).Error
		if err != nil {
			fmt.Println("保存用户关联组数据失败", err)
			return 0
		}

	}

	fmt.Println("迁移结算系统用户关联组数据完成")
	return 0
}
