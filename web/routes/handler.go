package routes

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ermites-io/passwd"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"hypebase/ent"
	"hypebase/ent/user"
)

type Handler struct {
	db *ent.Client
}

func NewHandler(db *ent.Client) *Handler {
	return &Handler{db: db}
}

func (h *Handler) Index(c *gin.Context) {
	sess := sessions.Default(c)

	sessionID := sess.Get("session_id")

	if sessionID == nil {
		c.HTML(http.StatusOK, "index.go.html", gin.H{"loggedIn": sessionID != ""})
	}
	c.Redirect(http.StatusFound, "/dashboard")
}

func (h *Handler) Register(c *gin.Context) {
	c.HTML(http.StatusOK, "register.go.html", nil)
}

func (h *Handler) RegisterForm(c *gin.Context) {
	type Form struct {
		Email    string `form:"email" binding:"required"`
		Password string `form:"pass" binding:"required"`
	}

	// session := sessions.Default(c)

	var form Form
	err := c.ShouldBind(&form)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("Failed to bind form: %e", err)})
		return
	}

	email, pass := strings.Trim(form.Email, " "), form.Password

	if email == "" || pass == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameters can't be empty"})
		return
	}

	rec, err := h.db.User.Query().Where(user.Email(email)).Only(c)
	if err != nil {
		switch err.(type) {
		case *ent.NotFoundError:
			break
		default:
			if email == rec.Email {
				c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed sign up due to database error"})
			return
		}
	}

	p, err := passwd.New(passwd.Argon2idDefault)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Password library failed to initialize"})
		return
	}

	hashedPass, err := p.Hash([]byte(pass))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Password library failed to hash password"})
		return
	}

	_, err = h.db.User.Create().
		SetEmail(email).
		SetPasswordHash(string(hashedPass)).
		Save(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user to database"})
	}

	c.Redirect(http.StatusFound, "/dashboard")
}

func (h *Handler) Login(c *gin.Context) {
	c.HTML(http.StatusOK, "login.go.html", nil)
}

func (h *Handler) LoginForm(c *gin.Context) {
	type Form struct {
		Email    string `form:"user" json:"user" binding:"required"`
		Password string `form:"pass" json:"pass" binding:"required"`
	}

	session := sessions.Default(c)

	var form Form
	err := c.ShouldBind(&form)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind form"})
		return
	}

	email, pass := strings.Trim(form.Email, " "), strings.Trim(form.Password, " ")

	if strings.Trim(email, " ") == "" || strings.Trim(pass, " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameters can't be empty"})
		return
	}

	rec, err := h.db.User.Query().Where(user.Email(email)).Only(c)
	if err != nil {
		switch err.(type) {
		case *ent.NotFoundError:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found with those credentialis"})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to authorize user in database"})
			return
		}
	}

	p, err := passwd.New(passwd.Argon2idDefault)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Password library failed to initialize"})
		return
	}

	err = p.Compare([]byte(rec.PasswordHash), []byte(pass))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Password library failed to hash password"})
		return
	}

	if rec.Email != email {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Emails do not match"})
		return
	}

	session.Set("sessionID", email)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully authenticated user"})
	c.Redirect(http.StatusFound, "/dashboard")
}

func (h *Handler) Logout(c *gin.Context) {
	session := sessions.Default(c)
	u := session.Get("sessionID")

	if u == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session token"})
		return
	}

	session.Delete("sessionID")
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
	c.Redirect(http.StatusFound, "/")
}

func (h *Handler) Dashboard(c *gin.Context) {
	c.HTML(http.StatusOK, "dashboard.go.html", gin.H{})
}

func (h *Handler) Services(c *gin.Context) {

}
