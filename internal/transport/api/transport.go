package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"rpgMonster/internal/core"
	"rpgMonster/internal/model"
)

const (
	authHeader = "Authorization"
)

type userCreateRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func SetupRouter(svc *core.Service) *gin.Engine {
	r := gin.Default()

	//// Just ping
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, svc.DoSomething())
	})

	r.GET("/", func(c *gin.Context) {
		//TODO default template for default page
		c.String(http.StatusOK, svc.GetTemplate())
	})

	// Get active tasks for current user
	r.GET("api/tasks", func(c *gin.Context) {
		userID, err := auth(c.GetHeader(authHeader), svc)
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
		userID, err := auth(c.GetHeader(authHeader), svc)
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
		userID, err := auth(c.GetHeader(authHeader), svc)
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
	r.POST("api/tasks", func(c *gin.Context) {
		var task model.Task

		if err := c.ShouldBindJSON(&task); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		task.Executor, _ = auth(c.GetHeader(authHeader), svc)

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
		userID, err := auth(c.GetHeader(authHeader), svc)
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
		userID, err := auth(c.GetHeader(authHeader), svc)
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
		userId, err := svc.CheckPassword(user.Login, user.Password)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		} else {
			c.String(http.StatusOK, userId)
		}
	})

	return r
}

func auth(header string, svc *core.Service) (userID string, err error) {
	if header == "" {
		return "", fmt.Errorf("empty header")
	}
	log, pass, err := ParseAuthHeader(header)
	if err != nil {
		return "", fmt.Errorf("invalid header")
	}
	userID, err = svc.CheckPassword(log, pass)
	if err != nil {
		return "", fmt.Errorf("invalid token")
	}
	return userID, nil
}

// extracts the user login and password from header string
// The header Authorization should contain "<user-id>:<password>".
func ParseAuthHeader(header string) (login, password string, err error) {
	parts := strings.Split(header, ":")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		err = fmt.Errorf("invalid Authorization header format")
		return
	}

	login = parts[0]
	password = parts[1]
	return
}
