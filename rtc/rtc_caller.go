package rtc

import (
	"github.com/openimsdk/protocol/rpccall"
	"google.golang.org/grpc"
)

func InitRtcService(conn *grpc.ClientConn) {
	SignalMessageAssembleCaller.SetConn(conn)
	SignalGetRoomByGroupIDCaller.SetConn(conn)
	SignalGetTokenByRoomIDCaller.SetConn(conn)
	SignalGetRoomsCaller.SetConn(conn)
	GetSignalInvitationInfoCaller.SetConn(conn)
	GetSignalInvitationInfoStartAppCaller.SetConn(conn)
	SignalSendCustomSignalCaller.SetConn(conn)
	GetSignalInvitationRecordsCaller.SetConn(conn)
	DeleteSignalRecordsCaller.SetConn(conn)
}

var (
	SignalMessageAssembleCaller           = rpccall.NewRpcCaller[SignalMessageAssembleReq, SignalMessageAssembleResp](RtcService_SignalMessageAssemble_FullMethodName)
	SignalGetRoomByGroupIDCaller          = rpccall.NewRpcCaller[SignalGetRoomByGroupIDReq, SignalGetRoomByGroupIDResp](RtcService_SignalGetRoomByGroupID_FullMethodName)
	SignalGetTokenByRoomIDCaller          = rpccall.NewRpcCaller[SignalGetTokenByRoomIDReq, SignalGetTokenByRoomIDResp](RtcService_SignalGetTokenByRoomID_FullMethodName)
	SignalGetRoomsCaller                  = rpccall.NewRpcCaller[SignalGetRoomsReq, SignalGetRoomsResp](RtcService_SignalGetRooms_FullMethodName)
	GetSignalInvitationInfoCaller         = rpccall.NewRpcCaller[GetSignalInvitationInfoReq, GetSignalInvitationInfoResp](RtcService_GetSignalInvitationInfo_FullMethodName)
	GetSignalInvitationInfoStartAppCaller = rpccall.NewRpcCaller[GetSignalInvitationInfoStartAppReq, GetSignalInvitationInfoStartAppResp](RtcService_GetSignalInvitationInfoStartApp_FullMethodName)
	SignalSendCustomSignalCaller          = rpccall.NewRpcCaller[SignalSendCustomSignalReq, SignalSendCustomSignalResp](RtcService_SignalSendCustomSignal_FullMethodName)
	GetSignalInvitationRecordsCaller      = rpccall.NewRpcCaller[GetSignalInvitationRecordsReq, GetSignalInvitationRecordsResp](RtcService_GetSignalInvitationRecords_FullMethodName)
	DeleteSignalRecordsCaller             = rpccall.NewRpcCaller[DeleteSignalRecordsReq, DeleteSignalRecordsResp](RtcService_DeleteSignalRecords_FullMethodName)
)
