package event

import (
	"reflect"
)

type Subscription struct {
	ec    *EventCenter
	id    int
	types []reflect.Type
	c     chan interface{}
}

func newSubscription(ec *EventCenter, types []reflect.Type, c chan interface{}) *Subscription {
	sub := Subscription{
		ec:    ec,
		id:    ec.maxId + 1,
		types: types,
		c:     c,
	}

	for _, t := range types {
		subscribes, ok := ec.subscribes[t]
		if !ok {
			subscribes = make(subscribeContainer)
			ec.subscribes[t] = subscribes
		}

		subscribes[sub.id] = &sub
	}

	ec.maxId++
	return &sub
}

func (sub *Subscription) Unsubscribe() {
	if sub.ec == nil {
		return
	}

	sub.ec.mu.Lock()
	defer sub.ec.mu.Unlock()

	for _, t := range sub.types {
		subscribes, ok := sub.ec.subscribes[t]
		if !ok {
			continue
		}

		delete(subscribes, sub.id)
	}

	sub.types = nil
}
