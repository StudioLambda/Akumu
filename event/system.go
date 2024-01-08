package event

import (
	"context"
	"io"
	"log/slog"
	"sync"
	"time"
)

type Config struct {
	timeout time.Duration
	logger  *slog.Logger
}

type System struct {
	Config
	subscribers map[Event][]Subscriber
	mutex       sync.RWMutex
}

type Builder func(*Config)

func WithLogger(logger *slog.Logger) Builder {
	return func(config *Config) {
		config.logger = logger
	}
}

func WithTimeout(timeout time.Duration) Builder {
	return func(config *Config) {
		config.timeout = timeout
	}
}

func New(builders ...Builder) *System {
	config := Config{
		timeout: 5 * time.Second,
		logger: slog.New(
			slog.NewTextHandler(io.Discard, nil),
		),
	}

	for _, builder := range builders {
		builder(&config)
	}

	return &System{
		Config:      config,
		subscribers: make(map[Event][]Subscriber),
		mutex:       sync.RWMutex{},
	}
}

func (system *System) Emit(event Event, payload any) error {
	system.mutex.RLock()
	defer system.mutex.RUnlock()

	ctx, cancel := context.WithTimeout(context.Background(), system.timeout)
	defer cancel()

	subscribers := system.subscribers[event]

	for _, subscriber := range subscribers {
		select {
		case subscriber <- Payload{Event: event, Data: payload}:
			continue
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	system.logger.Debug("event emitted", "event", event, "subscribers", len(subscribers))

	return nil
}

func (system *System) Subscribe(subscriber Subscriber, events ...Event) {
	system.mutex.Lock()
	defer system.mutex.Unlock()

	for _, event := range events {
		system.subscribers[event] = append(system.subscribers[event], subscriber)
	}

	system.logger.Debug("subscribed", "events", events)
}

func (system *System) Unsubscribe(subscriber Subscriber, events ...Event) {
	system.mutex.Lock()
	defer system.mutex.Unlock()

	for _, event := range events {
		for index, sub := range system.subscribers[event] {
			if sub == subscriber {
				system.subscribers[event] = append(
					system.subscribers[event][:index],
					system.subscribers[event][index+1:]...,
				)
			}
		}
	}

	system.logger.Debug("unsubscribed", "events", events)
}

func (system *System) IsSubscribed(subscriber Subscriber, events ...Event) bool {
	system.mutex.RLock()
	defer system.mutex.RUnlock()

	for _, event := range events {
		for _, sub := range system.subscribers[event] {
			if sub == subscriber {
				return true
			}
		}
	}

	return false
}
