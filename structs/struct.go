package structs

import "github.com/R1kkass/GoCloudGRPC/proto/chat"

type DataStreamConnect struct {
	ChatId int
	UserID uint
	Stream chat.ChatGreeter_StreamGetMessagesServer
	Chan   chan *chat.Message
}
