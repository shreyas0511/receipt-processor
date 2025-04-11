package models

type Points struct {
	Points int64
}

// Store the points calculated for a particular id
var PointsForId = map[string]Points{}
