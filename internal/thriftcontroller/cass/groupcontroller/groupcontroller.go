package groupcontroller

import (
	"github.com/jinzhu/gorm"
	"uims/internal/model"
	"uims/internal/service"
	"uims/pkg/db"
	thriftserver "uims/pkg/thrift/server"
)

type PushOrgData struct {
	GroupUUID    string        `json:"groupUUID"`
	Name         string        `json:"name"`
	ParentUUID   string        `json:"parentUUID"`
	Platform     int           `json:"platform"`
	UserAccounts []userAccount `json:"userAccounts"`
}

type userAccount struct {
	Account string `json:"account"`
	Name    string `json:"name"`
}

//保存结算系统组数据
func Create(c *thriftserver.Context) {
	var req PushOrgData
	var err error
	if err = c.ShouldBind(&req); err != nil {
		c.Response.BadParams(err)
		return
	}

	err = db.Def().Transaction(func(tx *gorm.DB) error {
		err, _ := saveCassGroupData(req, tx)
		if err != nil {
			return err
		}
		if err = saveCassGroupUserData(req, tx); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		c.Response.Error(err)
	} else {
		c.Response.Success(nil, "")
	}
	return
}

//保存用户信息
func saveCassGroupData(bizParamBody PushOrgData, db *gorm.DB) (err error, orgID int) {

	orgLevel, parentID := 0, 0
	if bizParamBody.ParentUUID != "" {
		var partentOrgInfo model.Org
		if err := db.Model(&partentOrgInfo).
			Where("org_code = ?", bizParamBody.ParentUUID).
			First(&partentOrgInfo).Error; err != nil {
			return err, 0
		}
		parentID = partentOrgInfo.ID
	}
	if parentID != 0 {
		orgLevel = 1
	}

	var clientInfo model.Client
	if err := db.Model(&clientInfo).
		Where("client_flag_code = ?", "VDK_CASS_BACK").
		First(&clientInfo).Error; err != nil {
		return err, 0
	}

	var orgInfo model.Org
	var saveOrg = model.Org{
		OrgCode:     bizParamBody.GroupUUID,
		ParentOrgID: parentID,
		ClientID:    clientInfo.ID,
		ClientAppID: clientInfo.AppId,
		OrgNameCN:   bizParamBody.Name,
		OrgLevel:    orgLevel,
	}

	if err := db.Where("org_code = ?", bizParamBody.GroupUUID).
		First(&orgInfo).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			if err := db.Create(&saveOrg).Error; err != nil {
				return err, saveOrg.ID
			}
		} else {
			return err, saveOrg.ID
		}
	} else {
		if err := db.Model(&saveOrg).
			Where("org_code = ?", bizParamBody.GroupUUID).
			Update(&saveOrg).Error; err != nil {
			return err, saveOrg.ID
		}
	}

	orgID = saveOrg.ID
	if saveOrg.ID == 0 {
		orgID = orgInfo.ID
	}

	return nil, orgID
}

//保存用户关联组信息
func saveCassGroupUserData(bizParamBody PushOrgData, tx *gorm.DB) (err error) {

	if len(bizParamBody.UserAccounts) == 0 {
		return nil
	}
	userOrgService := service.UserOrgService{}
	for _, v := range bizParamBody.UserAccounts {
		var userInfo model.User
		if err := tx.Model(&userInfo).
			Where("account = ?", v.Account).
			First(&userInfo).Error; err != nil {
			return err
		}
		var clientInfo model.Client
		if err := tx.Model(&clientInfo).
			Where("client_flag_code = ?", "VDK_CASS_BACK").
			First(&clientInfo).Error; err != nil {
			return err
		}
		var userRoleInfo model.UserRole
		if err := tx.Model(&userRoleInfo).
			Where("user_id = ? and client_id = ?", userInfo.ID, clientInfo.ID).
			First(&userRoleInfo).Error; err != nil {
			return err
		}
		var roleInfo model.Role
		if err := tx.Model(&roleInfo).
			Where("id = ?", userRoleInfo.RoleID).
			First(&roleInfo).Error; err != nil {
			return err
		}
		if err = userOrgService.SaveUserOrgData(userInfo.ID, v.Name, roleInfo.RoleNameEN, tx); err != nil {
			return err
		}
	}

	return nil
}
