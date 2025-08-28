package main

import (
	"file-fusion-rename/internal/app"
	"file-fusion-rename/internal/ui/dialog"
	"file-fusion-rename/internal/ui/menu"
	"file-fusion-rename/internal/ui/tray"

	"github.com/getlantern/systray"
)

type App struct {
	state  *app.State
	tray   *tray.Manager
	dialog *dialog.Manager
	menu   *menu.Manager
}

func NewApp() *App {
	return &App{
		state:  app.NewState(),
		tray:   tray.NewManager(),
		dialog: dialog.NewManager(),
	}
}

func (a *App) onReady() {
	a.tray.Initialize()
	a.menu = menu.NewManager(a.state, a.dialog, a.tray)
	a.menu.CreateMenus()
	a.tray.UpdateTitle(a.state.GetConvertDirection())
}

func (a *App) onExit() {
}

func main() {
	app := NewApp()
	systray.Run(app.onReady, app.onExit)
}
