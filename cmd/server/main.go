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
	e.Use(
		h.SessionMiddleware,
		middleware.RateLimiterWithConfig(config),
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

	auth := e.Group("/auth")
	auth.GET("/login", h.GetLoginPage)
	// auth.POST("/login", h.PostLogin)
	auth.GET("/signup", h.GetSignUpPage)
	// auth.POST("/signup", h.PostSignUp)
	auth.GET("/logout", h.GetLogOut)
	auth.GET("/oauth/:provider", h.GetOAuthFlow)
	auth.GET("/oauth/:provider/callback", h.GetOAuthCallback)

	articles := e.Group("/articles")
	articles.GET("", h.GetArticlesPage)
	articles.GET("/grid", h.GetArticlesGrid)
	article := articles.Group("/:slug")
	article.GET("", h.GetArticlePage)
	// articleID.GET("/comments", h.GetArticleComments)
	// articleID.POST("/comments", h.PostArticleComment)
	// articleID.DELETE("/comments/:commentID", h.DeleteArticleComment)

	editor := e.Group("/editor/:articleID", h.RoleMiddleware(internal.Superuser), h.URLIDMiddleware("articleID"))
	editor.GET("", h.GetArticleEditorPage)
	editor.DELETE("", h.DeleteArticle)
	editor.PATCH("/title", h.PatchArticleTitle)
	editor.PATCH("/description", h.PatchArticleDescription)
	editor.PATCH("/content", h.PatchArticleContent)
	editor.PATCH("/visibility", h.PatchArticleVisibility)
	// editor.PUT("/tags", h.PutArticleTags)

	// tags := e.Group("/tags")
	// tags.GET("", h.GetTags)
	// tags.POST("", h.PostTag)
	// tags.DELETE("/:tagID", h.DeleteTag)

	internal.GracefulShutdown(e, internal.Settings.Port)
}
