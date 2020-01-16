package ctxpoller

import (
	"context"
	"errors"
	"time"
)

type Poller interface {
	IsActive() bool
	Stop()
	Start() error
}

type ctxPoller struct {
	ctx      context.Context
	stopFunc context.CancelFunc
	active   bool
	interval time.Duration
	action   func(context.Context)
}

func DefaultPoller(action func(context.Context)) Poller {
	return &ctxPoller{
		interval: 5 * time.Second,
		action:   action,
		ctx:      context.TODO(),
	}
}

func NewPoller(ctx context.Context, action func(context.Context), interval time.Duration) Poller {
	return &ctxPoller{
		interval: interval,
		action:   action,
		ctx:      ctx,
	}
}

func (p *ctxPoller) IsActive() bool {
	return p.active
}

func (p *ctxPoller) Stop() {
	if !p.active {
		return
	}
	p.stopFunc()
	p.active = false
}

func (p *ctxPoller) Start() error {
	if p.active {
		return errors.New("Poller already started")
	}
	if p.interval < (5 * time.Second) {
		return errors.New("Invalid interval, should be greater or equal to 5 seconds")
	}

	ctx, cancel := context.WithCancel(p.ctx)
	p.ctx = ctx
	p.stopFunc = cancel
	p.active = true

	go p.poll()

	return nil
}

func (p *ctxPoller) poll() {
	for {
		select {
		case <-p.ctx.Done():
			return
		case <-time.After(p.interval):
			p.action(p.ctx)
		}
	}
}
