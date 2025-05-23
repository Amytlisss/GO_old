package models

type Animal struct {
	ID          int
	Name        string
	Type        string // dog, cat, etc.
	Breed       string
	Age         int
	Description string
	ImageURL    string
	Available   bool
}
