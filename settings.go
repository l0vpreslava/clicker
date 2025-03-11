package main

import (
	"encoding/json"
	"fmt"
	"os"
)

var PossibleScales = []float32{1, 2, 2.5, 3, 3.5, 4, 6}

func GetResolution(scaleIndex int) (int32, int32) {
	scale := PossibleScales[scaleIndex]
	w := scale * WorldWidth
	h := scale * WorldHeight
	return int32(w), int32(h)
}

type Settings struct {
	Borderless bool    `json:"borderless"`
	GameVolume float32 `json:"volume"`
	ScaleIndex int     `json:"scaleIndex"`
}

func NewSettigs() Settings {
	return Settings{
		Borderless: false,
		GameVolume: 0.50,
		ScaleIndex: 1,
	}
}

func LoadSettings() Settings {
	s, err := os.ReadFile(GetSettingsPath())
	if err != nil {
		fmt.Printf("Error reading settings: %s", err)
		return NewSettigs()
	}

	settings := Settings{}
	err = json.Unmarshal(s, &settings)
	if err != nil {
		fmt.Printf("Error unmarshaling settings: %s", err)
		return NewSettigs()
	}

	return settings
}

func (settings *Settings) Save() {
	s, err := json.Marshal(settings)
	if err != nil {
		fmt.Printf("Error marshaling settings: %s", err)
		return
	}

	os.MkdirAll(GetGameCacheDir(), os.ModePerm)
	err = os.WriteFile(GetSettingsPath(), s, os.ModePerm)
	if err != nil {
		fmt.Printf("Error saving settings: %s", err)
		return
	}
}

func GetSettingsPath() string {
	return fmt.Sprintf("%s/settings.json", GetGameCacheDir())
}
