package ui

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/goBookMarker/internal/app"
)

type SettingsPage struct {
	theme               *material.Theme
	state               *app.AppState
	list                widget.List
	syncEnabled         widget.Bool
	darkMode            widget.Bool
	saveButton          widget.Clickable
	logoutButton        widget.Clickable
	previousSyncEnabled bool
}

func NewSettingsPage(th *material.Theme, state *app.AppState) *SettingsPage {
	return &SettingsPage{
		theme: th,
		state: state,
		list: widget.List{
			List: layout.List{Axis: layout.Vertical},
		},
	}
}

func (p *SettingsPage) Layout(gtx layout.Context) layout.Dimensions {
	user := p.state.CurrentUser()
	if user == nil {
		return layout.Dimensions{}
	}

	// Handle sync toggle
	if p.syncEnabled.Value != p.previousSyncEnabled {
		user.SyncEnabled = p.syncEnabled.Value
		p.state.SaveUser(user)
		p.previousSyncEnabled = p.syncEnabled.Value
	}

	// Handle dark mode toggle
	if p.darkMode.Value != (user.Theme == "dark") {
		if p.darkMode.Value {
			user.Theme = "dark"
		} else {
			user.Theme = "light"
		}
		p.state.SaveUser(user)
	}

	// Handle save button
	if p.saveButton.Clicked(gtx) {
		p.state.SaveUser(user)
	}

	// Handle logout button
	if p.logoutButton.Clicked(gtx) {
		p.state.Logout()
	}

	return layout.UniformInset(unit.Dp(16)).Layout(gtx,
		func(gtx layout.Context) layout.Dimensions {
			return p.list.List.Layout(gtx, 1,
				func(gtx layout.Context, index int) layout.Dimensions {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							title := material.H6(p.theme, "Settings")
							return title.Layout(gtx)
						}),
						layout.Rigid(layout.Spacer{Height: unit.Dp(24)}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return p.layoutSection(gtx, "Account", func(gtx layout.Context) layout.Dimensions {
								return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
											layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
												label := material.Body1(p.theme, "Email")
												return label.Layout(gtx)
											}),
											layout.Rigid(func(gtx layout.Context) layout.Dimensions {
												value := material.Body1(p.theme, user.Email)
												value.Color = color.NRGBA{R: 128, G: 128, B: 128, A: 255}
												return value.Layout(gtx)
											}),
										)
									}),
								)
							})
						}),
						layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return p.layoutSection(gtx, "Preferences", func(gtx layout.Context) layout.Dimensions {
								return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										sw := material.Switch(p.theme, &p.darkMode, "Dark Mode")
										p.darkMode.Value = user.Theme == "dark"
										return sw.Layout(gtx)
									}),
									layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										sw := material.Switch(p.theme, &p.syncEnabled, "Enable Sync")
										p.syncEnabled.Value = user.SyncEnabled
										return sw.Layout(gtx)
									}),
								)
							})
						}),
						layout.Rigid(layout.Spacer{Height: unit.Dp(24)}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{}.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									btn := material.Button(p.theme, &p.saveButton, "Save Changes")
									return btn.Layout(gtx)
								}),
								layout.Rigid(layout.Spacer{Width: unit.Dp(16)}.Layout),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									btn := material.Button(p.theme, &p.logoutButton, "Logout")
									btn.Background = color.NRGBA{R: 255, G: 0, B: 0, A: 255}
									return btn.Layout(gtx)
								}),
							)
						}),
					)
				},
			)
		},
	)
}

func (p *SettingsPage) layoutSection(gtx layout.Context, title string, content func(gtx layout.Context) layout.Dimensions) layout.Dimensions {
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			rect := image.Rectangle{
				Max: image.Point{
					X: gtx.Constraints.Max.X,
					Y: gtx.Constraints.Max.Y,
				},
			}
			paint.FillShape(gtx.Ops,
				p.theme.Bg,
				clip.Stroke{
					Path:  clip.RRect{Rect: rect, NE: 8, NW: 8, SE: 8, SW: 8}.Path(gtx.Ops),
					Width: float32(unit.Dp(1)),
				}.Op(),
			)
			return layout.Dimensions{Size: gtx.Constraints.Max}
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(16)).Layout(gtx,
				func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							label := material.H6(p.theme, title)
							return label.Layout(gtx)
						}),
						layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
						layout.Rigid(content),
					)
				},
			)
		}),
	)
}
