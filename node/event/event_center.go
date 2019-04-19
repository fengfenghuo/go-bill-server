package event

import (
	"reflect"
	"sync"

	"github.com/op/go-logging"
)

var (
	log = logging.MustGetLogger("event")
)

type subscribeContainer map[int]*Subscription

type EventCenter struct {
	mu         sync.Mutex
	maxId      int
	subscribes map[reflect.Type]subscribeContainer
}

func NewEventCenter() *EventCenter {
	return &EventCenter{
		subscribes: map[reflect.Type]subscribeContainer{},
	}
}

func (ec *EventCenter) Subscribe(types []reflect.Type, c chan interface{}) *Subscription {
	ec.mu.Lock()
	defer ec.mu.Unlock()

	return newSubscription(ec, types, c)
}

func (ec *EventCenter) NotifyEvent(evt interface{}) {
	ec.mu.Lock()
	defer ec.mu.Unlock()

	if subscribes, ok := ec.subscribes[reflect.TypeOf(evt)]; ok {
		for _, subscribe := range subscribes {
			subscribe.c <- evt
		}
	}
}
