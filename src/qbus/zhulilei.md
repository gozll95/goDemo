
type Sender interface {
	Send(message *models.MessageModel) (err error)
}


type MessageModel struct {
	Id        bson.ObjectId `bson:"_id" json:"id"`
	AppId     string        `bson:"app_id" json:"app_id"`
	Subject   string        `bson:"subject" json:"subject"`
	From      string        `bson:"from" json:"from"`
	To        string        `bson:"to" json:"to"`
	UserId    uint32        `bson:"user_id"json:"user_id"`
	Content   string        `bson:"content" json:"content"`
	Type      MessageType   `bson:"type" json:"type"`
	Status    MessageStatus `bson:"status" json:"status"`
	UpdatedAt time.Time     `bson:"updated_at" json:"updated_at"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`

	isNewRecord bool `bson:"-" json:"-"`
}


Notification *send.Config
type Config struct {
	Host       string       `json:"host"`
	ClientId   string       `json:"client_id"`
	Provider   ProviderType `json:"provider"`
	Retry      int          `json:"retry"`
	Ratelimit  int          `json:"ratelimit"` // per minute
	Concurrent int          `json:"concurrent"`
	Timeout    int          `json:"timeout"`
}

type ProviderType string


# !++ manager.go-types
type MessageManager struct {
	Queue    *MessageQueue
	Sender   Sender
	SendChan chan *models.MessageModel

	Config *Config
	Logger gogo.Logger
}

type Config struct {
        Host        string        `json:"host"`
        ClientId    string        `json:"client_id"`
        Provider    ProviderType    `json:"provider"`
        Retry        int        `json:"retry"`
        Ratelimit    int        `json:"ratelimit"`
        Concurrent    int        `json:"concurrent"`
        Timeout        int        `json:"timeout"`
    }

type ProviderType string


type MessageQueue struct {
	First *MessageNode
	Tail  *MessageNode
	mutex sync.Mutex
}

type MessageNode struct {
	Message *models.MessageModel
	Next    *MessageNode
}

# !-- manager.go-types


bbbb里Sender
type Sender struct {
	Host      string
	Token     Token
	ClientId  string
	token     string
	expiredAt time.Time
	Client    *httpclient.HTTPClient
	Headers   map[string]string
}


认证这里:
func (config *Config) QTransport() *oauth.Transport {
	return &oauth.Transport{
		Config: &oauth.Config{
			ClientId:     config.Account.ClientId,
			ClientSecret: config.Account.ClientSecret,
			Scope:        "Scope",
			AuthURL:      "<AuthURL>",
			TokenURL:     config.Account.Host + "/oauth2/token",
			RedirectURL:  "<RedirectURL>",
		},
		// Transport: http.DefaultTransport,
	}
}


var sender send.Sender
//从配置文件里选出sender
switch config.Notification.Provider 
    - case send.Providerbbbb:
                - sender = bbbb.NewSender(config.Notification.Host, config.Notification.ClientId,
                    func() (token string, expiredAt time.Time, err error) {
                        // fetch admin token
                        adminToken, _, err := APP.xauth.ExchangeByPassword(config.Account.User, config.Account.Passwd)
                        if err != nil {
                            return "", expiredAt, err
                        }

                        expiredAt = time.Unix(adminToken.TokenExpiry, 0)

                        return adminToken.AccessToken, expiredAt, nil
                    })
                            - func NewSender(host, clientId string, token Token) *Sender 
                                - 	sender := &Sender{
                                        Host:     host,
                                        Token:    token,
                                        ClientId: clientId,
                                        Client:   httpclient.NewHTTPClient(nil),
                                        Headers:  make(map[string]string),
                                    }
                                - // add client id headers
                                - sender.Headers["Client-Id"] = clientId
                                - return sender


