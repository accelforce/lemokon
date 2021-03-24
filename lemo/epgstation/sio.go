package epgstation

import (
	"context"
	"fmt"
	"time"

	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
)

var (
	connected                 = make(chan struct{})
	Connected <-chan struct{} = connected

	statusUpdated                 = make(chan struct{})
	StatusUpdated <-chan struct{} = statusUpdated

	encodeUpdated                 = make(chan struct{})
	EncodeUpdated <-chan struct{} = encodeUpdated
)

func Streaming(ctx context.Context, host string, port int, secure bool) {
	disconnect := make(chan struct{})

	for {
		sio, err := gosocketio.Dial(gosocketio.GetUrl(host, port, secure), transport.GetDefaultWebsocketTransport())
		if err != nil {
			fmt.Printf("Cannot connect to EPGStation: %s\n", err)
			fmt.Println("Retry in 15 seconds...")
			<-time.After(15 * time.Second)
			continue
		}

		_ = sio.On(gosocketio.OnConnection, func(channel *gosocketio.Channel) {
			fmt.Printf("Connected to EPGStation: %s\n", channel.Id())
			connected <- struct{}{}
		})

		_ = sio.On(gosocketio.OnDisconnection, func(channel *gosocketio.Channel) {
			fmt.Printf("Disconnected from EPGStation: %s\n", channel.Id())
			disconnect <- struct{}{}
		})

		_ = sio.On(gosocketio.OnError, func(channel *gosocketio.Channel) {
			fmt.Printf("Error happened on Socket.IO\n")
		})

		_ = sio.On("updateStatus", func(channel *gosocketio.Channel) {
			statusUpdated <- struct{}{}
		})

		_ = sio.On("updateEncode", func(channel *gosocketio.Channel) {
			encodeUpdated <- struct{}{}
		})

		select {
		case <-disconnect:
			continue
		case <-ctx.Done():
			fmt.Println("Gracefully shutting down EPGStation Socket.IO connection...")
			go sio.Close()
			<-disconnect
			return
		}
	}
}
