package data

type ClientData struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Revenue int    `json:"revenue"`
	UkBased bool   `json:"ukBased"`
}

type ColumnData struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type TableData struct {
	ColumnData []ColumnData `json:"columnData"`
	ClientData []ClientData `json:"clientData"`
}
