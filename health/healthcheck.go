package health

import (
	"net"
	"net/url"
	"sync"
	"time"
)

// NewHealthCheck is the ProxyHealth constructor
func NewHealthCheck(origin *url.URL) *HealthCheck {
	h := &HealthCheck{
		origin:      origin,
		check:       defaultHealthCheckFunc,
		period:      defaultHealthCheckPeriod,
		cancel:      make(chan struct{}),
		isAvailable: defaultHealthCheckFunc(origin),
	}
	h.run()

	return h
}

// HealthCheck is looking after the proxy origin availability using either a set by
// HealthCheck.SetHealthCheck check function or the defaultHealthCheck func.
type HealthCheck struct {
	origin *url.URL

	mu          sync.Mutex
	check       func(addr *url.URL) bool
	period      time.Duration
	cancel      chan struct{}
	isAvailable bool
}

// IsAvailable returns whether the proxy origin was successfully connected at the last check time.
func (h *HealthCheck) IsAvailable() bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.isAvailable
}

// SetCheckFunc sets the passed check func as the algorithm of checking the origin availability and
// calls for it with interval defined with the passed period variable. The SetCheckFunc provides a
// concurrency save way of setting and replacing the current health check algorithm, so the Stop function
// shouldn't be called before the SetCheckFunc call.
func (h *HealthCheck) SetCheckFunc(check func(addr *url.URL) bool, period time.Duration) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.stop()
	h.check = check
	h.period = period
	h.cancel = make(chan struct{})
	h.isAvailable = h.check(h.origin)
	h.run()
}

// Stop gracefully stops the instance execution. Should be called when the instance work is no more needed.
func (h *HealthCheck) Stop() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.stop()
}

// run runs the check func in a new goroutine.
func (h *HealthCheck) run() {
	checkHealth := func() {
		h.mu.Lock()
		defer h.mu.Unlock()
		isAvailable := h.check(h.origin)
		h.isAvailable = isAvailable
	}

	go func() {
		t := time.NewTicker(h.period)
		for {
			select {
			case <-t.C:
				checkHealth()
			case <-h.cancel:
				t.Stop()
				return
			}
		}
	}()
}

// stop stops the currently rinning check func.
func (h *HealthCheck) stop() {
	if h.cancel != nil {
		h.cancel <- struct{}{}
		close(h.cancel)
		h.cancel = nil
	}
}

// defaultHealthCheckFunc is the default most simple check function
var defaultHealthCheckFunc = func(addr *url.URL) bool {
	conn, err := net.DialTimeout("tcp", addr.Host, defaultHealthCheckTimeout)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}

var (
	defaultHealthCheckTimeout = 10 * time.Second
	defaultHealthCheckPeriod  = 10 * time.Second
)
