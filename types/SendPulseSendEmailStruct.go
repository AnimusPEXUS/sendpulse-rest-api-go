package types

type (
	SendPulseSendEmailStructEmailAddr struct {
		Name  string `json:"name"`
		EMail string `json:"email"`
	}

	SendPulseSendEmailStructEmailAddrList []SendPulseSendEmailStructEmailAddr

	SendPulseSendEmailStruct struct {
		Html    *string                               `json:"html,omitempty"`
		Text    *string                               `json:"text,omitempty"`
		Subject string                                `json:"subject"`
		From    SendPulseSendEmailStructEmailAddr     `json:"from"`
		To      SendPulseSendEmailStructEmailAddrList `json:"to"`
	}
)
