package requests

type ClientStatusRequest struct {
	ID     int    `json:"id" form:"id" binding:"required" comment:"客户端ID"`
	Status string `json:"status" form:"status" binding:"required" comment:"默认N：未授权不可用；Y：已授权可用；F-被禁用"`
}
