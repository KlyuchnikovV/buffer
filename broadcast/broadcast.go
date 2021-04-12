package broadcast

import (
	"context"

	"github.com/KlyuchnikovV/buffer/messages"
)

type Broadcast struct {
	context.Context
	cancel context.CancelFunc

	Receiver chan messages.RenderMessage

	listeners []func(messages.RenderMessage)
}

func New(ctx context.Context) *Broadcast {
	return &Broadcast{
		Context:   ctx,
		Receiver:  make(chan messages.RenderMessage, 1000),
		listeners: make([]func(messages.RenderMessage), 0),
	}
}

func (b *Broadcast) AddListener(l func(messages.RenderMessage)) {
	b.listeners = append(b.listeners, l)
}

func (b *Broadcast) Start() error {
	if b.cancel != nil {
		return nil
	}
	b.Context, b.cancel = context.WithCancel(b.Context)
	for {
		select {
		case msg, ok := <-b.Receiver:
			if !ok {
				panic("err")
			}
			for _, listener := range b.listeners {
				go listener(msg)
			}
		case <-b.Context.Done():
			return nil
		}
	}
}

func (b *Broadcast) Stop() {
	if b.cancel == nil {
		return
	}
	b.cancel()
	b.cancel = nil
}
