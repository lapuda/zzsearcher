package store

type DataStore interface {
	Save(interface{}, string)
}
