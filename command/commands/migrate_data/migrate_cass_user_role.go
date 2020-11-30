package migrate_data

import (
	"fmt"
	"uims/command"
	"uims/command/commands/version"
	"uims/internal/model"
	"uims/pkg/color"
	"uims/pkg/db"
)

var CMDMigrateCassUserRole = &command.Command{
	UsageLine: "migrate:cass_user_role [command]",
	Short:     "迁移结算系统用户角色关联数据",
	Long:      `迁移结算系统已添加的用户角色关联数据`,
	PreRun:    func(cmd *command.Command, args []string) { version.ShowShortVersionBanner() },
	Run:       migrateCassUserRole,
}

type CassModelHasRole struct {
	RoleId    int    `gorm:"column:role_id" json:"role_id"`
	ModelType string `gorm:"column:model_type" json:"model_type"`
	ModelId   int    `gorm:"column:model_id" json:"model_id"`
}

func init() {
	command.CMD.Register(CMDMigrateCassUserRole)
}

func migrateCassUserRole(_ *command.Command, args []string) int {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(color.Red("连接结算系统数据库失败"))
		}
	}()
	var err error
	cassDBConn := db.Conn("cass")

	modelHasRole := []CassModelHasRole{}
	rows := cassDBConn.
		Table("vz_model_has_roles").
		Model(CassModelHasRole{})
	if len(args) > 0 {
		for i := 0; i < len(args); i++ {
			rows = rows.Where(args[i])
		}
	}
	rows = rows.Scan(&modelHasRole)

	if rows == nil {
		fmt.Println("获取结算系统用户数据失败")
		return 0
	}
	//go func() {
	//	fmt.Println("协程迁移数据")
	//	time.Sleep(5 * time.Minute)
	//}()
	for _, item := range modelHasRole {

		var userID int
		var userRelationType string
		var uimsRoleInfo model.Role

		//获取roleID获取角色信息
		var roleInfo CassRoles
		cassDBConn.Table("vz_roles").
			Model(CassRoles{}).
			Where("id = ?", item.RoleId).
			First(&roleInfo)
		db.Def().Table("uims_role").Where("role_name_en = ?", roleInfo.Name).First(&uimsRoleInfo)

		//根据modelID获取用户信息
		if item.ModelType == "App\\Models\\Mysql\\User" {
			var userInfo CassUser
			var uimsUserInfo model.User
			cassDBConn.Table("vz_user").
				Model(CassUser{}).
				Where("id = ?", item.ModelId).
				Scan(&userInfo)
			db.Def().Table("uims_user_auth").Where("account = ?", userInfo.Account).First(&uimsUserInfo)
			userID = uimsUserInfo.ID
			if userID == 0 {
				continue
			}
			userRelationType = "user"
		}

		if item.ModelType == "App\\Models\\Mysql\\Group" {
			var orgInfo CassGroup
			var uimsOrgInfo model.Org
			cassDBConn.Table("vz_group").
				Model(CassUser{}).
				Where("id = ?", item.ModelId).
				Scan(&orgInfo)
			db.Def().Table("uims_organization").Where("org_name_cn = ?", orgInfo.Name).First(&uimsOrgInfo)
			userID = uimsOrgInfo.ID
			userRelationType = "org"
		}

		var userRole model.UserRole
		db.Def().Table("uims_user_role").
			Where("user_id = ? AND role_id = ? AND client_id = ? AND user_relation_type = ?", userID, uimsRoleInfo.ID, uimsRoleInfo.ClientID, userRelationType).
			First(&userRole)
		if userRole.ID != 0 {
			continue
		}

		var saveModelHasRoleData = model.UserRole{
			UserID:           userID,
			RoleID:           uimsRoleInfo.ID,
			ClientID:         uimsRoleInfo.ClientID,
			UserRelationType: userRelationType,
		}
		err = db.Def().Create(&saveModelHasRoleData).Error
		if err != nil {
			fmt.Println("保存用户角色关联数据失败", err)
			return 0
		}
	}

	fmt.Println("迁移结算系统用户角色关联数据完成")
	return 0
}
