package db

type User struct {
	Name        string // For example John Doe (Last First names)
	Age         *float64
	PhoneNumber string
	Email       string
	Skype       string
	Location    Location
}
