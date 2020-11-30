package service

import (
	"errors"
	"github.com/jinzhu/gorm"
	"time"
	"uims/internal/model"
	"uims/internal/service/uuid"
	"uims/pkg/db"
	"uims/pkg/encryption"
	"uims/pkg/glog"
)

type UserService struct {
}

//结算系统请求保存数据的结构体
type PushUserData struct {
	OpenID    string `json:"openID"`
	Account   string `json:"account"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Wechat    string `json:"wechat"`
	WechatID  string `json:"wechatID"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	PasswdRaw string `json:"passwdRaw"`
	Status    int    `json:"status"`
	Type      int    `json:"type"`
	TaxerType int    `json:"taxerType"`
	TaxerNO   string `json:"taxerNO"`
	RoleName  string `json:"roleName"`
	GroupName string `json:"groupName"`
}

// 任务系统请求user数据结构体
type VzhuoUserData struct {
	Phone        string     `json:"phone"`
	Name         string     `json:"name"`
	Sex          string     `json:"sex"`
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

//更新用户数据
func (UserService) UpdateUser(userID int, status int) error {
	var user model.User
	err := db.Def().Find(&user, "id = ? and isdel = ?", userID, "N").Error
	if err != nil {
		return err
	}

	err = db.Def().Transaction(func(tx *gorm.DB) error {
		if status == 1 {
			user.Status = "Y"
		} else {
			user.Status = "N"
		}
		err = tx.Save(&user).Error
		if err != nil {
			return err
		}
		var clientData []model.Client
		err = tx.Where("status = ?", "Y").Find(&clientData).Error
		if err != nil {
			return err
		}
		//广播所有的客户端
		for _, clientInfo := range clientData {
			err = requestUpdateClientUser(tx, int(clientInfo.ID), user.ID, "update_user")
			if err != nil {
				glog.Channel("thrift").WithError(err).Error("推送修改用户数据失败：", clientInfo.ClientFlagCode, err)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (UserService) EmailIsExists(email string) bool {
	user := model.User{}
	db.Def().Select([]string{"id"}).Where(&model.User{Email: email}).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}

func PhoneIsExists(phone string) bool {
	var user model.User
	db.Def().Select([]string{"id"}).Where("phone = ?", phone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}

//通过account获取用户基础信息
func (UserService) GetUserInfoByAccount(pUser *model.User, account string) error {
	err := db.Def().Where("account = ?", account).First(pUser).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	return nil
}

//通过用户ID获取用户基础信息
func (UserService) GetUserInfoByUserID(pUser *model.User, user_id int) error {
	err := db.Def().Where("id = ?", user_id).First(pUser).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	return nil
}

//保存用户信息
func (UserService) SaveCassUserData(bizParamBody PushUserData, tx *gorm.DB) (err error, userID int) {
	userBussType := "normal"
	if bizParamBody.Type == 1 {
		userBussType = "business"
	}
	if bizParamBody.Type == 2 {
		userBussType = "settle_company"
	}
	taxerType := ""
	if bizParamBody.TaxerType == 0 {
		taxerType = "A"
	}

	status := "N"
	if bizParamBody.Status == 0 {
		status = "Y"
	} else {
		status = "N"
	}

	userCode := uuid.GenerateForUIMS().String()
	//手机号重复校验
	var userInfoByPhone model.User
	db.Def().Table("uims_user_auth").Where("phone = ? and account != ? ", bizParamBody.Phone, bizParamBody.Account).First(&userInfoByPhone)
	if userInfoByPhone.ID > 0 {
		return errors.New("手机号重复，请更换手机号"), userInfoByPhone.ID
	}

	var userBaseInfo model.User
	var saveUserAuth = model.User{
		OpenID:   bizParamBody.OpenID,
		UserType: "CASS",
		Account:  bizParamBody.Account,
		UserCode: userCode,
		NaCode:   "+86",
		Phone:    &bizParamBody.Phone,
		Email:    bizParamBody.Email,
		Status:   status,
	}

	if bizParamBody.PasswdRaw != "" {
		saveUserAuth.Passwd, _ = encryption.BcryptHash(bizParamBody.PasswdRaw)
	}

	if err := tx.Where("account = ?", bizParamBody.Account).
		First(&userBaseInfo).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			if err := tx.Create(&saveUserAuth).Error; err != nil {
				return err, saveUserAuth.ID
			}
		} else {
			return err, userBaseInfo.ID
		}
	} else {
		if err := tx.Model(&saveUserAuth).
			Where("account = ?", bizParamBody.Account).
			Update(&saveUserAuth).Error; err != nil {
			return err, saveUserAuth.ID
		}
	}

	var userInfo model.UserInfo
	var saveUserInfo = model.UserInfo{
		UserCode:     userCode,
		IsIdentify:   "N",
		Phone:        bizParamBody.Phone,
		NaCode:       saveUserAuth.NaCode,
		NameCn:       bizParamBody.Name,
		Nickname:     bizParamBody.Name,
		Sex:          "M",
		TaxerNo:      bizParamBody.TaxerNO,
		TaxerType:    taxerType,
		UserID:       saveUserAuth.ID,
		UserType:     saveUserAuth.UserType,
		UserBussType: userBussType,
	}

	if err := tx.Where("user_id = ?", userBaseInfo.ID).
		First(&userInfo).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			if err := tx.Create(&saveUserInfo).Error; err != nil {
				return err, saveUserAuth.ID
			}
		} else {
			return err, userInfo.ID
		}
	} else {
		if err := tx.Model(&saveUserInfo).
			Where("user_id = ?", userBaseInfo.ID).
			Update(&saveUserInfo).Error; err != nil {
			return err, saveUserAuth.ID
		}
	}

	userID = saveUserAuth.ID
	if saveUserAuth.ID == 0 {
		userID = userBaseInfo.ID
	}
	return nil, userID
}

// 设置用户密码, 加密类型为 encrypt_type = 0
func (UserService) SetPassword(userId int, password string) error {
	hash, err := encryption.BcryptHash(password)
	if err != nil {
		return err
	}
	err = db.Def().Model(&model.User{}).Where("id = ?", userId).Updates(map[string]interface{}{
		"salt":         "",
		"encrypt_type": 0,
		"passwd":       hash,
	}).Error
	if err != nil {
		return err
	}
	return nil
}

//请求更新子业务系统用户数据
func requestUpdateClientUser(tx *gorm.DB, clientId int, userId int, method string) error {
	type clientUser struct {
		OpenID         string `json:"open_id"`
		Status         string `json:"status"`
		Phone          string `json:"phone"`
		ClientFlagCode string `json:"client_flag_code"`
		OnlyItem       bool   `json:"only_item"`
	}

	var clientInfo model.Client
	err := tx.Where("id = ?", clientId).Find(&clientInfo).Error
	if err != nil {
		return err
	}
	//目前只有结算系统有需要回写数据
	//if clientInfo.ClientType != "CASS" {
	//	return nil
	//}
	var userInfo model.User
	err = tx.Where("id = ?", userId).Find(&userInfo).Error
	if err != nil {
		return err
	}

	var user = clientUser{
		OpenID:         userInfo.OpenID,
		Status:         userInfo.Status,
		Phone:          *userInfo.Phone,
		ClientFlagCode: clientInfo.ClientFlagCode,
		OnlyItem:       true,
	}
	resp := GetThriftClientServer().
		ClientInvoke(clientId, method, user)
	if !resp.OK() {
		return errors.New(resp.Err())
	}
	return nil
}

func IsExistByPhone(phone string) (bool, error) {
	return isExist(phone, "phone")
}

func IsExistByEmail(email string) (bool, error) {
	return isExist(email, "email")
}

func isExist(fieldV, field string) (bool, error) {
	user := model.User{}
	c := 0
	err := db.Def().
		Select([]string{"id"}).
		Table(user.TableName()).
		Where(field+" = ?", fieldV).
		Where("isdel = ?", "N").
		Where("status = ?", "Y").
		Count(&c).
		Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return false, err
		} else {
			return false, err
		}
	} else {
		if 0 == c {
			return false, nil
		} else {
			return true, nil
		}
	}
}

// 通过 open_id 查询 user_id
func (UserService) OpenID2UserID(openId string) (uint, error) {
	var user model.User
	err := db.Def().Select("id").Where("open_id = ?", openId).First(&user).Error
	if err != nil {
		return 0, err
	}
	return uint(user.ID), nil
}
