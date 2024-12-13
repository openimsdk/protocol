package push

import (
	"github.com/openimsdk/protocol/rpccall"
	"google.golang.org/grpc"
)

func InitPushMsgService(conn *grpc.ClientConn) {
	PushMsgCaller.SetConn(conn)
	DelUserPushTokenCaller.SetConn(conn)
}

var (
	PushMsgCaller          = rpccall.NewRpcCaller[PushMsgReq, PushMsgResp](PushMsgService_PushMsg_FullMethodName)
	DelUserPushTokenCaller = rpccall.NewRpcCaller[DelUserPushTokenReq, DelUserPushTokenResp](PushMsgService_DelUserPushToken_FullMethodName)
)
