package servercard

import "net/http"

// WellKnownPath is the standard endpoint for Server Card discovery.
const WellKnownPath = "/.well-known/mcp-server-card"

// Handler returns an http.Handler that serves the Server Card JSON
// at any path (typically mounted at WellKnownPath).
//
// The handler sets Content-Type, CORS, and caching headers per SEP-2127:
//   - Content-Type: application/json
//   - Access-Control-Allow-Origin: *
//   - Access-Control-Allow-Methods: GET
//   - Access-Control-Allow-Headers: Content-Type
//   - Cache-Control: public, max-age=3600
func Handler(card *ServerCard) http.Handler {
	// Pre-serialize the card once; it's static.
	data, err := card.JSON()
	if err != nil {
		// Programming error: the card should always be serializable.
		panic("servercard: failed to marshal card: " + err.Error())
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			setCORSHeaders(w)
			w.WriteHeader(http.StatusNoContent)
			return
		}
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		setCORSHeaders(w)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "public, max-age=3600")
		_, _ = w.Write(data)
	})
}

func setCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}
