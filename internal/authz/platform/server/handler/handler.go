package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	authz_uc "github.com/krlspj/go-jwt-auth/internal/authz/usecase"
	"github.com/krlspj/go-jwt-auth/internal/config"
)

type AuthzMidd struct {
	app     *config.AppConfig
	authzUc authz_uc.AuthzUsecase
	//jwtUc   jwt_uc.JwtUsecase
}

func NewAuthz(a *config.AppConfig, authzuc authz_uc.AuthzUsecase) *AuthzMidd {
	return &AuthzMidd{
		app:     a,
		authzUc: authzuc,
		//jwtUc:   jwtuc,
	}
}

// VerifyToken is a simple token verification
func (h *AuthzMidd) VeriyToken(c *gin.Context) {
	// get token from body
	//var token map[string]string
	//if err := c.BindJSON(&token); err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	return
	//}
	//t := token["token"]

	// get token from header "X-Access-Token"
	token := c.Request.Header.Get("X-Access-Token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("No Authorization header provided")})
		c.Abort()
		return
	}

	//fmt.Println("->", token)

	//claims, err := h.authzUc.CheckValidToken(token)
	err := h.authzUc.CheckValidToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	//log.Println("[AuthzMidd] [VerifyToken] claims", claims)
	//c.Set("claims", claims)

	// if Ok
	c.Next()

}

func (h *AuthzMidd) OnlyAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get token from header "X-Access-Token"
		token := c.Request.Header.Get("X-Access-Token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No Authorization header provided"})
			c.Abort()
			return
		}

		err := h.authzUc.IsAdmin(c.Request.Context(), token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		c.Next()

	}
}

//func WsJwtAuthorize() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		params := c.Request.URL.Query()
//		fmt.Println("urls params ->", params)
//		fmt.Println("url path:", c.Request.URL)
//		fmt.Println("clientIP:", c.ClientIP())
//		token := params.Get("token")
//		fmt.Println("token->", token)
//		if token == "" {
//			c.JSON(http.StatusBadRequest, gin.H{"error": "missing access token"})
//			c.Abort()
//			return
//		}
//
//		claims, err := authz.ValidateToken(token)
//		if err != "" {
//			log.Printf("[WsJwtAuthorize] error validating token: %v", err)
//			c.JSON(http.StatusUnauthorized, gin.H{"error": err})
//			c.Abort()
//			return
//		}
//
//		c.Set("userId", claims.UserId)
//		c.Set("clientIP", c.ClientIP())
//
//		c.Next()
//
//	}
//}
/*
func RefreshToken(auc *authz.UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		// get claims from context
		userId, ok := c.Get("userId")
		foundClaims := true
		if !ok {
			foundClaims = false
		}
		//email, ok := c.Get("email")
		//if !ok {
		//	foundClaims = false
		//}
		//username, ok := c.Get("username")
		//if !ok {
		//	foundClaims = false
		//}
		//roleId, ok := c.Get("role")
		//if !ok {
		//	foundClaims = false
		//}

		if !foundClaims {
			c.JSON(http.StatusConflict, gin.H{"error": "missing claims"})
			c.Abort()
			return
		}

		token, err := authz.GenerateToken(fmt.Sprint(userId))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// refresh token in db
		//user := models.User{
		//	Token: token,
		//}
		//err = auc.DB.PartialyUpdateUser(c.Request.Context(), fmt.Sprint(userId), user)
		//if err != nil {
		//	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		//	return
		//}
		// store refreshed token in context
		c.Set("refreshToken", token)
		c.Header("X-Access-Token", token)
	}
}
*/
//func VerifyAbilities(auc *authz.UseCase, abilities []string) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		uid, ok := c.Get("userId")
//		if !ok {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "userId not found"})
//		}
//		fmt.Println("userid context", uid)
//		allowAccess, info, err := auc.IsAuthorized(c.Request.Context(), uid.(string), abilities)
//		if err != nil {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//			c.Abort()
//			return
//		}
//		if !allowAccess {
//			c.JSON(http.StatusUnauthorized, gin.H{"error": info})
//			c.Abort()
//			return
//		}
//		c.Next()
//	}
//}
//
