package gbase

// PanicHandler is called for recovering from panics spawned internally to the library (and thus
// not recoverable by the caller's goroutine). Defaults to nil, which means panics are not recovered.
var PanicHandler func(interface{})
