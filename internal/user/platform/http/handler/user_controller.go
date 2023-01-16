package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/krlspj/go-jwt-auth/internal/user/domain"
	"github.com/krlspj/go-jwt-auth/internal/user/service"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(engine *gin.Engine, us service.UserService) {
	handler := UserHandler{
		userService: us,
	}
	// Register Routes
	userGroup := engine.Group("/v1/users")
	userGroup.GET("", handler.CheckHealth)
}

func (h UserHandler) CheckHealth(c *gin.Context) {
	var u domain.User
	u, err := h.userService.GetByID(c.Request.Context(), "12345")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("user", u)
	c.JSON(http.StatusOK, gin.H{"info": u})
	//c.JSON(http.StatusOK, gin.H{"info": "Endpoint working"})
}
