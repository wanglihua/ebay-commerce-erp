package xxgl_models

var MESSAGE_TAG_NORMAL MessageTagInfo = MessageTagInfo{Id: 1, Name: "普通消息"}
var MESSAGE_TAG_NOREPLY MessageTagInfo = MessageTagInfo{Id: 2, Name: "无需回复"}
var MESSAGE_TAG_OTHER MessageTagInfo = MessageTagInfo{Id: 3, Name: "其他"}

func Get_message_tag_by_id(tag_id int64) MessageTagInfo {
	if tag_id == MESSAGE_TAG_NORMAL.Id {
		return MESSAGE_TAG_NORMAL
	} else if tag_id == MESSAGE_TAG_NOREPLY.Id {
		return MESSAGE_TAG_NOREPLY
	} else if tag_id == MESSAGE_TAG_OTHER.Id {
		return MESSAGE_TAG_OTHER
	}

	return MessageTagInfo{}
}

type MessageTagInfo struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}
