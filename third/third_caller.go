package third

import (
	"github.com/openimsdk/protocol/rpccall"
	"google.golang.org/grpc"
)

func InitThird(conn *grpc.ClientConn) {
	PartLimitCaller.SetConn(conn)
	PartSizeCaller.SetConn(conn)
	InitiateMultipartUploadCaller.SetConn(conn)
	AuthSignCaller.SetConn(conn)
	CompleteMultipartUploadCaller.SetConn(conn)
	AccessURLCaller.SetConn(conn)
	InitiateFormDataCaller.SetConn(conn)
	CompleteFormDataCaller.SetConn(conn)
	DeleteOutdatedDataCaller.SetConn(conn)
	FcmUpdateTokenCaller.SetConn(conn)
	SetAppBadgeCaller.SetConn(conn)
	UploadLogsCaller.SetConn(conn)
	DeleteLogsCaller.SetConn(conn)
	SearchLogsCaller.SetConn(conn)
}

var (
	PartLimitCaller               = rpccall.NewRpcCaller[PartLimitReq, PartLimitResp](Third_PartLimit_FullMethodName)
	PartSizeCaller                = rpccall.NewRpcCaller[PartSizeReq, PartSizeResp](Third_PartSize_FullMethodName)
	InitiateMultipartUploadCaller = rpccall.NewRpcCaller[InitiateMultipartUploadReq, InitiateMultipartUploadResp](Third_InitiateMultipartUpload_FullMethodName)
	AuthSignCaller                = rpccall.NewRpcCaller[AuthSignReq, AuthSignResp](Third_AuthSign_FullMethodName)
	CompleteMultipartUploadCaller = rpccall.NewRpcCaller[CompleteMultipartUploadReq, CompleteMultipartUploadResp](Third_CompleteMultipartUpload_FullMethodName)
	AccessURLCaller               = rpccall.NewRpcCaller[AccessURLReq, AccessURLResp](Third_AccessURL_FullMethodName)
	InitiateFormDataCaller        = rpccall.NewRpcCaller[InitiateFormDataReq, InitiateFormDataResp](Third_InitiateFormData_FullMethodName)
	CompleteFormDataCaller        = rpccall.NewRpcCaller[CompleteFormDataReq, CompleteFormDataResp](Third_CompleteFormData_FullMethodName)
	DeleteOutdatedDataCaller      = rpccall.NewRpcCaller[DeleteOutdatedDataReq, DeleteOutdatedDataResp](Third_DeleteOutdatedData_FullMethodName)
	FcmUpdateTokenCaller          = rpccall.NewRpcCaller[FcmUpdateTokenReq, FcmUpdateTokenResp](Third_FcmUpdateToken_FullMethodName)
	SetAppBadgeCaller             = rpccall.NewRpcCaller[SetAppBadgeReq, SetAppBadgeResp](Third_SetAppBadge_FullMethodName)
	UploadLogsCaller              = rpccall.NewRpcCaller[UploadLogsReq, UploadLogsResp](Third_UploadLogs_FullMethodName)
	DeleteLogsCaller              = rpccall.NewRpcCaller[DeleteLogsReq, DeleteLogsResp](Third_DeleteLogs_FullMethodName)
	SearchLogsCaller              = rpccall.NewRpcCaller[SearchLogsReq, SearchLogsResp](Third_SearchLogs_FullMethodName)
)
