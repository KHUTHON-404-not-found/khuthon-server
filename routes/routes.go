package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/KHUTHON-404-not-found/khuthon-server/controllers"
	"github.com/KHUTHON-404-not-found/khuthon-server/middleware"
)

// SetupRoutes API 라우트 설정
func SetupRoutes(router *gin.Engine) {
	// 기본 경로
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to PlantApp API",
		})
	})

	// API 그룹
	api := router.Group("/api")
	{
		// 인증 관련 라우트
		auth := api.Group("/auth")
		{
			auth.POST("/login", controllers.Login)
			auth.POST("/register", controllers.CreateUser)
		}

		// 인증이 필요한 라우트
		protected := api.Group("/")
		protected.Use(middleware.JWTAuth())
		{
			// 사용자 관련 라우트
			users := protected.Group("/users")
			{
				users.GET("/", controllers.GetAllUsers)
				users.GET("/:id", controllers.GetUser)
				users.PUT("/:id", controllers.UpdateUser)
				users.DELETE("/:id", controllers.DeleteUser)
				users.GET("/email", controllers.GetUserByEmail)
				users.GET("/project/:project_id", controllers.GetUsersByProject)
			}

			// 프로젝트 관련 라우트
			projects := protected.Group("/projects")
			{
				projects.POST("/", controllers.CreateProject)
				projects.GET("/", controllers.GetAllProjects)
				projects.GET("/:id", controllers.GetProject)
				projects.PUT("/:id", controllers.UpdateProject)
				projects.DELETE("/:id", controllers.DeleteProject)
				projects.GET("/todo/:todo_id", controllers.GetProjectsByTodo)
			}

			// 할일 관련 라우트
			todos := protected.Group("/todos")
			{
				todos.POST("/", controllers.CreateTodo)
				todos.GET("/", controllers.GetAllTodos)
				todos.GET("/:id", controllers.GetTodo)
				todos.PUT("/:id", controllers.UpdateTodo)
				todos.DELETE("/:id", controllers.DeleteTodo)
				todos.GET("/date", controllers.GetTodosByDate)
				todos.PUT("/:id/complete", controllers.CompleteTodo)
			}

			// 일기 관련 라우트
			diaries := protected.Group("/diaries")
			{
				diaries.POST("/", controllers.CreateDiary)
				diaries.GET("/", controllers.GetAllDiaries)
				diaries.GET("/:id", controllers.GetDiary)
				diaries.PUT("/:id", controllers.UpdateDiary)
				diaries.DELETE("/:id", controllers.DeleteDiary)
				diaries.GET("/project/:project_id", controllers.GetDiariesByProject)
				diaries.GET("/date", controllers.GetDiariesByDate)
				//diaries.POST("/:id/photo", controllers.UploadDiaryPhoto)
				diaries.PATCH("/:id/photo-url", controllers.UpdateDiaryPhotoURL)
			}
		}
	}
}