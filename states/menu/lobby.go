package menu

import (
	"ebijam23/resources"
	"ebijam23/states"
	"ebijam23/states/game"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Lobby struct {
	clickSound      *resources.Sound
	items           []resources.MenuItem
	multiplayerItem *resources.ButtonItem
	joinItem        *resources.ButtonItem
	hostItem        *resources.ButtonItem
	backItem        *resources.TextItem
	lobbyItem       *resources.InputItem
	playerEntries   []*PlayerEntry
	overlay         game.Overlay
}

func (s *Lobby) Init(ctx states.Context) error {
	s.overlay.Init(ctx)
	//
	s.clickSound = ctx.Manager.GetAs("sounds", "click", (*resources.Sound)(nil)).(*resources.Sound)

	s.playerEntries = append(s.playerEntries, &PlayerEntry{
		player: game.NewLocalPlayer(),
	})

	for _, e := range s.playerEntries {
		e.Init(ctx)
	}

	s.multiplayerItem = &resources.ButtonItem{
		Text: ctx.L("Multiplayer"),
		X:    500,
		Y:    20,
		Callback: func() bool {
			s.clickSound.Play(1.0)
			for _, item := range s.items {
				if item == s.multiplayerItem {
					s.multiplayerItem.SetHidden(true)
					s.lobbyItem.SetHidden(false)
					s.joinItem.SetHidden(false)
					s.hostItem.SetHidden(false)
					s.playerEntries = append(s.playerEntries, &PlayerEntry{})
					s.playerEntries[len(s.playerEntries)-1].Init(ctx)
					break
				}
			}
			return false
		},
	}

	s.lobbyItem = &resources.InputItem{
		X:     350,
		Y:     20,
		Width: 150,
		Callback: func() bool {
			return false
		},
	}
	s.lobbyItem.SetHidden(true)

	s.joinItem = &resources.ButtonItem{
		Text: ctx.L("Host"),
		X:    450,
		Y:    20,
		Callback: func() bool {
			s.clickSound.Play(1.0)
			// TODO: Create network server:
			//   - Check if lobby is an address or ip
			//      - If so, begin directly hosting and wait for a client to connect.
			//			- If not, connect to magnet's matchmaker with the lobby as the advertisement and begin waiting for a client to connect.
			return false
		},
	}
	s.joinItem.SetHidden(true)

	s.hostItem = &resources.ButtonItem{
		Text: ctx.L("Join"),
		X:    450 + 50,
		Y:    20,
		Callback: func() bool {
			s.clickSound.Play(1.0)
			// TODO: Create network client:
			//   - Check if lobby is an address or ip
			//      - If so, directly connect to it
			//			- If not, connect to magnet's matchmaker and use the lobby as the target name. Wait for response, and if an ip:port, directly connect to it using the same socket.
			return false
		},
	}
	s.hostItem.SetHidden(true)

	s.backItem = &resources.TextItem{
		Text: ctx.L("Back"),
		X:    30,
		Y:    335,
		Callback: func() bool {
			s.clickSound.Play(1.0)
			ctx.StateMachine.PopState()
			return false
		},
	}
	s.items = append(s.items, s.backItem, s.multiplayerItem, s.lobbyItem, s.joinItem, s.hostItem)

	return nil
}

func (s *Lobby) Finalize(ctx states.Context) error {
	return nil
}

func (s *Lobby) Enter(ctx states.Context) error {
	return nil
}

func (s *Lobby) Update(ctx states.Context) error {
	s.overlay.Update(ctx)

	s.lobbyItem.Update()

	// Check for controller button hit to activate player 2.
	for _, gamepadID := range ebiten.AppendGamepadIDs(nil) {
		if inpututil.IsGamepadButtonJustPressed(gamepadID, ebiten.GamepadButton9) {
			pl := game.NewLocalPlayer()
			s.playerEntries[1].player = pl
			s.playerEntries[1].controllerId = int(gamepadID)
			s.playerEntries[1].SyncController(ctx)
			pl.GamepadID = int(gamepadID)
			// TODO: Stop network stuff and hide host/join.
			s.hostItem.SetHidden(true)
			s.joinItem.SetHidden(true)
			s.lobbyItem.SetHidden(true)
		}
	}

	x := -(len(s.playerEntries) - 1) * 150 / 2
	for i, e := range s.playerEntries {
		e.Update(ctx, float64(x+i*150))
	}

	x, y := ebiten.CursorPosition()
	for _, m := range s.items {
		m.CheckState(float64(x), float64(y))
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) {
		for _, m := range s.items {
			if m.Hovered() {
				if m.Activate() {
					return nil
				}
			}
		}
	}
	return nil
}

func (s *Lobby) Draw(ctx states.DrawContext) {
	ctx.Text.SetColor(color.White)
	for _, e := range s.playerEntries {
		e.Draw(ctx)
	}
	for _, m := range s.items {
		m.Draw(ctx)
	}
	s.overlay.Draw(ctx)
}
