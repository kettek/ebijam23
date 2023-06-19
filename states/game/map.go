package game

import (
	"ebijam23/resources"
	"ebijam23/states"
	"errors"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	ErrMissingMap = errors.New("missing map")
)

type Map struct {
	data    *resources.Map
	actors  []Actor
	bullets []*Bullet
}

func (s *World) TravelToMap(ctx states.Context, mapName string) error {
	mapData := ctx.Manager.GetAs("maps", mapName, (*resources.Map)(nil)).(*resources.Map)
	if mapData == nil {
		return ErrMissingMap
	}

	m := &Map{
		data: mapData,
	}

	for i, l := range m.data.Layers {
		for j, row := range l.Cells {
			for k, cell := range row {
				switch m.data.RuneMap[string(cell.Type)] {
				case "wall":
					cell.Sprite = resources.NewSprite(ctx.Manager.GetAs("images", "wall", (*ebiten.Image)(nil)).(*ebiten.Image))
					cell.Blocks = true
				case "floor":
					cell.Sprite = resources.NewSprite(ctx.Manager.GetAs("images", "floor", (*ebiten.Image)(nil)).(*ebiten.Image))
				case "empty":
					cell.Sprite = resources.NewSprite(ctx.Manager.GetAs("images", "empty", (*ebiten.Image)(nil)).(*ebiten.Image))
				default:
					cell.Sprite = resources.NewSprite(ctx.Manager.GetAs("images", m.data.RuneMap[string(cell.Type)], (*ebiten.Image)(nil)).(*ebiten.Image))
				}
				cell.Sprite.SetXY(
					cell.Sprite.Width()*float64(k),
					cell.Sprite.Height()*float64(j),
				)
				row[k] = cell
			}
			l.Cells[j] = row
		}
		m.data.Layers[i] = l
	}

	// Create actors.
	for _, a := range m.data.Actors {
		switch a.Type {
		case "spawner":
			spawner := CreateSpawner(float64(a.Spawn[0])*16, float64(a.Spawn[1])*16)
			m.actors = append(m.actors, spawner)
		}
	}

	s.activeMap = m

	// Move players over to new map.
	for _, p := range s.Players {
		p.Actor().SetXY(float64(m.data.Start[0]*16), float64(m.data.Start[1]*16))
		m.actors = append(m.actors, p.Actor())
	}

	return nil
}

func (m *Map) Draw(screen *ebiten.Image) {
	for i := len(m.data.Layers) - 1; i >= 0; i-- {
		l := m.data.Layers[i]
		for _, row := range l.Cells {
			for _, cell := range row {
				cell.Sprite.Draw(screen)
			}
		}
	}
	for _, a := range m.actors {
		a.Draw(screen)
	}
	for _, b := range m.bullets {
		b.Draw(screen)
	}
}

func (m *Map) Collides(s Shape) bool {
	// Get nearest cell to shape coordinates and check adjacent cells for collisions.
	x, y, _, _ := s.Bounds()
	x /= 16
	y /= 16
	x = math.Round(x)
	y = math.Round(y)
	z := 0

	check := func(x, y int) bool {
		if y >= 0 && int(y) < len(m.data.Layers[z].Cells) && x >= 0 && int(x) < len(m.data.Layers[z].Cells[int(y)]) {
			if m.data.Layers[z].Cells[int(y)][int(x)].Blocks {
				cell := m.data.Layers[z].Cells[int(y)][int(x)]
				if s.Collides(&RectangleShape{
					X:      cell.Sprite.X,
					Y:      cell.Sprite.Y,
					Width:  cell.Sprite.Width(),
					Height: cell.Sprite.Height(),
				}) {
					return true
				}
			}
		}
		return false
	}

	// TODO: Get whatever its called when you get the minimum distance to ensure contact but not intersection and return it so the caller can still potentially move.
	// lol
	return check(int(x), int(y)) ||
		check(int(x), int(y+1)) ||
		check(int(x), int(y-1)) ||
		check(int(x-1), int(y)) ||
		check(int(x-1), int(y+1)) ||
		check(int(x-1), int(y-1)) ||
		check(int(x+1), int(y)) ||
		check(int(x+1), int(y+1)) ||
		check(int(x+1), int(y-1))
}
