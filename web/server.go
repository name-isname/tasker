package web

import (
	"taskctl/core"
	"github.com/gin-gonic/gin"
	"strconv"
)

// StartServer starts the web server
func StartServer(port string) error {
	r := gin.Default()

	// API routes
	api := r.Group("/api/v1")
	{
		api.GET("/processes", getProcesses)
		api.POST("/processes", createProcess)
		api.PUT("/processes/:id/status", setProcessStatus)
		api.DELETE("/processes/:id", deleteProcess)
		api.GET("/processes/:id/logs", getLogs)
		api.POST("/processes/:id/logs", addLog)
		api.GET("/search", searchProcesses)
	}

	// Serve static files
	r.Static("/", "./frontend/dist")

	return r.Run(":" + port)
}

func getProcesses(c *gin.Context) {
	statusStr := c.Query("status")
	var status *core.ProcessStatus
	if statusStr != "" {
		s := core.ProcessStatus(statusStr)
		status = &s
	}

	processes, err := core.ListProcesses(status)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, processes)
}

func createProcess(c *gin.Context) {
	var req struct {
		Title       string              `json:"title" binding:"required"`
		Description string              `json:"description"`
		ParentID    *uint               `json:"parent_id"`
		Priority    core.ProcessPriority `json:"priority"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if req.Priority == "" {
		req.Priority = core.PriorityMedium
	}
	process, err := core.CreateProcess(req.Title, req.Description, req.ParentID, req.Priority)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, process)
}

func setProcessStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}

	var req struct {
		Status core.ProcessStatus `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = core.SetProcessStatus(uint(id), req.Status)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok"})
}

func deleteProcess(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}

	err = core.DeleteProcess(uint(id))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok"})
}

func getLogs(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}

	logs, err := core.GetLogs(uint(id))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, logs)
}

func addLog(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}

	var req struct {
		LogType core.LogType `json:"log_type" binding:"required"`
		Content string       `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	log, err := core.AddLog(uint(id), req.LogType, req.Content)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, log)
}

func searchProcesses(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(400, gin.H{"error": "missing query parameter 'q'"})
		return
	}

	processes, err := core.SearchProcesses(query)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, processes)
}
