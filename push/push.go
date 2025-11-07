package push

import "errors"

func (x *PushMsgReq) Check() error {
	if x.MsgData == nil {
		return errors.New("MsgData is empty")
	}
	if err := x.MsgData.Check(); err != nil {
		return err
	}
	if x.ConversationID == "" {
		return errors.New("ConversationID is empty")
	}
	return nil
}

func (x *DelUserPushTokenReq) Check() error {
	if x.UserID == "" {
		return errors.New("UserID is empty")
	}
	if x.PlatformID < 1 || x.PlatformID > 9 {
		return errors.New("PlatformID is invalid")
	}
	return nil
}
