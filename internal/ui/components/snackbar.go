package components

import (
	"image"
	"image/color"
	"time"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Snackbar struct {
	theme    *material.Theme
	message  string
	visible  bool
	dismiss  widget.Clickable
	showTime time.Time
}

func NewSnackbar(th *material.Theme) *Snackbar {
	return &Snackbar{
		theme: th,
	}
}

func (s *Snackbar) Show(message string) {
	s.message = message
	s.visible = true
	s.showTime = time.Now()
}

func (s *Snackbar) ShowMessage(message string) {
	s.Show(message)
}

func (s *Snackbar) ShowError(message string) {
	s.Show(message)
}

func (s *Snackbar) Layout(gtx layout.Context) layout.Dimensions {
	if !s.visible {
		return layout.Dimensions{}
	}

	// Auto-hide after 3 seconds
	if time.Since(s.showTime) > 3*time.Second {
		s.visible = false
		return layout.Dimensions{}
	}

	// Dismiss on click
	if s.dismiss.Clicked(gtx) {
		s.visible = false
		return layout.Dimensions{}
	}

	return layout.Stack{Alignment: layout.S}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(16)).Layout(gtx,
				func(gtx layout.Context) layout.Dimensions {
					return layout.Stack{}.Layout(gtx,
						layout.Expanded(func(gtx layout.Context) layout.Dimensions {
							bounds := image.Rectangle{
								Max: gtx.Constraints.Min,
							}

							paint.FillShape(gtx.Ops,
								color.NRGBA{R: 50, G: 50, B: 50, A: 230},
								clip.UniformRRect(bounds, 4).Op(gtx.Ops),
							)

							return layout.Dimensions{Size: bounds.Max}
						}),
						layout.Stacked(func(gtx layout.Context) layout.Dimensions {
							return layout.UniformInset(unit.Dp(12)).Layout(gtx,
								func(gtx layout.Context) layout.Dimensions {
									button := s.dismiss.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
										label := material.Body1(s.theme, s.message)
										label.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
										label.Alignment = text.Middle
										return label.Layout(gtx)
									})
									return button
								},
							)
						}),
					)
				},
			)
		}),
	)
}
