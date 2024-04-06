package entity

import (
	"math/rand"

	"github.com/EstralMC/GoMine/server/block/cube"
	"github.com/EstralMC/GoMine/server/entity/effect"
	"github.com/EstralMC/GoMine/server/internal/nbtconv"
	"github.com/EstralMC/GoMine/server/world"
	"github.com/EstralMC/GoMine/server/world/particle"
	"github.com/EstralMC/GoMine/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
)

// BottleOfEnchanting is a bottle that releases experience orbs when thrown.
type BottleOfEnchanting struct {
	transform
	age   int
	close bool

	owner world.Entity

	c *ProjectileComputer
}

// NewBottleOfEnchanting ...
func NewBottleOfEnchanting(pos mgl64.Vec3, owner world.Entity) *BottleOfEnchanting {
	b := &BottleOfEnchanting{
		owner: owner,
		c: &ProjectileComputer{&MovementComputer{
			Gravity:           0.07,
			Drag:              0.01,
			DragBeforeGravity: true,
		}},
	}
	b.transform = newTransform(b, pos)
	return b
}

// Type returns BottleOfEnchantingType.
func (b *BottleOfEnchanting) Type() world.EntityType {
	return BottleOfEnchantingType{}
}

// Glint returns true if the bottle should render with glint. It always returns true.
func (b *BottleOfEnchanting) Glint() bool {
	return true
}

// Tick ...
func (b *BottleOfEnchanting) Tick(w *world.World, current int64) {
	if b.close {
		_ = b.Close()
		return
	}
	b.mu.Lock()
	m, result := b.c.TickMovement(b, b.pos, b.vel, 0, 0, b.ignores)
	b.pos, b.vel = m.pos, m.vel
	b.mu.Unlock()

	b.age++
	m.Send()

	if m.pos[1] < float64(w.Range()[0]) && current%10 == 0 {
		b.close = true
		return
	}

	if result != nil {
		colour, _ := effect.ResultingColour(nil)
		w.AddParticle(m.pos, particle.Splash{Colour: colour})
		w.PlaySound(m.pos, sound.GlassBreak{})

		for _, orb := range NewExperienceOrbs(m.pos, rand.Intn(9)+3) {
			orb.SetVelocity(mgl64.Vec3{(rand.Float64()*0.2 - 0.1) * 2, rand.Float64() * 0.4, (rand.Float64()*0.2 - 0.1) * 2})
			w.AddEntity(orb)
		}

		b.close = true
	}
}

// ignores returns whether the BottleOfEnchanting should ignore collision with the entity passed.
func (b *BottleOfEnchanting) ignores(entity world.Entity) bool {
	_, ok := entity.(Living)
	return !ok || entity == b || (b.age < 5 && entity == b.owner)
}

// New creates a BottleOfEnchanting with the position, velocity, yaw, and pitch provided. It doesn't spawn the
// BottleOfEnchanting, only returns it.
func (b *BottleOfEnchanting) New(pos, vel mgl64.Vec3, owner world.Entity) world.Entity {
	bottle := NewBottleOfEnchanting(pos, owner)
	bottle.vel = vel
	return bottle
}

// Owner ...
func (b *BottleOfEnchanting) Owner() world.Entity {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.owner
}

// BottleOfEnchantingType is a world.EntityType for BottleOfEnchanting.
type BottleOfEnchantingType struct{}

func (BottleOfEnchantingType) EncodeEntity() string {
	return "minecraft:xp_bottle"
}
func (BottleOfEnchantingType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.125, 0, -0.125, 0.125, 0.25, 0.125)
}

func (BottleOfEnchantingType) DecodeNBT(data map[string]any) world.Entity {
	b := NewBottleOfEnchanting(nbtconv.MapVec3(data, "Pos"), nil)
	b.vel = nbtconv.MapVec3(data, "Motion")
	return b
}

func (BottleOfEnchantingType) EncodeNBT(e world.Entity) map[string]any {
	b := e.(*BottleOfEnchanting)
	return map[string]any{
		"Pos":    nbtconv.Vec3ToFloat32Slice(b.Position()),
		"Motion": nbtconv.Vec3ToFloat32Slice(b.Velocity()),
	}
}
