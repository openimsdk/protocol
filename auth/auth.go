package auth

import (
	"errors"

	"github.com/openimsdk/protocol/constant"
)

func (x *GetAdminTokenReq) Check() error {
	if x.UserID == "" {
		return errors.New("userID is empty")
	}
	return nil
}

func (x *ForceLogoutReq) Check() error {
	if x.UserID == "" {
		return errors.New("userID is empty")
	}
	if x.PlatformID > constant.AdminPlatformID || x.PlatformID < constant.IOSPlatformID {
		return errors.New("platformID is invalidate")
	}
	return nil
}

func (x *ParseTokenReq) Check() error {
	if x.Token == "" {
		return errors.New("userID is empty")
	}
	return nil
}

func (x *GetUserTokenReq) Check() error {
	if x.UserID == "" {
		errors.New("userID is empty")
	}

	if x.PlatformID == 0 {
		errors.New("platformID is empty")
	}
	return nil
}
