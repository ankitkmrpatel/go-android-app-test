package components

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
)

type ColorPicker struct {
	colors   []color.NRGBA
	buttons  []widget.Clickable
	selected int
}

func NewColorPicker() *ColorPicker {
	defaultColors := []color.NRGBA{
		{R: 244, G: 67, B: 54, A: 255},   // Red
		{R: 233, G: 30, B: 99, A: 255},   // Pink
		{R: 156, G: 39, B: 176, A: 255},  // Purple
		{R: 63, G: 81, B: 181, A: 255},   // Indigo
		{R: 33, G: 150, B: 243, A: 255},  // Blue
		{R: 0, G: 188, B: 212, A: 255},   // Cyan
		{R: 0, G: 150, B: 136, A: 255},   // Teal
		{R: 76, G: 175, B: 80, A: 255},   // Green
		{R: 139, G: 195, B: 74, A: 255},  // Light Green
		{R: 205, G: 220, B: 57, A: 255},  // Lime
		{R: 255, G: 235, B: 59, A: 255},  // Yellow
		{R: 255, G: 193, B: 7, A: 255},   // Amber
		{R: 255, G: 152, B: 0, A: 255},   // Orange
		{R: 255, G: 87, B: 34, A: 255},   // Deep Orange
		{R: 121, G: 85, B: 72, A: 255},   // Brown
		{R: 158, G: 158, B: 158, A: 255}, // Grey
	}

	cp := &ColorPicker{
		colors:   defaultColors,
		buttons:  make([]widget.Clickable, len(defaultColors)),
		selected: 0,
	}
	return cp
}

func (cp *ColorPicker) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{}.Layout(gtx,
				cp.layoutColors()...,
			)
		}),
	)
}

func (cp *ColorPicker) layoutColors() []layout.FlexChild {
	var children []layout.FlexChild

	for i := range cp.colors {
		i := i
		children = append(children,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return cp.colorButton(gtx, i)
			}),
			layout.Rigid(layout.Spacer{Width: unit.Dp(4)}.Layout),
		)
	}

	return children
}

func (cp *ColorPicker) colorButton(gtx layout.Context, index int) layout.Dimensions {
	if cp.buttons[index].Clicked(gtx) {
		cp.selected = index
	}

	size := image.Point{X: gtx.Dp(24), Y: gtx.Dp(24)}
	return layout.Stack{Alignment: layout.Center}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			button := cp.buttons[index].Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				bounds := image.Rectangle{Max: size}

				paint.FillShape(gtx.Ops,
					cp.colors[index],
					clip.UniformRRect(bounds, 4).Op(gtx.Ops),
				)

				if index == cp.selected {
					paint.FillShape(gtx.Ops,
						color.NRGBA{A: 40},
						clip.UniformRRect(bounds, 4).Op(gtx.Ops),
					)
				}

				return layout.Dimensions{Size: size}
			})
			return button
		}),
	)
}

func (cp *ColorPicker) Selected() color.NRGBA {
	return cp.colors[cp.selected]
}

func (cp *ColorPicker) SetSelected(c color.NRGBA) {
	for i, col := range cp.colors {
		if col == c {
			cp.selected = i
			return
		}
	}
}
