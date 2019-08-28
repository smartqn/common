package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/linnv/logx"
)

func HttpPostJson(targetUrl string, dataStruct interface{}, respStruct interface{}, client *http.Client) error {
	if client == nil {
		return fmt.Errorf("nil client")
	}

	postdata, err := json.Marshal(dataStruct)
	if err != nil {
		logx.Warnf("err: %+v\n", err)
		return err
	}

	logx.Warnf("bs: %s\n", postdata)
	req, err := http.NewRequest("POST", targetUrl, bytes.NewReader(postdata))
	if err != nil {
		logx.Warnf("err: %+v\n", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		logx.Warnf("err: %+v\n", err)
		return err
	}
	if resp.StatusCode != http.StatusOK {
		bs, _ := httputil.DumpRequest(req, true)
		bsResp, _ := httputil.DumpResponse(resp, true)
		logx.Warnf("req %s \n resp: %s", bs, bsResp)
	}

	bs, _ := httputil.DumpRequest(req, true)
	bsResp, _ := httputil.DumpResponse(resp, true)
	logx.Warnf("req %s \n resp: %s\n", bs, bsResp)

	if respStruct != nil {
		if err := json.NewDecoder(resp.Body).Decode(&respStruct); err != nil {
			logx.Warnf("err: %+v\n", err)
			return err
		}
		resp.Body.Close()
	}
	return nil
}
