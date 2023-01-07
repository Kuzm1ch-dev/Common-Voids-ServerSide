package game

import "image/color"

type AppearanceData struct {
	Sex       uint8      `json:"sex"`
	HairType  uint8      `json:"hairType"`
	HairColor color.RGBA `json:"hairColor"`
	SkinColor color.RGBA `json:"skinColor"`
}
