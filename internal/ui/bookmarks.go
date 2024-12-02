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

type BookmarksPage struct {
	theme           *material.Theme
	state           *app.AppState
	list            widget.List
	searchBar       widget.Editor
	bookmarkActions map[string]*BookmarkActions
}

type BookmarkActions struct {
	favorite *widget.Clickable
	edit     *widget.Clickable
	delete   *widget.Clickable
	share    *widget.Clickable
}

func NewBookmarksPage(th *material.Theme, state *app.AppState) *BookmarksPage {
	return &BookmarksPage{
		theme: th,
		state: state,
		list: widget.List{
			List: layout.List{Axis: layout.Vertical},
		},
		searchBar: widget.Editor{
			SingleLine: true,
			Submit:     true,
		},
		bookmarkActions: make(map[string]*BookmarkActions),
	}
}

func (p *BookmarksPage) getBookmarkActions(id string) *BookmarkActions {
	if actions, exists := p.bookmarkActions[id]; exists {
		return actions
	}
	actions := &BookmarkActions{
		favorite: new(widget.Clickable),
		edit:     new(widget.Clickable),
		delete:   new(widget.Clickable),
		share:    new(widget.Clickable),
	}
	p.bookmarkActions[id] = actions
	return actions
}

func (p *BookmarksPage) Layout(gtx layout.Context) layout.Dimensions {
	if p.searchBar.Submit {
		p.state.Search(p.searchBar.Text())
		p.searchBar.Submit = false
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(16)).Layout(gtx,
				func(gtx layout.Context) layout.Dimensions {
					ed := material.Editor(p.theme, &p.searchBar, "Search bookmarks...")
					return ed.Layout(gtx)
				},
			)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return p.layoutBookmarks(gtx)
		}),
	)
}

func (p *BookmarksPage) layoutBookmarks(gtx layout.Context) layout.Dimensions {
	bookmarks := p.state.GetBookmarks()
	if len(bookmarks) == 0 {
		return p.layoutEmptyState(gtx)
	}

	return p.list.List.Layout(gtx, len(bookmarks),
		func(gtx layout.Context, index int) layout.Dimensions {
			return p.layoutBookmarkItem(gtx, &bookmarks[index])
		})
}

func (p *BookmarksPage) layoutEmptyState(gtx layout.Context) layout.Dimensions {
	return layout.Center.Layout(gtx,
		func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					icon := material.IconButton(p.theme, nil, icons.AddIcon, "Add bookmark") //material.Icon(icons.BookmarkIcon)
					icon.Color = p.theme.Fg
					icon.Size = unit.Dp(48)
					return icon.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					text := material.Body1(p.theme, "No bookmarks found")
					text.Color = p.theme.Fg
					return text.Layout(gtx)
				}),
			)
		},
	)
}

func (p *BookmarksPage) layoutBookmarkItem(gtx layout.Context, bookmark *models.Bookmark) layout.Dimensions {
	actions := p.getBookmarkActions(bookmark.ID)

	// Handle actions
	if actions.favorite.Clicked(gtx) {
		bookmark.IsFavorite = !bookmark.IsFavorite
		p.state.SaveBookmark(bookmark)
	}
	if actions.edit.Clicked(gtx) {
		p.state.EditBookmark(bookmark)
	}
	if actions.delete.Clicked(gtx) {
		p.state.DeleteBookmark(bookmark)
	}
	if actions.share.Clicked(gtx) {
		p.state.ShareBookmark(bookmark)
	}

	return layout.UniformInset(unit.Dp(16)).Layout(gtx,
		func(gtx layout.Context) layout.Dimensions {
			return layout.Stack{}.Layout(gtx,
				layout.Expanded(func(gtx layout.Context) layout.Dimensions {
					rect := image.Rectangle{
						Min: image.Point{X: 0, Y: 0},
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
									return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
										layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
											title := material.H6(p.theme, bookmark.Title)
											title.Color = p.theme.Fg
											return title.Layout(gtx)
										}),
										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											return p.layoutActions(gtx, actions, bookmark)
										}),
									)
								}),
								layout.Rigid(layout.Spacer{Height: unit.Dp(4)}.Layout),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									url := material.Body2(p.theme, bookmark.URL)
									url.Color = color.NRGBA{R: 128, G: 128, B: 128, A: 255}
									return url.Layout(gtx)
								}),
								layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return p.layoutTags(gtx, bookmark.Tags)
								}),
							)
						},
					)
				}),
			)
		},
	)
}

func (p *BookmarksPage) layoutActions(gtx layout.Context, actions *BookmarkActions, bookmark *models.Bookmark) layout.Dimensions {
	return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			btn := material.IconButton(p.theme, actions.favorite, icons.FavoriteIcon, "Favorite")
			if bookmark.IsFavorite {
				btn.Color = color.NRGBA{R: 255, G: 0, B: 0, A: 255}
			}
			return btn.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.IconButton(p.theme, actions.edit, icons.EditIcon, "Edit").Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.IconButton(p.theme, actions.delete, icons.DeleteIcon, "Delete").Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.IconButton(p.theme, actions.share, icons.ShareIcon, "Share").Layout(gtx)
		}),
	)
}

func (p *BookmarksPage) layoutTags(gtx layout.Context, tags []string) layout.Dimensions {
	if len(tags) == 0 {
		return layout.Dimensions{}
	}

	return layout.Flex{Alignment: layout.Start}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			var children []layout.FlexChild
			for _, tag := range tags {
				tag := tag // capture for closure
				children = append(children,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return p.tagChip(gtx, tag)
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(4)}.Layout),
				)
			}
			return layout.Flex{}.Layout(gtx, children...)
		}),
	)
}

func (p *BookmarksPage) tagChip(gtx layout.Context, tag string) layout.Dimensions {
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			rect := image.Rectangle{
				Min: image.Point{X: 0, Y: 0},
				Max: image.Point{
					X: gtx.Constraints.Max.X,
					Y: gtx.Constraints.Max.Y,
				},
			}
			paint.FillShape(gtx.Ops,
				color.NRGBA{R: 230, G: 230, B: 230, A: 255},
				clip.Stroke{
					Path:  clip.RRect{Rect: rect, NE: 16, NW: 16, SE: 16, SW: 16}.Path(gtx.Ops),
					Width: float32(unit.Dp(1)),
				}.Op(),
			)
			return layout.Dimensions{Size: gtx.Constraints.Max}
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(6)).Layout(gtx,
				func(gtx layout.Context) layout.Dimensions {
					label := material.Body2(p.theme, tag)
					label.Color = p.theme.Fg
					return label.Layout(gtx)
				},
			)
		}),
	)
}
