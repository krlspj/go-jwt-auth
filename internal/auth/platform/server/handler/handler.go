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
	"github.com/krlspj/go-jwt-auth/internal/config"
	jwt_uc "github.com/krlspj/go-jwt-auth/internal/jwt/service"
)

type AuthHandler struct {
	app         *config.AppConfig
	authService service.AuthService
	jwtService  jwt_uc.JwtService
	validate    *validator.Validate
}

//func NewHanldersRepo(a *config.AppConfig, render render_service.RenderService) *HandlerRepo {
//	return &HandlerRepo{
//		app: a,
//		rs:  render,
//	}
//}

func NewAuthHanlderRepo(
	a *config.AppConfig,
	as service.AuthService,
	jwtS jwt_uc.JwtService,
	v *validator.Validate,
) *AuthHandler {

	return &AuthHandler{
		app:         a,
		authService: as,
		jwtService:  jwtS,
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

	u, err := h.authService.LoginUser(c.Request.Context(), user.Name, user.Password)
	if err != nil {
		c.JSON(http.StatusForbidden, err.Error())
		return
	}

	// if correct credentials generate token
	token, err := h.jwtService.GenerateToken(h.app.TokenLifetime, u.ID(), u.Name())
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	u.SetToken(token)
	c.Header("X-Access-Token", token)
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

//type jwtClaims struct {
//	Username string `json:"username,omitemtpy"`
//	UserId   string `json:"userId,omitempty"`
//	jwt.StandardClaims
//}
//
//// generateToken returns the generated token and an error
//func generateToken(userId, username string) (string, error) {
//
//	claims := &jwtClaims{
//		Username: username,
//		UserId:   userId,
//		//Email:    email,
//		//RoleId: role,
//		StandardClaims: jwt.StandardClaims{
//			ExpiresAt: time.Now().Local().Add(1 * time.Minute).Unix(),
//		},
//	}
//	secretKey := "this is my secret"
//	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secretKey))
//
//	//refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(secret_key))
//
//	if err != nil {
//		log.Panic(err)
//		return "", err
//	}
//
//	return token, nil
//}
//
//func validateToken(signedToken string) error {
//	token, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
//		// Don't forget to validate the alg is what you expect:
//		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
//		}
//
//		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
//		hmacSampleSecret := "this is my secret"
//
//		return []byte(hmacSampleSecret), nil
//	})
//
//	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
//		fmt.Println("claims >", claims)
//	} else {
//		fmt.Println(err)
//		return err
//	}
//	return nil
//}
//
