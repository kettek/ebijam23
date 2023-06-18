package game

type ActorActions struct {
	Actor   Actor
	Actions []Action
}

type BulletActions struct {
	Bullet  *Bullet
	Actions []Action
}

type Action interface {
}

type ActionMove struct {
	X, Y float64
}

type ActionReflect struct {
	X, Y float64
}

type ActionDeflect struct {
	Direction float64
}

type ActionSpawnBullets struct {
	Bullets []*Bullet
}

type ActionFindNearestActor struct {
	Actor Actor
}
