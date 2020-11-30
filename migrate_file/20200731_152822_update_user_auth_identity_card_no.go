package migrate_file

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"uims/internal/model"
	"uims/pkg/db"
)

type UpdateUserAuthIdCard struct {
}

func (UpdateUserAuthIdCard) Key() string {
	return "20200731_152822_update_user_auth_identity_card_no"
}

func (UpdateUserAuthIdCard) Up() (err error) {
	err = db.Def().Table("uims_user_info").Where("identity_card_no = ? or identity_card_no = ?", "网络营销", "[\"会计\"").Update("identity_card_no", "").Error
	if err != nil {
		fmt.Println(err)
	}
	//重复身份证号的，写入临时表,再删除
	var userIDCardCommonData []model.UserInfo
	err = db.Def().Table("uims_user_info").Raw("SELECT * from uims_user_info where identity_card_no in (select identity_card_no from uims_user_info WHERE identity_card_no != '' group by identity_card_no having count(identity_card_no) > 1)").Scan(&userIDCardCommonData).Error
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range userIDCardCommonData {
		userInfoCursor := model.UserInfoCursor{
			UserInfoID:            v.ID,
			UserID:                v.UserID,
			IsIdentify:            v.IsIdentify,
			UserCode:              v.UserCode,
			UserType:              v.UserType,
			UserBussType:          v.UserBussType,
			NameEn:                v.NameEn,
			NameCn:                v.NameCn,
			NameCnAlias:           v.NameCnAlias,
			NameAbbrPy:            v.NameAbbrPy,
			NameFullPy:            v.NameFullPy,
			IdentityCardNo:        *v.IdentityCardNo,
			NaCode:                v.NaCode,
			Phone:                 v.Phone,
			LandlinePhone:         v.LandlinePhone,
			Sex:                   v.Sex,
			Birthday:              v.Birthday,
			Nickname:              v.Nickname,
			TaxerType:             v.TaxerType,
			TaxerNo:               v.TaxerNo,
			HeaderImgURL:          v.HeaderImgURL,
			IdentityCardPersonImg: v.IdentityCardPersonImg,
			IdentityCardEmblemImg: v.IdentityCardEmblemImg,
			Isdel:                 v.Isdel,
			CommonModel:           v.CommonModel,
		}
		err = db.Def().Transaction(func(tx *gorm.DB) error {
			err = tx.Table("uims_user_auth").Where("id = ?", v.UserID).Update("status", "N").Error
			if err != nil {
				return err
			}
			err = tx.Table("uims_user_info_cursor").Create(&userInfoCursor).Error
			if err != nil {
				return err
			}
			err = tx.Where("id = ?", v.ID).Delete(&model.UserInfo{}).Error
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	//userData := []model.UserInfo{}
	//err = db.Def().Table("uims_user_info").Where("identity_card_no != ?", "").Scan(&userData).Error
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//for _, userInfo := range userData {
	//	err = db.Def().Table("uims_user_auth").Where("id = ?", userInfo.UserID).Update("identity_card_no", userInfo.IdentityCardNo).Error
	//	if err != nil {
	//		fmt.Println("同步失败：", userInfo.ID, userInfo.IdentityCardNo)
	//	}
	//}
	return
}

func (UpdateUserAuthIdCard) Down() (err error) {
	return
}
