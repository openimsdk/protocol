package msggateway

import "errors"

func (x *OnlinePushMsgReq) Check() error {
	if x.MsgData == nil {
		return errors.New("MsgData is empty")
	}
	if err := x.MsgData.Check(); err != nil {
		return err
	}
	if x.PushToUserID == "" {
		return errors.New("PushToUserID is empty")
	}
	return nil
}

func (x *OnlineBatchPushOneMsgReq) Check() error {
	if x.MsgData == nil {
		return errors.New("MsgData is empty")
	}
	if err := x.MsgData.Check(); err != nil {
		return err
	}
	if x.PushToUserIDs == nil {
		return errors.New("PushToUserIDs is empty")
	}
	return nil
}

func (x *GetUsersOnlineStatusReq) Check() error {
	if x.UserIDs == nil {
		return errors.New("UserIDs is empty")
	}
	return nil
}

func (x *KickUserOfflineReq) Check() error {
	if x.PlatformID < 1 || x.PlatformID > 9 {
		return errors.New("PlatformID is invalid")
	}
	if x.KickUserIDList == nil {
		return errors.New("KickUserIDList is empty")
	}
	return nil
}

func (x *MultiTerminalLoginCheckReq) Check() error {
	if x.PlatformID < 1 || x.PlatformID > 9 {
		return errors.New("PlatformID is invalid")
	}
	if x.UserID == "" {
		return errors.New("UserID is empty")
	}
	if x.Token == "" {
		return errors.New("token is empty")
	}
	return nil
}
