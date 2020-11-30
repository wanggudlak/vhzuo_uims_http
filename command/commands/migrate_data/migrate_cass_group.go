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

var CMDMigrateCassGroup = &command.Command{
	UsageLine: "migrate:cass_group [command]",
	Short:     "迁移结算系统组数据",
	Long:      `迁移结算系统已添加的组数据`,
	PreRun:    func(cmd *command.Command, args []string) { version.ShowShortVersionBanner() },
	Run:       migrateCassGroup,
}

type CassGroup struct {
	ID          int    `gorm:"column:ID" json:"ID"`
	GroupUUID   string `gorm:"column:groupUUID" json:"groupUUID"`
	Name        string `gorm:"column:name" json:"name"`
	UserID      int    `gorm:"column:userID" json:"userID"`
	ParentID    int    `gorm:"column:parentID" json:"parentID"`
	Platform    int    `gorm:"column:platform" json:"platform"`
	Description string `gorm:"column:description" json:"description"`
	CreatedAt   int    `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt   int    `gorm:"column:updatedAt" json:"updatedAt"`
}

func init() {
	command.CMD.Register(CMDMigrateCassGroup)
}

func migrateCassGroup(_ *command.Command, args []string) int {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(color.Red("连接结算系统数据库失败"))
		}
	}()
	var err error
	cassDBConn := db.Conn("cass")

	cassGroup := []CassGroup{}
	rows := cassDBConn.
		Table("vz_group").
		Model(CassGroup{})

	if len(args) > 0 {
		for i := 0; i < len(args); i++ {
			rows = rows.Where(args[i])
		}
	}
	rows = rows.Scan(&cassGroup)

	if rows == nil {
		fmt.Println("获取结算系统组数据失败")
		return 0
	}
	//go func() {
	//	fmt.Println("协程迁移数据")
	//	time.Sleep(5 * time.Minute)
	//}()
	var clientInfo model.Client
	if err := db.Def().Model(&clientInfo).
		Where("client_flag_code = ?", "VDK_CASS_BACK").
		First(&clientInfo).Error; err != nil {
		return 0
	}

	for _, item := range cassGroup {
		if db.Def().Where("org_name_cn = ?", item.Name).Take(&model.Org{}).Error == nil {
			fmt.Println("该数据已存在，请勿重复添加:", item.Name)
			continue
		}
		var createdTime, updatedTime time.Time
		if item.CreatedAt == 0 {
			createdTime = time.Time{}
		} else {
			createdTime = tool.BigIntConvertTime(item.CreatedAt)
		}
		if item.UpdatedAt == 0 {
			updatedTime = time.Time{}
		} else {
			updatedTime = tool.BigIntConvertTime(item.UpdatedAt)
		}

		var orgLevel int
		if item.ParentID == 0 {
			orgLevel = 0
		} else {
			orgLevel = 1
		}

		var saveOrg = model.Org{
			OrgCode:     item.GroupUUID,
			ParentOrgID: item.ParentID,
			ClientID:    clientInfo.ID,
			ClientAppID: clientInfo.AppId,
			OrgNameCN:   item.Name,
			OrgLevel:    orgLevel,
			CommonModel: &model.CommonModel{
				CreatedAt: createdTime,
				UpdatedAt: updatedTime,
			},
		}
		err = db.Def().Create(&saveOrg).Error
		if err != nil {
			fmt.Println("保存角色基本数据失败", err)
			return 0
		}

	}

	fmt.Println("迁移结算系统组数据完成")
	return 0
}
