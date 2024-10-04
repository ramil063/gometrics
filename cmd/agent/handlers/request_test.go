package handlers

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type RequestMock struct {
	Request request
}

type ClientMock struct {
	Client client
}

func (c ClientMock) SendPostRequest(url string) error {
	_, err := c.NewRequest("POST", url)
	return err
}

func (c ClientMock) NewRequest(method string, url string) (*http.Request, error) {
	return httptest.NewRequest(method, url, nil), nil
}

func TestNewRequest(t *testing.T) {
	tests := []struct {
		name string
		want RequestInterface
	}{
		{"create new request", RequestMock{}.Request},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewRequest())
		})
	}
}

func Test_request_SendMetrics(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"send request"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := request{}
			assert.NoError(t, r.SendMetrics(ClientMock{}, 5))
		})
	}
}

func TestNewClient(t *testing.T) {
	tests := []struct {
		name string
		want ClientInterface
	}{
		{"create new request", ClientMock{}.Client},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewClient(), "NewClient()")
		})
	}
}

func Test_client_SendPostRequest(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{"send post", args{url: MainUrl + "/update/gauge/metric1/100"}, assert.NoError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ClientMock{}
			tt.wantErr(t, c.SendPostRequest(tt.args.url), fmt.Sprintf("SendPostRequest(%v)", tt.args.url))
		})
	}
}
