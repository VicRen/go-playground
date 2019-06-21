package addr

import (
	"net"
	"reflect"
	"testing"
)

func TestAddrToBytes(t *testing.T) {
	type args struct {
		addr string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "happy_case_v4_1",
			args: args{"127.0.0.1:80"},
			want: []byte{1, 127, 0, 0, 1, 0, 80},
		},
		{
			name: "happy_case_v4_2",
			args: args{"127.0.0.1:1080"},
			want: []byte{1, 127, 0, 0, 1, 0x04, 0x38},
		},
		{
			name: "happy_case_v6_1",
			args: args{"[1414:1414:1414:1414:1414:1414:1414:1414]:80"},
			want: []byte{0x04, 0x14, 0x14, 0x14, 0x14, 0x14, 0x14, 0x14, 0x14, 0x14, 0x14, 0x14, 0x14, 0x14, 0x14, 0x14, 0x14, 0x00, 0x50},
		},
		{
			name: "happy_case_v6_2",
			args: args{"[1414:1414:1414:1414:1414:1414:1414:1414]:1080"},
			want: []byte{0x04, 0x14, 0x14, 0x14, 0x14, 0x14, 0x14, 0x14, 0x14, 0x14, 0x14, 0x14, 0x14, 0x14, 0x14, 0x14, 0x14, 0x04, 0x38},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MarshalAddr(tt.args.addr)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddrToBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddrToBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnmarshalAddr(t *testing.T) {
	type args struct {
		addrBuf []byte
	}
	tests := []struct {
		name         string
		args         args
		wantAddrType uint8
		wantIp       net.IP
		wantPort     uint16
		wantErr      bool
	}{
		{
			name: "happy_case_v4_1",
			args: args{
				[]byte{1, 127, 0, 0, 1, 0, 80},
			},
			wantAddrType: 1,
			wantIp:       net.IP{127, 0, 0, 1},
			wantPort:     uint16(80),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAddrType, gotIp, gotPort, err := UnmarshalAddr(tt.args.addrBuf)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalAddr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotAddrType != tt.wantAddrType {
				t.Errorf("UnmarshalAddr() gotAddrType = %v, want %v", gotAddrType, tt.wantAddrType)
			}
			if !reflect.DeepEqual(gotIp, tt.wantIp) {
				t.Errorf("UnmarshalAddr() gotIp = %v, want %v", gotIp, tt.wantIp)
			}
			if gotPort != tt.wantPort {
				t.Errorf("UnmarshalAddr() gotPort = %v, want %v", gotPort, tt.wantPort)
			}
		})
	}
}
