package models

type AlertaEmergencia struct {
	ID            int    `json:"id"`
	Gravedad      int    `json:"gravedad"`
	Requerimiento string `json:"requerimiento"`
	Estado        string `json:"estado"`
}
