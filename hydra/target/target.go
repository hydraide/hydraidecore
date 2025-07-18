package target

import "github.com/hydraide/hydraidecore/hydra/swamp/treasure"

type Targets struct {
	TargetSwamps map[string]*Target `json:"targets"`
}

type Target struct {
	TargetSwampName string                           `json:"tsn"`
	EventTypes      map[treasure.TreasureStatus]bool `json:"et"`
}
