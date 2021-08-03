package routes

import (
	"net/http"
	"net/mail"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"hypebase/ent/user"
)

func (h *Handler) ValidateEmail(c *gin.Context) {
	type Req struct {
		Email string `form:"email" json:"email" binding:"required"`
	}

	var data Req
	err := c.ShouldBind(&data)
	email := data.Email
	log.Debug().Msgf("validating email: %s", email)

	_, err = mail.ParseAddress(email)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"value": false, "message": "Email address is not valid"})
		return
	}

	isExist, err := h.db.User.Query().Where(user.Email(email)).Exist(c)
	if err != nil || isExist {
		c.JSON(http.StatusOK, gin.H{"value": false, "message": "Email address already in use"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"value": true, "message": ""})
	return
}
