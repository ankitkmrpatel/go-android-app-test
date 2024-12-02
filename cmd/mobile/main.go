//go:build android || ios
// +build android ios

package main

import (
	"log"
	"os"

	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget/material"

	"gioui.org/app"
	"gioui.org/app/system"

	appState "github.com/goBookMarker/internal/app"
	"github.com/goBookMarker/internal/storage"
	"github.com/goBookMarker/internal/ui"
)

func main() {
	go func() {
		w := app.NewWindow(
			app.Title("GoBookMarker"),
			app.Size(unit.Dp(360), unit.Dp(640)),
			app.MinSize(unit.Dp(300), unit.Dp(500)),
		)
		if err := run(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(w *app.Window) error {
	// Initialize theme with default font collection
	th := material.NewTheme()
	th.Shaper = gofont.Collection()

	// Initialize storage
	db, err := storage.NewSQLiteDB()
	if err != nil {
		return err
	}
	defer db.Close()

	// Initialize application state
	state := appState.NewAppState()

	// Initialize UI
	ui := ui.NewUI(th, state)

	// Create operation list for window
	var ops op.Ops

	// Event loop
	for {
		e := w.NextEvent()

		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err

		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)

			// Layout UI
			ui.Layout(gtx)

			// Render frame
			e.Frame(gtx.Ops)

		case system.StageEvent:
			if e.Stage >= system.StageRunning {
				// App is visible, trigger initial data load
				state.LoadInitialData()
			}
		}
	}
}
