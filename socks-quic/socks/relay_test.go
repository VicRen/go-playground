package socks

import (
	"errors"
	"net"
	"reflect"
	"testing"
)

func TestNewRelay(t *testing.T) {
	type args struct {
		srcAddr *net.UDPAddr
		data    []byte
	}
	tt := []struct {
		name    string
		args    args
		want    *Relay
		wantErr error
	}{
		{
			"too_short_1",
			args{
				srcAddr: &net.UDPAddr{IP: net.IP{}, Port: 1000},
			},
			nil,
			errors.New("bad request: too short"),
		},
		{
			"too_short_2",
			args{
				&net.UDPAddr{IP: net.IP{}, Port: 1000},
				[]byte{0x00},
			},
			nil,
			errors.New("bad request: too short"),
		},
		{
			"bad_request_eof",
			args{
				&net.UDPAddr{IP: net.IP{}, Port: 1000},
				[]byte{
					0x00, 0x00, // RSV
					0x00, // FRAG
					0x01,
				},
			},
			nil,
			errors.New("bad request: EOF"),
		},
		{
			"bad_source",
			args{
				nil,
				[]byte{
					0x00, 0x00, // RSV
					0x00,                   // FRAG
					0x01,                   // ATYP
					0x00, 0x00, 0x00, 0x00, // ADDR
					0x00, 0x00, // PORT
					0x00, // DATA
				},
			},
			nil,
			errors.New("bad request: bad source"),
		},
		{
			"happy_case",
			args{
				&net.UDPAddr{IP: net.IP{}, Port: 1000},
				[]byte{
					0x00, 0x00, // RSV
					0x00,                   // FRAG
					0x01,                   // ATYP
					0x00, 0x00, 0x00, 0x00, // ADDR
					0x00, 0x00, // PORT
					0x00, // DATA
				},
			},
			&Relay{
				0x00,
				&net.UDPAddr{IP: net.IP{}, Port: 1000},
				&AddrSpec{IP: net.IP{0, 0, 0, 0}, Port: 0},
				[]byte{0x00},
			},
			nil,
		},
		{
			"happy_case_2",
			args{
				&net.UDPAddr{IP: net.IP{127, 0, 0, 1}, Port: 1000},
				[]byte{
					0x00, 0x00, // RSV
					0x00,                   // FRAG
					0x01,                   // ATYP
					0x01, 0x01, 0x01, 0x01, // ADDR
					0x00, 0x10, // PORT
					0x00, // DATA
				},
			},
			&Relay{
				0x00,
				&net.UDPAddr{IP: net.IP{127, 0, 0, 1}, Port: 1000},
				&AddrSpec{IP: net.IP{1, 1, 1, 1}, Port: 16},
				[]byte{0x00},
			},
			nil,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got, err := NewRelay(tc.args.srcAddr, tc.args.data)
			if !reflect.DeepEqual(err, tc.wantErr) {
				t.Errorf("NewRelay() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("NewRelay() = %v, want %v", got, tc.want)
			}
		})
	}
}
