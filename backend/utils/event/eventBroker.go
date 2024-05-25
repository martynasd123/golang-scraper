package event

import (
	"fmt"
	"sync"
)

// StateBroker Holds a registry of state broadcasters (see stateBroadcaster.go) and associates them
// with corresponding topics. A topic can be any datatype, as long as it is comparable
type StateBroker[Topic comparable, Data any] struct {
	broadcasters map[Topic]*StateBroadcaster[Data]
	mu           sync.Mutex
}

func CreateStateBroker[Topic comparable, Data any]() *StateBroker[Topic, Data] {
	return &StateBroker[Topic, Data]{
		broadcasters: make(map[Topic]*StateBroadcaster[Data]),
		mu:           sync.Mutex{},
	}
}

func (broker *StateBroker[Topic, Data]) AddStateBroadcaster(topic Topic) (*StateBroadcaster[Data], error) {
	broker.mu.Lock()
	defer broker.mu.Unlock()
	if _, exists := broker.broadcasters[topic]; exists {
		return nil, fmt.Errorf("state broadcaster already exists for given topic: %v", topic)
	} else {
		broadcaster := CreateStateBroadcaster[Data]()
		broker.broadcasters[topic] = broadcaster
		return broadcaster, nil
	}
}

func (broker *StateBroker[Topic, Data]) GetStateBroadcaster(topic Topic) (*StateBroadcaster[Data], error) {
	broker.mu.Lock()
	defer broker.mu.Unlock()
	if value, exists := broker.broadcasters[topic]; exists {
		return value, nil
	} else {
		return nil, fmt.Errorf("state broadcaster does not exist for topic: %v", topic)
	}
}

func (broker *StateBroker[Topic, Data]) DeleteStateBroadcaster(topic Topic) error {
	broker.mu.Lock()
	defer broker.mu.Unlock()
	if _, exists := broker.broadcasters[topic]; exists {
		delete(broker.broadcasters, topic)
	} else {
		return fmt.Errorf("state broadcaster does not exist for topic: %v", topic)
	}
	return nil
}
