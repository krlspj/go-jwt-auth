package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/krlspj/go-jwt-auth/internal/auth/platform/storage/inmemory"
	"github.com/krlspj/go-jwt-auth/internal/auth/platform/storage/mongodb"
	"github.com/krlspj/go-jwt-auth/internal/auth/service"
)

type AuthHandler struct {
	// app *config.AppConfig
	// rs  render_service.RenderService
	authService service.AuthService
	validate    *validator.Validate
}

//func NewHanldersRepo(a *config.AppConfig, render render_service.RenderService) *HandlerRepo {
//	return &HandlerRepo{
//		app: a,
//		rs:  render,
//	}
//}

func NewAuthHanlderRepo(as service.AuthService, v *validator.Validate) *AuthHandler {
	return &AuthHandler{
		authService: as,
		validate:    v, // validator.New(),
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

func (h *AuthHandler) FindUser(c *gin.Context) {
	id := c.Param("id")
	user, err := h.authService.FindUser(c.Request.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, mongodb.ErrBadMethod):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})

		case errors.Is(err, inmemory.ErrNotUserFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

		default:
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})

		}
		return
	}
	fmt.Println("user:", user)
	c.JSON(http.StatusOK, gin.H{"user": toUserResp(user)})

}

func (h *AuthHandler) FindUserByField(c *gin.Context) {
	x := make(map[string]string)
	if err := c.BindJSON(&x); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(x) != 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("Only the first json value will be used for the search, remove the others").Error()})
		return
	}

	keys := make([]string, 0, len(x))
	for k := range x {
		keys = append(keys, k)
	}
	user, err := h.authService.FindUserByField(c.Request.Context(), keys[0], x[keys[0]])
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusFound, toUserResp(user))

}

func (h *AuthHandler) Login(c *gin.Context) {
	//user := new(userHttp)
	var user userReq
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.validateReq(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	fmt.Println("user", user)

	u, err := h.authService.LoginUser(c.Request.Context(), user.Name, user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, gin.H{"info": "correct credentials", "user": toUserResp(u)})
}
func (h *AuthHandler) CreateNewUser(c *gin.Context) {
	var user userReq
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//err := user.ValidateUser()
	err := h.validateReq(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err = h.authService.CreateUser(c.Request.Context(), user.toDomainUser())
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, "")
}

func (h *AuthHandler) validateReq(req userReq) error {
	if err := h.validate.Struct(req); err != nil {
		return err
	}
	return nil
}
