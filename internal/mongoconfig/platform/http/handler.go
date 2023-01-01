package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/krlspj/go-jwt-auth/internal/mongoconfig/service"
)

type DBConfigHanlder struct {
	configService service.DBConfigService
}

func NewDBConfigHanlder(cs service.DBConfigService) *DBConfigHanlder {
	return &DBConfigHanlder{
		configService: cs,
	}
}

func (h DBConfigHanlder) UpdateConfig(c *gin.Context) {
}
