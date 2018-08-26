package rpc

import (
	"context"
	"sync"
)

type ID string

type Subscription struct {
	ID        ID
	namespace string
	err       chan error // closed on unsubscribe
}

func (s *Subscription) Err() <-chan error {
	return s.err
}

type Notifier struct {
	subMu    sync.RWMutex // guards active and inactive maps
	active   map[ID]*Subscription
	inactive map[ID]*Subscription
	receiver func(id ID, data interface{}) error
}

type notifierKey struct{}

func NewNotifier(receiver func(id ID, data interface{}) error) *Notifier {
	return &Notifier{
		active:   make(map[ID]*Subscription),
		inactive: make(map[ID]*Subscription),
		receiver: receiver,
	}
}

func AttachNotifier(ctx context.Context, notifier *Notifier) context.Context {
	return context.WithValue(ctx, notifierKey{}, notifier)
}

func NotifierFromContext(ctx context.Context) *Notifier {
	n, ok := ctx.Value(notifierKey{}).(*Notifier)
	if ok {
		return n
	} else {
		return nil
	}
}

func (n *Notifier) CreateSubscription() *Subscription {
	s := &Subscription{ID: NewID(), err: make(chan error)}
	n.subMu.Lock()
	n.inactive[s.ID] = s
	n.subMu.Unlock()
	return s
}

func (n *Notifier) Notify(id ID, data interface{}) error {
	n.subMu.RLock()
	defer n.subMu.RUnlock()

	// log.Debugf("notify %s, %v", id, data)

	_, active := n.active[id]
	if active {
		if err := n.receiver(id, data); err != nil {
			return err
		}
	}
	return nil
}

func (n *Notifier) Unsubscribe(id ID) error {
	n.subMu.Lock()
	defer n.subMu.Unlock()
	if s, found := n.active[id]; found {
		close(s.err)
		delete(n.active, id)
		return nil
	}
	return &ErrSubscriptionNotFound{}
}

func (n *Notifier) Activate(id ID) {
	n.subMu.Lock()
	defer n.subMu.Unlock()

	if sub, found := n.inactive[id]; found {
		n.active[id] = sub
		delete(n.inactive, id)
	}
}
