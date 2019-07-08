package main

import (
	"io/ioutil"
	"log"

	sendpulse "github.com/AnimusPEXUS/sendpulse-rest-api-go"
	"github.com/AnimusPEXUS/sendpulse-rest-api-go/types"
)

func main() {

	log.Print("getting token..")

	m, err := sendpulse.NewSendPulse("", "")
	if err != nil {
		log.Fatal(err)
	}

	log.Print("got token ", m.Token())

	log.Print("sending message..")

	email := &types.SendPulseSendEmailStruct{
		// Html: &[]string{"test"}[0],
		Subject: "test subject",
		Text:    &[]string{"test"}[0],
		To: types.SendPulseSendEmailStructEmailAddrList{
			types.SendPulseSendEmailStructEmailAddr{
				"lundovsky",
				"a@guidsy.travel",
			},
		},
		From: types.SendPulseSendEmailStructEmailAddr{
			"Guidsy Service EMail",
			"noreply@guidsy.travel",
		},
	}

	res, err := m.SmtpEmailsPost(email)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("sending result ", res)

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("sending result body ", string(b))

}
