package models
type RecursoClinico struct {
	ID             int    `json:"id"`
	ClinicaID      int    `json:"clinica_id"`
	TipoMaquina    string `json:"tipo_maquina"`
	EstaDisponible bool   `json:"esta_disponible"`
}