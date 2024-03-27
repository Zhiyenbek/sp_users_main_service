package models

type Recruiter struct {
	PublicID        string     `json:"public_id"`
	CompanyPublicID string     `json:"company_public_id"`
	FirstName       string     `json:"first_name"`
	LastName        string     `json:"last_name"`
	Photo           string     `json:"photo"`
	Company         *Company   `json:"company"`
	Positions       []Position `json:"positions"`
}

type Position struct {
	PublicID string `json:"public_id"`
	Name     string `json:"name"`
	Status   int    `json:"status"`
}
