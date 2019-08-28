package cryptometa

import (
	"testing"

	"smartqn/algorithm/rsa"

	"github.com/golang/protobuf/proto"
	"github.com/linnv/logx"
)

func TestLicense_GetMacAddress(t *testing.T) {
	lic := &License{}
	lic.SecondsLimit = *proto.Int64(22)
	lic.SecondsUsed = *proto.Int64(21)
	lic.LicenseID = *proto.String("ssmac")
	bs, err := proto.Marshal(lic)
	if err != nil {
		t.Fatal(err.Error())
	}
	logx.Debugf("bs: %s\n", bs)

	cryptStr := rsa.CrypteMsg(bs)
	logx.Debugf("cryptStr: %s\n", cryptStr)

	nbs := rsa.DeCrypteMsg(cryptStr)
	var newLic License
	err = proto.UnmarshalMerge(nbs, &newLic)
	if err != nil {
		t.Fatal(err.Error())
	}
	logx.Debugf("newLic: %s\n", newLic.String())
}
