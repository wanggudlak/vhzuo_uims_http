package user_controller

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	responses2 "uims/internal/controllers/responses"
	requests2 "uims/internal/controllers/user_controller/requests"
	"uims/internal/model"
	"uims/internal/service"
	"uims/internal/service/uuid"
	"uims/pkg/db"
	encry "uims/pkg/encryption"
)

// @Summary 创建新用户
// @Produce  json
// @Param user body requests.UserStoreRequest true "注册信息"
// @Success 200 {object} responses.Response
// @Router /api/users [post]
func Create(c *gin.Context) {
	var request requests2.UserStoreRequest
	var err error
	if err = c.ShouldBindJSON(&request); err != nil {
		responses2.BadReq(c, err)
		return
	}

	if service.PhoneIsExists(request.Phone) {
		responses2.Failed(c, "Phone has been used", nil)
		return
	}
	//todo 假设是微桌平台则需要进行加盐等加密处理方式(这里暂时写个假的处理)
	var salt string
	var pwd string

	switch request.EncryptType {
	// uims密码加密方式
	case 0:
		pwd, err = encry.BcryptHash(request.Passwd)
		if err != nil {
			responses2.Error(c, err)
		}
	// vzhuo业务系统加密方式
	case 1:
		pwd, salt = encry.DefaultPBKDF2Options.GeneratePasswdPBKDF2Key([]byte(request.Passwd), []byte{})

	// 结算系统加密方式
	case 2:
		pwd, err = encry.BcryptHash(request.Passwd)
		if err != nil {
			responses2.Error(c, err)
		}
	}

	u := model.User{
		UserType:    "VDK",
		Account:     request.Account,
		UserCode:    uuid.GenerateForUIMS().String(),
		NaCode:      "+86",
		Phone:       &request.Phone,
		Email:       request.Email,
		Salt:        salt,
		EncryptType: request.EncryptType,
		Passwd:      pwd,
	}

	err = db.Def().Create(&u).Error
	if err != nil {
		responses2.Error(c, err)
		return
	}

	responses2.Success(c, "success", u)
}

func List(c *gin.Context) {
	var req requests2.UserListRequest
	var err error
	if err = c.ShouldBind(&req); err != nil {
		responses2.BadReq(c, err)
		return
	}
	// 进行分页数据处理和转换
	page := req.Page
	pagesize := req.PageSize
	if page == 0 {
		page = 1
	}
	if pagesize == 0 {
		pagesize = 10
	}
	page = page - 1
	var users []model.User
	var total int
	// 获取用户总数
	q := db.Def().
		Model(&model.User{}).
		Select("count(id)").
		Where("isdel = ?", "N")
	if req.Phone != "" {
		q = q.Where("phone like ?", "%"+req.Phone+"%")
	}
	if req.Email != "" {
		q = q.Where("email like ?", "%"+req.Email+"%")
	}
	err = q.Count(&total).Error
	if err != nil {
		responses2.Error(c, errors.Wrap(err, "查询用户总数失败"))
		return
	}
	q = db.Def().
		Offset(page * pagesize).
		Limit(pagesize).
		Order("updated_at desc")
	if req.Phone != "" {
		q = q.Where("phone like ?", "%"+req.Phone+"%")
	}
	if req.Email != "" {
		q = q.Where("email like ?", "%"+req.Email+"%")
	}
	err = q.Find(&users, "isdel = ?", "N").Error
	if err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			responses2.Error(c, errors.Wrap(err, "查询用户列表失败"))
			return
		}
	}
	body := make(map[string]interface{})
	var data_list []interface{}

	for _, user := range users {
		user_map := make(map[string]interface{})
		user_map["id"] = user.ID
		user_map["open_id"] = user.OpenID
		if user.Phone == nil {
			user_map["phone"] = ""
		} else {
			user_map["phone"] = *user.Phone
		}
		user_map["email"] = user.Email
		user_map["user_code"] = user.UserCode
		user_map["na_code"] = user.NaCode
		user_map["encrypt_type"] = user.EncryptType
		user_map["status"] = user.Status
		user_map["created_at"] = user.CreatedAt
		user_map["updated_at"] = user.UpdatedAt
		// 取出user 所拥有的的角色id 列表
		var userRoleList []model.UserRole
		err = db.Def().Where("user_id = ?", user.ID).Find(&userRoleList).Error
		if err != nil {
			if !gorm.IsRecordNotFoundError(err) {
				responses2.Error(c, errors.Wrap(err, "查询用户角色列表失败"))
				return
			}
		}
		var roleIdList []int
		for _, userRole := range userRoleList {
			roleIdList = append(roleIdList, userRole.RoleID)
		}
		// 获取所有的角色数据
		var roleList []model.Role
		err = db.Def().Where("id in (?)", roleIdList).Find(&roleList).Error
		if err != nil {
			if !gorm.IsRecordNotFoundError(err) {
				responses2.Error(c, errors.Wrap(err, "查询用户角色数据失败"))
				return
			}
		}
		user_map["role_list"] = roleList
		data_list = append(data_list, user_map)
	}
	body["data"] = data_list
	body["total"] = total
	//log.Print(users)
	responses2.Success(c, "success", body)
}

// 用户状态解禁或者冻结
func Update(c *gin.Context) {
	var req requests2.UserStatusRequest
	var err error
	if err = c.ShouldBind(&req); err != nil {
		responses2.BadReq(c, err)
		return
	}
	err = service.UserService{}.UpdateUser(req.Id, req.Type)
	if err != nil {
		responses2.BadReq(c, err)
		return
	}
	responses2.Success(c, "success", nil)
}

func Get(c *gin.Context) {
	id := c.Query("id")
	var err error

	var user model.User
	err = db.Def().Where("id = ?", id).First(&user).Error
	if err != nil {
		responses2.Error(c, err)
		return
	}
	user_map := make(map[string]interface{})
	// 取出user 所拥有的的角色id 列表
	var user_role_list []model.UserRole
	err = db.Def().Where("user_id = ?", user.ID).Find(&user_role_list).Error
	if err != nil {
		responses2.Error(c, err)
		return
	}
	var role_id_list []int
	for _, user_role := range user_role_list {
		role_id_list = append(role_id_list, user_role.RoleID)
	}
	// 获取所有的角色数据
	var role_list []model.Role
	err = db.Def().Where("id in (?)", role_id_list).Find(&role_list).Error
	if err != nil {
		responses2.Error(c, err)
		return
	}
	user_map["id"] = user.ID
	user_map["open_id"] = user.OpenID
	if user.Phone == nil {
		user_map["phone"] = ""
	} else {
		user_map["phone"] = *user.Phone
	}
	user_map["email"] = user.Email
	user_map["user_code"] = user.UserCode
	user_map["na_code"] = user.NaCode
	user_map["encrypt_type"] = user.EncryptType
	user_map["status"] = user.Status
	user_map["role_list"] = role_list
	responses2.Success(c, "success", user_map)
}
