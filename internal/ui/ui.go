package ui

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/goBookMarker/internal/app"
)

type UI struct {
	theme     *material.Theme
	state     *app.AppState
	nav       *NavigationPage
	home      *HomePage
	bookmarks *BookmarksPage
	tags      *TagsPage
	settings  *SettingsPage
}

type navItem struct {
	name  string
	icon  *widget.Icon
	click *widget.Clickable
}

func NewUI(th *material.Theme, state *app.AppState) *UI {
	ui := &UI{
		theme: th,
		state: state,
	}

	// Initialize navigation and pages
	ui.nav = NewNavigationPage(th, state)
	ui.home = NewHomePage(th, state)
	ui.bookmarks = NewBookmarksPage(th, state)
	ui.tags = NewTagsPage(th, state)
	ui.settings = NewSettingsPage(th, state)

	return ui
}

func (ui *UI) Layout(gtx layout.Context) layout.Dimensions {
	// Determine if navigation is at top or bottom
	isTopNav := ui.state.CurrentUser() != nil && ui.state.CurrentUser().NavPosition == "top"

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if isTopNav {
				return ui.nav.Layout(gtx)
			}
			return layout.Dimensions{}
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return ui.layoutCurrentPage(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if !isTopNav {
				return ui.nav.Layout(gtx)
			}
			return layout.Dimensions{}
		}),
	)
}

func (ui *UI) layoutCurrentPage(gtx layout.Context) layout.Dimensions {
	switch ui.state.GetCurrentPage() {
	case "home":
		return ui.home.Layout(gtx)
	case "bookmarks":
		return ui.bookmarks.Layout(gtx)
	case "tags":
		return ui.tags.Layout(gtx)
	case "settings":
		return ui.settings.Layout(gtx)
	default:
		return ui.home.Layout(gtx)
	}
}
