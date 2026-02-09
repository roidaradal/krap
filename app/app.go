package app

// Initializer function with name
type Initializer struct {
	Fn   func() error
	Name string
}
