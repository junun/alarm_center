package repo


type DingTalkRepository interface {
	GenerateClient(name string) (*DingTalkClient, error)
	SendMessage(name string, c *DingTalkClient, request *RobotSendRequest) (*RobotSendResponse, error)
	GetPushUrl(c *DingTalkClient) (string, error)
	Sign(timestamp string, c *DingTalkClient) (string, error)
	SendQueueDingTalkMsg()
}

type DingTalkClient struct {
	Webhook string
	Secret  string
	Keyword string
}

type Text struct {
	Content string `json:"content"`
}

type Markdown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type Link struct {
	Title      string `json:"title"`
	Text       string `json:"text"`
	PicUrl     string `json:"picUrl"`
	MessageUrl string `json:"messageUrl"`
}

type At struct {
	AtMobiles []string `json:"atMobiles"`
	IsAtAll   bool     `json:"isAtAll"`
}

type RobotSendRequest struct {
	MsgType  string   `json:"msgtype"`
	Text     Text     `json:"text"`
	Markdown Markdown `json:"markdown"`
	Link     Link     `json:"link"`
	At       At       `json:"at"`
}

type RobotSendResponse struct {
	ErrMsg  string `json:"errmsg"`
	ErrCode int64  `json:"errcode"`
}

