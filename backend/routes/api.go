package routes

import (
	"apoyo-anonimo/backend/handlers" // Nuestras reglas
	"github.com/gin-gonic/gin"
)

// SetupAPIRoutes configura todos los caminos de nuestra API
func SetupAPIRoutes(router *gin.Engine) {
	api := router.Group("/api") // Creamos un grupo de caminos que empiezan con /api
	{
		// Rutas para Grupos
		api.POST("/groups", handlers.CreateGroup)        // Para crear un grupo
		api.GET("/groups", handlers.GetGroups)           // Para ver todos los grupos
		api.GET("/groups/:id", handlers.GetGroupByID)    // Para ver un grupo específico

		// Rutas para Mensajes y Chat
		api.POST("/messages", handlers.SendMessage)       // Para enviar un mensaje
		api.GET("/groups/:id/messages", handlers.GetMessagesByGroup) // Para ver mensajes de un grupo
		api.GET("/ws/groups/:id", handlers.HandleWebSocket) // ¡Para el chat en tiempo real!
	}
}
