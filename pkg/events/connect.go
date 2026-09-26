package events

import "strings"

// Connect opens the event backbone this platform was generated for (nats).
// `memory://` (or an empty URL) selects the in-process bus, for tests and local runs.
func Connect(url, name string) (Bus, error) {
	if url == "" || strings.HasPrefix(url, "memory://") {
		return NewMemory(), nil
	}
	return ConnectNATS(url, name)
}
