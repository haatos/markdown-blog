package main

import (
	"github.com/haatos/markdown-blog/internal"
	"github.com/haatos/markdown-blog/internal/data"
	"github.com/haatos/markdown-blog/internal/handler"
	"github.com/haatos/markdown-blog/internal/templates"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	internal.ReadDotenv()
	internal.Settings = internal.NewSettings()

	rdb := data.InitDatabase(true)
	defer rdb.Close()
	rwdb := data.InitDatabase(false)
	defer rwdb.Close()

	data.RunMigrations(rwdb)

	h := handler.NewHandler(rdb, rwdb)
	handler.NewSecureCookie()

	e := echo.New()
	e.Renderer = templates.NewRenderer()
	e.HTTPErrorHandler = handler.ErrorHandler
	loggerFormat := "${method} ${uri} [${status}] (${latency_human}) | ${short_file}:${line} | ${message}\n"
	e.Logger.SetHeader(loggerFormat)

	config := internal.GetRateLimiterConfig()
	e.Use(middleware.RateLimiterWithConfig(config))

	e.Use(
		middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: loggerFormat,
		}),
		middleware.GzipWithConfig(middleware.GzipConfig{
			Skipper: middleware.DefaultSkipper,
			Level:   3,
		}),
	)

	e.Static("/", "public")

	e.GET("", h.GetIndexPage)

	// auth := e.Group("/auth")
	// auth.GET("/login", h.GetLoginPage)
	// auth.POST("/login", h.PostLogin)
	// auth.GET("/signup", h.GetSignUpPage)
	// auth.POST("/signup", h.PostSignUp)
	// auth.GET("/oauth/:provider", h.GetOAuthFlow)
	// auth.GET("/oauth/:provider/callback", h.GetOAuthCallback)

	// articles := e.Group("/articles")
	// articles.GET("", h.GetArticlesPage)
	// articleID := articles.Group("/:slug")
	// articleID.GET("", h.GetArticlePage)
	// articleID.GET("/comments", h.GetArticleComments)
	// articleID.POST("/comments", h.PostArticleComment)
	// articleID.DELETE("/comments/:commentID", h.DeleteArticleComment)

	// editor := e.Group("/editor", h.EnsureSuperuserMiddleware)
	// editor.GET("/:articleID", h.GetArticleEditorPage)
	// editor.PATCH("/:articleID/title", h.PatchArticleTitle)
	// editor.PATCH("/:articleID/description", h.PatchArticleDescription)
	// editor.PATCH("/:articleID/content", h.PatchArticleContent)
	// editor.PATCH("/:articleID/visibility", h.PatchArticleVisibility)
	// editor.PUT("/:articleID/tags", h.PutArticleTags)

	// tags := e.Group("/tags")
	// tags.GET("", h.GetTags)
	// tags.POST("", h.PostTag)
	// tags.DELETE("/:tagID", h.DeleteTag)

	internal.GracefulShutdown(e, internal.Settings.Port)
}
