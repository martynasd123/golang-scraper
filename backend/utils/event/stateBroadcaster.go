package event

import (
	"github.com/barkimedes/go-deepcopy"
	"log"
)

type Subscriber[T any] struct {
	Data chan T
	Done chan struct{}
}

// StateBroadcaster is a type of broadcaster, where the latest update value is remembered. When a subscriber is registered,
// the latest value is sent. A deep copy of the data is always made to ensure no concurrent modification+read
type StateBroadcaster[T any] struct {
	subscribers      []Subscriber[T]
	addSubscriber    chan Subscriber[T]
	removeSubscriber chan Subscriber[T]
	publisher        chan T
	stateUpdates     chan T
	newestState      T
}

func CreateStateBroadcaster[T any]() *StateBroadcaster[T] {
	return &StateBroadcaster[T]{
		subscribers:      make([]Subscriber[T], 0),
		addSubscriber:    make(chan Subscriber[T]),
		removeSubscriber: make(chan Subscriber[T]),
		publisher:        make(chan T),
		stateUpdates:     make(chan T, 10),
		newestState:      *new(T),
	}
}

func (broadcaster *StateBroadcaster[T]) notifyListeners() {
	for _, subscriber := range broadcaster.subscribers {
		subscriber.Data <- broadcaster.newestState
	}
}

func (broadcaster *StateBroadcaster[T]) closeListeners() {
	for _, subscriber := range broadcaster.subscribers {
		close(subscriber.Data)
	}
}

func (broadcaster *StateBroadcaster[T]) closeListener(subscriber Subscriber[T]) {
	for i, sub := range broadcaster.subscribers {
		if sub == subscriber {
			close(subscriber.Data)
			broadcaster.subscribers = append(broadcaster.subscribers[:i], broadcaster.subscribers[i+1:]...)
		}
	}
}

// Start function starts the broadcaster. It is the responsibility of the caller to call End when function
// is not needed anymore
func (broadcaster *StateBroadcaster[T]) Start(data T) {
	broadcaster.newestState = data
	go func() {
		defer broadcaster.closeListeners()
		for {
			select {
			case subscriber := <-broadcaster.addSubscriber:
				subscriber.Data <- broadcaster.newestState
				broadcaster.subscribers = append(broadcaster.subscribers, subscriber)
			case subscriber := <-broadcaster.removeSubscriber:
				broadcaster.closeListener(subscriber)
			case update, ok := <-broadcaster.stateUpdates:
				if !ok {
					return
				}
				updateCopy, err := deepcopy.Anything(update)
				if err != nil {
					log.Fatalf("failed to copy state: %v", err)
				}
				broadcaster.newestState = updateCopy.(T)
				broadcaster.notifyListeners()
			}
		}
	}()
}

// Listen starts listening for updates. Caller must write to done channel when updates are no longer needed,
// unless the data channel has been closed. It is the responsibility of the caller
// to ensure that the done channel is closed after it is no longer needed.
// Returns:
//
//	data (<-chan T): The channel through which data is to be sent
//	done (chan<- struct{}): Channel through which to send signal when updates are no longer needed
func (broadcaster *StateBroadcaster[T]) Listen() (data <-chan T, done chan<- struct{}) {
	subscriber := Subscriber[T]{Data: make(chan T, 1), Done: make(chan struct{})}
	go func() {
		select {
		case _, ok := <-subscriber.Done:
			if ok {
				// Subscriber unsubscribed from updates
				broadcaster.removeSubscriber <- subscriber
			} else {
				// Broadcaster has been closed
			}
		}
	}()
	broadcaster.addSubscriber <- subscriber
	return subscriber.Data, subscriber.Done
}

// End the broadcaster. This unregisters all listeners and closes relevant channels.
func (broadcaster *StateBroadcaster[T]) End() {
	close(broadcaster.stateUpdates)
}

// Publish data to the channel
func (broadcaster *StateBroadcaster[T]) Publish(data T) {
	broadcaster.stateUpdates <- data
}
