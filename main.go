package main

import (
	"ebijam23/states"
	"ebijam23/states/menu"
	"embed"
	"io/fs"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/kettek/go-multipath/v2"
	"golang.org/x/image/font/sfnt"
)

//go:embed assets/*
var embedFS embed.FS

func main() {
	game := &Game{}

	// Allow loading from filesystem.
	game.Manager.files.InsertFS(os.DirFS("assets"), multipath.FirstPriority)

	// Also allow loading from embedded filesystem.
	sub, err := fs.Sub(embedFS, "assets")
	if err != nil {
		panic(err)
	}
	game.Manager.files.InsertFS(sub, multipath.LastPriority)

	if err := game.Manager.Setup(); err != nil {
		panic(err)
	}

	// Might as well load all assets up front (for now -- might not want to with music later).
	if err := game.Manager.LoadAll(); err != nil {
		panic(err)
	}

	// Set our locale.
	game.Localizer.manager = &game.Manager
	game.Localizer.SetLocale("ja")

	// Initialize game fields as necessary.
	if err := game.Init(); err != nil {
		panic(err)
	}

	// Ensure we have our font.
	if f := game.Manager.GetAs("fonts", "x12y16pxMaruMonica", (*sfnt.Font)(nil)).(*sfnt.Font); f == nil {
		panic("missing font")
	} else {
		game.Text.SetFont(f)
		game.Text.Utils().SetCache8MiB()
	}

	// Initialize audio.
	audio.NewContext(44100)

	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("ebijam23")
	ebiten.SetCursorMode(ebiten.CursorModeHidden)

	// Push the state.
	game.PushState(&menu.Menu{})

	if err := ebiten.RunGame(game); err != nil {
		if err == states.ErrQuitGame {
			return
		}
		panic(err)
	}
}
