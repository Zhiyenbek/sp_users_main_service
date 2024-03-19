package models

type Company struct {
	ID          int    `json:"-"`
	PublicID    string `json:"public_id"`
	Name        string `json:"name"`
	Logo        string `json:"logo"`
	Description string `json:"description"`
}
