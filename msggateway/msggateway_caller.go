package msggateway

import (
	"github.com/openimsdk/protocol/rpccall"
	"google.golang.org/grpc"
)

func InitMsgGateway(conn *grpc.ClientConn) {
	OnlinePushMsgCaller.SetConn(conn)
	GetUsersOnlineStatusCaller.SetConn(conn)
	OnlineBatchPushOneMsgCaller.SetConn(conn)
	SuperGroupOnlineBatchPushOneMsgCaller.SetConn(conn)
	KickUserOfflineCaller.SetConn(conn)
	MultiTerminalLoginCheckCaller.SetConn(conn)
}

var (
	OnlinePushMsgCaller                   = rpccall.NewRpcCaller[OnlinePushMsgReq, OnlinePushMsgResp](MsgGateway_OnlinePushMsg_FullMethodName)
	GetUsersOnlineStatusCaller            = rpccall.NewRpcCaller[GetUsersOnlineStatusReq, GetUsersOnlineStatusResp](MsgGateway_GetUsersOnlineStatus_FullMethodName)
	OnlineBatchPushOneMsgCaller           = rpccall.NewRpcCaller[OnlineBatchPushOneMsgReq, OnlineBatchPushOneMsgResp](MsgGateway_OnlineBatchPushOneMsg_FullMethodName)
	SuperGroupOnlineBatchPushOneMsgCaller = rpccall.NewRpcCaller[OnlineBatchPushOneMsgReq, OnlineBatchPushOneMsgResp](MsgGateway_SuperGroupOnlineBatchPushOneMsg_FullMethodName)
	KickUserOfflineCaller                 = rpccall.NewRpcCaller[KickUserOfflineReq, KickUserOfflineResp](MsgGateway_KickUserOffline_FullMethodName)
	MultiTerminalLoginCheckCaller         = rpccall.NewRpcCaller[MultiTerminalLoginCheckReq, MultiTerminalLoginCheckResp](MsgGateway_MultiTerminalLoginCheck_FullMethodName)
)
