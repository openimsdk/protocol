// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package user

import (
	"errors"
	"fmt"

	"github.com/openimsdk/protocol/util/datautil"
)

func (x *GetAllUserIDReq) Check() error {
	if x.Pagination == nil {
		return errors.New("pagination is empty")
	}
	if x.Pagination.PageNumber < 1 {
		return errors.New("pageNumber is invalid")
	}
	return nil
}

func (x *AccountCheckReq) Check() error {
	if x.CheckUserIDs == nil {
		return errors.New("CheckUserIDs is empty")
	}
	return nil
}

func (x *GetDesignateUsersReq) Check() error {
	if x.UserIDs == nil {
		return errors.New("UserIDs is empty")
	}
	return nil
}

func (x *UpdateUserInfoReq) Check() error {
	if x.UserInfo == nil {
		return errors.New("UserInfo is empty")
	}
	if x.UserInfo.UserID == "" {
		return errors.New("UserID is empty")
	}
	return nil
}

func (x *SetGlobalRecvMessageOptReq) Check() error {
	if x.GlobalRecvMsgOpt > 2 || x.GlobalRecvMsgOpt < 0 {
		return errors.New("GlobalRecvMsgOpt is invalid")
	}
	if x.UserID == "" {
		return errors.New("UserID is empty")
	}
	return nil
}

func (x *SetConversationReq) Check() error {
	if err := x.Conversation.Check(); err != nil {
		return err
	}
	if x.NotificationType < 1 || x.NotificationType > 3 {
		return errors.New("NotificationType is invalid")
	}
	return nil
}
func (x *SetRecvMsgOptReq) Check() error {
	if x.OwnerUserID == "" {
		return errors.New("OwnerUserID is empty")
	}
	if x.ConversationID == "" {
		return errors.New("ConversationID is empty")
	}
	if x.RecvMsgOpt < 0 || x.RecvMsgOpt > 2 {
		return errors.New("RecvMsgOpt is invalid")
	}
	if x.NotificationType < 1 || x.NotificationType > 3 {
		return errors.New("NotificationType is invalid")
	}
	return nil
}

func (x *GetConversationReq) Check() error {
	if x.OwnerUserID == "" {
		return errors.New("OwnerUserID is empty")
	}
	if x.ConversationID == "" {
		return errors.New("ConversationID is empty")
	}
	return nil
}

func (x *GetConversationsReq) Check() error {
	if x.OwnerUserID == "" {
		return errors.New("OwnerUserID is empty")
	}
	if x.ConversationIDs == nil {
		return errors.New("ConversationIDs is empty")
	}
	return nil
}

func (x *GetAllConversationsReq) Check() error {
	if x.OwnerUserID == "" {
		return errors.New("OwnerUserID is empty")
	}
	return nil
}

func (x *BatchSetConversationsReq) Check() error {
	if x.OwnerUserID == "" {
		return errors.New("OwnerUserID is empty")
	}
	if x.Conversations == nil {
		return errors.New("ConversationIDs is empty")
	}
	if x.NotificationType < 1 || x.NotificationType > 3 {
		return errors.New("NotificationType is invalid")
	}
	return nil
}

func (x *GetPaginationUsersReq) Check() error {
	if x.Pagination == nil {
		return errors.New("pagination is empty")
	}
	if x.Pagination.PageNumber < 1 {
		return errors.New("pageNumber is invalid")
	}
	return nil
}

func (x *UserRegisterReq) Check() error {
	if x.Users == nil {
		return errors.New("users are empty")
	}
	for _, u := range x.Users {
		switch {
		case u == nil:
			return errors.New("user is empty")
		case !datautil.IsLegalUserID(u.UserID):
			return errors.New("userID is legal")
		}
	}
	for _, u := range x.Users {
		if u.Nickname == "" {
			return errors.New("user name is empty")
		}
	}

	return nil
}

func (x *GetGlobalRecvMessageOptReq) Check() error {
	if x.UserID == "" {
		return errors.New("UserID is empty")
	}
	return nil
}

func (x *UserRegisterCountReq) Check() error {
	if x.Start <= 0 {
		return errors.New("start is invalid")
	}
	if x.End <= 0 {
		return errors.New("end is invalid")
	}
	return nil
}

func (x *SubscribeOrCancelUsersStatusReq) Check() error {
	if x.UserID == "" {
		return errors.New("UserID is empty")
	}
	if x.UserIDs == nil {
		return errors.New("subscription User-list is empty")
	}
	if x.Genre <= 0 || x.Genre >= 3 {
		return errors.New("invalid subscription type parameter")
	}
	return nil
}

func (x *GetUserStatusReq) Check() error {
	if len(x.UserIDs) == 0 {
		return errors.New("UserIDs is empty")
	}
	return nil
}

func (x *GetSubscribeUsersStatusReq) Check() error {
	if x.UserID == "" {
		return errors.New("UserID is empty")
	}
	return nil
}

func (x *ProcessUserCommandAddReq) Check() error {
	if x.UserID == "" {
		return errors.New("userID is empty")
	}
	if x.Type == 0 {
		return errors.New("type is not specified")
	}
	if x.Uuid == "" {
		return errors.New("UUID is empty")
	}
	return nil
}
func (x *ProcessUserCommandDeleteReq) Check() error {
	if x.UserID == "" {
		return errors.New("userID is empty")
	}
	if x.Type == 0 {
		return errors.New("type is not specified")
	}
	if x.Uuid == "" {
		return errors.New("UUID is empty")
	}
	return nil
}
func (x *ProcessUserCommandUpdateReq) Check() error {
	if x.UserID == "" {
		return errors.New("userID is empty")
	}
	if x.Type == 0 {
		return errors.New("type is not specified")
	}
	if x.Uuid == "" {
		return errors.New("UUID is empty")
	}
	return nil
}
func (x *ProcessUserCommandGetReq) Check() error {
	if x.UserID == "" {
		return errors.New("userID is empty")
	}
	if x.Type == 0 {
		return errors.New("type is not specified")
	}
	return nil
}
func (x *ProcessUserCommandGetAllReq) Check() error {
	if x.UserID == "" {
		return errors.New("userID is empty")
	}
	return nil
}
func (x *AddNotificationAccountReq) Check() error {
	if x.NickName == "" {
		return errors.New("nickName is empty")
	}
	//if x.FaceURL == "" {
	//	return errors.New("faceURL is empty")
	//}
	return nil
}

func (x *UpdateNotificationAccountInfoReq) Check() error {
	if x.UserID == "" {
		return errors.New("userID is empty")
	}

	if x.FaceURL == "" && x.NickName == "" {
		return errors.New("faceURL and nickName is empty at the same")
	}
	return nil
}

func (x *SearchNotificationAccountReq) Check() error {
	if x.Pagination == nil {
		return errors.New("pagination is empty")
	}
	if x.Pagination.PageNumber < 1 {
		return errors.New("pageNumber is invalid")
	}
	return nil
}

func (x *UpdateUserInfoExReq) Check() error {
	if x.UserInfo == nil {
		return errors.New("userInfo is empty")
	}
	if x.UserInfo.UserID == "" {
		return errors.New("userID is empty")
	}
	return nil
}

func (x *GetDesignateUsersResp) Format() any {
	if len(x.UsersInfo) > 20 {
		return fmt.Sprintf("len is %v", len(x.UsersInfo))
	}
	return x
}
