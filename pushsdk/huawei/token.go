package huawei

import (
  "sync"
  "time"
  "net/url"
  "encoding/json"
  "github.com/alecthomas/log4go"
  "errors"
)

const (
  maxRequestTokenRetry = 3
)

type tokenData struct {
  AccessToken string `json:"access_token"`
  ExpireIn    int64  `json:"expires_in"`
  ExpireAt    int64  `json:"expires_at"`
  Scope       string `json:"scope"`
  Error       int32  `json:"error"`
  ErrorDesc   string `json:"error_description"`
}

type token struct {
  data  tokenData
  lock  *sync.RWMutex
  timer *time.Timer
}

func newToken() *token {
  data := tokenData{AccessToken: "", ExpireAt: 0}
  lock := new(sync.RWMutex)
  return &token{data, lock, nil}
}

func (t *token) refresh(cliId string, cliSecret string) {
  t.lock.Lock()
  defer t.lock.Unlock()
  nowSecs := time.Now().Unix()
  t.doRefresh(nowSecs, cliId, cliSecret)
}

func (t *token) get(cliId string, cliSecret string) (string, error) {
  t.lock.Lock()
  defer t.lock.Unlock()
  nowSecs := time.Now().Unix()
  if t.data.AccessToken != "" && t.data.ExpireAt > nowSecs {
    return t.data.AccessToken, nil
  } else if err := t.doRefresh(nowSecs, cliId, cliSecret); err != nil {
    return "", err
  } else {
    return t.data.AccessToken, nil
  }
}

func (t *token) expireImmediately() {
  t.lock.Lock()
  defer t.lock.Unlock()
  t.data.AccessToken = ""
  if t.timer != nil {
    t.timer.Stop()
    t.timer = nil
  }
}

func (t *token) doRefresh(nowSecs int64, cliId string, cliSecret string) error {
  f := url.Values{}
  f.Add("client_id", cliId)
  f.Add("client_secret", cliSecret)
  f.Add("grant_type", "client_credentials")
  bytes, err := httpPost(accessTokenApi, f, maxRequestTokenRetry)
  if err != nil {
    return err
  }
  var newToken tokenData
  err = json.Unmarshal(bytes, &newToken)
  if err != nil {
    return err
  }
  newToken.ExpireAt = nowSecs + newToken.ExpireIn
  t.data = newToken
  if t.timer != nil {
    t.timer.Stop()
  }
  if t.data.Error != 0 {
    t.data.AccessToken = ""
    return errors.New(t.data.ErrorDesc)
  }
  expireTime := time.Second * time.Duration(newToken.ExpireIn)
  t.timer = time.AfterFunc(expireTime, func() {
    t.lock.Lock()
    defer t.lock.Unlock()
    t.data.AccessToken = ""
  })
  log4go.Debug("huawei: new token '%v', expired in '%v' secs\n",
      t.data.AccessToken, t.data.ExpireIn)
  return nil
}

var accessToken = newToken()
