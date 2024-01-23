package game

import (
	"dungeon/internal/animation"
	"dungeon/internal/numerics"
	"github.com/hajimehoshi/ebiten/v2"
)

type Projectile struct {
	// Direction is the direction the projectile is moving in.
	Direction numerics.Vec2

	*Object
}

func NewProjectile(src *Object, direction numerics.Vec2, img *animation.Image) *Projectile {
	aabb := NewAABB(numerics.ZeroVec2(), img)

	obj := &Object{
		Image:       map[Orientation]*animation.Image{All: img},
		Op:          &ebiten.DrawImageOptions{},
		Position:    src.Position,
		Center:      src.Center,
		Velocity:    numerics.OneVec2().MulScalar(2),
		AABB:        aabb,
		Orientation: All,
	}

	return &Projectile{Direction: direction, Object: obj}
}

func (p *Projectile) Step() {
	diff := p.Direction.Mul(p.Velocity)
	p.UpdatePosition(diff)
}
