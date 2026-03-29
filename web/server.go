package web

import (
	"io"
	"io/fs"
	"net/http"
	"strings"
	"taskctl/core"
	"github.com/gin-gonic/gin"
	"strconv"
)

// StartServer starts the web server
func StartServer(port string) error {
	r := gin.Default()
	r.RedirectTrailingSlash = false

	// API routes - register these first
	api := r.Group("/api/v1")
	{
		api.GET("/processes", getProcesses)
		api.GET("/processes/:id", getProcess)
		api.POST("/processes", createProcess)
		api.PUT("/processes/:id/status", setProcessStatus)
		api.DELETE("/processes/:id", deleteProcess)
		api.GET("/processes/:id/logs", getLogs)
		api.POST("/processes/:id/logs", addLog)
		api.GET("/search", searchProcesses)
	}

	// Create HTTP file server for embedded frontend
	fsys, err := fs.Sub(FrontendFS, "frontend/dist")
	if err != nil {
		return err
	}

	// Serve all non-API routes through the file server
	r.NoRoute(func(c *gin.Context) {
		// Check if this is an API route that wasn't found
		if len(c.Request.URL.Path) >= 5 && c.Request.URL.Path[:5] == "/api/" {
			c.JSON(404, gin.H{"error": "not found"})
			return
		}

		// Try to serve the file from embedded filesystem
		path := strings.TrimPrefix(c.Request.URL.Path, "/")
		if path == "" {
			path = "index.html"
		}

		// Check if file exists and serve it
		file, err := fsys.Open(path)
		if err == nil {
			defer file.Close()
			stat, _ := file.Stat()
			http.ServeContent(c.Writer, c.Request, path, stat.ModTime(), file.(interface {
				io.ReadSeeker
			}))
			return
		}

		// File doesn't exist, return index.html for SPA routing
		indexFile, _ := fsys.Open("index.html")
		defer indexFile.Close()
		stat, _ := indexFile.Stat()
		http.ServeContent(c.Writer, c.Request, "index.html", stat.ModTime(), indexFile.(interface {
			io.ReadSeeker
		}))
	})

	return r.Run(":" + port)
}

func getProcess(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}

	process, err := core.GetProcess(uint(id))
	if err != nil {
		c.JSON(404, gin.H{"error": "process not found"})
		return
	}
	c.JSON(200, process)
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
		Reason string             `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Use ChangeProcessState for transactional consistency with TUI
	err = core.ChangeProcessState(uint(id), req.Status, req.Reason)
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
