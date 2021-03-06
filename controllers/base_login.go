package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/itcloudy/gin-base-framework/common"
	"github.com/itcloudy/gin-base-framework/middles"
	"github.com/itcloudy/gin-base-framework/models"
	"github.com/itcloudy/gin-base-framework/services"
)

// @tags  用户登录
// @Description 用户登录
// @Summary 用户登录
// @Accept  json
// @Produce  json
// @Param name query string true "用户名"
// @Param password query string true "密码"
// @Success 200 {string} json "{"code":200,"data":{"token":"Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjc1MzIxNjI5NTQsIk5hbWUiOiJhZG1pbiIsIlJvbGUiOm51bGwsIlVzZXJJZCI6MSwiSXNBZG1pbiI6dHJ1ZX0.YzPovX4xP6PPlZV9UGPKWgoLfGL8hnnC01j3L-k6f56mJuds7UDL--3Nts_P6RhuOQgzv7BL7hh6CJdSJdopjzE4A4HmEsq80_DN7cazuFE6gzA2ZfVLI7jslnWcJmVHVPTfu8_57NScdfxCDX_nFbbZUWjzDy7iT5L5zLXrBvg","user":{"id":1,"name":"admin","alias":"","email":"admin@block.vc","password":"","roles":[],"openid":"","active":true,"is_admin":true}},"message":"success"}"
// @Router /login [post]
func Login(c *gin.Context) {
	var (
		bindUser models.User
	)
	err := c.ShouldBindWith(&bindUser, binding.JSON)
	user, err, code := services.CheckUser(bindUser.Name, bindUser.Password)
	if err == nil {
		var roleList []string
		for _, role := range user.Roles {
			roleList = append(roleList, fmt.Sprintf("role_%d", role.ID))
		}
		token := middles.GenerateJWT(user.Name, roleList, user.ID, user.IsAdmin)
		var data map[string]interface{}
		data = make(map[string]interface{})
		data["user"] = user
		data["token"] = token
		common.GenResponse(c, common.SUCCESSED, data, "success")

	} else {
		common.GenResponse(c, code, nil, err.Error())
	}

}

// @tags  用户登录
// @Description 用户微信登录
// @Summary 用户微信登录
// @Accept  json
// @Produce  json
// @Param openid query string true "微信openid"
// @Param name query string false "用户名"
// @Param head query string false "用户头像"
// @Success 200 {string} json "{"code":200,"data":{"token":"Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjc1MzIxNjMwODMsIk5hbWUiOiJhZG1pbiIsIlJvbGUiOm51bGwsIlVzZXJJZCI6MiwiSXNBZG1pbiI6ZmFsc2V9.HZq5jBw4-ZQipQPnq0K7Ei0_LvaRXZGNgKqLoFnhV_vpfQupmddsDMZbiI_Yy0Zhd7J7AvRGDXMfVwW9-TidsDrux6-L4KQWIV0Mrlj4SXgW13HvMSXW0XzHYQBxiai61AeJx4VmQR84s2lI5hmKuiVOpsyOZAduJoO1K26b8X4","user":{"id":2,"name":"admin","alias":"","email":"","password":"","roles":[],"openid":"admin","active":true,"is_admin":false}},"message":"success"}"
// @Router /login/wechat [post]
func LoginWechat(c *gin.Context) {
	var (
		bindUser models.User
		code     int
		user     *models.User
		err      error
	)
	err = c.ShouldBindWith(&bindUser, binding.JSON)
	if bindUser.OpenId == "" {
		code = common.OPEN_ID_IS_EMPITY
		common.GenResponse(c, code, nil, "openid is empty")
		return
	}
	_, err, code = services.GetOpenId(bindUser.OpenId)
	if err != nil {
		common.GenResponse(c, code, nil, "openid is not exist")
		return
	}
	name := bindUser.Name
	head := bindUser.Head
	user, err, code = services.GetUserByOpenId(&bindUser)
	if err == nil {
		_, err, code = services.UpdateUser(bindUser.OpenId, name, head)
		if err != nil {
			common.GenResponse(c, code, nil, err.Error())
			return
		}
		var roleList []string
		for _, role := range user.Roles {
			roleList = append(roleList, fmt.Sprintf("role_%d", role.ID))
		}
		token := middles.GenerateJWT(user.Name, roleList, user.ID, user.IsAdmin)
		var data map[string]interface{}
		data = make(map[string]interface{})
		data["user"] = user
		data["token"] = token
		common.GenResponse(c, common.SUCCESSED, data, "success")
	} else {
		common.GenResponse(c, code, nil, err.Error())
	}
}
