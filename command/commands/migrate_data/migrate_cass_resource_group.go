package migrate_data

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"uims/command"
	"uims/command/commands/version"
	"uims/internal/model"
	"uims/internal/service"
	"uims/pkg/color"
	"uims/pkg/db"
	"uims/pkg/tool"
)

var CMDMigrateCassResourceGroup = &command.Command{
	UsageLine: "migrate:cass_resource_group",
	Short:     "迁移结算系统资源组数据",
	Long:      `迁移结算系统资源组数据到UIMS系统`,
	PreRun:    func(cmd *command.Command, args []string) { version.ShowShortVersionBanner() },
	Run:       migrateCassResourceGroup,
}

// cass
type CassroleHasPermissions struct {
	PermissionId int `gorm:"column:permission_id" json:"permission_id"`
	RoleId       int `gorm:"column:role_id" json:"role_id"`
}

func init() {
	command.CMD.Register(CMDMigrateCassResourceGroup)
}

func migrateCassResourceGroup(*command.Command, []string) int {
	// role resource group
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(color.Red("连接结算系统数据库失败"))
		}
	}()
	var err error
	cassDBConn := db.Conn("cass")

	cassRoles := []CassRoles{}
	rows := cassDBConn.
		Table("vz_roles").
		Model(CassRoles{}).
		Scan(&cassRoles)

	if rows == nil {
		fmt.Println(color.Red("获取结算系统角色数据失败"))
		return 1
	}

	for _, item := range cassRoles {
		var (
			permissionId    []int
			roleResourceIDs []int
			clientFlagCode  string
		)

		if item.Platform == 2 {
			clientFlagCode = Front
		} else {
			clientFlagCode = Back
		}

		clientId := service.GetClientService().GetClientId(map[string]interface{}{"client_flag_code": clientFlagCode})

		if clientId == 0 {
			fmt.Println(color.Red("获取结算系统客户端 clientId 失败"))
			return 1
		}

		// select role and permission map
		cassDBConn.Table("vz_role_has_permissions").
			Where("role_id = ?", item.ID).
			Pluck("permission_id", &permissionId)

		if permissionId != nil {
			cassPermissions := []CassPermissions{}
			rows := cassDBConn.Table("vz_permissions").
				Model(CassPermissions{}).
				Where(permissionId).
				Scan(&cassPermissions)

			if rows != nil {
				for _, v := range cassPermissions {
					roleMap := make(map[string]interface{})
					roleMap["client_id"] = clientId
					roleMap["org_id"] = OrgId
					roleMap["res_name_en"] = v.Name
					roleMap["res_name_cn"] = v.Title
					resourceIDs, err := getResourceIds(roleMap)

					if resourceIDs != nil && err == nil {
						roleResourceIDs = append(roleResourceIDs, resourceIDs...)
					}
				}
			}
		}

		err = saveResourceGroup(clientId, item.Name, item.Title, roleResourceIDs)
		if err != nil {
			fmt.Println(color.Red("保存结算系统角色资源组数据失败"))
			return 1
		}
	}

	// front resource group
	resFront, msgFront := initResourceGroup(Front, FrontEn, FrontCn)
	if !resFront {
		fmt.Println(color.Red(msgFront))
		return 1
	}

	// back resource group
	resBack, msgBack := initResourceGroup(Back, BackEn, BackCn)
	if !resBack {
		fmt.Println(color.Red(msgBack))
		return 1
	}

	fmt.Println(color.Green("migrate cass resource group success."))
	return 0
}

// init front and back resource group
func initResourceGroup(clientFlagCode string, resGroupEn string, resGroupCn string) (bool, string) {
	clientId := service.GetClientService().GetClientId(map[string]interface{}{"client_flag_code": clientFlagCode})

	if clientId == 0 {
		return false, "获取结算系统客户端 clientId 失败"
	}

	resourceIDs, err := getResourceIds(map[string]interface{}{
		"client_id": clientId,
		"org_id":    OrgId,
	})

	if err != nil {
		return false, fmt.Sprintf("获取资源点 resourceIDs 数据失败: %+v", err)
	}

	err = saveResourceGroup(clientId, resGroupEn, resGroupCn, resourceIDs)
	if err != nil {
		return false, "保存结算系统资源组数据失败"
	}

	return true, "初始化资源组数据成功"
}

// get resource ids
func getResourceIds(whereMaps interface{}) ([]int, error) {
	var (
		resourceData []model.Resource
		resourceIDs  []int
		err          error
	)

	err = db.Def().
		Select("id").
		Where("isdel = ?", "N").
		Where(whereMaps).
		Find(&resourceData).
		Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	for _, item := range resourceData {
		resourceIDs = append(resourceIDs, int(item.ID))
	}

	return resourceIDs, nil
}

// save data resource group
func saveResourceGroup(clientId uint, resGroupEn string, resGroupCn string, resourceIDs []int) (err error) {
	err = db.Def().
		Select("id").
		Where("client_id = ?", clientId).
		Where("org_id = ?", OrgId).
		Where("isdel = ?", "N").
		Where("res_group_en = ?", resGroupEn).
		Where("res_group_cn = ?", resGroupCn).
		First(&model.ResourceGroup{}).
		Error

	// 资源组存在，无需添加
	if err != nil {
		resourceGroup := model.ResourceGroup{
			ResGroupCode: tool.GenXid(),
			ResGroupEn:   resGroupEn,
			ResGroupCn:   resGroupCn,
			ResGroupType: "DEFAULT",
			ClientId:     clientId,
			OrgId:        OrgId,
		}

		if resourceIDs != nil {
			resourceOfCurr := new(model.ResourceOfCurr)
			resourceOfCurr.ResourceIDs = resourceIDs

			resourceGroup.ResOfCurr = resourceOfCurr
		}

		err = db.Def().Create(&resourceGroup).Error
	}

	return
}
