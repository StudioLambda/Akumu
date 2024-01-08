package event_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/studiolambda/akumu/event"
)

func TestCanCreateSystem(t *testing.T) {
	system := event.New()

	require.NotNil(t, system)
}

func TestCanSubscribeToEvents(t *testing.T) {
	t.Run("subscribe to one event", func(t *testing.T) {
		system := event.New()
		subscriber := make(event.Subscriber)
		e := event.Event("test.event")

		system.Subscribe(subscriber, e)

		require.True(t, system.IsSubscribed(subscriber, e))
	})

	t.Run("subscribe to multiple events", func(t *testing.T) {
		system := event.New()
		subscriber := make(event.Subscriber)
		e1 := event.Event("test.event1")
		e2 := event.Event("test.event2")

		system.Subscribe(subscriber, e1, e2)

		require.True(t, system.IsSubscribed(subscriber, e1))
		require.True(t, system.IsSubscribed(subscriber, e2))
	})

	t.Run("subscribe to same event multiple times", func(t *testing.T) {
		system := event.New()
		subscriber := make(event.Subscriber)
		e := event.Event("test.event")

		system.Subscribe(subscriber, e)
		system.Subscribe(subscriber, e)

		require.True(t, system.IsSubscribed(subscriber, e))
	})

	t.Run("subscribe to same event multiple times with different subscribers", func(t *testing.T) {
		system := event.New()
		subscriber1 := make(event.Subscriber)
		subscriber2 := make(event.Subscriber)
		e := event.Event("test.event")

		system.Subscribe(subscriber1, e)
		system.Subscribe(subscriber2, e)

		require.True(t, system.IsSubscribed(subscriber1, e))
		require.True(t, system.IsSubscribed(subscriber2, e))
	})

	t.Run("dont subscribe", func(t *testing.T) {
		system := event.New()
		subscriber := make(event.Subscriber)
		e := event.Event("test.event")

		require.False(t, system.IsSubscribed(subscriber, e))
	})
}

func TestCanUnsubscribeFromEvents(t *testing.T) {
	t.Run("unsubscribe from one event", func(t *testing.T) {
		system := event.New()
		subscriber := make(event.Subscriber)
		e := event.Event("test.event")

		system.Subscribe(subscriber, e)
		system.Unsubscribe(subscriber, e)

		require.False(t, system.IsSubscribed(subscriber, e))
	})

	t.Run("unsubscribe from multiple events", func(t *testing.T) {
		system := event.New()
		subscriber := make(event.Subscriber)
		e1 := event.Event("test.event1")
		e2 := event.Event("test.event2")

		system.Subscribe(subscriber, e1, e2)
		system.Unsubscribe(subscriber, e1, e2)

		require.False(t, system.IsSubscribed(subscriber, e1))
		require.False(t, system.IsSubscribed(subscriber, e2))
	})

	t.Run("unsubscribe from same event multiple times", func(t *testing.T) {
		system := event.New()
		subscriber := make(event.Subscriber)
		e := event.Event("test.event")

		system.Subscribe(subscriber, e)
		system.Unsubscribe(subscriber, e)
		system.Unsubscribe(subscriber, e)

		require.False(t, system.IsSubscribed(subscriber, e))
	})

	t.Run("unsubscribe from same event multiple times with different subscribers", func(t *testing.T) {
		system := event.New()
		subscriber1 := make(event.Subscriber)
		subscriber2 := make(event.Subscriber)
		e := event.Event("test.event")

		system.Subscribe(subscriber1, e)
		system.Subscribe(subscriber2, e)
		system.Unsubscribe(subscriber1, e)
		system.Unsubscribe(subscriber2, e)

		require.False(t, system.IsSubscribed(subscriber1, e))
		require.False(t, system.IsSubscribed(subscriber2, e))
	})

	t.Run("dont unsubscribe", func(t *testing.T) {
		system := event.New()
		subscriber := make(event.Subscriber)
		e := event.Event("test.event")

		system.Unsubscribe(subscriber, e)

		require.False(t, system.IsSubscribed(subscriber, e))
	})
}
