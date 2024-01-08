package event

type Event string

type Payload struct {
	Event Event
	Data  any
}

type Subscriber chan<- Payload
