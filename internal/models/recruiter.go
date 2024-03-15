package models

type Recruiter struct {
	ID              int
	PublicID        string
	CompanyPublicID string
	FirstName       string
	LastName        string
	Photo           string
	Company         *Company
	Positions       []Position
}

type Company struct {
	ID          int
	PublicID    string
	Name        string
	Description string
}

type Position struct {
	ID       int
	PublicID string
	Name     string
	Status   int
}
