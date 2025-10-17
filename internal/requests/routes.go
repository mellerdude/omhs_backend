package requests

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, controller *RequestController) {
	r.POST("/:database/:collection", controller.Create)
	r.GET("/:database/:collection/:id", controller.Get)
	r.PUT("/:database/:collection/:id", controller.Update)
	r.DELETE("/:database/:collection/:id", controller.Delete)
	r.GET("/:database/:collection", controller.GetAll)
}
