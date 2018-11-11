package web

import (
	"fmt"
	"net/http"
	"strings"

	"jaredpearson.com/dbweb/data"
)

type SetDetailPage struct {
	pageTitle  string
	userInfo   UserInfo
	name       string
	miniatures []*data.Miniature
}

func (setDetailPage SetDetailPage) PageTitle() string {
	return setDetailPage.pageTitle
}
func (setDetailPage SetDetailPage) UserInfo() UserInfo {
	return setDetailPage.userInfo
}
func (setDetailPage SetDetailPage) Name() string {
	return setDetailPage.name
}
func (setDetailPage SetDetailPage) Miniatures() []*data.Miniature {
	return setDetailPage.miniatures
}

func newSetDetailPage(r *http.Request, set data.MiniatureSet, minis []*data.Miniature) MainLayoutData {
	userInfo, _ := UserInfoFromRequest(r)
	return SetDetailPage{
		pageTitle:  fmt.Sprintf("%s Set", set.Name()),
		userInfo:   userInfo,
		name:       set.Name(),
		miniatures: minis,
	}
}

func ShowSetDetailPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.NotFound(w, r)
		return
	}

	pathParts := strings.Split(r.URL.Path[1:], "/")
	if len(pathParts) <= 1 {
		http.NotFound(w, r)
		return
	}

	// assume that the second part of the path is the mini ID
	miniID := pathParts[1]

	s, exists := data.GetMiniatureSetByID(miniID)
	if exists != nil {
		http.NotFound(w, r)
		return
	}

	minis, err := data.GetMiniaturesBySet(s.ID())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pageModel := newSetDetailPage(r, *s, minis)

	ShowTemplateInMainLayout(w, r, "setDetail", pageModel)
}
