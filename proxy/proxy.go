package proxy

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
	"github.com/Llifintsefv/-GoPrometheusProxy/server"
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

func init() {
	targetURLString := os.Getenv("PROXY_TARGET_SERVER")
	if targetURLString == "" {
		targetURLString = defailServerPort
	}
	var err error
	targetURL,err = url.Parse(targetURLString)
	if err != nil {
		log.Fatalf("Invalid proxy target url: %v",err)
	}
	serverPort := os.Getenv("PROXY_SERVER_PORT")
	if serverPort == "" {
		serverPort = defailServerPort
	}
}


func RunProxyServer(ctx context.Context) {
	http.HandleFunc("/",handlerRequest)
	srv := server.NewServer(http.HandlerFunc(handlerRequest),serverPort)
	if err := srv.Run(ctx);err != nil {
		log.Fatalf("failed to run proxy server %v",err)
	}
}

func handlerRequest(w http.ResponseWriter, r *http.Request){
	startReq := time.Now()
	targetRequestURL := targetURL.ResolveReference(r.URL)
	log.Println(r.Method,targetRequestURL.String())
	proxyReq,err := http.NewRequestWithContext(r.Context(),r.Method,targetRequestURL.String(),r.Body)
	if err != nil {
		http.Error(w,"error create proxy request",http.StatusInternalServerError)
		return
	}
	proxyReq.Header = r.Header.Clone()
	
}


func getRequestWithRetry(req *http.Request) (*http.Response,error) {
	var resp *http.Response
	var err error
	for _,backoff := range backoffSchedule {
		resp,err = client.Do(req)
		if err == nil {
			return resp,nil
		}
		log.Printf("request error: %v\n",err)
		log.Printf("Retruing in %v \n",backoff)
		time.Sleep(backoff)
	}
	return nil,err
}
