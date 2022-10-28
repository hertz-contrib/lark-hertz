# lark-hertz (This is a community driven project)

an [oapi-sdk-go](https://github.com/larksuite/oapi-sdk-go) extension package that integrates the hertz web framework

## Installation

```bash
go get github.com/hertz-contrib/lark-hertz
```

## Usage

```go
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
```

## License

This project is under Apache License. See the [LICENSE](./LICENSE-APACHE) file for the full license text.