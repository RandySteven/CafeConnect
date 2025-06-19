package messages

type EmailMessage struct {
	Receiver string `json:"receiver"`
	Header   string `json:"header"`
	Content  string `json:"email"`
}
