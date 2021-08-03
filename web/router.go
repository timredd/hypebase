package web

import (
	"io/fs"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"hypebase/web/routes"
)

type Router struct {
	*gin.Engine
}

func NewRouter() *Router {
	return &Router{Engine: gin.Default()}
}

func (r *Router) Route(h *routes.Handler) {
	r.GET("/", h.Index)
	r.GET("/register", h.Register)
	r.GET("/login", h.Login)

	r.POST("/register", h.RegisterForm)
	r.POST("/login", h.LoginForm)

	auth := r.Group("/")
	auth.Use(AuthRequired())
	{
		auth.POST("/logout", h.Logout)

		auth.GET("/dashboard", h.Dashboard)
		auth.GET("/services", h.Services)
	}

	system := r.Group("/sys")
	{
		system.POST("/validate_email", h.ValidateEmail)
	}
}

func (r *Router) AddHTMLFiles(rootDir string) error {
	var htmlFiles []string
	err := filepath.WalkDir(rootDir,
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if filepath.Ext(path) == ".html" {
				htmlFiles = append(htmlFiles, path)
			}
			return nil
		},
	)

	r.LoadHTMLFiles(htmlFiles...)
	return err
}

func (r *Router) AddStaticFiles(rootDir string) {
	r.Static("/static", rootDir)
}
