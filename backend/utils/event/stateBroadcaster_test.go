package event

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type State struct {
	Value int
}

func TestStateBroadcaster_InitialStateSentToSubscriber(t *testing.T) {
	broadcaster := CreateStateBroadcaster[State]()
	initialState := State{Value: 1}
	broadcaster.Start(initialState)
	defer broadcaster.End()

	dataCh, doneCh := broadcaster.Listen()

	select {
	case state := <-dataCh:
		assert.Equal(t, initialState, state)
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for initial state")
	}
	close(doneCh)
}

func TestStateBroadcaster_UpdatesSentToSubscribers(t *testing.T) {
	broadcaster := CreateStateBroadcaster[State]()
	initialState := State{Value: 1}
	broadcaster.Start(initialState)
	defer broadcaster.End()

	dataCh1, doneCh1 := broadcaster.Listen()
	dataCh2, doneCh2 := broadcaster.Listen()

	// Wait for initial states
	<-dataCh1
	<-dataCh2

	updatedState := State{Value: 2}
	broadcaster.Publish(updatedState)

	select {
	case state := <-dataCh1:
		assert.Equal(t, updatedState, state)
	case <-time.After(1 * time.Second):
		t.Fatal("timeout waiting for updated state")
	}

	select {
	case state := <-dataCh2:
		assert.Equal(t, updatedState, state)
	case <-time.After(1 * time.Second):
		t.Fatal("timeout waiting for updated state")
	}

	close(doneCh1)
	close(doneCh2)
}

func TestStateBroadcaster_UnsubscribedChannelsAreClosed(t *testing.T) {
	broadcaster := CreateStateBroadcaster[State]()
	initialState := State{Value: 1}
	broadcaster.Start(initialState)
	defer broadcaster.End()

	dataCh, doneCh := broadcaster.Listen()

	// Wait for initial state
	<-dataCh

	// Indicate that we are done
	doneCh <- struct{}{}

	select {
	case _, ok := <-dataCh:
		assert.False(t, ok, "data channel should be closed")
	case <-time.After(1 * time.Second):
		t.Fatal("timeout waiting for data channel to close")
	}

	close(doneCh)
}

func TestStateBroadcaster_EndClosesAllChannels(t *testing.T) {
	broadcaster := CreateStateBroadcaster[State]()
	initialState := State{Value: 1}
	broadcaster.Start(initialState)

	dataCh1, doneCh1 := broadcaster.Listen()
	dataCh2, doneCh2 := broadcaster.Listen()

	// Wait for initial states
	<-dataCh1
	<-dataCh2

	broadcaster.End()

	select {
	case _, ok := <-dataCh1:
		assert.False(t, ok, "data channel 1 should be closed")
	case <-time.After(1 * time.Second):
		t.Fatal("timeout waiting for data channel 1 to close")
	}

	select {
	case _, ok := <-dataCh2:
		assert.False(t, ok, "data channel 2 should be closed")
	case <-time.After(1 * time.Second):
		t.Fatal("timeout waiting for data channel 2 to close")
	}

	close(doneCh1)
	close(doneCh2)
}
