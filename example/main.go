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

	"github.com/chyroc/lark"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/lark-hertz"
)

func main() {
	// 1: init lark client
	larkCli := lark.New(
		lark.WithAppCredential("", ""),
		lark.WithEventCallbackVerify("", ""),
	)

	// 2: register lark msg callback login
	// handle text, file, image and send response
	larkCli.EventCallback.HandlerEventV2IMMessageReceiveV1(func(ctx context.Context, cli *lark.Lark, schema string, header *lark.EventHeaderV2, event *lark.EventV2IMMessageReceiveV1) (string, error) {
		content, err := lark.UnwrapMessageContent(event.Message.MessageType, event.Message.Content)
		if err != nil {
			return "", err
		}
		switch event.Message.MessageType {
		case lark.MsgTypeText:
			_, _, err = cli.Message.Reply(event.Message.MessageID).SendText(ctx, fmt.Sprintf("got text: %s", content.Text.Text))
		case lark.MsgTypeFile:
			_, _, err = cli.Message.Reply(event.Message.MessageID).SendText(ctx, fmt.Sprintf("got file: %s, key: %s", content.File.FileName, content.File.FileKey))
		case lark.MsgTypeImage:
			_, _, err = cli.Message.Reply(event.Message.MessageID).SendText(ctx, fmt.Sprintf("got image: %s", content.Image.ImageKey))
		}
		return "", err
	})

	// 3: init hertz server
	h := server.Default()
	h.POST("/api/lark_callback", lark_hertz.ListenCallback(larkCli, lark_hertz.WithIgnoreCheckSignature(true)))

	h.Spin()

	// 4: deploy server to cloud, and set lark callback url to `<host>/api/lark_callback`
}
