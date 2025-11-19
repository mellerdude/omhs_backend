package kanban

import "github.com/gin-gonic/gin"

const BasePath = "/kanban"

func RegisterRoutes(r *gin.RouterGroup, controller *KanbanController) {
	group := r.Group(BasePath)
	{
		group.GET("", controller.GetKanban)
		group.POST("", controller.CreateKanban)
		group.PUT("", controller.UpdateKanban)
	}
}
