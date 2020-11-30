package migrate_data

import (
	"fmt"
	"time"
	"uims/command"
	"uims/command/commands/version"
	"uims/internal/model"
	"uims/pkg/db"
)

var CMDInitRoleResMap = &command.Command{
	UsageLine: "init:role_res [command]",
	Short:     "初始化角色资源组数据",
	Long:      `初始化uims角色资源组数据`,
	PreRun:    func(cmd *command.Command, args []string) { version.ShowShortVersionBanner() },
	Run:       initRoleResMap,
}

func init() {
	command.CMD.Register(CMDInitRoleResMap)
}

func initRoleResMap(_ *command.Command, args []string) int {

	//go func() {
	//	fmt.Println("协程迁移数据")
	//	time.Sleep(5 * time.Minute)
	//}()

	var resGroup []model.ResourceGroup
	rows := db.Def().
		Table("uims_res_group")
	if len(args) > 0 {
		for i := 0; i < len(args); i++ {
			rows = rows.Where(args[i])
		}
	}
	rows = rows.Scan(&resGroup)
	if rows == nil {
		fmt.Println("获取结算系统资源组数据失败")
		return 0
	}

	for _, item := range resGroup {

		var uimsRoleInfo model.Role
		db.Def().Table("uims_role").Where("role_name_en = ?", item.ResGroupEn).First(&uimsRoleInfo)
		if uimsRoleInfo.ID == 0 {
			continue
		}

		var roleResMap model.RoleResMap
		db.Def().Table("uims_role_res_map").
			Where("role_id = ? AND res_grp_id = ? AND client_id = ?", uimsRoleInfo.ID, int(item.ID), uimsRoleInfo.ClientID).
			First(&roleResMap)
		if roleResMap.ID != 0 {
			continue
		}
		var saveRoleRes = model.RoleResMap{
			RoleID:       uimsRoleInfo.ID,
			ResGrpID:     int(item.ID),
			ClientID:     uimsRoleInfo.ClientID,
			OrgId:        uimsRoleInfo.OrgID,
			StartValidAt: time.Time{},
			IsDel:        "N",
			ForgetAt:     time.Time{},
		}
		err := db.Def().Create(&saveRoleRes).Error
		if err != nil {
			fmt.Println("保存角色关联资源组数据失败", err)
			return 0
		}

	}
	fmt.Println("初始化角色资源组数据完成")
	return 0
}
