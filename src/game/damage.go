package game

const (
	Physical int = 0
	Magic        = 1
	Pure         = 2
)

type Damage struct {
	Value float32
	Type  int
}
