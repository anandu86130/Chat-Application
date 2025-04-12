package chat

type Client struct{
	ID string
	MsgChan chan string
}