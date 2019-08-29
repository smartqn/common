// Package mobile
package mobile

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"strings"

	"time"

	"github.com/smartqn/common/config"
	"github.com/smartqn/common/util"

	"github.com/linnv/logx"
	uuid "github.com/satori/go.uuid"
)

var conf = config.Config()

type Mobile struct {
	Mobile            string `json:"Mobile"`
	Name              string `json:"Name,omitempty"`
	Ext               string `json:"Ext"`
	SMSChannel        string `json:"SMSChannel"`
	TextMsg           string `json:"TextMsg"`
	MessageTemplateID string `json:"MessageTemplateID,omitempty"`
	MessageQnzsID     string `json:"MessageQnzsID,omitempty"`

	FlowID       string `json:"FlowID,omitempty"`
	SessionID    string `json:"SessionID"`
	EnterpriseID string `json:"EnterpriseID"`
}

func NewMobile(sessID, eID, flowID, mobile, msg, tplID string) *Mobile {
	flowID = strings.TrimPrefix(flowID, eID+"_")
	return &Mobile{
		Mobile:            mobile,
		MessageTemplateID: tplID,
		FlowID:            flowID,
		SessionID:         sessID,
		EnterpriseID:      eID,
		TextMsg:           msg,
	}
}

const NotifyTypeMobile = byte('1')
const MAX_MSG_TEX_LEN = 60

func (m *Mobile) GetDetail() (instance interface{}, kind byte, err error) {
	return m, NotifyTypeMobile, nil
}

func (m *Mobile) Send() error {
	if m == nil {
		return nil
	}
	const splitQuota = ";;"
	msgUUID := uuid.NewV4()
	qnzsMsgID := m.EnterpriseID + splitQuota + msgUUID.String()
	if len(qnzsMsgID) > MAX_MSG_TEX_LEN {
		qnzsMsgID = qnzsMsgID[:MAX_MSG_TEX_LEN]
	}
	m.MessageQnzsID = qnzsMsgID
	logx.Debugfln(" send msg: %+v , mobile:%s\n , smsChannel", m.TextMsg, m.Mobile, m.SMSChannel)
	resp, err := Send(m.TextMsg, m.Mobile, m.SMSChannel, qnzsMsgID)
	if err != nil {
		logx.Warnf(" err: %+v\n", err)
		return err
	}
	_ = resp
	logx.Debugf("send mobile msg %s\n", m.TextMsg)
	return nil
}

type ReqMsgQn struct {
	Clientid     string `json:"clientid"`
	Password     string `json:"password"`
	Mobile       string `json:"mobile"`
	Content      string `json:"content"`
	Extend       string `json:"extend"`
	UID          string `json:"uid"`
	CompressType string `json:"compress_type"`
}

type RespMsgQn struct {
	ComporessType string `json:"comporess_type"`
	Data          []struct {
		Code   int    `json:"code"`
		Fee    int    `json:"fee"`
		Mobile string `json:"mobile"`
		Msg    string `json:"msg"`
		Sid    string `json:"sid"`
		UID    string `json:"uid"`
	} `json:"data"`
	TotalFee int `json:"total_fee"`
}

func (r *RespMsgQn) GetMsgID() string {
	if r == nil {
		return ""
	}
	for _, v := range r.Data {
		return v.Sid
	}
	return ""
}

func (r *RespMsgQn) GetMsgQnzsID() string {
	if r == nil {
		return ""
	}
	for _, v := range r.Data {
		return v.UID
	}
	return ""
}

const NoCompress = "1"
const SMSIndustryChan = "SMS_INDUSTRY"

type QnMsg struct {
}

//@TODO use interface
// func (qm *QnMsg) Send(c, mobile, qnzsID string) (respStruct RespMsgQn, err error) {
// 	return Send(c, mobile, qnzsID)
// }

func smsDispatch(smsChannel string) (string, string) {
	if smsChannel == SMSIndustryChan {
		return conf.SMSIndustryAccount, conf.SMSIndustryPassword
	}

	return conf.SMSMarketAccount, conf.SMSMarketPassword
}

func Send(c, mobile, smsChannel, qnzsID string) (respStruct RespMsgQn, err error) {
	if conf == nil {
		logx.Warnln("smartQn conf not init")
		return
	}
	reqUrl := conf.SmsSrvAddr + conf.MsgSendApi
	account, passw := smsDispatch(smsChannel)

	if qnzsID == "" {
		qnzsID = util.RandNumber(10)
	}
	reqDataStruct := ReqMsgQn{
		Clientid:     account,
		Password:     util.GetMd5(passw),
		Mobile:       mobile,
		Content:      c,
		UID:          qnzsID,
		CompressType: NoCompress,
	}
	// respStruct := RespMsgQn{}

	postdata, err := json.Marshal(reqDataStruct)
	if err != nil {
		logx.Warnf("err: %+v\n", err)
		return
	}

	logx.Debugf("postdata: \n%s\n", postdata)
	req, err := http.NewRequest("POST", reqUrl, bytes.NewReader(postdata))
	if err != nil {
		logx.Warnf("err: %+v\n", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Do(req)
	if err != nil {
		logx.Warnf("err: %+v\n", err)
		return
	}

	bs, _ := httputil.DumpRequest(req, true)
	bsResp, _ := httputil.DumpResponse(resp, true)
	logx.Debugf("\nreq %s\n%s \n\nresp: %s", bs, postdata, bsResp)

	if err = json.NewDecoder(resp.Body).Decode(&respStruct); err != nil {
		logx.Warnf("err: %+v\n", err)
		return
	}
	resp.Body.Close()

	return
}
