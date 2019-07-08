package sendpulse

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/AnimusPEXUS/sendpulse-rest-api-go/types"
)

func (self *SendPulse) SmtpEmailsPost(email *types.SendPulseSendEmailStruct) (*http.Response, error) {

	var s string

	{
		b, err := json.Marshal(email)
		if err != nil {
			return nil, err
		}
		s = string(b)
	}

	if self.Debug {
		log.Print("email json ", s)
	}

	resp, err := self.SendRequest(
		"smtp/emails",
		"POST",
		&url.Values{"email": []string{s}},
		true,
		false,
	)
	return resp, err
}
