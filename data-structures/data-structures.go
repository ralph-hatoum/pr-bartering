package datastructures

import "time"

type NodeScore struct {
	NodeIP string
	Score  float64
}

type PeerStorageUse struct {
	NodeIP        string
	StorageAtNode float64
}

type NodeRatio struct {
	NodeIP string
	Ratio  float64
}

type FulfilledRequest struct {
	CID  string
	Peer string
}

type StorageRequest struct {
	CID      string
	FileSize float64
}

type StorageRequestTimed struct {
	CID             string
	DurationMinutes int64
}

type StorageRequestTimedAccepted struct {
	CID      string
	Deadline time.Time
}

type FilesAtPeers struct {
	Peer  string
	Files []string
}

type ScoreVariationScenario struct {
	Scenario  string
	Variation float64
}
