package api

import "time"

type result struct {
	Result string `json:"result"`
	Error  string `json:"error,omitempty"`
}

// TestDesc describes one ping/http test
type TestDesc struct {
	Groups      []string  `json:"cat"`
	Description string    `json:"desc"`
	Favourite   bool      `json:"fav"`
	Slaves      []string  `json:"slaves"`
	AutoAdded   time.Time `json:"auto-added,omitempty"`
	AS          int64     `json:"as"`
}

// ConfigInfo type represent information about current configuration of tests
type ConfigInfo struct {
	Counter int64 `json:"counter"`
	Ping    struct {
		IPs       map[string]TestDesc `json:"ips"`
		Timeout   float32             `json:"timeout"`
		Interval  float32             `json:"interval"`
		Slowdown  int                 `json:"slowdown"`
		SlowEvery int                 `json:"slowEvery"`
		LastIP    string              `json:"lastIP"`
	} `json:"ping"`
	HTTP struct {
		URLs     map[string]TestDesc `json:"urls"`
		Timeout  float32             `json:"timeout"`
		Interval float32             `json:"interval"`
	}
}
