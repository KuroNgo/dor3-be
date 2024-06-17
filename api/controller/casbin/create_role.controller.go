package casbin_controller

import (
	"clean-architecture/internal/casbin"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddRole(ctx *gin.Context) {
	var data RoleData
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, "can not get data")
		return
	}

	if data.API == nil {
		data.API = append(data.API, "http://localhost:8080")
	}

	// Add policy rules
	for _, api := range data.API {
		for _, method := range data.Method {
			// Assuming rbac is already initialized and AddPolicy is defined
			_, err := casbin.Rbac.AddPolicy(data.Role, api, method)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "ERROR"})
				return
			}
		}
	}

	ctx.JSON(http.StatusCreated, "success added role: "+data.Role)
}
