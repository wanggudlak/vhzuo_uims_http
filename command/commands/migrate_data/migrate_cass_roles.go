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
)

var CMDMigrateCassRoles = &command.Command{
	UsageLine: "migrate:cass_roles [command]",
	Short:     "迁移结算系统角色数据",
	Long:      `迁移结算系统已添加的角色数据`,
	PreRun:    func(cmd *command.Command, args []string) { version.ShowShortVersionBanner() },
	Run:       migrateCassRoles,
}

type CassRoles struct {
	ID          int       `gorm:"column:id" json:"id"`
	Name        string    `gorm:"column:name" json:"name"`
	GuardName   string    `gorm:"column:guard_name" json:"guard_name"`
	Platform    int       `gorm:"column:platform" json:"platform"`
	Title       string    `gorm:"column:title" json:"title"`
	Description string    `gorm:"column:description" json:"description"`
	CreatedAt   time.Time `gorm:"column:updated_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
	GroupID     int       `gorm:"column:groupID" json:"groupID"`
}

func init() {
	command.CMD.Register(CMDMigrateCassRoles)
}

func migrateCassRoles(_ *command.Command, args []string) int {
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
		Model(CassRoles{})

	if len(args) > 0 {
		for i := 0; i < len(args); i++ {
			rows = rows.Where(args[i])
		}
	}
	rows = rows.Scan(&cassRoles)

	if rows == nil {
		fmt.Println("获取结算系统角色数据失败")
		return 0
	}
	//go func() {
	//	fmt.Println("协程迁移数据")
	//	time.Sleep(5 * time.Minute)
	//}()
	for _, item := range cassRoles {
		if db.Def().Where("role_name_en = ?", item.Name).Take(&model.Role{}).Error == nil {
			fmt.Println("该数据已存在，请勿重复添加:", item.Name)
			continue
		}
		var client_flag_code string
		var clientID uint
		if item.Platform == 2 {
			client_flag_code = Front
		} else {
			client_flag_code = Back
		}
		clientID = service.GetClientService().GetClientId(map[string]interface{}{"client_flag_code": client_flag_code})
		var saveRole = model.Role{
			ClientID:   int(clientID),
			RoleType:   "C",
			RoleCode:   "",
			RoleNameEN: item.Name,
			RoleNameCN: item.Title,
			IsDel:      "N",
			CommonModel: &model.CommonModel{
				CreatedAt: item.CreatedAt,
				UpdatedAt: item.UpdatedAt,
			},
		}
		err = db.Def().Create(&saveRole).Error
		if err != nil {
			fmt.Println("保存角色基本数据失败", err)
			return 0
		}

	}

	fmt.Println("迁移结算系统角色数据完成")
	return 0
}
