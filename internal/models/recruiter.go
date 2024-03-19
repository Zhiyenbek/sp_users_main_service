package models

type Recruiter struct {
	PublicID        string
	CompanyPublicID string
	FirstName       string
	LastName        string
	Photo           string
	Company         *Company
	Positions       []Position
}

type Position struct {
	PublicID string
	Name     string
	Status   int
}
