package api

type Address struct {
	Locations [][][]Location
	Street    string
	Building  int
	Flat      int
}
