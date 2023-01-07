package game

import (
	"github.com/ByteArena/box2d"
	"server/src/game/buffs"
)

type Entity struct {
	Name     string
	Collider box2d.B2Body
	Stats    EntityStat
	Buffs    []buffs.Buff
}

func (e *Entity) SubHealth(value float32) {
	e.Stats.CurrentHealth -= value
}

func (e *Entity) AddHealth(value float32) {
	e.Stats.CurrentHealth += value
}

func (e *Entity) SubMana(value float32) {
	e.Stats.CurrentMana -= value
}

func (e *Entity) AddMana(value float32) {
	e.Stats.CurrentMana += value
}

func (e *Entity) GetDamage(damage Damage) {
	if damage.Type == Physical {
		e.Stats.CurrentHealth -= damage.Value * (1 - e.Stats.PhysicalResistance)
		return
	}
	if damage.Type == Magic {
		e.Stats.CurrentHealth -= damage.Value * (1 - e.Stats.MagicResistance)
		return
	}
	if damage.Type == Pure {
		e.Stats.CurrentHealth -= damage.Value
		return
	}
}

func (e *Entity) CheckHealth() {
	if e.Stats.CurrentHealth <= 0 {

	}
	if e.Stats.CurrentHealth >= e.Stats.MaxHealth {
		e.Stats.CurrentHealth = e.Stats.MaxHealth
	}
}

func (e *Entity) SetPosition(pos box2d.B2Vec2) {
	e.Collider.SetTransform(pos, e.Collider.GetAngle())
}

func (e *Entity) SetAngle(angle float64) {
	e.Collider.SetTransform(e.Collider.GetPosition(), angle)
}
