package reload

import "context"

// Notifier knows how to trigger a reload process.
type Notifier interface {
	Notify(ctx context.Context) (string, error)
}

// NotifierFunc is a helper to create notifiers from functions.
type NotifierFunc func(ctx context.Context) (string, error)

// Notify satisfies Notifier interface.
func (n NotifierFunc) Notify(ctx context.Context) (string, error) { return n(ctx) }

// NotifierChan is a helper to create notifiers from channels.
//
// Note: Closing the channel is not safe, as the channel will be reused and read
// from it multiple times for each notification.
type NotifierChan <-chan string

// Notify satisfies Notifier interface.
func (n NotifierChan) Notify(ctx context.Context) (string, error) { return <-n, nil }
