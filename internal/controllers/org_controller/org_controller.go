package org_controller

import (
	"github.com/gin-gonic/gin"
	requests2 "uims/internal/controllers/org_controller/requests"
	responses2 "uims/internal/controllers/responses"
	"uims/internal/service"
)

// @Summary 获取客户端信息
// @Produce  json
// @Param id query int true "ID"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/client [GET]
func Detail(c *gin.Context) {
	var request requests2.OrgDetailRequest

	if err := c.ShouldBindQuery(&request); err != nil {
		responses2.Error(c, err)
		return
	}

	//查询组织是否存在
	if !service.GetOrgService().ExistOrgByID(request.ID) {
		responses2.Failed(c, "org_id  not exist", nil)
	}

	//查询org详细信息
	org, e := service.GetOrgService().GetOrgByID(request.ID)

	if e != nil {
		responses2.Error(c, e)
		return
	}

	responses2.Success(c, "success", org)
}
