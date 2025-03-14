package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/samuelIkoli/GODM/controllers"
)

type Route struct {
	Controller *controller.Controller
}

func (r *Route) RegisterRoutes(router *gin.Engine) {
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "Deployed and Running"})
	})

	router.GET("/task", controller.Task)
	router.GET("/test", controller.Test)
	router.GET("/ping", controller.Ping)
	router.GET("/get", controller.GetMessage)
	router.GET("/users", r.Controller.UpdateUsers)
	router.GET("/gemini", r.Controller.AIResponse)
	router.GET("/gemini-embeddings", r.Controller.GenerateEmbeddings)
	router.GET("/search", r.Controller.VectorSearch)
	router.GET("/rag-search", r.Controller.RaggedResponse)
}

func NewRoute(controller *controller.Controller) *Route {
	return &Route{
		Controller: controller,
	}
}