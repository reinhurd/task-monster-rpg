package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"rpgMonster/internal/core"
	"rpgMonster/internal/model"
)

const (
	authHeader      = "Authorization"
	tempTokenCookie = "token"
)

type userCreateRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func SetupRouter(svc *core.Service) *gin.Engine {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		//TODO default template for default page
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(svc.GetTemplate()))
	})

	// Get active tasks for current user
	r.GET("api/tasks", func(c *gin.Context) {
		userID, err := auth(svc, c)
		if err != nil || userID == "" {
			c.String(http.StatusUnauthorized, "Unauthorized")
			return
		}
		tasks, err := svc.GetListTasksByUserID(context.Background(), userID)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		} else {
			c.JSON(http.StatusOK, tasks)
		}
	})

	// Get task by Id. Need to check rights
	r.GET("api/tasks/:id", func(c *gin.Context) {
		id := c.Param("id")
		userID, err := auth(svc, c)
		if err != nil || userID == "" {
			c.String(http.StatusUnauthorized, "Unauthorized")
			return
		}
		task, err := svc.GetTask(context.Background(), id, userID)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		} else {
			c.JSON(http.StatusOK, task)
		}
	})

	// Create a new task from GPT
	r.GET("api/tasks/create/gpt", func(c *gin.Context) {
		userID, err := auth(svc, c)
		if err != nil || userID == "" {
			c.String(http.StatusUnauthorized, "Unauthorized")
			return
		}
		req := c.Query("req")
		task, err := svc.CreateTaskFromGPTByRequest(req, userID)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		} else {
			c.JSON(http.StatusOK, task)
		}
	})

	// Create a new task manually
	r.POST("api/tasks/create", func(c *gin.Context) {
		var task model.Task

		if err := c.ShouldBindJSON(&task); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		task.Executor, _ = auth(svc, c)

		err := svc.CreateTask(context.Background(), &task)
		if err != nil {
			responseText := err.Error()
			c.String(http.StatusInternalServerError, responseText)
		} else {
			responseText := "Task created"
			c.String(http.StatusOK, responseText)
		}
	})

	r.POST("api/tasks/:id", func(c *gin.Context) {
		id := c.Param("id")
		var task model.Task
		if err := c.ShouldBindJSON(&task); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userID, err := auth(svc, c)
		if err != nil || userID == "" {
			c.String(http.StatusUnauthorized, "Unauthorized")
			return
		}
		err = svc.UpdateTask(context.Background(), &model.Task{
			BizId:       id,
			Title:       task.Title,
			Description: task.Description,
			Executor:    task.Executor,
			Reviewer:    task.Reviewer,
			Completed:   task.Completed,
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   time.Now(),
			Deadline:    task.Deadline,
			Tags:        task.Tags,
		})
		if err != nil {
			responseText := err.Error()
			c.String(http.StatusInternalServerError, responseText)
		} else {
			responseText := "Updated Task ID:" + id
			c.String(http.StatusOK, responseText)
		}
	})

	r.PUT("api/tasks/:id/status", func(c *gin.Context) {
		id := c.Param("id")
		userID, err := auth(svc, c)
		if err != nil || userID == "" {
			c.String(http.StatusUnauthorized, "Unauthorized")
			return
		}
		err = svc.UpdateTask(context.Background(), &model.Task{
			BizId:     id,
			Completed: true,
		})
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		} else {
			c.String(http.StatusOK, "Task updated")
		}
	})

	//create user
	r.POST("api/users", func(c *gin.Context) {
		var user userCreateRequest
		err := c.ShouldBindBodyWithJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		_, err = svc.CreateNewUser(user.Login, user.Password)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		} else {
			c.String(http.StatusOK, "User created")
		}
	})

	//auth prestep - if success, frontend can set log:pass in header - but maybe better to use token?
	r.POST("api/users/login", func(c *gin.Context) {
		var user userCreateRequest
		err := c.ShouldBindBodyWithJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userId, token, err := svc.CheckPassword(user.Login, user.Password)
		if err != nil || userId == "" {
			c.String(http.StatusInternalServerError, err.Error())
		} else {
			c.JSON(http.StatusOK, gin.H{tempTokenCookie: token, "success": true})
			//set token to cookie
			c.SetCookie(tempTokenCookie, token, 3600000, "/", "localhost", false, true)
		}
	})

	return r
}

func auth(svc *core.Service, c *gin.Context) (userID string, err error) {
	token, err := c.Cookie(tempTokenCookie)
	if err != nil {
		return "", fmt.Errorf("no token")
	}
	if token == "" {
		return "", fmt.Errorf("empty cookie")
	}
	userID, err = svc.GetUserByTempToken(token)
	if err != nil {
		return "", fmt.Errorf("invalid token")
	}
	return userID, nil
}