APP.messageManager = send.NewMessageManager(sender, config.Notification, APP.Logger())
    - func NewMessageManager(sender Sender, config *Config, logger gogo.Logger) *MessageManager 
        - 	manager := &MessageManager{
                Queue: &MessageQueue{
                    First: nil,
                    Tail:  nil,
                },
                Sender:   sender,
                SendChan: make(chan *models.MessageModel, config.Ratelimit+1),
                Config:   config,
                Logger:   logger,
            }
        //开启goroutine dispatch(调度)
        - go manager.dispatch()
            - func (manager *MessageManager) dispatch() 
                -  ticker := time.Tick(time.Second * 60)
                - for now := range ticker
                    - for 循环
                        - message := manager.Queue.FetchMessage()
                            - func (queue *MessageQueue) FetchMessage() (message *models.MessageModel)
                                - queue.mutex.Lock()
                                - defer queue.mutex.Unlock()
                                - if queue.IsEmpty()
                                         - func (queue *MessageQueue) IsEmpty() bool
                                            - if queue.First == nil || queue.Tail == nil 
                                                - return true
                                    - return nil
                                //单向链表式
                                - message = queue.First.Message
                                - queue.First = queue.First.Next
                            
                        - manager.SendChan <- message
                        - 
        - for i := 0; i < config.Concurrent; i++
            - go manager.sending()
                    - func (manager *MessageManager) sending() 
                        - for 循环
                            - message := <-manager.SendChan
                            //这里是interface method
                            - err := manager.Sender.Send(message)
                            - err = message.Save()
        - return manager


# !+mail_message.go types
type MailMessage struct {
	Uid     uint32   `json:"uid"`
	To      []string `json:"to"`
	Cc      []string `json:"cc"`
	Bcc     []string `json:"bcc"`
	Subject string   `json:"subject"`
	Content string   `json:"content"`
}

# !-mail_message.go types
type MailMessage struct {
	Uid     uint32   `json:"uid"`
	To      []string `json:"to"`
	Cc      []string `json:"cc"`
	Bcc     []string `json:"bcc"`
	Subject string   `json:"subject"`
	Content string   `json:"content"`
}

type MailResult struct {
	Code int    `json:"code"`
	Oid  string `json:"oid"`
	Msg  string `json:"msg"`
}


# !+ morse的send.go
func (s *Sender) Send(message *models.MessageModel) (err error)
    - err = s.EnsureToken()
            - func (s *Sender) EnsureToken() (err error)
                 - if time.Now().Add(-10*time.Second).After(s.expiredAt) || s.token == ""
                - s.token, s.expiredAt, err = s.Token()
                - s.refreshTokenHeader()
                        - func (s *Sender) refreshTokenHeader()
                            - authorization := fmt.Sprintf("Bearer %s", s.token)
                            - s.Headers["Authorization"] = authorization
    - var res MailResult
    - switch message.Type
        - case models.MessageTypeMail:
            - mailMessage := NewMailMessage(message)
                        - func NewMailMessage(message *models.MessageModel) *MailMessage
                            - 	return &MailMessage{
                                    Uid:     message.UserId,
                                    To:      []string{message.To},
                                    Subject: message.Subject,
                                    Content: message.Content,
                                }
            - _, err = s.Client.DoJSONWithHeaders(http.MethodPost, s.SendMailUrl(), s.Headers, mailMessage, &res)
                    - 


# !- morse的send.go

# !+ controller层message.go
//省略若干代码
//这里send api就是将消息render之后加入队列
err = APP.MessageManager().SendMessage(message)
        - func (manager *MessageManager) SendMessage(message *models.MessageModel) (err error) 
            - err = message.Save()
            - manager.Queue.AddMessage(message)
                    - func (queue *MessageQueue) AddMessage(message *models.MessageModel) 
                        - queue.mutex.Lock()
                        - defer queue.mutex.Unlock()
                        - 	node := &MessageNode{
                                Message: message,
                                Next:    nil,
                            }
                        - // empty queue
                        -	if queue.IsEmpty() {
                            queue.First = node
                            queue.Tail = node
                            return
                        }
                        - queue.Tail.Next = node
                        - queue.Tail = node


# !- controller层message.go