package models

type Data struct {
	MoveObject MoveObject `json:"moveObject"`
}

type BcsData struct {
	MoveObject BcsMoveObject `json:"moveObject"`
}

type MoveObject struct {
	Type              string                 `json:"type"`
	HasPublicTransfer bool                   `json:"hasPublicTransfer"`
	Fields            map[string]interface{} `json:"fields"`
}

type BcsMoveObject struct {
	Type              string `json:"type"`
	HasPublicTransfer bool   `json:"hasPublicTransfer"`
	Version           int64  `json:"version"`
}

type Object struct {
	Package  PackageObject  `json:"package"`
}
type PackageObject struct {
	Disassembled map[string]interface{} `json:"disassembled"`
}
