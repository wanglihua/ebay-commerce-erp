package xxgl_models

var MESSAGE_TAG_NORMAL MessageTagInfo = MessageTagInfo{Id: 1, Name: "普通消息"}
var MESSAGE_TAG_NOREPLY MessageTagInfo = MessageTagInfo{Id: 2, Name: "无需回复"}
var MESSAGE_TAG_OTHER MessageTagInfo = MessageTagInfo{Id: 3, Name: "其他"}

type MessageTagInfo struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}
