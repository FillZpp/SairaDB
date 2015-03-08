package meta

const (
	Read  = 1
	Write = 2
	Alter = 4
	
)

type User struct {
	Name string
	GlobalAuthority uint
}

