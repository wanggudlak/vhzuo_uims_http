package service

import (
	"errors"
	"github.com/jinzhu/gorm"
	"uims/internal/model"
	"uims/pkg/db"
)

type UserRoleService struct {
}

// 查询用户关联的角色
func (UserRoleService) GetUserRoles(maps interface{}) ([]int, error) {
	var (
		userRoles []model.UserRole
		err       error
		RoleIDS   []int
	)

	err = db.Def().Where(maps).Order("id desc").Find(&userRoles).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	for _, v := range userRoles {
		RoleIDS = append(RoleIDS, v.RoleID)
	}

	return RoleIDS, nil
}

// 查询用户是否已经关联过角色
func (UserRoleService) ExistUserRole(dict map[string]int) bool {

	var userRole model.UserRole
	err := db.Def().Select("id").
		Where("user_id = ? and role_id = ?", dict["user_id"], dict["role_id"]).
		First(&userRole).Error
	if err != nil {
		return false
	}
	return true
}

// 为用户添加角色关联关系
func (UserRoleService) AddUserRole(dict map[string]int) error {
	userRole := model.UserRole{
		ClientID: dict["client_id"],
		RoleID:   dict["role_id"],
		UserID:   dict["user_id"],
	}

	err := db.Def().Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&userRole).Error; err != nil {
			return err
		}
		err := requestUpdateClientUserRole(userRole.ClientID, userRole.UserID, userRole.RoleID, "add_client_user_role")
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// 为用户删除角色关联关系
func (UserRoleService) DeleteUserRole(dict map[string]int) error {
	userRole := model.UserRole{
		ClientID: dict["client_id"],
		RoleID:   dict["role_id"],
		UserID:   dict["user_id"],
		//CreatedAt: time.Now(),
	}
	err := db.Def().Transaction(func(tx *gorm.DB) error {
		if err := tx.Unscoped().
			Where("role_id = ? and user_id = ?", userRole.RoleID, userRole.UserID).
			Delete(&userRole).Error; err != nil {
			return err
		}
		err := requestUpdateClientUserRole(userRole.ClientID,
			userRole.UserID, userRole.RoleID, "delete_client_user_role")
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

//保存用户关联角色
func (UserRoleService) SaveUserRoleData(userID int, roleName string, userRelationType string, tx *gorm.DB) error {

	var roleInfo model.Role

	err := db.Def().Table("uims_role").Where("role_name_en = ? ", roleName).First(&roleInfo).Error
	if err != nil {
		return err
	}

	var UserRoleInfo model.UserRole
	var userRole = model.UserRole{
		UserID:           userID,
		RoleID:           roleInfo.ID,
		ClientID:         roleInfo.ClientID,
		UserRelationType: userRelationType,
	}
	if err := tx.Where("user_id = ? and role_id = ? and client_id = ? and user_relation_type = ?", userID, roleInfo.ID, roleInfo.ClientID, userRelationType).
		First(&UserRoleInfo).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			if err := tx.Create(&userRole).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		if err := tx.Model(&userRole).
			Where("id = ?", UserRoleInfo.ID).
			Update(&userRole).Error; err != nil {
			return err
		}
	}

	return nil
}

//请求更新子业务系统用户角色数据
func requestUpdateClientUserRole(clientId int, userId int, roleId int, method string) error {
	type clientUserRole struct {
		OpenID     string `json:"open_id"`
		RoleNameEN string `json:"role_name_en"`
		OnlyItem   bool   `json:"only_item"`
	}
	clientInfo, err := ClientService{}.GetClientByID(clientId)
	if clientInfo == nil || err != nil {
		return err
	}
	var userInfo model.User
	err = UserService{}.GetUserInfoByUserID(&userInfo, userId)
	if err != nil {
		return err
	}

	//目前只有结算系统有需要回写数据
	//if userInfo.UserType == "MP" {
	//	return errors.New("该用户无法使用客户端为" + clientInfo.ClientName + "的角色数据")
	//}

	roleInfo, err := RoleService{}.GetRole(roleId)
	if err != nil {
		return err
	}
	var addUserRole = clientUserRole{
		OpenID:     userInfo.OpenID,
		RoleNameEN: roleInfo.RoleNameEN,
		OnlyItem:   true,
	}
	resp := GetThriftClientServer().
		ClientInvoke(clientId, method, addUserRole)
	if !resp.OK() {
		return errors.New(resp.Err())
	}
	return nil
}
