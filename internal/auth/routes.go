package auth

import "github.com/gin-gonic/gin"

const BasePath = "/auth"

func RegisterRoutes(r *gin.RouterGroup, controller *AuthController) {
	group := r.Group(BasePath)
	{
		group.POST("/register", controller.Register)
		group.POST("/login", controller.Login)
		group.POST("/reset-password", controller.ResetPassword)
		group.POST("/change-password", controller.ChangePassword)
	}
}
