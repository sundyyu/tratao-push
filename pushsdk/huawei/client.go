package huawei

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alecthomas/log4go"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	maxPushRetry = 3
)

const (
	accessTokenApi = "https://login.cloud.huawei.com/oauth2/v2/token"
	pushApi        = "https://api.push.hicloud.com/pushsend.do"
	pushSvcName    = "openpush.message.api.send"
)

var (
	UnknownServerError    = errors.New("huawei: unknown server error")
	RateLimitError        = errors.New("huawei: rate limit")
	InvalidSessionError   = errors.New("huawei: invalid session")
	SessionExpiredError   = errors.New("huawei: session expired")
	InvalidParameterError = errors.New("huawei: invalid parameter")
)

type Client struct {
	cliId     string
	cliSecret string
}

func NewClient(cliId string, cliSecret string) *Client {
	return &Client{cliId, cliSecret}
}

func (c *Client) Push(n *Notification) (*Response, error) {
	params := n.Form()
	params, err := c.defaultParams(params)
	if err != nil {
		return nil, err
	}
	params.Add("nsp_svc", pushSvcName)
	nspCtx, _ := json.Marshal(map[string]interface{}{
		"ver":   "1",
		"appId": c.cliId,
	})
	nspCtxStr := url.QueryEscape(string(nspCtx))
	postUrl := fmt.Sprintf("%s?nsp_ctx=%s", pushApi, nspCtxStr)
	bytes, err := httpPost(postUrl, params, maxPushRetry)
	if err != nil {
		return nil, err
	}
	var result Response
	err = json.Unmarshal(bytes, &result)
	if err == InvalidSessionError || err == SessionExpiredError {
		accessToken.expireImmediately()
		return c.Push(n)
	} else if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) defaultParams(params url.Values) (url.Values, error) {
	t, err := accessToken.get(c.cliId, c.cliSecret)
	if err != nil {
		return params, err
	}
	nspTs := strconv.FormatInt(time.Now().Unix(), 10)
	params.Add("nsp_ts", nspTs)
	params.Add("access_token", t)
	return params, nil
}

func httpPost(url string, form url.Values, maxRetry int) ([]byte, error) {
	var result []byte
	var req *http.Request
	var resp *http.Response
	var err error
	req, err = http.NewRequest("POST", url, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type",
		"application/x-www-form-urlencoded; charset=UTF-8")
	client := &http.Client{}
	try := 0
tryAgain:
	resp, err = client.Do(req)
	if err != nil {
		log4go.Warn("huawei: post error: %v, try count: %v", err, try)
		try += 1
		if try < maxRetry {
			goto tryAgain
		}
		return nil, err
	}

	// if resp.StatusCode != 200 {
	// 	fmt.Println(url, resp)
	// }
	// if resp.StatusCode == 503 {
	// 	return nil, RateLimitError
	// } else if resp.StatusCode != 200 {
	// 	return nil, UnknownServerError
	// }
	// nspStatus := resp.Header.Get("NSP_STATUS")
	// if nspStatus == "6" {
	// 	return nil, SessionExpiredError
	// } else if nspStatus == "102" {
	// 	return nil, InvalidSessionError
	// } else if nspStatus != "" && nspStatus != "0" {
	// 	return nil, fmt.Errorf("huawei: nsp_status %v", nspStatus)
	// }

	defer resp.Body.Close()
	result, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	str := string(result)
	str, err = strconv.Unquote(str)
	if err != nil {
		str = string(result)
	}
	return []byte(str), nil
}
