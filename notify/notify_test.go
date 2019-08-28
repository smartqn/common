package notify

import (
	"reflect"
	"testing"

	"github.com/linnv/logx"
)

func TestNewNotify(t *testing.T) {
	a := []byte("abcd")
	b := make([]byte, len(a)+1)
	b[0] = '9'
	copy(b[1:], a)
	logx.Warnf("b: %s\n", b)
	bs := []byte(`1{"Mobile":"Mobile1","Name":"Name2"}`)
	n, err := NewNotify(bs)
	if err != nil {
		logx.Warnf("err: %+v\n", err)
	}
	logx.Debugf("n: %+v\n", n)
	n.Send()
	return

	type args struct {
		bs []byte
	}
	tests := []struct {
		name    string
		args    args
		wantN   Notifier
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotN, err := NewNotify(tt.args.bs)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewNotify() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotN, tt.wantN) {
				t.Errorf("NewNotify() = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}
