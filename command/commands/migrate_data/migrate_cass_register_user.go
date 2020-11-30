package migrate_data

import (
	"fmt"
	"time"
	"uims/command"
	"uims/command/commands/version"
	"uims/internal/model"
	"uims/internal/service/uuid"
	"uims/pkg/color"
	"uims/pkg/db"
	"uims/pkg/tool"
)

var CMDMigrateCassRegisterUser = &command.Command{
	UsageLine: "migrate:cass_register_user [command]", //eg:./uims migrate:cass_register_user name=\"北京鲁拓科技有限公司\"
	Short:     "迁移结算系统注册用户数据",
	Long:      `迁移结算系统已注册的用户数据`,
	PreRun:    func(cmd *command.Command, args []string) { version.ShowShortVersionBanner() },
	Run:       migrateCassRegisterUser,
}

type CassUser struct {
	ID             int    `gorm:"column:ID" json:"ID"`
	OpenID         string `gorm:"column:openID" json:"openID"`
	Account        string `gorm:"column:account" json:"account"`
	Name           string `gorm:"column:name" json:"name"`
	Email          string `gorm:"column:email" json:"email"`
	Wechat         string `gorm:"column:wechat" json:"wechat"`
	Password       string `gorm:"column:password" json:"password"`
	RemitSK        string `gorm:"column:remitSK" json:"remitSK"`
	RoleK          string `gorm:"column:roleK" json:"roleK"`
	RoleCN         string `gorm:"column:roleCN" json:"roleCN"`
	Phone          string `gorm:"column:phone" json:"phone"`
	Status         int    `gorm:"column:status" json:"status"`
	SmsReceiveType int    `gorm:"column:smsReceiveType" json:"smsReceiveType"`
	WechatID       string `gorm:"column:wechatID" json:"wechatID"`
	CreatedAt      int64  `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt      int64  `gorm:"column:updatedAt" json:"updatedAt"`
	IsDel          int    `gorm:"column:isDel" json:"isDel"`
	BcUUID         string `gorm:"column:bcUUID" json:"bcUUID"`
	ParentID       int    `gorm:"column:parentID" json:"parentID"`
	GroupID        int    `gorm:"column:groupID" json:"groupID"`
}

type BusinessCustomer struct {
	ID   int `gorm:"column:ID" json:"ID"`
	UID  int `gorm:"column:uid" json:"uid"`
	Type int `gorm:"column:type" json:"type"`
}

func init() {
	command.CMD.Register(CMDMigrateCassRegisterUser)
}

func migrateCassRegisterUser(_ *command.Command, args []string) int {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(color.Red(fmt.Sprintf("连接结算系统数据库失败: %+v", err)))
		}
	}()
	var err error
	cassDBConn := db.Conn("cass")

	cassUsers := []CassUser{}
	rows := cassDBConn.
		Table("vz_user").
		Model(CassUser{})

	if len(args) > 0 {
		for i := 0; i < len(args); i++ {
			rows = rows.Where(args[i])
		}
	}
	rows = rows.Scan(&cassUsers)

	if rows == nil {
		fmt.Println("获取结算系统用户数据失败")
		return 0
	}

	//go func() {
	//	fmt.Println("协程迁移数据")
	//	time.Sleep(5 * time.Minute)
	//}()
	for _, item := range cassUsers {

		if db.Def().Where("account = ?", item.Account).Take(&model.User{}).Error == nil {
			fmt.Println("该数据已存在，请勿重复添加:", item.Account)
			continue
		}
		var status, isDel string
		if item.Status == 0 {
			status = "Y"
		} else {
			status = "N"
		}

		if item.IsDel == 0 {
			isDel = "N"
		} else {
			isDel = "Y"
		}

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

		//获取用户类型
		var bussUInfo BusinessCustomer
		var userBussType string
		cassDBConn.Table("vz_business_customer").
			Where("uid = ?", item.ID).
			Scan(&bussUInfo)
		userBussType = "back"
		if bussUInfo.Type == 1 {
			userBussType = "business"
		}
		if bussUInfo.Type == 2 {
			userBussType = "settle_company"
		}

		userCode := uuid.GenerateForUIMS().String()
		//操作员
		if item.ParentID != 0 {
			userBussType = "business"
		}

		var user = model.User{
			UserType:    "CASS",
			OpenID:      item.OpenID,
			Account:     item.Account,
			UserCode:    userCode,
			NaCode:      "+86",
			Phone:       &item.Phone,
			Email:       item.Email,
			Salt:        "",
			EncryptType: 0,
			Passwd:      item.Password,
			Status:      status,
			Isdel:       isDel,
			CommonModel: &model.CommonModel{
				CreatedAt: createdTime,
				UpdatedAt: updatedTime,
			},
		}
		err = db.Def().Create(&user).Error
		if err != nil {
			fmt.Println("保存用户基本数据失败", err)
			return 0
		}

		var userInfo = model.UserInfo{
			UserCode:     userCode,
			IsIdentify:   "N",
			Isdel:        isDel,
			Phone:        item.Phone,
			NaCode:       user.NaCode,
			NameCn:       item.Name,
			Nickname:     item.Name,
			Sex:          "M",
			TaxerType:    "A",
			UserID:       user.ID,
			UserType:     user.UserType,
			UserBussType: userBussType,
			CommonModel: &model.CommonModel{
				CreatedAt: createdTime,
				UpdatedAt: updatedTime,
			},
		}
		err = db.Def().Create(&userInfo).Error
		if err != nil {
			fmt.Println("保存用户额外数据失败", err)
			return 0
		}
		// 创建 user_wechat 数据
		if item.WechatID != "" {
			userWeChat := model.UserWeChat{
				UserId:       uint(user.ID),
				WeChatId:     1,
				WeChatOpenId: item.WechatID,
				CommonModel:  nil,
			}
			err = db.Def().Create(&userWeChat).Error
			if err != nil {
				fmt.Printf("保存 user_wechat 失败: [%+v] \n", err)
				return 0
			}
		}
	}

	fmt.Println("迁移结算系统注册用户数据完成")
	return 0
}
