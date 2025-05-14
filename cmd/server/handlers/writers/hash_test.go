package writers

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHashWriter(t *testing.T) {

	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{"test 1", "*writers.hashWriter", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var rw http.ResponseWriter
			got := NewHashWriter(rw, []byte("test"), "test")
			assert.Equalf(t, tt.want, reflect.ValueOf(got).Type().String(), "NewCompressReader()")
		})
	}
}

func Test_hashWriter_CreateSha256(t *testing.T) {
	type fields struct {
		w    http.ResponseWriter
		key  string
		body []byte
	}
	var w http.ResponseWriter
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"test1 ", fields{w: w, key: "test", body: []byte("test")}, "88cd2108b5347d973cf39cdf9053d7dd42704876d8c9a9bd8e2d168259d3ddf7"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hw := &hashWriter{
				w:    tt.fields.w,
				key:  tt.fields.key,
				body: tt.fields.body,
			}
			assert.Equal(t, tt.want, hw.CreateSha256())
		})
	}
}

func Test_hashWriter_Header(t *testing.T) {
	type fields struct {
		w    http.ResponseWriter
		key  string
		body []byte
	}
	f := fields{
		w:    httptest.NewRecorder(),
		key:  "test",
		body: []byte("test"),
	}

	tests := []struct {
		name   string
		fields fields
		want   http.Header
	}{
		{"tast 1", f, http.Header{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hw := &hashWriter{
				w:    tt.fields.w,
				key:  tt.fields.key,
				body: tt.fields.body,
			}

			assert.Equalf(t, tt.want, hw.Header(), "error in header")
		})
	}
}

func Test_hashWriter_Write(t *testing.T) {
	type fields struct {
		w    http.ResponseWriter
		key  string
		body []byte
	}
	w := httptest.NewRecorder()
	type args struct {
		p []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{"test 1", fields{w: w, key: "test", body: []byte("test")}, args{p: []byte("test")}, 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hw := &hashWriter{
				w:    tt.fields.w,
				key:  tt.fields.key,
				body: tt.fields.body,
			}
			got, err := hw.Write(tt.args.p)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_hashWriter_WriteHeader(t *testing.T) {
	type fields struct {
		w    http.ResponseWriter
		key  string
		body []byte
	}
	w := httptest.NewRecorder()
	type args struct {
		statusCode int
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantHeader string
	}{
		{"test 1", fields{w: w, key: "test", body: []byte("test")}, args{http.StatusOK}, "88cd2108b5347d973cf39cdf9053d7dd42704876d8c9a9bd8e2d168259d3ddf7"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hw := &hashWriter{
				w:    tt.fields.w,
				key:  tt.fields.key,
				body: tt.fields.body,
			}
			hw.WriteHeader(tt.args.statusCode)
			assert.Equal(t, tt.wantHeader, hw.Header().Get("HashSHA256"))
		})
	}
}
