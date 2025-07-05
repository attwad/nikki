package model

type Date struct {
	Year   int64
	Month  int64
	Day    int64
	Things []Thing
}

type Thing struct {
	ID   int64
	What string
}

type NewThing struct {
	Year  int64
	Month int64
	Day   int64
	What  string
}
