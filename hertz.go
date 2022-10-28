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
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/larksuite/oapi-sdk-go/v3/card"
	"github.com/larksuite/oapi-sdk-go/v3/event"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
)

func doProcess(c context.Context, ctx *app.RequestContext, reqHandler larkevent.IReqHandler, options ...larkevent.OptionFunc) {
	eventReq, err := translate(c, &ctx.Request)
	if err != nil {
		ctx.Response.SetStatusCode(consts.StatusInternalServerError)
		ctx.Response.SetBodyString(err.Error())
		return
	}

	eventResp := reqHandler.Handle(c, eventReq)

	err = write(c, &ctx.Response, eventResp)
	if err != nil {
		reqHandler.Logger().Error(c, fmt.Sprintf("write resp result error:%s", err.Error()))
	}
}

func NewCardActionHandlerFunc(cardActionHandler *larkcard.CardActionHandler, options ...larkevent.OptionFunc) func(c context.Context, ctx *app.RequestContext) {
	cardActionHandler.InitConfig(options...)
	return func(c context.Context, ctx *app.RequestContext) {
		doProcess(c, ctx, cardActionHandler, options...)
	}
}

func NewEventHandlerFunc(eventDispatcher *dispatcher.EventDispatcher, options ...larkevent.OptionFunc) func(c context.Context, ctx *app.RequestContext) {
	eventDispatcher.InitConfig(options...)
	return func(c context.Context, ctx *app.RequestContext) {
		doProcess(c, ctx, eventDispatcher, options...)
	}
}

func write(c context.Context, resp *protocol.Response, eventResp *larkevent.EventResp) error {
	resp.SetStatusCode(eventResp.StatusCode)
	for k, vs := range eventResp.Header {
		for _, v := range vs {
			resp.Header.Add(k, v)
		}
	}

	if len(eventResp.Body) > 0 {
		resp.SetBody(eventResp.Body)
		return nil
	}
	return nil
}

func translate(ctx context.Context, req *protocol.Request) (*larkevent.EventReq, error) {
	headers := make(map[string][]string)
	req.Header.VisitAll(func(key, value []byte) {
		keyStr := string(key)
		valueStr := string(value)
		headers[keyStr] = append(headers[keyStr], valueStr)
	})

	eventReq := &larkevent.EventReq{
		Header: headers,
		Body:   req.Body(),
	}

	return eventReq, nil
}
