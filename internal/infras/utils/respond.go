package utils

import "github.com/gin-gonic/gin"

func JsonRespond(code int, message string, data interface{}, c *gin.Context) {
	c.JSON(code, gin.H{
		"code"	: code,
		"message": message,
		"data"   : data,
	})
}
