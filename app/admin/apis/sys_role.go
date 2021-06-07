package apis

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-admin-team/go-admin-core/sdk"
	"go-admin/app/admin/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/app/admin/service"
	"go-admin/app/admin/service/dto"
	"go-admin/common/global"
)

type SysRole struct {
	api.Api
}

// GetPage
// @Summary 角色列表数据
// @Description Get JSON
// @Tags 角色/Role
// @Param roleName query string false "roleName"
// @Param status query string false "status"
// @Param roleKey query string false "roleKey"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/role [get]
// @Security Bearer
func (e SysRole) GetPage(c *gin.Context) {
	s := service.SysRole{}
	req := dto.SysRoleSearch{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.Form).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	list := make([]models.SysRole, 0)
	var count int64

	err = s.GetPage(&req, &list, &count)
	if err != nil {
		e.Error(500, err, "查询失败")
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get
// @Summary 获取Role数据
// @Description 获取JSON
// @Tags 角色/Role
// @Param roleId path string false "roleId"
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/role/{id} [get]
// @Security Bearer
func (e SysRole) Get(c *gin.Context) {
	s := service.SysRole{}
	req := dto.SysRoleById{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, nil).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	var object models.SysRole

	err = s.Get(&req, &object)
	if err != nil {
		e.Error(http.StatusUnprocessableEntity, err, "查询失败")
		return
	}

	e.OK(object, "查看成功")
}

// Insert
// @Summary 创建角色
// @Description 获取JSON
// @Tags 角色/Role
// @Accept  application/json
// @Product application/json
// @Param data body dto.SysRoleControl true "data"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/role [post]
// @Security Bearer
func (e SysRole) Insert(c *gin.Context) {
	s := service.SysRole{}
	req := dto.SysRoleControl{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.JSON).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	// 设置创建人
	req.CreateBy = user.GetUserId(c)
	if req.Status == "" {
		req.Status = "2"
	}
	cb := sdk.Runtime.GetCasbinKey(c.Request.Host)
	err = s.Insert(&req, cb)
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, "创建失败")
		return
	}
	_, err = global.LoadPolicy(c)
	if err != nil {
		e.Error(500, err, "")
		return
	}
	e.OK(req.GetId(), "创建成功")
}

// Update 修改用户角色
// @Summary 修改用户角色
// @Description 获取JSON
// @Tags 角色/Role
// @Accept  application/json
// @Product application/json
// @Param data body dto.SysRoleControl true "body"
// @Success 200 {string} string	"{"code": 200, "message": "修改成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "修改失败"}"
// @Router /api/v1/role/{id} [put]
// @Security Bearer
func (e SysRole) Update(c *gin.Context) {
	s := service.SysRole{}
	req := dto.SysRoleControl{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, nil, binding.JSON).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	cb := sdk.Runtime.GetCasbinKey(c.Request.Host)

	req.SetUpdateBy(user.GetUserId(c))

	err = s.Update(&req, cb)
	if err != nil {
		e.Logger.Error(err)
		return
	}
	_, err = global.LoadPolicy(c)
	if err != nil {
		e.Error(http.StatusInternalServerError, err, "")
		return
	}
	e.OK(req.GetId(), "更新成功")
}

// Delete
// @Summary 删除用户角色
// @Description 删除数据
// @Tags 角色/Role
// @Param data body dto.SysRoleById true "body"
// @Success 200 {string} string	"{"code": 200, "message": "删除成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "删除失败"}"
// @Router /api/v1/role [delete]
// @Security Bearer
func (e SysRole) Delete(c *gin.Context) {
	s := new(service.SysRole)
	req := dto.SysRoleById{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.JSON).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	err = s.Remove(&req)
	if err != nil {
		e.Logger.Error(err)
		e.Error(http.StatusInternalServerError, err, "")
		return
	}
	_, err = global.LoadPolicy(c)
	if err != nil {
		e.Error(http.StatusInternalServerError, err, "")
		return
	}
	e.OK(req.GetId(), "删除成功")
}

// Update2DataScope 更新角色数据权限
func (e SysRole) Update2DataScope(c *gin.Context) {
	s := &service.SysRole{}
	req := dto.RoleDataScopeReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.JSON, nil).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	data := &models.SysRole{
		RoleId:    req.RoleId,
		DataScope: req.DataScope,
		DeptIds:   req.DeptIds,
	}
	data.UpdateBy = user.GetUserId(c)
	err = s.UpdateDataScope(data)
	if err != nil {
		e.Error(http.StatusInternalServerError, err, "")
		return
	}
	e.OK(nil, "操作成功")
}