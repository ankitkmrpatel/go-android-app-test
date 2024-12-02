package ui

import (
	"fmt"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/goBookMarker/internal/app"
	"github.com/goBookMarker/internal/models"
	"github.com/goBookMarker/internal/ui/components"
	"github.com/goBookMarker/internal/ui/icons"
)

type TagsPage struct {
	theme     *material.Theme
	state     *app.AppState
	list      widget.List
	searchBar widget.Editor
	toast     *components.Snackbar

	// Tag management
	addTag struct {
		button      widget.Clickable
		name        component.TextField
		color       *components.ColorPicker
		description component.TextField
		parentID    string
		visible     bool
	}
	editTag struct {
		tag         *models.Tag
		name        component.TextField
		color       *components.ColorPicker
		description component.TextField
		parentID    string
		visible     bool
		save        widget.Clickable
		cancel      widget.Clickable
	}
	deleteConfirm struct {
		tag     *models.Tag
		visible bool
		confirm widget.Clickable
		cancel  widget.Clickable
	}

	// Tag groups
	addGroup struct {
		button  widget.Clickable
		name    component.TextField
		visible bool
	}
	editGroup struct {
		group   *models.TagGroup
		name    component.TextField
		visible bool
		save    widget.Clickable
		cancel  widget.Clickable
	}

	// Drag and drop
	draggedTag   *models.Tag
	draggedGroup *models.TagGroup
	dropTarget   struct {
		tagID    string
		groupID  string
		position string // "before", "after", "into"
	}

	// Batch operations
	selectedTags map[string]bool
	batchOps     struct {
		visible bool
		delete  widget.Clickable
		group   widget.Clickable
		export  widget.Clickable
	}

	// Import/Export
	importBtn widget.Clickable
	exportBtn widget.Clickable

	filteredTags []models.Tag
	tagGroups    []models.TagGroup

	// ...
}

func NewTagsPage(th *material.Theme, state *app.AppState) *TagsPage {
	tp := &TagsPage{
		theme: th,
		state: state,
		list:  widget.List{List: layout.List{Axis: layout.Vertical}},
		searchBar: widget.Editor{
			SingleLine: true,
			Submit:     true,
		},
		selectedTags: make(map[string]bool),
		toast:        components.NewSnackbar(th),
	}
	tp.filteredTags = state.GetTags()
	tp.tagGroups = state.GetTagGroups()
	return tp
}

func (tp *TagsPage) Layout(gtx layout.Context) layout.Dimensions {
	// Handle search
	// for _, e := range tp.searchBar.Events() {
	// 	if _, ok := e.(widget.SubmitEvent); ok {
	// 		tp.filteredTags = tp.state.SearchTags(tp.searchBar.Text())
	// 	}
	// }

	if tp.searchBar.Submit {
		tp.filteredTags = tp.state.SearchTags(tp.searchBar.Text())
	}

	// Handle import/export
	if tp.importBtn.Clicked(gtx) {
		go tp.handleImport()
	}
	if tp.exportBtn.Clicked(gtx) {
		go tp.handleExport()
	}

	// Handle batch operations
	if tp.batchOps.visible {
		if tp.batchOps.delete.Clicked(gtx) {
			tp.handleBatchDelete()
		}
		if tp.batchOps.group.Clicked(gtx) {
			tp.handleBatchGroup()
		}
		if tp.batchOps.export.Clicked(gtx) {
			tp.handleBatchExport()
		}
	}

	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return tp.layoutToolbar(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if len(tp.selectedTags) > 0 {
						return tp.layoutBatchOperations(gtx)
					}
					return layout.Dimensions{}
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if tp.addTag.visible {
						return tp.layoutAddTag(gtx)
					}
					return layout.Dimensions{}
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if tp.editTag.visible {
						return tp.layoutEditTag(gtx)
					}
					return layout.Dimensions{}
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if tp.deleteConfirm.visible {
						return tp.layoutDeleteConfirm(gtx)
					}
					return layout.Dimensions{}
				}),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					if len(tp.filteredTags) == 0 && len(tp.tagGroups) == 0 {
						return tp.layoutEmptyState(gtx)
					}
					return tp.layoutTagsAndGroups(gtx)
				}),
			)
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return tp.toast.Layout(gtx)
		}),
	)
}

