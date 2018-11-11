package web

import (
	"net/http"
)

type UserInfo struct {
	Username string
}

func UserInfoFromRequest(r *http.Request) (UserInfo, bool) {
	var userInfo UserInfo
	if r.Context().Value(AuthUserToken) != nil {
		userInfo = r.Context().Value(AuthUserToken).(UserInfo)
		return userInfo, true
	}
	return UserInfo{}, false
}
