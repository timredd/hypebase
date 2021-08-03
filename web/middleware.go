package web

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// func Session() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		session := sessions.Default(c)
//
// 		sessionID := session.Get("sessionID")
// 		if sessionID == nil {
// 			id := uuid.New()
// 			session.Set("sessionID", id)
// 			_ = session.Save()
// 			return
// 		}
// 	}
// }

func Push() gin.HandlerFunc {
	return func(c *gin.Context) {
		pusher := c.Writer.Pusher()
		if pusher != nil {
			// TODO: walk through files to load individually, wildcard usage not possible
			err := pusher.Push("/static/*", nil)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("pushing file failed: %e", err)})
				return
			}
		}

	}
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		sessionID := session.Get("id")
		if sessionID == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		c.Next()
	}
}
