package servercard

import "net/http"

// Attach builds a Server Card from opts and returns an http.Handler
// that serves it at WellKnownPath.
//
// Returns an error if opts fails validation (missing required fields or
// invalid name format).
//
// Typical usage with a standard http.ServeMux:
//
//	handler, err := servercard.Attach(opts)
//	if err != nil {
//		log.Fatal(err)
//	}
//	mux.Handle(servercard.WellKnownPath, handler)
func Attach(opts Options) (http.Handler, error) {
	card, err := Build(opts)
	if err != nil {
		return nil, err
	}
	return Handler(card), nil
}
