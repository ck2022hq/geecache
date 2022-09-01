package geecache

import (
	"reflect"
	"testing"
)

func TestByteView_String(t *testing.T) {
	type fields struct {
		b []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"case1",
			fields{b: []byte("aaaaa")},
			"aaaaa",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			byteView := ByteView{
				b: tt.fields.b,
			}
			if got := byteView.String(); got != tt.want {
				t.Errorf("ByteView.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestByteView_Len(t *testing.T) {
	type fields struct {
		b []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			"case1",
			fields{b: []byte("aaaaa")},
			5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			byteView := ByteView{
				b: tt.fields.b,
			}
			if got := byteView.Len(); got != tt.want {
				t.Errorf("ByteView.Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestByteView_ByteSlice(t *testing.T) {
	type fields struct {
		b []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			"case1",
			fields{b: []byte("aaaaa")},
			[]byte("aaaaa"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			byteView := ByteView{
				b: tt.fields.b,
			}
			if got := byteView.ByteSlice(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ByteView.ByteSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
