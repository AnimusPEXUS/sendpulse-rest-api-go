package sendpulse

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/AnimusPEXUS/sendpulse-rest-api-go/types"
)

const SENDPULSE_URI = `https://api.sendpulse.com`

type SendPulse struct {
	Debug bool

	client_id             string
	client_secret         string
	token                 *types.TokenResponse
	token_expiration_time time.Time

	signin_error error

	http_client *http.Client
}

func NewSendPulse(client_id, client_secret string) (*SendPulse, error) {
	self := &SendPulse{
		Debug:                 false,
		client_id:             client_id,
		client_secret:         client_secret,
		token_expiration_time: time.Now(),
		// signin_error:          nil,
		http_client: &http.Client{},
	}
	return self, nil
}

func (self *SendPulse) IsTokenUpdateRequired() bool {
	return time.Now().After(self.token_expiration_time)
}

func (self *SendPulse) UpdateToken(force_renew bool) error {
	if force_renew || self.IsTokenUpdateRequired() {
		return self.updateToken()
	}
	return nil
}

func (self *SendPulse) Token() *types.TokenResponse {
	return self.token
}

func (self *SendPulse) SendRequest(
	path, method string,
	params *url.Values,
	use_token bool,
	force_token_renew bool,
) (*http.Response, error) {
	err := self.UpdateToken(force_token_renew)
	if err != nil {
		return nil, err
	}
	return self.sendRequest(path, method, params, use_token)
}

func (self *SendPulse) sendRequest(
	path, method string,
	params *url.Values,
	use_token bool,
) (*http.Response, error) {

	method = strings.ToUpper(method)

	req, err := http.NewRequest(
		method,
		fmt.Sprintf("%s/%s", SENDPULSE_URI, path),
		nil,
	)
	if err != nil {
		return nil, err
	}

	if use_token && self.token != nil {
		req.Header.Set(
			"Authorization",
			fmt.Sprintf(
				"%s %s",
				self.token.TokenType,
				self.token.AccessToken,
			),
		)
	}

	if method == "POST" && params != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		b := params.Encode()
		if self.Debug {
			log.Print("body ", b)
		}
		req.Body = ioutil.NopCloser(strings.NewReader(b))
	}

	if self.Debug {
		log.Print("header ", req.Header)
	}

	return self.http_client.Do(req)
}

func (self *SendPulse) updateToken() error {

	values := &url.Values{
		"grant_type":    []string{"client_credentials"},
		"client_id":     []string{self.client_id},
		"client_secret": []string{self.client_secret},
	}

	request_time := time.Now()

	resp, err := self.sendRequest("oauth/access_token", "POST", values, false)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		if self.Debug {
			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			log.Print("bad response body ", string(data))
		}
		return errors.New(fmt.Sprintf("http failure. code: %d", resp.StatusCode))
	}

	resp_form := &types.TokenResponse{}

	{
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		err = json.Unmarshal(data, resp_form)
		if err != nil {
			return err
		}
	}

	self.token = resp_form

	self.token_expiration_time = request_time.Add(time.Duration(self.token.ExpiresIn) * time.Second)

	return nil
}
