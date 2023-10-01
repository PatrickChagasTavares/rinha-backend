package httpRouter

import (
	"context"
	"net/http"
)

type (
	Router interface {
		Server(port string) error
		ServeHTTP(w http.ResponseWriter, req *http.Request)
		Get(path string, f HandlerFunc)
		Post(path string, f HandlerFunc)
		Put(path string, f HandlerFunc)
		Delete(paht string, f HandlerFunc)
		ParseHandler(h http.HandlerFunc) HandlerFunc
	}

	HandlerFunc func(ctx Context)

	Context interface {
		// Context returns the request's context. To change the context, use
		// Clone or WithContext.
		//
		// The returned context is always non-nil; it defaults to the
		// background context.
		//
		// For outgoing client requests, the context controls cancellation.
		//
		// For incoming server requests, the context is canceled when the
		// client's connection closes, the request is canceled (with HTTP/2),
		// or when the ServeHTTP method returns.
		Context() context.Context
		// Header is an intelligent shortcut for c.Writer.Header().Set(key, value).
		// It writes a header in the response.
		// If value == "", this method removes the header `c.Writer.Header().Del(key)`
		SetHeader(key, value string)
		// JSON serializes the given struct as JSON into the response body.
		// It also sets the Content-Type as "application/json".
		JSON(statusCode int, data any)
		// String writes the given string into the response body.
		String(statusCode int, value string)
		// Bind checks the Method and Content-Type to select a binding engine automatically,
		// Depending on the "Content-Type" header different bindings are used, for example:
		//
		//	"application/json" --> JSON binding
		//	"application/xml"  --> XML binding
		//
		// It parses the request's body as JSON if Content-Type == "application/json" using JSON or XML as a JSON input.
		// It decodes the json payload into the struct specified as a pointer.
		// It writes a 400 error and sets Content-Type header "text/plain" in the response if input is not valid.
		Decode(data any) error
		GetResponseWriter() http.ResponseWriter
		GetRequestReader() *http.Request
		// Query returns the keyed url query value if it exists,
		// otherwise it returns an empty string `("")`.
		// It is shortcut for `c.Request.URL.Query().Get(key)`
		//
		//	    GET /path?id=1234&name=Manu&value=
		//		   c.Query("id") == "1234"
		//		   c.Query("name") == "Manu"
		//		   c.Query("value") == ""
		//		   c.Query("wtf") == ""
		GetQuery(param string) string
		// Param returns the value of the URL param.
		// It is a shortcut for c.Params.ByName(key)
		//
		//	router.GET("/user/:id", func(c *gin.Context) {
		//	    // a GET request to /user/john
		//	    id := c.Param("id") // id == "/john"
		//	    // a GET request to /user/john/
		//	    id := c.Param("id") // id == "/john/"
		//	})
		GetParam(param string) string
		Validate(input any) error
	}
)
