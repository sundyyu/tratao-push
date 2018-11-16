package xiaomi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	maxNumRegIds = 1000
	maxPostRetry = 3
)

const (
	officialApiHost = "https://api.xmpush.xiaomi.com"
	// officialApiHost = "https://api.xmpush.global.xiaomi.com/"
	sandboxApiHost = "https://sandbox.xmpush.xiaomi.com"
)

const (
	officialRegIdUrl = "/v3/message/regid"
	sandboxRegIdUrl  = "/v2/message/regid"
)

// -----------------------------------------------------------------------------
// Section: Client
// -----------------------------------------------------------------------------
type Client struct {
	apiHost   string
	appSecret string
	pkgNames  []string
	httpCli   *http.Client
}

func NewClient(appSecret string, pkgNames []string) *Client {
	return &Client{
		apiHost:   officialApiHost,
		appSecret: appSecret,
		pkgNames:  pkgNames,
		httpCli:   &http.Client{},
	}
}

func (c *Client) UseOfficialApi() {
	c.apiHost = officialApiHost
}

func (c *Client) UseSandboxApi() {
	c.apiHost = sandboxApiHost
}

func (c *Client) Send(msg *Message, regId string) (*Response, error) {
	params := c.buildSendParams(msg, regId)
	uri := officialRegIdUrl
	if c.apiHost == sandboxApiHost {
		uri = sandboxRegIdUrl
	}
	bytes, err := c.doPost(c.apiHost+uri, params)
	if err != nil {
		return nil, err
	}
	var result Response
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) MultiSend(msg *Message, regIds []string) (*Response, error) {
	if len(regIds) == 0 || len(regIds) > maxNumRegIds {
		return nil, fmt.Errorf("xiaomi: invalid regIds, length: %v", len(regIds))
	}
	return c.Send(msg, strings.Join(regIds, ","))
}

func (m *Client) buildSendParams(msg *Message, regId string) url.Values {
	form := m.defaultForm(msg)
	form.Add("registration_id", regId)
	return form
}

func (m *Client) defaultForm(msg *Message) url.Values {
	f := msg.Form()
	if len(m.pkgNames) > 0 {
		f.Add("restricted_package_name", strings.Join(m.pkgNames, ","))
	}
	return f
}

func (c *Client) doPost(url string, form url.Values) ([]byte, error) {
	var result []byte
	var req *http.Request
	var res *http.Response
	var err error
	req, err = http.NewRequest("POST", url, strings.NewReader(form.Encode()))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.Header.Set("Authorization", "key="+c.appSecret)
	tryTime := 0
tryAgain:
	res, err = c.httpCli.Do(req)
	if err != nil {
		fmt.Println("xiaomi push post err:", err, tryTime)
		tryTime += 1
		if tryTime < maxPostRetry {
			goto tryAgain
		}
		return nil, err
	}
	if res.Body == nil {
		return nil, errors.New("xiaomi: empty response")
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("network error %v", res.StatusCode)
	}
	result, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return result, nil
}
