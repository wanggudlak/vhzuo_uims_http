package migrate_data

import (
	"fmt"
	"time"
	"uims/command"
	"uims/command/commands/version"
	"uims/internal/model"
	"uims/internal/service"
	"uims/pkg/color"
	"uims/pkg/db"
	"uims/pkg/tool"
)

var CMDMigrateCassResource = &command.Command{
	UsageLine: "migrate:cass_resource",
	Short:     "迁移结算系统资源点数据",
	Long:      `迁移结算系统资源点数据到UIMS系统`,
	PreRun:    func(cmd *command.Command, args []string) { version.ShowShortVersionBanner() },
	Run:       migrateCassResource,
}

type CassPermissions struct {
	ID              int       `gorm:"column:id" json:"id"`
	Name            string    `gorm:"column:name" json:"name"`
	GuardName       string    `gorm:"column:guard_name" json:"guard_name"`
	Title           string    `gorm:"column:title" json:"title"`
	CreatedAt       time.Time `gorm:"column:updated_at" json:"created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at" json:"updated_at"`
	resDataLocation *model.LocationData
}

var (
	OrgId   uint = 0
	Front        = "VDK_CASS_FRONT"
	FrontEn      = "vzhuo_front"
	FrontCn      = "结算系统前台"
	Back         = "VDK_CASS_BACK"
	BackEn       = "vzhuo_back"
	BackCn       = "结算系统后台"
)

func init() {
	command.CMD.Register(CMDMigrateCassResource)
}

func migrateCassResource(*command.Command, []string) int {
	// front resource
	resFront, msgFront := saveResource(Front, FrontEn)
	if !resFront {
		fmt.Println(color.Red(msgFront))
		return 1
	}

	// back resource
	resBack, msgBack := saveResource(Back, BackEn)
	if !resBack {
		fmt.Println(color.Red(msgBack))
		return 1
	}

	fmt.Println(color.Green("migrate cass resource success."))
	return 0
}

// save data resource
func saveResource(clientFlagCode string, guardName string) (bool, string) {
	clientId := service.GetClientService().GetClientId(map[string]interface{}{"client_flag_code": clientFlagCode})
	if clientId == 0 {
		return false, "获取结算系统客户端ID失败"
	}

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(color.Red("连接结算系统数据库失败"))
		}
	}()
	cassDBConn := db.Conn("cass")

	cassPermissions := []CassPermissions{}
	rows := cassDBConn.
		Table("vz_permissions").
		Where("guard_name = ?", guardName).
		Model(CassPermissions{}).
		Scan(&cassPermissions)

	if rows == nil {
		return false, "获取结算系统资源点数据失败"
	}

	for _, item := range cassPermissions {
		// 资源点存在，执行下一个
		err := db.Def().
			Select("id").
			Where("client_id = ?", clientId).
			Where("org_id = ?", OrgId).
			Where("isdel = ?", "N").
			Where("res_name_en = ?", item.Name).
			Where("res_name_cn = ?", item.Title).
			First(&model.Resource{}).
			Error

		if err == nil {
			continue
		}

		saveData := model.Resource{
			ClientId:     clientId,
			OrgId:        OrgId,
			ResCode:      tool.GenXid(),
			ResFrontCode: tool.GenXid(),
			ResType:      "A",
			ResSubType:   "AM",
			ResNameEn:    item.Name,
			ResNameCn:    item.Title,
			CommonModel: &model.CommonModel{
				CreatedAt: item.CreatedAt,
				UpdatedAt: item.UpdatedAt,
			},
		}

		// json data
		if item.resDataLocation != nil {
			locationData := new(model.LocationData)
			locationData.DataBase = item.resDataLocation.DataBase
			locationData.Table = item.resDataLocation.Table
			locationData.Status = item.resDataLocation.Status

			saveData.ResDataLocation = locationData
		}

		err = db.Def().Create(&saveData).Error
		if err != nil {
			return false, "保存资源点基本数据失败"
		}

	}

	return true, "保存资源点基本数据成功"
}
