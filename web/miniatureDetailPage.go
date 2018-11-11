package web

import (
	"fmt"
	"net/http"
	"strings"

	"jaredpearson.com/dbweb/data"
)

type MiniatureDetailPage struct {
	pageTitle       string
	userInfo        UserInfo
	ID              string
	Name            string
	Lineage         string
	Aspect          string
	SpawnCost       string
	AspectCost      string
	Power           string
	Defense         string
	Life            string
	Abilities       string
	FlavorText      string
	CollectorNumber string
	SetCode         string
	Set             string
	Rarity          string
	NextMiniURL     string
	PrevMiniURL     string
}

func (page MiniatureDetailPage) UserInfo() UserInfo {
	return page.userInfo
}
func (page MiniatureDetailPage) PageTitle() string {
	return page.pageTitle
}

func newMiniatureDetailPage(r *http.Request, miniature *data.Miniature) (page MiniatureDetailPage) {
	userInfo, _ := UserInfoFromRequest(r)

	page.pageTitle = miniature.Name()
	page.userInfo = userInfo
	page.ID = miniature.ID()
	page.Name = miniature.Name()
	page.Lineage = emptyToDash(miniature.Lineage())
	page.Aspect = emptyToDash(miniature.Aspect())
	page.SpawnCost = emptyToDash(miniature.SpawnCost())
	page.AspectCost = emptyToDash(miniature.AspectCost())
	page.Power = emptyToDash(miniature.Power())
	page.Defense = emptyToDash(miniature.Defense())
	page.Life = emptyToDash(miniature.Life())
	page.Abilities = emptyToDash(miniature.Abilities())
	page.FlavorText = emptyToDash(miniature.FlavorText())
	page.CollectorNumber = emptyToDash(miniature.CollectorNumber())
	page.Rarity = emptyToDash(miniature.Rarity())
	miniSet, setExists := miniature.Set()
	if setExists {
		page.SetCode = miniSet.ID()
		page.Set = miniSet.Name()
	} else {
		page.Set = "Unknown"
	}
	if len(miniature.NextMiniID()) > 0 {
		page.NextMiniURL = fmt.Sprintf("/miniature/" + miniature.NextMiniID())
	}
	if len(miniature.PrevMiniID()) > 0 {
		page.PrevMiniURL = fmt.Sprintf("/miniature/" + miniature.PrevMiniID())
	}
	return
}

func emptyToDash(value string) string {
	if len(strings.TrimSpace(value)) == 0 {
		return "-"
	}
	return strings.TrimSpace(value)
}

func ShowMiniatureDetailPage(w http.ResponseWriter, r *http.Request) {
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

	m, err := data.GetMiniatureByID(miniID)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	pageModel := newMiniatureDetailPage(r, m)

	ShowTemplateInMainLayout(w, r, "miniDetail", pageModel)
}
