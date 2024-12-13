package auth

import (
	"github.com/openimsdk/protocol/rpccall"
	"google.golang.org/grpc"
)

func InitAuth(conn *grpc.ClientConn) {
	GetAdminTokenCaller.SetConn(conn)
	GetUserTokenCaller.SetConn(conn)
	ForceLogoutCaller.SetConn(conn)
	ParseTokenCaller.SetConn(conn)
	InvalidateTokenCaller.SetConn(conn)
	KickTokensCaller.SetConn(conn)
}

var (
	GetAdminTokenCaller   = rpccall.NewRpcCaller[GetAdminTokenReq, GetAdminTokenResp](Auth_GetAdminToken_FullMethodName)
	GetUserTokenCaller    = rpccall.NewRpcCaller[GetUserTokenReq, GetUserTokenResp](Auth_GetUserToken_FullMethodName)
	ForceLogoutCaller     = rpccall.NewRpcCaller[ForceLogoutReq, ForceLogoutResp](Auth_ForceLogout_FullMethodName)
	ParseTokenCaller      = rpccall.NewRpcCaller[ParseTokenReq, ParseTokenResp](Auth_ParseToken_FullMethodName)
	InvalidateTokenCaller = rpccall.NewRpcCaller[InvalidateTokenReq, InvalidateTokenResp](Auth_InvalidateToken_FullMethodName)
	KickTokensCaller      = rpccall.NewRpcCaller[KickTokensReq, KickTokensResp](Auth_KickTokens_FullMethodName)
)
