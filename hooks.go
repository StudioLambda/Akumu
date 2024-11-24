package akumu

// OnErrorKey is used in the [http.Request]'s context
// to store the hook or callback that will run whenever
// a server error is found.
type OnErrorKey struct{}

// OnErrorHook stores the current handler that will take
// care of handling a server error in case it's found.
//
// The `error` is always a [ErrServer] and it's joined with
// either of those:
//   - [ErrServerWriter]
//   - [ErrServerBody]
//   - [ErrServerStream]
//   - [ErrServerDefault]
//
// To "extend" the hook, you can make use of composition
// to also call the previous function in the new one, creating
// a chain of handlers. Keep in mind you will need to check
// for `nil` if that was the case.
type OnErrorHook func(error)
