package usercontroller

import (
	"github.com/jinzhu/gorm"
	"uims/internal/service"
	"uims/pkg/db"
	thriftserver "uims/pkg/thrift/server"
)

//保存结算系统用户数据
func Create(c *thriftserver.Context) {
	var req service.PushUserData
	var err error
	if err := c.ShouldBind(&req); err != nil {
		c.Response.BadParams(err)
		return
	}
	userService := service.UserService{}
	userRoleService := service.UserRoleService{}
	userOrgService := service.UserOrgService{}
	err = db.Def().Transaction(func(tx *gorm.DB) error {
		err, userID := userService.SaveCassUserData(req, tx)
		if err != nil {
			return err
		}
		if err = userRoleService.SaveUserRoleData(userID, req.RoleName, "user", tx); err != nil {
			return err
		}
		if err = userOrgService.SaveUserOrgData(userID, req.GroupName, req.RoleName, tx); err != nil {
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
