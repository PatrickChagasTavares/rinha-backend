package httpRouter

import (
	"context"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/patrickchagastavares/rinha-backend/pkg/validator"
)

type (
	ginRouter struct {
		router *gin.Engine
	}
)

func NewGinRouter() Router {
	if os.Getenv("env") != "local" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	router.Use(
		// Set the content type default = application/json
		setContentType("application/json"),
		gin.Recovery(),
	)

	return &ginRouter{
		router: router,
	}
}

func setContentType(contentType string) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Content-Type", contentType)
		ctx.Next()
	}
}

func (r *ginRouter) Get(path string, f HandlerFunc) {
	r.router.GET(path, func(ctx *gin.Context) {
		f(newGinContext(ctx))
	})
}

func (r *ginRouter) Post(path string, f HandlerFunc) {
	r.router.POST(path, func(ctx *gin.Context) {
		f(newGinContext(ctx))
	})
}

func (r *ginRouter) Put(path string, f HandlerFunc) {
	r.router.PUT(path, func(ctx *gin.Context) {
		f(newGinContext(ctx))
	})
}

func (r *ginRouter) Delete(path string, f HandlerFunc) {
	r.router.DELETE(path, func(ctx *gin.Context) {
		f(newGinContext(ctx))
	})
}

func (r *ginRouter) Server(port string) error {
	return http.ListenAndServe(port, r.router)
}

func (r *ginRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}

func (m *ginRouter) ParseHandler(h http.HandlerFunc) HandlerFunc {
	return func(c Context) {
		h(c.GetResponseWriter(), c.GetRequestReader())
	}
}

type ginContext struct {
	r *gin.Context
	v validator.Validator
}

func newGinContext(ctx *gin.Context) Context {
	return &ginContext{
		r: ctx,
		v: validator.New(),
	}
}

func (c *ginContext) Context() context.Context {
	return c.r.Request.Context()
}

func (c *ginContext) SetHeader(key, value string) {
	c.r.Header(key, value)
}

func (c *ginContext) String(statusCode int, value string) {
	c.r.Header("Content-Type", "text/plain")
	c.r.String(statusCode, value)
}

func (c *ginContext) JSON(statusCode int, data any) {
	c.r.JSON(statusCode, data)
}

func (c *ginContext) Decode(data any) error {
	return c.r.Bind(&data)
}

func (c *ginContext) GetQuery(query string) string {
	return c.r.Query(query)
}

func (c *ginContext) GetParam(param string) string {
	return c.r.Param(param)
}

func (c *ginContext) GetResponseWriter() http.ResponseWriter {
	return c.r.Writer
}

func (c *ginContext) GetRequestReader() *http.Request {
	return c.r.Request
}

func (c *ginContext) Validate(input any) error {
	if err := c.v.Validate(input); err != nil {
		return err.ToHttpErr()
	}
	return nil
}
