package models

type Points struct {
	Points int64
}

// Map of {id: CalculatedPoints} to avoid recalculation
var PointsForId = map[string]Points{}
