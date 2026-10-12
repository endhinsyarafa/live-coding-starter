package response

import "github.com/gin-gonic/gin"

func SuccessWithMeta(c *gin.Context, statusCode int, message string, data any, meta any) {
	c.JSON(statusCode, gin.H{
		"success": true,
		"message": message,
		"data":    data,
		"meta":    meta,
	})
}
