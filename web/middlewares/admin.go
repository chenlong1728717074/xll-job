package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xll-job/utils"
	"xll-job/web/dto"
)

func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		value, exists := c.Get("claims")
		claims := value.(*utils.CustomClaims)
		if !exists || claims.Role == 2 {
			c.JSON(http.StatusOK, dto.NewResponse(http.StatusForbidden, "暂无权限", ""))
			c.Abort()
		}
		c.Next()
	}
}
