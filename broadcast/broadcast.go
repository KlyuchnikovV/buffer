package broadcast

import (
	"context"

	"github.com/KlyuchnikovV/edigode-cli/types"
)

type Broadcast struct {
	context.Context
	cancel context.CancelFunc

	Receiver chan types.Message

	listeners []func(types.Message)
}

func New(ctx context.Context) *Broadcast {
	return &Broadcast{
		Context:   ctx,
		Receiver:  make(chan types.Message, 1000),
		listeners: make([]func(types.Message), 0),
	}
}

func (b *Broadcast) AddListener(l func(types.Message)) {
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
