package utils

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func NewReverseProxy(target string) http.Handler {
	url, err := url.Parse(target)
	if err != nil {
		log.Fatal("Error parsing URL:", err)
	}
	return httputil.NewSingleHostReverseProxy(url)
}