func (tp *TagsPage) layoutTagsAndGroups(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			// Render tag groups
			return tp.layoutTagGroups(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			// Render individual tags
			return tp.layoutIndividualTags(gtx)
		}),
	)
}

func (tp *TagsPage) layoutTagGroups(gtx layout.Context) layout.Dimensions {
	if len(tp.tagGroups) == 0 {
		return layout.Dimensions{}
	}

	list := &layout.List{
		Axis: layout.Vertical,
	}
	return list.Layout(gtx, len(tp.tagGroups), func(gtx layout.Context, index int) layout.Dimensions {
		group := tp.tagGroups[index]
		return layout.Inset{
			Top:    unit.Dp(8),
			Bottom: unit.Dp(8),
			Left:   unit.Dp(16),
			Right:  unit.Dp(16),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return material.Body1(tp.theme, group.Name).Layout(gtx)
		})
	})
}

func (tp *TagsPage) layoutIndividualTags(gtx layout.Context) layout.Dimensions {
	if len(tp.filteredTags) == 0 {
		return layout.Dimensions{}
	}

	list := &layout.List{
		Axis: layout.Vertical,
	}
	return list.Layout(gtx, len(tp.filteredTags), func(gtx layout.Context, index int) layout.Dimensions {
		tag := tp.filteredTags[index]
		return layout.Inset{
			Top:    unit.Dp(8),
			Bottom: unit.Dp(8),
			Left:   unit.Dp(16),
			Right:  unit.Dp(16),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return material.Body1(tp.theme, tag.Name).Layout(gtx)
		})
	})
}

func (tp *TagsPage) layoutToolbar(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			ed := material.Editor(tp.theme, &tp.searchBar, "Search tags...")
			return ed.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			btn := material.IconButton(tp.theme, &tp.importBtn, icons.ImportIcon, "Import")
			return btn.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			btn := material.IconButton(tp.theme, &tp.exportBtn, icons.ExportIcon, "Export")
			return btn.Layout(gtx)
		}),
	)
}

func (tp *TagsPage) handleImport() {
	// TODO: Implement tag import logic
	// This might involve:
	// 1. Opening a file dialog to select import file
	// 2. Parsing the imported file (JSON/CSV)
	// 3. Validating and adding new tags
	// 4. Updating the state and UI
	tp.toast.Show("Tag import functionality not yet implemented")
}

func (tp *TagsPage) handleExport() {
	// Export all tags
	tagIDs := make([]string, 0, len(tp.state.GetTags()))
	for _, tag := range tp.state.GetTags() {
		tagIDs = append(tagIDs, tag.ID)
	}

	exportPath, err := tp.state.ExportTags(tagIDs)
	if err != nil {
		tp.toast.ShowError("Failed to export tags: " + err.Error())
		return
	}
	tp.toast.ShowMessage("Tags exported successfully to " + exportPath)
}

func (tp *TagsPage) handleBatchDelete() {
	// Collect the IDs of selected tags
	tagIDs := make([]string, 0, len(tp.selectedTags))
	for tagID := range tp.selectedTags {
		tagIDs = append(tagIDs, tagID)
	}

	// Perform batch delete operation
	err := tp.state.DeleteTags(tagIDs)
	if err != nil {
		tp.toast.Show(fmt.Sprintf("Error deleting tags: %v", err))
		return
	}

	// Refresh tags after deletion
	tp.filteredTags = tp.state.GetTags()
	tp.selectedTags = make(map[string]bool)
	tp.batchOps.visible = false
	tp.toast.Show(fmt.Sprintf("%d tags deleted", len(tagIDs)))
}

