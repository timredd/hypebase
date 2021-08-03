package web

import (
	"io"
	"time"

	"github.com/gin-contrib/logger"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"hypebase/ent"
	"hypebase/web/routes"
)

type Server struct {
	db *ent.Client
}

func NewServer(db *ent.Client) *Server {
	return &Server{db: db}
}

func (s *Server) Run() error {
	r := NewRouter()
	r.Use(logger.SetLogger(
		logger.WithUTC(true),
		logger.WithLogger(func(c *gin.Context, out io.Writer, latency time.Duration) zerolog.Logger {
			return zerolog.New(out).With().
				Str("path", c.Request.URL.Path).
				Dur("latency", latency).
				Logger()
		}),
	))
	r.Use(gin.Recovery())
	r.Use(sessions.Sessions("session", cookie.NewStore([]byte("secret"))))

	r.LoadHTMLGlob("web/views/**/*.go.html")

	h := routes.NewHandler(s.db)
	r.Route(h)
	r.AddStaticFiles("web/static")

	err := r.Run(":8081")
	if err != nil {
		return err
	}
	return nil
}
