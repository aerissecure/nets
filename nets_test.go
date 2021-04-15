package nets

import (
	"net"
	"reflect"
	"testing"
)

func TestIPInc(t *testing.T) {
	tests := []struct {
		name string
		ip   net.IP
		want net.IP
	}{
		{
			name: "1.1.1.1",
			ip:   net.IP{1, 1, 1, 1},
			want: net.IP{1, 1, 1, 2},
		},
		{
			name: "1.1.1.255",
			ip:   net.IP{1, 1, 1, 255},
			want: net.IP{1, 1, 2, 0},
		},
		{
			name: "255.255.255.255",
			ip:   net.IP{255, 255, 255, 255},
			want: net.IP{0, 0, 0, 0},
		},
		{
			name: "0000:0000:0000:0000:0000:0000:0000:0000",
			ip:   net.ParseIP("0000:0000:0000:0000:0000:0000:0000:0000"),
			want: net.ParseIP("0000:0000:0000:0000:0000:0000:0000:0001"),
		},
		{
			name: "0000:0000:0000:0000:0000:0000:0000:ffff",
			ip:   net.ParseIP("0000:0000:0000:0000:0000:0000:0000:ffff"),
			want: net.ParseIP("0000:0000:0000:0000:0000:0000:0001:0000"),
		},
		{
			name: "ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff",
			ip:   net.ParseIP("ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff"),
			want: net.ParseIP("0000:0000:0000:0000:0000:0000:0000:0000"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.ip
			IPInc(got)
			t.Log("got:", []byte(got))

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IPAdd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIPAdd(t *testing.T) {
	type args struct {
		ip     net.IP
		offset int
	}
	tests := []struct {
		name string
		args args
		want net.IP
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IPAdd(tt.args.ip, tt.args.offset); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IPAdd() = %v, want %v", got, tt.want)
			}
		})
	}
}

var resultIP net.IP

func BenchmarkInc(b *testing.B) {
	ip := net.IPv4(1, 1, 1, 1)
	for i := 1; i <= b.N; i++ {
		IPInc(ip)
	}
	b.Log("IP:", ip, "N:", b.N)
}

func BenchmarkAdd1(b *testing.B) {
	ip := net.IPv4(1, 1, 1, 1)
	for i := 1; i < b.N; i++ {
		ip = IPAdd(ip, 1)
	}
	b.Log("IP:", ip, "N:", b.N)
}

func BenchmarkAdd(b *testing.B) {
	ip := net.IPv4(1, 1, 1, 1)
	ip = IPAdd(ip, b.N)
	b.Log("IP:", ip, "N:", b.N)
}

func BenchmarkNetRange(b *testing.B) {
	_, ipnet, _ := net.ParseCIDR("1.1.1.0/16")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s, e := NetRange(*ipnet)
		resultIP = s
		resultIP = e
	}
}

func TestNetRange(t *testing.T) {
	newIPNet := func(s string) net.IPNet {
		_, ipnet, _ := net.ParseCIDR(s)
		return *ipnet
	}
	tests := []struct {
		name      string
		ipnet     net.IPNet
		wantFirst net.IP
		wantLast  net.IP
	}{
		{
			name:      "2002:0000:0000:1234:abcd:ffff:c0a8:0100/124",
			ipnet:     newIPNet("2002:0000:0000:1234:abcd:ffff:c0a8:0100/124"),
			wantFirst: net.ParseIP("2002:0000:0000:1234:abcd:ffff:c0a8:0100"),
			wantLast:  net.ParseIP("2002:0000:0000:1234:abcd:ffff:c0a8:010f"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := NetRange(tt.ipnet)
			if !reflect.DeepEqual(got, tt.wantFirst) {
				t.Errorf("NetRange() got = %v, want %v", got, tt.wantFirst)
			}
			if !reflect.DeepEqual(got1, tt.wantLast) {
				t.Errorf("NetRange() got1 = %v, want %v", got1, tt.wantLast)
			}
		})
	}
}
