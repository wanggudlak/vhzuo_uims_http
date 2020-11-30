package esigncontroller

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"strings"
	"time"
	"uims/internal/model"
	"uims/internal/service/uuid"
	"uims/pkg/db"
	"uims/pkg/randc"
	"uims/pkg/storage"
	thriftserver "uims/pkg/thrift/server"
)

// 接收电签通知
// 可能手机号会为空
func NotifyESign(c *thriftserver.Context) {
	var req NotifyESignReq
	if err := c.ShouldBind(&req); err != nil {
		c.Response.BadParams(err)
		return
	}
	// 判断身份证号或手机号是否已注册
	var phoneUser model.User
	var idCardUser model.User
	var personImgPath string
	var emblemImgPath string
	err := db.Def().Transaction(func(tx *gorm.DB) error {
		var err error
		saveIdCardImgPath := "id_cards/" + time.Now().Format("20060102")
		if req.IdentityCardPersonImgBase64 != "" {
			personImgPath, err = storage.Storage.StoreBase64RandomName(saveIdCardImgPath, req.IdentityCardPersonImgBase64)
			if err != nil {
				return errors.Wrap(err, "保存头像面身份证照失败")
			}
		}
		if req.IdentityCardEmblemImgBase64 != "" {
			emblemImgPath, err = storage.Storage.StoreBase64RandomName(saveIdCardImgPath, req.IdentityCardEmblemImgBase64)
			if err != nil {
				return errors.Wrap(err, "保存国徽面身份证照失败")
			}
		}
		if req.Phone != "" {
			err = db.Def().Where("phone = ?", req.Phone).First(&phoneUser).Error
			if err != nil {
				if !gorm.IsRecordNotFoundError(err) {
					return errors.Wrap(err, "查找 user 失败")
				}
				// 手机号未找到
			}
		}
		err = db.Def().Where("identity_card_no = ?", req.IdCard).First(&idCardUser).Error
		if err != nil {
			if !gorm.IsRecordNotFoundError(err) {
				return errors.Wrap(err, "查找 user 失败")
			}
			// 身份证未找到
		}
		if phoneUser.ID == 0 && idCardUser.ID == 0 {
			var user model.User
			var userInfo model.UserInfo
			// 手机号身份证均未使用, 进行注册
			if req.Phone != "" {
				user.Phone = &req.Phone
			}
			user.OpenID = strings.ToUpper(randc.UUID())
			user.UserType = "ESIGN"
			user.UserCode = uuid.GenerateForUIMS().String()
			user.NaCode = "+86"
			user.IdentityCardNo = &req.IdCard
			if err = db.Def().Create(&user).Error; err != nil {
				return errors.Wrap(err, "创建 user 失败")
			}
			userInfo.IdentityCardNo = &req.IdCard
			userInfo.IdentityCardPersonImg = personImgPath
			userInfo.IdentityCardEmblemImg = emblemImgPath
			userInfo.UserID = user.ID
			userInfo.NaCode = user.NaCode
			userInfo.Phone = req.Phone
			userInfo.IsIdentify = "Y"
			userInfo.NameCn = req.Name
			userInfo.UserCode = user.UserCode
			userInfo.UserType = user.UserType
			if err = db.Def().Create(&userInfo).Error; err != nil {
				return errors.Wrap(err, "创建 user_info 失败")
			}
		} else if phoneUser.ID == 0 {
			// 用户未设置手机号或已保存的手机号和传入手机号不一致
			if idCardUser.Phone != nil && *idCardUser.Phone != req.Phone {
				return errors.New(fmt.Sprintf("传入的手机号 %s 和系统中已存在的手机号 %s 不一致", req.Phone, *idCardUser.Phone))
			}
			idCardUser.NaCode = "+86"
			if req.Phone != "" {
				idCardUser.Phone = &req.Phone
			}
			if err = db.Def().Save(&idCardUser).Error; err != nil {
				return errors.Wrap(err, "保存 user 失败")
			}
			if err = db.Def().Model(&model.UserInfo{}).Where("user_id = ?", idCardUser.ID).Updates(&model.UserInfo{
				Phone:                 req.Phone,
				NaCode:                idCardUser.NaCode,
				IdentityCardPersonImg: personImgPath,
				IdentityCardEmblemImg: emblemImgPath,
			}).Error; err != nil {
				return errors.Wrap(err, "更新 user_info 失败")
			}
		} else if idCardUser.ID == 0 {
			// 未设置身份证号
			if phoneUser.IdentityCardNo != nil && *phoneUser.IdentityCardNo != req.IdCard {
				return errors.New(fmt.Sprintf("传入的身份证号 %s 和系统中已存在的身份证号 %s 不一致", req.IdCard, phoneUser.IdentityCardNo))
			}
			phoneUser.IdentityCardNo = &req.IdCard
			if err = db.Def().Save(&phoneUser).Error; err != nil {
				return errors.Wrap(err, "保存 user 失败")
			}
			if err = db.Def().Model(&model.UserInfo{}).Where("user_id = ?", phoneUser.ID).Updates(&model.UserInfo{
				IsIdentify:            "Y",
				IdentityCardNo:        &req.IdCard,
				IdentityCardPersonImg: personImgPath,
				IdentityCardEmblemImg: emblemImgPath,
			}).Error; err != nil {
				return errors.Wrap(err, "更新 user_info 失败")
			}

		} else {
			if phoneUser.ID != idCardUser.ID {
				return errors.New("手机号与身份证号绑定在不同用户")
			} else {
				// 只更新照片数据
				if err = db.Def().Model(&model.UserInfo{}).Where("user_id = ?", phoneUser.ID).Updates(&model.UserInfo{
					IsIdentify:            "Y",
					IdentityCardPersonImg: personImgPath,
					IdentityCardEmblemImg: emblemImgPath,
				}).Error; err != nil {
					return errors.Wrap(err, "更新 user_info 失败")
				}
			}
		}
		return nil
	})
	if err != nil {
		c.Response.Error(err)
		return
	}
	// 已注册则更新信息
	// 未注册则创建用户
	c.Response.Success(NotifyESignResp{}, "ok")
	return
}
