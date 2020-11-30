package user_info_controller

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"time"
	responses2 "uims/internal/controllers/responses"
	requests2 "uims/internal/controllers/user_info_controller/requests"
	"uims/internal/model"
	"uims/pkg/db"
	//"uims/db"
	//"uims/internal/model"
	"github.com/cao-guang/pinyin"
)

// 创建用户资料数据
func Create(c *gin.Context) {
	// 创建模型绑定参数对象
	var request requests2.UserInfoCreateRequest
	var err error
	if err = c.ShouldBind(&request); err != nil {
		responses2.Error(c, err)
		return
	}
	//中文转换为拼音
	pinyin.LoadingPYFileName("conf/pinyin.txt")
	name_full_py, err := pinyin.To_Py(request.NameCn, "", "") //默认造型： hanzipinyin
	name_abbr_py := string([]byte(name_full_py)[:1])

	birthdayTime, err := time.ParseInLocation("2006-01-02", request.Birthday, time.Local)
	if err != nil {
		responses2.Error(c, err)
		return
	}

	//创建数据
	user_info_obj := model.UserInfo{
		Birthday:              &birthdayTime,
		HeaderImgURL:          request.HeaderImgURL,
		IdentityCardEmblemImg: request.IdentityCardEmblemImg,
		IdentityCardPersonImg: request.IdentityCardPersonImg,
		IdentityCardNo:        &request.IdentityCardNo,
		NameCn:                request.NameCn,
		NameEn:                request.NameEn,
		Nickname:              request.Nickname,
		NameAbbrPy:            name_abbr_py,
		NameFullPy:            name_full_py,
		UserID:                request.UserID,
	}

	err = db.Def().Create(&user_info_obj).Error
	if err != nil {
		responses2.Error(c, err)
		return
	}
	responses2.Success(c, "success", user_info_obj)

}

func Update(c *gin.Context) {
	// 模型绑定校验
	var request requests2.UserInfoCreateRequest
	var err error
	if err = c.ShouldBind(&request); err != nil {
		responses2.Error(c, err)
		return
	}
	var user_info model.UserInfo
	err = db.Def().Where("user_id = ?", request.UserID).First(&user_info).Error
	if err != nil {
		responses2.Error(c, err)
		return
	}

	//中文转换为拼音
	pinyin.LoadingPYFileName("conf/pinyin.txt")
	name_full_py, err := pinyin.To_Py(request.NameCn, "", "") //默认造型： hanzipinyin
	if err != nil {
		responses2.Error(c, err)
		return
	}
	name_abbr_py := string([]byte(name_full_py)[:1])
	// 转换json时间类型
	birthdayTime, err := time.ParseInLocation("2006-01-02", request.Birthday, time.Local)
	if err != nil {
		responses2.Error(c, err)
		return
	}

	user_info.Nickname = request.Nickname
	user_info.NameFullPy = name_full_py
	user_info.NameAbbrPy = name_abbr_py
	user_info.NameEn = request.NameEn
	user_info.NameCn = request.NameCn
	user_info.HeaderImgURL = request.HeaderImgURL
	user_info.Birthday = &birthdayTime
	user_info.Sex = request.Sex
	//user_info.UpdatedAt = time.Now()

	// TODO RPC远程同步到微桌业务系统
	tx := db.Def().Begin()

	err = tx.Save(&user_info).Error
	if err != nil {
		tx.Rollback()
		responses2.Error(c, err)
		return
	}
	type RpcParams struct {
		phone              string
		nickname           string
		realname           string
		sex                string
		birthday           string
		id_card_person_img string
		id_card_emblem_img string
		id_no              string
		avatar             string
	}

	//var rpc_params model.UserInfo
	//var rpc_params = RpcParams{user_info.Phone,
	//	user_info.Nickname,
	//	user_info.NameCn, user_info.Sex,
	//	user_info.Birthday.Format("2006-01-02 15:04:05"),
	//	user_info.IdentityCardPersonImg,
	//	user_info.IdentityCardEmblemImg,
	//	user_info.IdentityCardNo,
	//	user_info.HeaderImgURL,
	//}
	//rpc_params.Phone = user_info.Phone
	//rpc_params.name = user_info.NameCn
	//rpc_params.birthday = user_info.Birthday.Format("2006-01-02 15:04:05")

	//response, err := service.GetThriftClientServer().
	//	VzhuoUserThriftClientInvoke("update_user_info", user_info)
	//fmt.Printf("%s\n", response)
	//
	//if response == nil && err != nil {
	//	tx.Rollback()
	//	responses2.Error(c, err)
	//	return
	//}

	tx.Commit()

	responses2.Success(c, "success", user_info)
}

func Get(c *gin.Context) {
	// 进行模型参数校验
	var request requests2.UserDetailRequest
	var err error
	if err = c.ShouldBind(&request); err != nil {
		responses2.Error(c, err)
		return
	}
	var user_base model.User
	var user_info model.UserInfo

	err = db.Def().First(&user_base,
		"id = ? and isdel = ? and status = ?", request.Id, "N", "Y").Error
	if err != nil {
		responses2.Error(c, err)
		return
	}

	err = db.Def().First(&user_info,
		"user_id = ? and isdel = ?", request.Id, "N").Error
	if err == gorm.ErrRecordNotFound {
		responses2.Success(c, "success", nil)
		return
	}
	if err != nil {
		responses2.Error(c, err)
		return
	}
	body := map[string]interface {
	}{}
	body["phone"] = *user_base.Phone
	if user_base.Phone == nil {
		body["phone"] = ""
	} else {
		body["phone"] = *user_base.Phone
	}
	body["email"] = user_base.Email
	body["user_code"] = user_base.UserCode
	body["user_info"] = user_info

	responses2.Success(c, "success", body)
}
