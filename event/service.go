package event

type Service interface {
	Emit(event Event, payload any) error
	Subscribe(subscriber Subscriber, events ...Event)
	Unsubscribe(subscriber Subscriber, events ...Event)
	IsSubscribed(subscriber Subscriber, events ...Event) bool
}
