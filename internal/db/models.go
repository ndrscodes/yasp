package db

type System struct {
	Id   uint
	Name string
}

func NewSystem(name string) System {
	return System{
		Name: name,
	}
}
