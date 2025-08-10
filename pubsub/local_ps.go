package pubsub

import (
	"context"
	"github.com/coderconquerer/social-todo/common"
	"log"
	"sync"
)

// A local pubsub run locally (in-mem)
// It has a queue as a buffer channel at it's core and many groups of subscribers
// To send messages with a specific topic for many subscribers in a group
type localPubSub struct {
	name         string
	messageQueue chan *Message
	mapChannel   map[Topic][]chan *Message
	locker       *sync.RWMutex
}

const defaultMessageQueueSize = 1000

func NewLocalPubsub(name string) *localPubSub {
	pb := &localPubSub{
		name:         name,
		messageQueue: make(chan *Message, defaultMessageQueueSize),
		mapChannel:   make(map[Topic][]chan *Message),
		locker:       new(sync.RWMutex),
	}

	return pb
}

// Publish sends a new message to the specified topic.
// The message is pushed into the pubsub's internal message queue.
// It runs the actual sending in a separate goroutine to avoid blocking the caller.
func (ps *localPubSub) Publish(ctx context.Context, channel Topic, data *Message) error {
	// Set the message's channel (channel) field so receivers know which channel it belongs to
	data.SetChannel(channel)

	// Run the publishing logic asynchronously in its own goroutine
	go func() {
		// Ensure that if this goroutine panics, it will be recovered and not crash the program
		defer common.Recovery()

		// Send the message into the pubsub's message queue
		ps.messageQueue <- data

		// Log the publishing event for debugging/monitoring
		log.Printf("New message published on channel: %s - with data: %s\n", channel, data.String())
	}()

	// No error handling here — always return nil
	return nil
}

// Subscribe registers a new subscriber to a given topic.
// It returns:
//   - a receive-only channel (<-chan *Message) for delivering messages
//   - an unsubscribe function to remove the subscriber and close the channel
func (ps *localPubSub) Subscribe(ctx context.Context, channel Topic) (ch <-chan *Message, unsubscribe func()) {
	// Create a new message channel for this subscriber
	c := make(chan *Message)

	ps.locker.Lock()
	if val, ok := ps.mapChannel[channel]; ok {
		// If channel already has subscribers, append the new channel
		val = append(ps.mapChannel[channel], c)
		ps.mapChannel[channel] = val
	} else {
		// If channel has no subscribers yet, initialize the slice with this channel
		ps.mapChannel[channel] = []chan *Message{c}
	}
	ps.locker.Unlock()

	// Return the channel and an unsubscribe function
	return c, func() {
		log.Println("...Unsubscribe")

		// Remove this subscriber from the channel
		if chans, ok := ps.mapChannel[channel]; ok {
			for i := range chans {
				if chans[i] == c {
					// Remove element at index i from chans
					chans = append(chans[:i], chans[i+1:]...)

					// Update the map with the new list of channels
					ps.locker.Lock()
					ps.mapChannel[channel] = chans
					ps.locker.Unlock()

					// Close the subscriber's channel
					close(c)
					break
				}
			}
		}
	}
}

// run starts a goroutine that listens for messages from the message queue
// and delivers them to all subscribers of the corresponding topic.
func (ps *localPubSub) run() error {
	go func() {
		// Recover from panics in this goroutine so they don't crash the process.
		defer common.Recovery()

		for {
			// Wait until a message is received from the queue.
			mess := <-ps.messageQueue
			log.Println("Message dequeued:", mess.String())

			ps.locker.RLock()
			// Look up subscribers for the message's channel/topic.
			if subs, ok := ps.mapChannel[mess.Channel()]; ok {
				// For each subscriber channel...
				for i := range subs {
					// Deliver the message asynchronously to avoid blocking.
					go func(c chan *Message) {
						// Recover from panics in the subscriber delivery.
						defer common.Recovery()
						c <- mess // Send the message to the subscriber's channel.
					}(subs[i])
				}
			}
			ps.locker.RUnlock()
			// If there are no subscribers, message is just dropped (commented out code suggests possible requeueing).
		}
	}()
	return nil
}

func (l *localPubSub) Name() string {
	return l.name
}

func (l *localPubSub) GetPrefix() string {
	return l.name
}

func (l *localPubSub) Get() interface{} {
	return l
}

func (l *localPubSub) InitFlags() {
}

func (l *localPubSub) Configure() error {
	return nil
}

func (l *localPubSub) Run() error {
	return l.run()
}

func (l *localPubSub) Stop() <-chan bool {
	c := make(chan bool)
	go func() {
		c <- true
	}()
	return c
}
