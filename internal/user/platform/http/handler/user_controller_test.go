package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/krlspj/go-jwt-auth/internal/user/domain"
	"github.com/krlspj/go-jwt-auth/internal/user/service/mocks"
	"github.com/stretchr/testify/mock"
)

func TestUserCheckHealth(t *testing.T) {

	router := gin.Default()

	t.Run("StatusOk", func(t *testing.T) {
		uResp := domain.User{
			ID:   "1234567890",
			Name: "John Doe",
		}
		mockGetUserByID := new(mocks.UserService)
		mockGetUserByID.On("GetByID",
			mock.AnythingOfType("*context.emptyCtx"),
			mock.AnythingOfType("string"),
		).Return(uResp, nil)

		rec := httptest.NewRecorder()

		//router := gin.Default() // Declared outside the Run function

		// Option 1 - Create handler with routes inside
		NewUserHandler(router, mockGetUserByID)

		// Option 2 - Create handler and then create the pointing route
		//handler := UserHandler{
		//	userService: mockGetUserByID,
		//}
		//router.GET("/v1/users", handler.CheckHealth)

		//req, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, "/v1/users", nil)
		req, err := http.NewRequest(http.MethodGet, "/v1/users", nil)
		if err != nil {
			t.Error(err)
		}
		router.ServeHTTP(rec, req)
		//		fmt.Println("req ->", req)
		//		fmt.Println("resp ->", rec.Body)
		type resp struct {
			Info domain.User `json:"info"`
		}
		var rsp resp
		err = json.Unmarshal(rec.Body.Bytes(), &rsp)
		if err != nil {
			t.Error(err)
		}
		if rsp.Info.ID != "1234567890" {
			t.Error("Bad ID. Expected 1234567890, got:", rsp.Info.ID)
		}
		if rsp.Info.Name != "John Doe" {
			t.Error("Bad Name. Expected John Doe, got:", rsp.Info.Name)
		}

	})
}
