package models

type Candidate struct {
	PublicID        *string     `json:"public_id"`
	FirstName       *string     `json:"first_name"`
	LastName        *string     `json:"last_name"`
	CurrentPosition *string     `json:"current_position"`
	Resume          *string     `json:"resume"`
	Bio             *string     `json:"bio"`
	Skills          []*string   `json:"skills"`
	Interviews      []Interview `json:"interviews,omitempty"`
	Education       *string     `json:"education"`
}

type Interview struct {
	ID       int                    `json:"id"`
	PublicID string                 `json:"public_id"`
	Results  map[string]interface{} `json:"results"`
}
