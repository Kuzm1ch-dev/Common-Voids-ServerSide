package game

type Entity struct {
	Name  string
	Stats EntityStat
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
