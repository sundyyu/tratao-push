package meizu

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/alecthomas/log4go"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

const (
	maxPostRetry = 3
)

const (
	apiHostUrl             = "http://server-api-push.meizu.com"
	apiSendNotificationUri = "/garcia/api/server/push/varnished/pushByPushId"
)

type Client struct {
	appId     string
	appSecret string
}

func NewClient(appId string, appSecret string) *Client {
	return &Client{appId, appSecret}
}

func (c *Client) Push(msg *NotificationMsg, ids []string) (
	*Response, error) {
	params := url.Values{}
	params.Add("appId", c.appId)
	params.Add("pushIds", strings.Join(ids, ","))
	msgStr, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	params.Add("messageJson", string(msgStr))
	sign := c.sign(params)
	params.Add("sign", sign)
	bytes, err := httpPost(apiHostUrl+apiSendNotificationUri, params,
		maxPostRetry)
	if err != nil {
		return nil, err
	}
	fmt.Printf("meizu: response '%v'\n", string(bytes))
	var resp Response
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) sign(params url.Values) string {
	var keys []string
	var parameterStr string

	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, value := range keys {
		parameterStr += value + "=" + params[value][0]
	}
	parameterStr += c.appSecret

	md5Ctx := md5.New()
	md5Ctx.Write([]byte(parameterStr))
	signStr := hex.EncodeToString(md5Ctx.Sum(nil))

	return signStr
}

func httpPost(url string, form url.Values, maxRetry int) ([]byte, error) {
	var result []byte
	var req *http.Request
	var resp *http.Response
	var err error
	req, err = http.NewRequest("POST", url, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	try := 0
tryAgain:
	resp, err = client.Do(req)
	if err != nil {
		log4go.Warn("meizu: post error: %v, try count: %v", err, try)
		try += 1
		if try < maxRetry {
			goto tryAgain
		}
		return nil, err
	}
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
