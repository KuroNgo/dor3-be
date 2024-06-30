package casbin_route

import (
	casbin_controller "clean-architecture/api/controller/casbin"
	"clean-architecture/internal/casbin"
	"github.com/gin-gonic/gin"
)

func CasbinRouter(group *gin.RouterGroup) {
	r := casbin.SetUp()
	cbGroup := group.Group("/casbin")
	cbGroup.POST("/add", casbin_controller.AddRole)
	cbGroup.POST("/add/user", casbin_controller.AddRoleForUser)
	cbGroup.POST("/add/role/api", casbin_controller.AddRoleForAPI)
	cbGroup.POST("/add/api/role", casbin_controller.AddAPIForRole)
	cbGroup.DELETE("/delete", casbin_controller.DeleteRole)
	cbGroup.DELETE("/delete/user", casbin_controller.DeleteRoleForUser)
	cbGroup.DELETE("/delete/role/api", casbin_controller.AddRoleForAPI)
	cbGroup.DELETE("/delete/api/role", casbin_controller.AddAPIForRole)
	err := r.SavePolicy()
	if err != nil {
		return
	}
}
