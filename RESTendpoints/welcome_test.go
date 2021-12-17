package RESTendpoints

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWelcome(t *testing.T) {
	type args struct {
		r      *http.Request
		m      http.HandlerFunc
		status int
	}
	tests := []struct {
		name string
		args args
	}{
		{"First time is ok", args{get_request(t, "GET", "/", "127.0.0.1"), Welcome, http.StatusOK}},
		{"Second time is ok", args{get_request(t, "GET", "/","127.0.0.1"), Welcome, http.StatusOK}},
		{"Third time is forbidden", args{ get_request(t, "GET", "/","127.0.0.1"), Welcome, http.StatusForbidden}},
		{"Keep rejecting", args{get_request(t, "GET", "/","127.0.0.1"), Welcome, http.StatusForbidden}},
		{"New IP is alllowed", args{get_request(t, "GET", "/","10.0.0.10"), Welcome, http.StatusOK}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			handler := tt.args.m
			handler.ServeHTTP(rr, tt.args.r)
			assert.Equal(t, tt.args.status, rr.Code)
		})
	}
}

func get_request(t *testing.T, method string, url string, ip string) *http.Request {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.RemoteAddr = ip+":443"
	return req
}

