package sdkws

import (
	"errors"
	"fmt"

	"github.com/openimsdk/protocol/constant"
)

func (x *MsgData) Check() error {
	if x.SendID == "" {
		return errors.New("sendID is empty")
	}
	if x.Content == nil {
		return errors.New("content is empty")
	}
	if x.SessionType == constant.NotificationChatType && x.ContentType != constant.OANotification ||
		x.SessionType != constant.NotificationChatType && x.ContentType == constant.OANotification {
		return errors.New("notification msg must have correct session type and content type")
	}
	return nil
}

func (x *RequestPagination) Check() error {
	if x == nil {
		return errors.New("pagination is nil")
	}
	if x.PageNumber < 1 {
		return errors.New("pageNumber is invalid")
	}
	if x.ShowNumber < 1 {
		return errors.New("showNumber is invalid")
	}
	return nil
}

func (x *GetMaxSeqResp) Format() any {
	if len(x.MaxSeqs) > 50 {
		return fmt.Sprintf("len is %v", len(x.MaxSeqs))
	}
	if len(x.MinSeqs) > 50 {
		return fmt.Sprintf("len is %v", len(x.MinSeqs))
	}
	return x
}
