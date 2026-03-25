package web

import (
	"taskctl/core"
	"github.com/gin-gonic/gin"
)

// StartServer starts the web server
func StartServer(port string) error {
	r := gin.Default()

	// API routes
	api := r.Group("/api/v1")
	{
		api.GET("/tasks", getTasks)
		api.POST("/tasks", createTask)
		api.PUT("/tasks/:id/complete", completeTask)
		api.DELETE("/tasks/:id", deleteTask)
	}

	// Serve static files
	r.Static("/", "./frontend/dist")

	return r.Run(":" + port)
}

func getTasks(c *gin.Context) {
	tasks, err := core.ListTasks()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, tasks)
}

func createTask(c *gin.Context) {
	var req struct {
		Title string `json:"title"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	task, err := core.AddTask(req.Title)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, task)
}

func completeTask(c *gin.Context) {
	// TODO: Implement
	c.JSON(200, gin.H{"status": "ok"})
}

func deleteTask(c *gin.Context) {
	// TODO: Implement
	c.JSON(200, gin.H{"status": "ok"})
}
