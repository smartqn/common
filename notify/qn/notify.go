package qn

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/linnv/logx"
	uuid "github.com/satori/go.uuid"
	"github.com/smartqn/common/apistruct"
	"github.com/smartqn/common/config"
	"github.com/smartqn/common/util"
)

var httpClient = &http.Client{
	Timeout: time.Second * 10,
}

const NOTIFY_TYPE_QN_ROBOT = byte('2')

//1转坐席2发短信3信息采集
const NOTIFY_TYPE_ACTION_AGENT = byte(1)
const NOTIFY_TYPE_ACTION_MSG = byte(2)
const NOTIFY_TYPE_ACTION_COLLECT = byte(3)

const ROBOT = "robot"
const CUSTOMER = "cus"

type QnNotify struct {
	Items      map[string]string `json:"items"`
	NotifyType byte              `json:"notifyType"`

	Mobile       string `json:"mobile"`
	EnterPriseID string `json:"enterpriseID"`
	FlowID       string `json:"flowID"`
	SmsTplID     string `json:"smsTplID"`
	QnMsgID      string `json:"qnMsgID"`
	SessionID    string `json:"sessionID"`

	DialogHistory []*apistruct.Msg `json:"DialogHistory"`
}

type chatRecord struct {
	User string `json:"user"`
	Msg  string `json:"msg"`
}

type respData struct {
	Result string `json:"result"`
	Desc   string `json:"desc"`
	Object string `json:"object"`
}

type ReqRobotNotify struct {
	EntID          string `json:"entID"`
	RequestID      string `json:"requestID"`
	Timestamp      string `json:"timestamp"`
	Mac            string `json:"mac"`
	SessionId      string `json:"sessionId"`
	SerialNo       string `json:"serialNo"`
	DeviceNumber   string `json:"deviceNumber"`
	OptType        int    `json:"type"`
	CustomerNum    string `json:"customerNum"`
	Content        string `json:"content"`
	SkillId        string `json:"skillId"`
	SkillName      string `json:"skillName"`
	AgentId        string `json:"agentId"`
	UserData       string `json:"userData"`
	MakeCallTime   string `json:"makeCallTime"`
	CoverAgentTime string `json:"coverAgentTime"`
	TotalDuration  string `json:"totalDuration"`
	CusID          string `json:"cusID"`
	CusHostNum     string `json:"cusHostNum"`
	CusSex         string `json:"cusSex"`
	CusName        string `json:"cusName"`
	CusArea        string `json:"cusArea"`
	//ChatRecord string `json:"chatRecord"`
	ChatRecord []chatRecord `json:"chatRecord"`
}

type RespRobotNotify struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

//type markData struct {
//	GroupID string `json:"groupId"`
//	GroupName string `json:"groupName"`
//}

func parseIDByItems(items map[string]string, mayKeys []string) (eid string) {
	for _, v := range mayKeys {
		enterpriseID, ok := items[v]
		if ok && enterpriseID != "" {
			eid = enterpriseID
		}
	}
	return
}

func sessionIDTransform(sessionID string) (result string) {
	const SESSION_ID_MIN_SPLIT_NUM = 3
	idList := strings.Split(sessionID, ":")
	if len(idList) < SESSION_ID_MIN_SPLIT_NUM {
		return
	}

	str := idList[0]
	if str == "" {
		logx.Warnf("invalid session_id: %s\n", sessionID)
		return
	}

	bn := util.Hex2BigInt(str)
	result = bn.String()
	return
}

