package lark_hertz

import (
	"bytes"
	"context"
	"net/http"

	"github.com/chyroc/lark"
	"github.com/cloudwego/hertz/pkg/app"
)

type Option struct {
	IgnoreCheckSignature bool
}

func WithIgnoreCheckSignature(b bool) func(option *Option) {
	return func(option *Option) {
		option.IgnoreCheckSignature = b
	}
}

// ListenCallback listen lark callback, return hertz app.HandlerFunc
func ListenCallback(cli *lark.Lark, options ...func(option *Option)) app.HandlerFunc {
	opt := new(Option)
	for _, v := range options {
		v(opt)
	}
	return func(ctx context.Context, c *app.RequestContext) {
		header := http.Header{}
		c.Request.Header.VisitAll(func(key, value []byte) {
			header.Add(string(key), string(value))
		})
		body := bytes.NewReader(c.Request.Body())
		if !opt.IgnoreCheckSignature {
			cli.EventCallback.ListenSecurityCallback(ctx, header, body, &hertzResponseWriter{c})
		} else {
			cli.EventCallback.ListenCallback(ctx, body, &hertzResponseWriter{c})
		}
	}
}

type hertzResponseWriter struct {
	ctx *app.RequestContext
}

func (r *hertzResponseWriter) Header() http.Header {
	header := http.Header{}
	r.ctx.Response.Header.VisitAll(func(key, value []byte) {
		header.Add(string(key), string(value))
	})
	return header
}

func (r *hertzResponseWriter) Write(i []byte) (int, error) {
	r.ctx.Response.SetBodyRaw(i)
	return len(i), nil
}

func (r *hertzResponseWriter) WriteHeader(statusCode int) {
	r.ctx.Response.SetStatusCode(statusCode)
}
