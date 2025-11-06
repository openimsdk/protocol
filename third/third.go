package third

import (
	"errors"

	"github.com/openimsdk/protocol/constant"
)

func (x *FcmUpdateTokenReq) Check() error {
	if x.PlatformID > constant.AdminPlatformID || x.PlatformID < constant.IOSPlatformID {
		return errors.New("platformID is invalidate")
	}
	if x.FcmToken == "" {
		return errors.New("FcmToken is empty")
	}
	if x.Account == "" {
		return errors.New("account is empty")
	}
	return nil
}

func (x *SetAppBadgeReq) Check() error {
	if x.UserID == "" {
		return errors.New("UserID is empty")
	}
	return nil
}

func (x *InitiateMultipartUploadReq) Check() error {
	if x.UrlPrefix == "" {
		return errors.New("UrlPrefix is empty")
	}
	return nil
}

func (x *CompleteMultipartUploadReq) Check() error {
	if x.UrlPrefix == "" {
		return errors.New("UrlPrefix is empty")
	}
	return nil
}

func (x *CompleteFormDataReq) Check() error {
	if x.UrlPrefix == "" {
		return errors.New("UrlPrefix is empty")
	}
	return nil
}

func (x *DeleteOutdatedDataReq) Check() error {
	if x.Limit <= 0 {
		return errors.New("limit must be greater than 0")
	}
	if len(x.ObjectGroup) == 0 {
		return errors.New("ObjectGroup is empty")
	}
	return nil
}
