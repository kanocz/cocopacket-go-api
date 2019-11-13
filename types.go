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

// GroupConfig == settings for group :)
type GroupConfig struct {
	IsPublic         bool            `json:"isPublic"`
	LossTreshold     float32         `json:"lossTreshold"`
	LatencyThreshold float32         `json:"latencyTreshold"`
	PushNotify       map[string]bool `json:"pushNotify"`
	IsAutoGroup      bool            `json:"isAutoGroup"`
	AGNetwork        string          `json:"agNetwork"`
	AGCount          int             `json:"agCount"`
	AGSlaves         []string        `json:"agSlaves"`
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
	Groups map[string]GroupConfig `json:"groups"`
}

// AvgChunk is part of GroupStatsData struct
type AvgChunk struct {
	Count   int     `json:"count"`   // total count of tests
	Loss    int     `json:"loss"`    // count of lost tests
	Latency float32 `json:"latency"` // total SUMM of latency
}

// GroupStatsData is result of GroupStats api call
type GroupStatsData struct {
	Ping map[string]map[int64]*AvgChunk
	HTTP map[string]map[int64]*AvgChunk
}
