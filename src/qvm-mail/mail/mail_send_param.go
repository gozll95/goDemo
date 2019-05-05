package mail

type MailSendParam struct {
	Host       string           `json:"host"`
	AccessKey  string           `json:"access_key"`
	SecretKey  string           `json:"secret_key"`
	TemplateId string           `json:"template_id"`
	Template   *MessageTemplate `json:"template"`
	Data       string           `json:"data"`
	Provider   *string          `json:"provider,omitempty"`
}

type MessageShowParam struct {
	Host      string `json:"host"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	MessageId string `json:"message_id"`
}

type MessageTemplate struct {
	To          string `json:"to"`
	Subject     string `json:"subject"`
	Content     string `json:"content"`
	MessageType string `json:"message_type"`
	RenderType  string `json:"render_type"`
}

type TemplateData struct {
	User  *UserInfo  `json:"user"`
	User1 *User1Info `json:"user1"`
	SMS   *SMSInfo   `json:"sms"`
}

type UserInfo struct {
	Email string `json:"email"`
}

type User1Info struct {
	A string `json:"a"`
}

type SMSInfo struct {
	Phone string `json:"phone"`
}
