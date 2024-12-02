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
	"github.com/goBookMarker/internal/models"
	"github.com/goBookMarker/internal/ui/icons"
)

type HomePage struct {
	theme     *material.Theme
	state     *app.AppState
	searchBar widget.Editor
	addButton *widget.Clickable
	list      widget.List
}

func NewHomePage(th *material.Theme, state *app.AppState) *HomePage {
	return &HomePage{
		theme:     th,
		state:     state,
		addButton: new(widget.Clickable),
		searchBar: widget.Editor{
			SingleLine: true,
			Submit:     true,
		},
		list: widget.List{
			List: layout.List{Axis: layout.Vertical},
		},
	}
}

func (h *HomePage) Layout(gtx layout.Context) layout.Dimensions {
	if h.searchBar.Submit {
		h.state.Search(h.searchBar.Text())
		h.searchBar.Submit = false
	}

	if h.addButton.Clicked(gtx) {
		h.state.ShowAddBookmark()
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(16)).Layout(gtx,
				func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{}.Layout(gtx,
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							ed := material.Editor(h.theme, &h.searchBar, "Search bookmarks...")
							return ed.Layout(gtx)
						}),
						layout.Rigid(layout.Spacer{Width: unit.Dp(16)}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							btn := material.IconButton(h.theme, h.addButton, icons.AddIcon, "Add")
							btn.Size = unit.Dp(36)
							btn.Background = h.theme.ContrastBg
							btn.Color = h.theme.ContrastFg
							return btn.Layout(gtx)
						}),
					)
				},
			)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return h.layoutRecentBookmarks(gtx)
		}),
	)
}

func (h *HomePage) layoutRecentBookmarks(gtx layout.Context) layout.Dimensions {
	bookmarks := h.state.GetBookmarks()
	if len(bookmarks) == 0 {
		return h.layoutEmptyState(gtx)
	}

	return h.list.List.Layout(gtx, len(bookmarks),
		func(gtx layout.Context, index int) layout.Dimensions {
			return h.layoutBookmarkItem(gtx, &bookmarks[index])
		})
}

func (h *HomePage) layoutEmptyState(gtx layout.Context) layout.Dimensions {
	return layout.Center.Layout(gtx,
		func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					icon := material.IconButton(h.theme, h.addButton, icons.BookmarkIcon, "Add bookmark")
					icon.Color = h.theme.Fg
					icon.Size = unit.Dp(48)
					return icon.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					text := material.Body1(h.theme, "No bookmarks yet")
					text.Color = h.theme.Fg
					return text.Layout(gtx)
				}),
			)
		},
	)
}

func (h *HomePage) layoutBookmarkItem(gtx layout.Context, bookmark *models.Bookmark) layout.Dimensions {
	return layout.UniformInset(unit.Dp(16)).Layout(gtx,
		func(gtx layout.Context) layout.Dimensions {
			return layout.Stack{}.Layout(gtx,
				layout.Expanded(func(gtx layout.Context) layout.Dimensions {
					rect := image.Rectangle{
						Max: image.Point{
							X: gtx.Constraints.Max.X,
							Y: gtx.Constraints.Max.Y,
						},
					}
					rr := clip.RRect{
						Rect: rect,
						NE:   8, NW: 8, SE: 8, SW: 8,
					}
					paint.FillShape(gtx.Ops,
						h.theme.Bg,
						clip.Stroke{
							Path:  rr.Path(gtx.Ops),
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
									title := material.H6(h.theme, bookmark.Title)
									title.Color = h.theme.Fg
									return title.Layout(gtx)
								}),
								layout.Rigid(layout.Spacer{Height: unit.Dp(4)}.Layout),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									url := material.Body2(h.theme, bookmark.URL)
									url.Color = color.NRGBA{R: 128, G: 128, B: 128, A: 255}
									return url.Layout(gtx)
								}),
							)
						},
					)
				}),
			)
		},
	)
}
