package api_struct

type Ctx interface {
	Persist() error
	Load() error
}
