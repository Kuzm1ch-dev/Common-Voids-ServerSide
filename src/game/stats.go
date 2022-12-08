package game

type EntityStat struct {
	CurrentHealth      float32
	MaxHealth          float32
	HealthRegeneration float32
	CurrentMana        float32
	MaxMana            float32
	ManaRegeneration   float32
	PhysicalResistance float32
	MagicResistance    float32
}

type ArmorStat struct {
	PhysicalResistance float32
	MagicResistance    float32
	AdditionalHealth   float32
	HealthRegeneration float32
	AdditionalMana     float32
	ManaRegeneration   float32
}

type WeaponStat struct {
	Damage Damage
}