func newReqRobotNotify(items map[string]string, diaLogHistory []*apistruct.Msg) (req *ReqRobotNotify, err error) {
	//MD5(entID+entSecret+timestamp+requestID)
	logx.Debugf("items: %+v\n", items)
	if items == nil {
		logx.Warnln("item nil ")
		return req, fmt.Errorf("item nil ")
	}
	entID := items["ent_id"]
	entSecret := items["ent_secret"]
	timestamp := time.Now().Format("20060102150405")
	msgUUID, err := uuid.NewV4()
	if err != nil {
		logx.Warnf("create uuid err: %s\n", err.Error())
		return req, err
	}
	requestID := msgUUID.String()
	mac := util.GetMd5(entID + entSecret + timestamp + requestID)
	sessionID := items["session_id"]
	callOut := items["@NUMBER"]
	content := items["_real_msg"]
	dn := items["dn"]
	custHostNum := items["custHostNum"]
	custID := items["custID"]
	agentID := items["agent_id"]
	groupID := items["groupId"]
	groupName := items["groupName"]

	result := sessionIDTransform(sessionID)
	if result != "" {
		sessionID = result
	}

	//chatRecordStr := items["chatRecord"]

	/*	mobile := items["_real_mobile"]
		smsTplID := items["msg_tpl_id"]

		eIDKey, eIDDefaultKey := config.Config().EnterpriseStrEx, config.Config().EnterpriseStr
		flowIDKey, flowIDDefaultKey := config.Config().FlowIdStrEx, config.Config().FlowIdStr
		enterPriseID := parseIDByItems(items, []string{eIDKey, eIDDefaultKey})
		flowID := parseIDByItems(items, []string{flowIDKey, flowIDDefaultKey})*/

	//logx.Debugf("robot chats: %+v\n", chatRecordStr)
	//#TODO
	//err = json.Unmarshal([]byte(chatRecordStr), chatRecord)
	//if err != nil {
	//	logx.Warnf("chatRecord json.Unmarshal err: %s\n", err.Error())
	//	return
	//}
	req = &ReqRobotNotify{}
	req.Mac = mac
	req.EntID = entID
	req.RequestID = requestID
	req.Timestamp = timestamp
	req.SessionId = sessionID
	req.DeviceNumber = dn
	req.SerialNo = callOut
	req.CusHostNum = custHostNum
	req.CusID = custID
	req.AgentId = agentID
	req.SkillId = groupID
	req.SkillName = groupName

	req.Content = content
	req.ChatRecord = make([]chatRecord, 0)
	for _, v := range diaLogHistory {
		chatRec := []chatRecord{{CUSTOMER, v.Question}, {ROBOT, v.Answer}}
		req.ChatRecord = append(req.ChatRecord, chatRec...)
	}
	logx.Debugf("robot chats: %+v\n", req.ChatRecord)

	return req, nil
}

func NewQnNotify(notifyType byte, items map[string]string) *QnNotify {
	one := &QnNotify{
		NotifyType: notifyType,
		Items:      items,
	}
	return one
}

func (m *QnNotify) GetDetail() (instance interface{}, kind byte, err error) {
	if m == nil {
		return
	}
	return m, NOTIFY_TYPE_QN_ROBOT, nil
}

func (m *QnNotify) Send() (isSuc bool, noRetry bool, err error) {
	if m == nil {
		logx.Warnln("nil QnNotify obj ")
		return false, true, fmt.Errorf("nil QnNotify obj ")
	}
	reqDataStruct, err := newReqRobotNotify(m.Items, m.DialogHistory)
	if err != nil {
		logx.Warnf("newReqRobotNotify err: %s\n", err.Error())
		return false, true, err
	}

	reqDataStruct.OptType = int(m.NotifyType)
	m.QnMsgID = reqDataStruct.RequestID
	m.SessionID = reqDataStruct.SessionId
	m.Mobile = m.Items["_real_mobile"]
	m.SmsTplID = m.Items["msg_tpl_id"]
	eIDKey, eIDDefaultKey := config.Config().EnterpriseStrEx, config.Config().EnterpriseStr
	flowIDKey, flowIDDefaultKey := config.Config().FlowIdStrEx, config.Config().FlowIdStr
	m.EnterPriseID = parseIDByItems(m.Items, []string{eIDKey, eIDDefaultKey})
	m.FlowID = parseIDByItems(m.Items, []string{flowIDKey, flowIDDefaultKey})
	if m.EnterPriseID == "" {
		logx.Warnf("can't get any enterpriseID from items: %+v\n", m.Items)
	}
	if m.FlowID == "" {
		logx.Warnf("can't get any flowID from items: %+v\n", m.Items)
	}

	logx.Debugf("reqDataStruct: %+v\n", reqDataStruct)
	logx.Debugf("NotifyType: %+v\n", m)
	if reqDataStruct == nil {
		logx.Warnln("reqDataStruct nil")
		return false, true, fmt.Errorf("reqDataStruct nil\n")
	}
	respStruct, err := reqDataStruct.Send()
	if err != nil {
		logx.Warnf("QnNotify QnRobot send err: %s\n", err.Error())
		return false, false, err
	}

	if respStruct.Code != 0 {
		logx.Warnf("QnNotify QnRobot err, code: %d, msg: %s\n", respStruct.Code, respStruct.Msg)
		return false, true, fmt.Errorf("QnNotify QnRobot err, code: %d, msg: %s\n", respStruct.Code, respStruct.Msg)
	}

	return true, true, nil
}

func (r *ReqRobotNotify) Send() (respStruct RespRobotNotify, err error) {
	reqUrl := config.Config().NotifyRobotAddr + config.Config().RobotFinishCallReq
	postdata, err := json.Marshal(r)
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

	resp, err := httpClient.Do(req)
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
