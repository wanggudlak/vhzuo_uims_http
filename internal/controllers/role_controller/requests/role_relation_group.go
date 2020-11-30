package requests

// 新建角色资源组关系
type GroupNewRequest struct {
	ID      int    `json:"id" form:"client_id" binding:"required"`
	GroupID int    `json:"group_id" form:"group_id" binding:"-"`
	Forget  string `json:"forget" form:"-" binding:"" example:"2020-05-11 20:00:45" comment:"过期时间"`
}

// 删除角色资源组关系
type RoleGroupMapRequest struct {
	RoleID  int `json:"role_id" form:"role_id" binding:"required" comment:"角色id"`
	GroupID int `json:"group_id" form:"group_id" binding:"required" comment:"资源组id"`
}