func (tp *TagsPage) handleBatchGroup() {
	// Collect the IDs of selected tags
	tagIDs := make([]string, 0, len(tp.selectedTags))
	for tagID := range tp.selectedTags {
		tagIDs = append(tagIDs, tagID)
	}

	// Perform batch group operation
	err := tp.state.GroupTags(tagIDs)
	if err != nil {
		tp.toast.Show(fmt.Sprintf("Error grouping tags: %v", err))
		return
	}

	// Refresh tags after grouping
	tp.filteredTags = tp.state.GetTags()
	tp.selectedTags = make(map[string]bool)
	tp.batchOps.visible = false
	tp.toast.Show(fmt.Sprintf("%d tags grouped", len(tagIDs)))
}

func (tp *TagsPage) handleBatchExport() {
	// Collect the IDs of selected tags
	tagIDs := make([]string, 0, len(tp.selectedTags))
	for tagID := range tp.selectedTags {
		tagIDs = append(tagIDs, tagID)
	}

	// Perform batch export operation
	exportPath, err := tp.state.ExportTags(tagIDs)
	if err != nil {
		tp.toast.ShowError("Failed to export tags: " + err.Error())
		return
	}
	tp.toast.ShowMessage("Tags exported successfully to " + exportPath)
}

func (tp *TagsPage) layoutEmptyState(gtx layout.Context) layout.Dimensions {
	label := material.Body1(tp.theme, "No tags or tag groups found. Create your first tag or group!")
	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{
			Top:    unit.Dp(20),
			Bottom: unit.Dp(20),
			Left:   unit.Dp(16),
			Right:  unit.Dp(16),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return label.Layout(gtx)
		})
	})
}

func (tp *TagsPage) layoutEditTag(gtx layout.Context) layout.Dimensions {
	th := tp.theme
	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Top: unit.Dp(16), Bottom: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return material.Body1(th, "Edit Tag").Layout(gtx)
				})
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return tp.editTag.name.Layout(gtx, th, "Tag Name")
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return tp.editTag.description.Layout(gtx, th, "Description")
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return tp.editTag.color.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						btn := material.Button(th, &tp.editTag.save, "Save")
						return btn.Layout(gtx)
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						btn := material.Button(th, &tp.editTag.cancel, "Cancel")
						return btn.Layout(gtx)
					}),
				)
			}),
		)
	})
}

func (tp *TagsPage) layoutAddTag(gtx layout.Context) layout.Dimensions {
	th := tp.theme

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(16), Bottom: unit.Dp(16), Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceBetween}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return tp.addTag.name.Layout(gtx, th, "Tag Name")
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return tp.addTag.description.Layout(gtx, th, "Description (Optional)")
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return tp.addTag.color.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{}.Layout(gtx,
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								btn := material.Button(th, &tp.addTag.button, "Add Tag")
								return btn.Layout(gtx)
							}),
						)
					}),
				)
			})
		}),
	)
}

func (tp *TagsPage) layoutDeleteConfirm(gtx layout.Context) layout.Dimensions {
	if !tp.deleteConfirm.visible {
		return layout.Dimensions{}
	}

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Top:    unit.Dp(16),
				Bottom: unit.Dp(16),
				Left:   unit.Dp(16),
				Right:  unit.Dp(16),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return material.Body1(tp.theme, fmt.Sprintf("Are you sure you want to delete the tag '%s'?", tp.deleteConfirm.tag.Name)).Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis: layout.Horizontal,
			}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					btn := material.Button(tp.theme, &tp.deleteConfirm.confirm, "Confirm")
					return btn.Layout(gtx)
				}),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					btn := material.Button(tp.theme, &tp.deleteConfirm.cancel, "Cancel")
					return btn.Layout(gtx)
				}),
			)
		}),
	)
}

func (tp *TagsPage) layoutBatchOperations(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					btn := material.Button(tp.theme, &tp.batchOps.delete, "Delete")
					return btn.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					btn := material.Button(tp.theme, &tp.batchOps.group, "Group")
					return btn.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					btn := material.Button(tp.theme, &tp.batchOps.export, "Export")
					return btn.Layout(gtx)
				}),
			)
		}),
	)
}

// ... rest of the code remains the same ...
