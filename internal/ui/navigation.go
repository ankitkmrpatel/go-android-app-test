package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/goBookMarker/internal/app"
	"github.com/goBookMarker/internal/ui/icons"
)

type NavigationPage struct {
	theme    *material.Theme
	state    *app.AppState
	list     widget.List
	navItems []NavItem
	selected int
}

type NavItem struct {
	Label    string
	Icon     *widget.Icon
	Click    widget.Clickable
	Selected bool
}

func NewNavigationPage(th *material.Theme, state *app.AppState) *NavigationPage {
	return &NavigationPage{
		theme: th,
		state: state,
		list: widget.List{
			List: layout.List{Axis: layout.Vertical},
		},
	}
}

func (p *NavigationPage) Layout(gtx layout.Context) layout.Dimensions {
	user := p.state.CurrentUser()
	if user == nil {
		return layout.Dimensions{}
	}

	// Convert string position to int
	p.selected = 0 // Default to first item
	if user.NavPosition == "home" {
		p.selected = 0
	} else if user.NavPosition == "bookmarks" {
		p.selected = 1
	} else if user.NavPosition == "settings" {
		p.selected = 2
	}

	// Update nav items based on user preferences
	p.navItems = []NavItem{
		{Label: "Home", Icon: icons.HomeIcon},
		{Label: "Bookmarks", Icon: icons.BookmarkIcon},
		{Label: "Settings", Icon: icons.SettingsIcon},
	}

	return layout.UniformInset(unit.Dp(8)).Layout(gtx,
		func(gtx layout.Context) layout.Dimensions {
			return p.list.List.Layout(gtx, len(p.navItems),
				func(gtx layout.Context, index int) layout.Dimensions {
					return p.layoutNavItem(gtx, &p.navItems[index], index)
				})
		},
	)
}

func (p *NavigationPage) layoutNavItem(gtx layout.Context, item *NavItem, index int) layout.Dimensions {
	if item.Click.Clicked(gtx) {
		p.selected = index
		switch index {
		case 0:
			p.state.SetCurrentPage("home")
		case 1:
			p.state.SetCurrentPage("bookmarks")
		case 2:
			p.state.SetCurrentPage("settings")
		}
	}

	item.Selected = index == p.selected
	return layout.UniformInset(unit.Dp(8)).Layout(gtx,
		func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					iconButton := material.IconButton(p.theme, &item.Click, item.Icon, item.Label)
					if item.Selected {
						iconButton.Color = p.theme.ContrastBg
					}
					return iconButton.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Width: unit.Dp(16)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					label := material.Body1(p.theme, item.Label)
					if item.Selected {
						label.Color = p.theme.ContrastBg
					}
					return label.Layout(gtx)
				}),
			)
		},
	)
}
