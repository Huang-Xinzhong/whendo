package app

import (
	"context"
	"runtime"

	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// NewMenu returns the application menu.
func NewMenu(ctx context.Context) *menu.Menu {
	appMenu := menu.NewMenu()

	if runtime.GOOS == "darwin" {
		m := appMenu.AddSubmenu("WhenDo")
		m.AddText("About", nil, func(_ *menu.CallbackData) {
			wailsruntime.EventsEmit(ctx, "menu:about")
		})
		m.AddSeparator()
		m.AddText("Quit", keys.CmdOrCtrl("q"), func(_ *menu.CallbackData) {
			wailsruntime.Quit(ctx)
		})
	}

	editMenu := appMenu.AddSubmenu("Edit")
	editMenu.AddText("Undo", keys.CmdOrCtrl("z"), nil)
	editMenu.AddText("Redo", keys.CmdOrCtrl("shift+z"), nil)
	editMenu.AddSeparator()
	editMenu.AddText("Cut", keys.CmdOrCtrl("x"), nil)
	editMenu.AddText("Copy", keys.CmdOrCtrl("c"), nil)
	editMenu.AddText("Paste", keys.CmdOrCtrl("v"), nil)
	editMenu.AddText("Select All", keys.CmdOrCtrl("a"), nil)

	return appMenu
}
