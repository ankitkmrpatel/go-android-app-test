package ui

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip" // Ensure this is exactly this import
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/goBookMarker/internal/app"
)

type OnboardingPage struct {
	theme     *material.Theme
	state     *app.AppState
	pages     []onboardingStep
	current   int
	next      widget.Clickable
	skip      widget.Clickable
	googleBtn widget.Clickable
	msBtn     widget.Clickable
	skipBtn   widget.Clickable
}

type onboardingStep struct {
	title       string
	description string
	image       image.Image // You'll need to add actual images
}

func NewOnboardingPage(th *material.Theme, state *app.AppState) *OnboardingPage {
	return &OnboardingPage{
		theme: th,
		state: state,
		pages: []onboardingStep{
			{
				title:       "Welcome to GoBookMarker",
				description: "Your personal bookmark manager for all your links and images",
			},
			{
				title:       "Share & Save",
				description: "Share links and images directly from other apps to save them instantly",
			},
			{
				title:       "Sync Across Devices",
				description: "Optionally sync your bookmarks with Google Drive or OneDrive",
			},
		},
	}
}

func (p *OnboardingPage) Layout(gtx layout.Context) layout.Dimensions {
	// Handle navigation
	if p.next.Clicked(gtx) {
		if p.current < len(p.pages)-1 {
			p.current++
		} else {
			p.state.SetCurrentPage("auth")
		}
	}
	if p.skip.Clicked(gtx) {
		p.state.SetCurrentPage("home")
	}

	// Layout
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Flexed(0.7, func(gtx layout.Context) layout.Dimensions {
			return p.layoutContent(gtx)
		}),
		layout.Flexed(0.3, func(gtx layout.Context) layout.Dimensions {
			return p.layoutNavigation(gtx)
		}),
	)
}

func (p *OnboardingPage) layoutContent(gtx layout.Context) layout.Dimensions {
	page := p.pages[p.current]

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Top:    unit.Dp(32),
				Bottom: unit.Dp(16),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				title := material.H4(p.theme, page.title)
				title.Alignment = text.Middle
				title.Color = p.theme.Palette.ContrastBg
				return title.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Bottom: unit.Dp(32),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				desc := material.Body1(p.theme, page.description)
				desc.Alignment = text.Middle
				return desc.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if p.current == len(p.pages)-1 {
				return p.layoutAuthButtons(gtx)
			}
			return layout.Dimensions{}
		}),
	)
}

func (p *OnboardingPage) layoutAuthButtons(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			btn := material.Button(p.theme, &p.googleBtn, "Continue with Google")
			btn.Background = color.NRGBA{R: 66, G: 133, B: 244, A: 255}
			return layout.Inset{Bottom: unit.Dp(8)}.Layout(gtx, btn.Layout)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			btn := material.Button(p.theme, &p.msBtn, "Continue with Microsoft")
			btn.Background = color.NRGBA{R: 0, G: 120, B: 212, A: 255}
			return layout.Inset{Bottom: unit.Dp(16)}.Layout(gtx, btn.Layout)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			btn := material.Button(p.theme, &p.skipBtn, "Skip for now")
			btn.Color = p.theme.Palette.ContrastBg
			btn.Background = color.NRGBA{A: 0}
			return btn.Layout(gtx)
		}),
	)
}

func (p *OnboardingPage) layoutNavigation(gtx layout.Context) layout.Dimensions {
	return layout.Flex{
		Axis:      layout.Horizontal,
		Alignment: layout.Middle,
		Spacing:   layout.SpaceBetween,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if p.current < len(p.pages)-1 {
				skip := material.Button(p.theme, &p.skip, "Skip")
				skip.Background = color.NRGBA{A: 0}
				skip.Color = p.theme.Palette.ContrastBg
				return skip.Layout(gtx)
			}
			return layout.Dimensions{}
		}),
		layout.Rigid(p.layoutDots),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if p.current < len(p.pages)-1 {
				next := material.Button(p.theme, &p.next, "Next")
				return next.Layout(gtx)
			}
			return layout.Dimensions{}
		}),
	)
}

func (p *OnboardingPage) layoutDots(gtx layout.Context) layout.Dimensions {
	return layout.Flex{
		Axis:      layout.Horizontal,
		Alignment: layout.Middle,
		Spacing:   layout.SpaceEvenly,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return p.layoutDot(gtx, 0)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return p.layoutDot(gtx, 1)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return p.layoutDot(gtx, 2)
		}),
	)
}

func (p *OnboardingPage) layoutDot(gtx layout.Context, index int) layout.Dimensions {
	size := int(8)
	if index == p.current {
		size = 12
	}

	centerX, centerY := size/2, size/2 // Center of the circle

	defer clip.RRect{
		Rect: image.Rect(
			centerX-size/2, // Left
			centerY-size/2, // Top
			centerX+size/2, // Right
			centerY+size/2, // Bottom
		),
		NE: size / 2, // Equal radii for all corners
		NW: size / 2,
		SE: size / 2,
		SW: size / 2,
	}.Push(gtx.Ops).Pop()

	color := p.theme.Palette.ContrastBg
	if index == p.current {
		color = p.theme.Palette.ContrastFg
	}

	paint.ColorOp{Color: color}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	return layout.Dimensions{
		Size: image.Point{X: int(size), Y: int(size)},
	}
}
