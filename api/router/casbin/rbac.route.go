package casbin_route

import (
	casbin_controller "clean-architecture/api/controller/casbin"
	"clean-architecture/internal/casbin"
	"github.com/gin-gonic/gin"
)

func CasbinRouter(group *gin.RouterGroup) {
	r := casbin.SetUp()
	cbGroup := group.Group("/casbin")
	cbGroup.POST("/add-role", casbin_controller.AddRole)
	err := r.SavePolicy()
	if err != nil {
		return
	}
}
