package concurrency

import (
	"net"
	"sync"
)

// Run executes a long-running task.
// longRunningTask is a function that didn't take in any cancellation method.
// That's not great but we can't do anything about it.
// One thing we can do is to return a done channel to notify when the task is done and run longRunningTask in a goroutine.
func Run(longRunningTask func()) {
	longRunningTask()
}

// ConsumeChannel needs some work to be done.
// When the channel closes, if main exists, it's possible that not all messages will be handled.
func ConsumeChannel[T any](ch <-chan T, handle func(t T)) {
	// Rule 1: if nobody has created a way to stop the operation
	// return a cancel() function yourself.
	// The easiest way to do this is to create a context with cancel.
	// Rule 2: Return a done channel to notify when the operation is done.
	// Rule 3: ALWAYS close the done channel, always call cancel().

	for t := range ch {
		handle(t)
	}
}

// HandleConcurrently takes a channel and a handler function.
// HandleConcurrently needs fixing
// If the channel closes, and somehow main exists, it's possible that not all goroutines will resume.
// We need a way to stop not just the function, and every goroutine.(separately)
func HandleConcurrently[T any](ch <-chan T, handle func(t T)) {
	for t := range ch {
		go func() {
			handle(t)
		}()
	}
}

// Serve is not easily testable. Can you tell why?
// Can we shutdown the goroutines?
func Serve(l net.Listener, h func(c net.Conn)) error {
	var wg sync.WaitGroup
	var conn net.Conn
	var err error
	for {
		conn, err = l.Accept()
		if err != nil {
			break
		}
		wg.Add(1)
		go func(c net.Conn) {
			defer wg.Done()
			h(c)
		}(conn)
	}
	wg.Wait()
	return err
}
