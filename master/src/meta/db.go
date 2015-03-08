package meta

type Table struct {
	Name string
	key string
	Column map[string]string
}

type NameSpace struct {
	Name string
	Tables map[string]Table
}



