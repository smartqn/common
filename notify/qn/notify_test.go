package qn

import (
	"io/ioutil"
	"smartqn/common/apistruct"
	"smartqn/config"
	"testing"
)

func TestQnNotify_Send(t *testing.T) {
	return
	bs, err := ioutil.ReadFile("/home/yangjc/qnzsPro/smartOutCall/src/MessageProxy/config/config.yaml")
	if err != nil {
		t.Fatalf(err.Error())
		return
	}

	config.Init(bs)
	items := make(map[string]string)
	items["ent_id"] = "2019052801"
	items["ent_secret"] = "79bac4dc9b09416eaf606370bd463c34"
	items["session_id"] = "7304106625105657903"
	items["@NUMBER"] = "TEL:99048DFD17F3B8QQ78322959761"
	items["_real_msg"] = "hello world"
	items["dn"] = "10000040"
	items["custHostNum"] = "173******09"
	items["custID"] = "123456789012345678"
	items["agentID"] = "testAgentID"
	items["groupID"] = "testGroupID"
	items["groupName"] = "testGroupName"

	items[config.Config().EnterpriseStrEx] = "10000001"
	items[config.Config().FlowIdStr] = "20190830"

	var diaLogHistory []*apistruct.Msg
	msg := &apistruct.Msg{}
	msg.Question = "what?"
	msg.Answer = "yes"
	diaLogHistory = append(diaLogHistory, msg)

	notify := NewQnNotify(NOTIFY_TYPE_ACTION_MSG, items)
	notify.DialogHistory = diaLogHistory

	isSuc, noRetry, err := notify.Send()
	if err != nil {
		t.Fatalf(err.Error())
		return
	}
	if !isSuc || !noRetry {
		t.Fatalf(err.Error())
		return
	}

	return
}
