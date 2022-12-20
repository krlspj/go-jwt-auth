package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/krlspj/go-jwt-auth/internal/auth/service"
)

type AuthHandler struct {
	// app *config.AppConfig
	// rs  render_service.RenderService
	authService service.AuthService
}

//func NewHanldersRepo(a *config.AppConfig, render render_service.RenderService) *HandlerRepo {
//	return &HandlerRepo{
//		app: a,
//		rs:  render,
//	}
//}

func NewAuthHanlderRepo(as service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: as,
	}
}

func (h *AuthHandler) CheckHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fmt.Println("CheckHandler")
		ctx.String(http.StatusOK, "everything is ok!")
	}

}

func (h *AuthHandler) FindUsersC(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"users": "found"})
}

func (h *AuthHandler) FindUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := h.authService.FindUsers(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		uRestList, err := toUserRestList(users)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, gin.H{"users": uRestList})
	}
}
