package proxy

import (
	"net/http"
	"net/url"
	"time"
)

const (
	defaultTargetUrl = "http://localhost:8081"
	defailServerPort = "8080"
)

var (
	targetURL *url.URL
	serverPort string
	backoffSchedule = []time.Duration{
		1* time.Second,
		3 * time.Second,
		10 * time.Second,
	}
	client = http.Client{Timeout: 30 * time.Second}
)

