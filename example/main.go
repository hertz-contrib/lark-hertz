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

package main

import (
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app/server"
	lark_hertz "github.com/hertz-contrib/lark-hertz"
	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	larkcontact "github.com/larksuite/oapi-sdk-go/v3/service/contact/v3"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

func main() {
	// create event handler
	handler := dispatcher.NewEventDispatcher("v", "1212121212").OnP2MessageReceiveV1(func(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
		fmt.Println(larkcore.Prettify(event))
		fmt.Println(event.RequestId())
		return nil
	}).OnP2MessageReadV1(func(ctx context.Context, event *larkim.P2MessageReadV1) error {
		fmt.Println(larkcore.Prettify(event))
		fmt.Println(event.RequestId())
		return nil
	}).OnP2UserCreatedV3(func(ctx context.Context, event *larkcontact.P2UserCreatedV3) error {
		fmt.Println(larkcore.Prettify(event))
		fmt.Println(event.RequestId())
		return nil
	})

	// create card action handler.
	cardHandler := larkcard.NewCardActionHandler("v", "", func(ctx context.Context, cardAction *larkcard.CardAction) (interface{}, error) {
		fmt.Println(larkcore.Prettify(cardAction))

		// return card
		// return getCard(),nil

		// return custom resp
		// return getCustomResp(),nil

		// return nil
		return nil, nil
	})

	// register handler
	h := server.Default(server.WithHostPorts(":9999"))

	h.POST("/webhook/event", lark_hertz.NewEventHandlerFunc(handler))
	h.POST("/webhook/card", lark_hertz.NewCardActionHandlerFunc(cardHandler))

	// start server
	h.Spin()
}

//func getCard() *larkcard.MessageCard {
//	// config
//	config := larkcard.NewMessageCardConfig().
//		WideScreenMode(true).
//		EnableForward(true).
//		UpdateMulti(false).
//		Build()
//
//	// CardUrl
//	cardLink := larkcard.NewMessageCardURL().
//		PcUrl("https://open.feishu.com").
//		IoSUrl("https://open.feishu.com").
//		Url("https://open.feishu.com").
//		AndroidUrl("https://open.feishu.com").
//		Build()
//
//	// header
//	header := larkcard.NewMessageCardHeader().
//		Template("turquoise").
//		Title(larkcard.NewMessageCardPlainText().
//			Content("").
//			Build()).
//		Build()
//
//	// Elements
//	divElement := larkcard.NewMessageCardDiv().
//		Fields([]*larkcard.MessageCardField{larkcard.NewMessageCardField().
//			Text(larkcard.NewMessageCardLarkMd().
//				Content("**🕐 time：**\\n2021-02-23 20:17:51").
//				Build()).
//			IsShort(true).
//			Build()}).
//		Build()
//
//	content := "✅ " + "name" + "already deal"
//	processPersonElement := larkcard.NewMessageCardDiv().
//		Fields([]*larkcard.MessageCardField{larkcard.NewMessageCardField().
//			Text(larkcard.NewMessageCardLarkMd().
//				Content(content).
//				Build()).
//			IsShort(true).
//			Build()}).
//		Build()
//
//	// card message body
//	messageCard := larkcard.NewMessageCard().
//		Config(config).
//		Header(header).
//		Elements([]larkcard.MessageCardElement{divElement, processPersonElement}).
//		CardLink(cardLink).
//		Build()
//
//	return messageCard
//}
//
//func getCustomResp() interface{} {
//	body := make(map[string]interface{})
//	body["content"] = "hello"
//
//	i18n := make(map[string]string)
//	i18n["zh_cn"] = "你好"
//	i18n["en_us"] = "hello"
//	i18n["ja_jp"] = "こんにちは"
//	body["i18n"] = i18n
//
//	resp := larkcard.CustomResp{
//		StatusCode: 400,
//		Body:       body,
//	}
//	return &resp
//}
