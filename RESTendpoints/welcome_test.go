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
		m      WelcomeModel
		status int
	}
	tests := []struct {
		name string
		args args
	}{
		{"First time is ok", args{getRequest(t, "GET", "/", "127.0.0.1"), WelcomeModel{}, http.StatusOK}},
		{"Second time is ok", args{getRequest(t, "GET", "/", "127.0.0.1"), WelcomeModel{}, http.StatusOK}},
		{"Third time is forbidden", args{getRequest(t, "GET", "/", "127.0.0.1"), WelcomeModel{}, http.StatusForbidden}},
		{"Keep rejecting", args{getRequest(t, "GET", "/", "127.0.0.1"), WelcomeModel{}, http.StatusForbidden}},
		{"New IP is alllowed", args{getRequest(t, "GET", "/", "10.0.0.10"), WelcomeModel{}, http.StatusOK}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			tt.args.m.Welcome(rr, tt.args.r)
			assert.Equal(t, tt.args.status, rr.Code)
		})
	}
}

func getRequest(t *testing.T, method string, url string, ip string) *http.Request {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.RemoteAddr = ip + ":443"
	return req
}
