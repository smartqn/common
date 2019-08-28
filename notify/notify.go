package notify

import (
	"encoding/json"
	"fmt"

	"github.com/linnv/logx"

	"smartqn/common/notify/mobile"
)

// type Notify struct {
// 	Type     string `json:"Type"`
// 	Notifier []byte
// }

type Notifier interface {
	Send() error
	GetDetail() (instance interface{}, kind byte, err error)
}

var ERR_INVALIDNOTIFIER = fmt.Errorf("invalid notify type")

const KEY_NOTIFY = "QNZS_NTF"

//NewNotify implements parse specified notifiter by format
//$Type$JsonBody
//e.g. `1{"Mobile":"Mobile1","Name":"Name2"}`
func NewNotify(bs []byte) (n Notifier, err error) {
	if len(bs) < 2 {
		err = ERR_INVALIDNOTIFIER
		return
	}
	const notifyTypeIndex = 0
	notifyType := bs[notifyTypeIndex]
	logx.Debugf("notifyType: %s\n", string(notifyType))

	switch notifyType {
	case mobile.NotifyTypeMobile:
		//@TODO use pool
		// logx.Debugf("mobile: %s\n", bs[notifyTypeIndex:])
		var mobile = new(mobile.Mobile)
		err = json.Unmarshal(bs[notifyTypeIndex+1:], &mobile)
		return mobile, err
	}

	return
}
