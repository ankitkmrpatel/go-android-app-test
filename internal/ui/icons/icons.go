package icons

import (
	"gioui.org/widget"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

var (
	AddIcon      = mustIcon(icons.ContentAdd)
	BookmarkIcon = mustIcon(icons.ActionBookmark)
	DeleteIcon   = mustIcon(icons.ActionDelete)
	EditIcon     = mustIcon(icons.ImageEdit)
	ExportIcon   = mustIcon(icons.FileCloudUpload)
	ImportIcon   = mustIcon(icons.FileCloudDownload)
	FavoriteIcon = mustIcon(icons.ActionFavorite)
	SearchIcon   = mustIcon(icons.ActionSearch)
	SettingsIcon = mustIcon(icons.ActionSettings)
	ShareIcon    = mustIcon(icons.SocialShare)
	TagIcon      = mustIcon(icons.ActionLabel)
	HomeIcon     = mustIcon(icons.ActionHome)
)

func mustIcon(data []byte) *widget.Icon {
	icon, err := widget.NewIcon(data)
	if err != nil {
		panic(err)
	}
	return icon
}
