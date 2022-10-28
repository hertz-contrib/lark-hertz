/*
 * Copyright 2022 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package lark_hertz

import (
	"bytes"
	"context"
	"net/http"

	"github.com/chyroc/lark"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/adaptor"
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
			cli.EventCallback.ListenSecurityCallback(ctx, header, body, adaptor.GetCompatResponseWriter(&c.Response))
		} else {
			cli.EventCallback.ListenCallback(ctx, body, adaptor.GetCompatResponseWriter(&c.Response))
		}
	}
}
