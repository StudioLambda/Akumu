package akumu

type OnErrorKey struct{}

type OnErrorHook func(ServerError)
