// Package timeouts resolves the SDK's default connection read timeout for a
// service.
//
// Generated clients call GetServiceReadTimeout for themselves and apply the
// result with aws/transport/http's BuildableClient.WithReadTimeout.
package timeouts

import (
	"os"
	"sync"
	"time"
)

// DefaultReadTimeout is the SDK's default read timeout for a service with no
// entry in serviceInactivityTimeoutMillis.
const DefaultReadTimeout = 5 * time.Minute

// enableReadTimeoutEnvVar opts a public build in to the default read timeout.
// Internal builds are gated by enableReadTimeout2026 instead.
const enableReadTimeoutEnvVar = "AWS_ENABLE_DEFAULT_READ_TIMEOUT_2026"

var enableFromEnv = sync.OnceValue(func() bool {
	return os.Getenv(enableReadTimeoutEnvVar) == "true"
})

// GetServiceReadTimeout reports the SDK's default read timeout for a service,
// and whether one applies.
func GetServiceReadTimeout(serviceID string) (time.Duration, bool) {
	if enableReadTimeout2026 {
		if len(readTimeout2026Rollout) > 0 && !readTimeout2026Rollout[serviceID] {
			return 0, false
		}
	} else if !enableFromEnv() {
		return 0, false
	}

	ms, ok := serviceInactivityTimeoutMillis[serviceID]
	if !ok {
		return DefaultReadTimeout, true
	}
	if ms < 0 {
		return 0, false
	}

	return time.Duration(ms) * time.Millisecond, true
}
