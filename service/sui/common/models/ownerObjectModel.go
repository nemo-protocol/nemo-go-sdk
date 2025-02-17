package models

type CommonOnChainDataResp struct {
	MoveObject FieldsMoveObject `json:"moveObject"`
}

type FieldsMoveObject struct {
	Fields interface{} `json:"fields"`
}
