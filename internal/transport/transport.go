package transport

import (
	"context"
	"net/http"
	"rpgMonster/internal/core"
	"rpgMonster/internal/tasks"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	//// Just ping
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Get active tasks for current user
	r.GET("api/tasks", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Get task by Id. Need to check rights
	r.GET("api/tasks/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.String(http.StatusOK, "Task Id: "+id)
	})

	// Create a new task
	r.POST("api/tasks", func(c *gin.Context) {
		var task tasks.Task

		if err := c.ShouldBindJSON(&task); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		task.Executor, _ = core.GetCurrentUserID(c.Request.Header)

		err := tasks.CreateTask(context.Background(), &task)
		if err != nil {
			responseText := err.Error()
			c.String(http.StatusInternalServerError, responseText)
		} else {
			responseText := "Task created"
			c.String(http.StatusOK, responseText)
		}
	})

	r.PUT("api/tasks/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.String(http.StatusOK, "Task Id: "+id)
	})

	r.PUT("api/tasks/:id/status", func(c *gin.Context) {
		id := c.Param("id")
		err := tasks.UpdateTask(context.Background(), &tasks.Task{
			BizId:     id,
			Completed: true,
		})
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		} else {
			c.String(http.StatusOK, "Task updated")
		}
	})

	//// User logic
	r.POST("api/users", func(c *gin.Context) {
		var user struct {
			Login    string `json:"login"`
			Password string `json:"password"`
		}
		err := c.ShouldBindBodyWithJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		_, err = core.CreateNewUser(user.Login, user.Password)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		} else {
			c.String(http.StatusOK, "User created")
		}
	})

	r.POST("api/users/login", func(c *gin.Context) {
		var user struct {
			Login    string `json:"login"`
			Password string `json:"password"`
		}
		err := c.ShouldBindBodyWithJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userId, err := core.CheckPassword(user.Login, user.Password)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		} else {
			c.String(http.StatusOK, userId)
		}
	})

	return r
}
