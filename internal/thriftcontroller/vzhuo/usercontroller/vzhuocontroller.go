package usercontroller

import (
	"github.com/cao-guang/pinyin"
	"github.com/jinzhu/gorm"
	"time"
	"uims/internal/model"
	"uims/internal/service/uuid"
	"uims/pkg/db"
	"uims/pkg/encryption"
	thriftserver "uims/pkg/thrift/server"
)

// 任务系统请求user数据结构体
type VzhuoUserData struct {
	Phone        string     `json:"phone"`
	Name         string     `json:"name"`
	Sex          int        `json:"sex"`
	Nickname     string     `json:"nickname"`
	HeaderImgURL string     `json:"header_img_url"`
	Birthday     *time.Time `json:"birthday" format:"ISO 8601"`
}

// 任务系统的修改身份证照片等资料
type VzhuoUserIdCardData struct {
	Phone          string `json:"phone"`
	IdCard1        string `json:"id_card_1"`        // 人像面
	IdCard2        string `json:"id_card_2"`        // 国徽面
	IdentityCardNo string `json:"identity_card_no"` //身份证号

}

// 任务系统请求重置密码数据结构体
type VzhuoPasswordData struct {
	Phone    string `json:"phone"`
	Password string `json:"password"` //密码
}

// 任务系统创建的user同步到uims,根据手机号
type VzhuoUser struct {
	Phone    string `json:"phone"`
	Nickname string `json:"nickname"`
	OpenId   string `json:"open_id"`
}

//修改用户信息
func UpdateUserInfo(c *thriftserver.Context) {
	var req VzhuoUserData
	var err error
	if err = c.ShouldBind(&req); err != nil {
		c.Response.BadParams(err)
		return
	}
	var user_info model.UserInfo

	err = db.Def().Where("phone = ?", req.Phone).First(&user_info).Error
	if err != nil {
		c.Response.Error(err)
		return
	}
	user_info.Phone = req.Phone
	user_info.NameCn = req.Name
	//中文转换为拼音
	pinyin.LoadingPYFileName("conf/pinyin.txt")
	name_full_py, err := pinyin.To_Py(req.Name, "", "") //默认造型： hanzipinyin
	if err != nil {
		c.Response.Error(err)
		return
	}
	name_abbr_py := string([]byte(name_full_py)[:1])
	user_info.NameAbbrPy = name_abbr_py
	user_info.HeaderImgURL = req.HeaderImgURL
	user_info.Birthday = req.Birthday
	user_info.Nickname = req.Nickname
	//user_info.Sex =
	var user_sex string
	if req.Sex == 1 {
		user_sex = "M"
	} else {
		user_sex = "F"
	}
	user_info.Sex = user_sex
	err = db.Def().Save(&user_info).Error
	if err != nil {
		c.Response.Error(err)
		return
	}

	c.Response.Success(nil, "")
	return
}

// 保存微桌系统用户身份证信息
func UpdateUserIdentity(c *thriftserver.Context) {
	var req VzhuoUserIdCardData
	var err error
	if err = c.ShouldBind(&req); err != nil {
		c.Response.BadParams(err)
		return
	}
	var user_info model.UserInfo
	err = db.Def().Where("phone = ?", req.Phone).First(&user_info).Error
	if err != nil {
		c.Response.BadParams(err)
		return
	}
	var userInfo_by_identity model.UserInfo
	err = db.Def().Where("identity_card_no = ? and phone != ?", req.IdentityCardNo, req.Phone).First(&userInfo_by_identity).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		c.Response.BadParams(err)
		return
	}

	if userInfo_by_identity.ID != 0 {
		c.Response.Fail("身份证号已被绑定，请更换身份证号")
		return
	}
	user_info.IdentityCardNo = &req.IdentityCardNo
	user_info.IdentityCardPersonImg = req.IdCard1
	user_info.IdentityCardEmblemImg = req.IdCard2
	user_info.IsIdentify = "Y"
	err = db.Def().Transaction(func(tx *gorm.DB) error {
		err = tx.Save(&user_info).Error
		if err != nil {
			return err
		}
		err = tx.Table("uims_user_auth").Where("id = ?", user_info.UserID).Update("identity_card_no", req.IdentityCardNo).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		c.Response.BadParams(err)
		return
	}
	c.Response.Success(nil, "")
	return
}

//修改用户密码
func UpdateUserPassword(c *thriftserver.Context) {
	var req VzhuoPasswordData
	var err error
	if err = c.ShouldBind(&req); err != nil {
		c.Response.BadParams(err)
		return
	}
	var user model.User

	err = db.Def().Where("phone = ?", req.Phone).First(&user).Error
	if err != nil {
		c.Response.Error(err)
		return
	}

	// 生成密码
	user.Passwd, err = encryption.BcryptHash(req.Password)
	user.EncryptType = 0
	if err != nil {
		c.Response.Error(err)
		return
	}

	err = db.Def().Save(&user).Error
	if err != nil {
		c.Response.Error(err)
		return
	}

	c.Response.Success(nil, "")
	return
}

// 根据微桌后台系统推送的手机号创建user
func SaveVzhuoUser(c *thriftserver.Context) {
	var req VzhuoUser
	var err error
	if err = c.ShouldBind(&req); err != nil {
		c.Response.BadParams(err)
		return
	}
	var uims_user model.User
	err = db.Def().Where("phone = ?", req.Phone).First(&uims_user).Error
	if err == nil {
		c.Response.Fail("用户已存在")
		return
	}
	userCode := uuid.GenerateForUIMS().String()

	user := model.User{
		Phone:    &req.Phone,
		OpenID:   req.OpenId,
		UserCode: userCode,
		UserType: "MP",
	}

	err = db.Def().Create(&user).Error
	if err != nil {
		//c.Response.Error(err)
		c.Response.Fail("uims创建用户失败")
		return
	}

	user_info := model.UserInfo{
		UserID:   uims_user.ID,
		Phone:    req.Phone,
		Nickname: req.Nickname,
		UserCode: userCode,
		UserType: "MP",
	}
	err = db.Def().Create(&user_info).Error
	if err != nil {
		//c.Response.Error(err)
		c.Response.Fail("uims创建用户失败")
		return
	}
	c.Response.Success(nil, "")
	return
}
