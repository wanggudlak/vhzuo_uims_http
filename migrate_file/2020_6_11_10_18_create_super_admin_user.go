package migrate_file

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"uims/internal/model"
	"uims/internal/service/uuid"
	"uims/pkg/db"
	"uims/pkg/tool"
)

type CreateSuperAdminUserMigrate struct {
}

func (CreateSuperAdminUserMigrate) Key() string {
	return "2020_6_11_10_18_create_super_admin_user"
}

func (CreateSuperAdminUserMigrate) Up() (err error) {
	passwd := "yxtkUIMS123"
	hashPwd, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	userCode := uuid.GenerateForUIMS().String()
	openID := tool.GenerateRandStrWithMath(32)
	phone := "17852000001"
	saveBase := model.User{
		OpenID:      openID,
		UserType:    "UIMS",
		Account:     "uims_super_admin",
		UserCode:    userCode,
		NaCode:      "+86",
		Phone:       &phone,
		Email:       "vzhuo@vzhuo.com",
		Salt:        "",
		EncryptType: 0,
		Passwd:      string(hashPwd),
		Status:      "Y",
		Isdel:       "N",
		//CommonModel: &model.CommonModel{
		//	CreatedAt: time.Now(),
		//	UpdatedAt: time.Now(),
		//},
	}
	err = db.Def().Create(&saveBase).Error
	if err != nil {
		fmt.Println("保存超级管理员失败")
		return
	}
	saveInfo := model.UserInfo{
		UserID:                saveBase.ID,
		IsIdentify:            "N",
		UserCode:              userCode,
		UserType:              "VDK",
		UserBussType:          "uims",
		NameEn:                "",
		NameCn:                "uims超级管理员",
		Nickname:              "uims超级管理员",
		NameCnAlias:           "",
		NameAbbrPy:            "",
		NameFullPy:            "",
		IdentityCardNo:        nil,
		NaCode:                "+86",
		Phone:                 "17852000001",
		LandlinePhone:         "",
		Sex:                   "M",
		TaxerType:             "",
		TaxerNo:               "",
		HeaderImgURL:          "",
		IdentityCardPersonImg: "",
		IdentityCardEmblemImg: "",
		Isdel:                 "N",
	}
	err = db.Def().Create(&saveInfo).Error
	if err != nil {
		fmt.Println("保存超级管理员详细数据失败")
		return
	}

	return
}

func (CreateSuperAdminUserMigrate) Down() (err error) {
	db.Def().Where("account = ?", "uims_super_admin").Delete(&model.User{})
	return
}
