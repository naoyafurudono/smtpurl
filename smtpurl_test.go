package smtpurl

import (
	"net/smtp"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		raw        string
		wantServer string
		wantAuth   smtp.Auth
		wantErr    bool
	}{
		{
			"smtp://username:password@smtp.example.com:25",
			"smtp.example.com:25",
			smtp.PlainAuth("", "username", "password", "smtp.example.com"),
			false,
		},
		{
			"smtp://username:password@smtp.example.com",
			"smtp.example.com:25",
			smtp.PlainAuth("", "username", "password", "smtp.example.com"),
			false,
		},
		{
			"smtp://username:password@smtp.example.com:587",
			"smtp.example.com:587",
			smtp.PlainAuth("", "username", "password", "smtp.example.com"),
			false,
		},
		{
			"smtp://username;AUTH=PLAIN:password@smtp.example.com",
			"smtp.example.com:25",
			smtp.PlainAuth("", "username", "password", "smtp.example.com"),
			false,
		},
		{
			"smtp://username;AUTH=CRAM-MD5:password@smtp.example.com",
			"smtp.example.com:25",
			smtp.CRAMMD5Auth("username", "password"),
			false,
		},
	}
	for _, tt := range tests {
		gotServer, gotAuth, err := Parse(tt.raw)
		if err != nil {
			if !tt.wantErr {
				t.Errorf("got err %v\n", err)
			}
			continue
		}
		if tt.wantErr {
			t.Error("want err\n")
		}
		if gotServer != tt.wantServer {
			t.Errorf("got %v\nwant %v", gotServer, tt.wantServer)
		}
		if !reflect.DeepEqual(gotAuth, tt.wantAuth) {
			t.Errorf("got %v\nwant %v", gotAuth, tt.wantAuth)
		}
	}
}
