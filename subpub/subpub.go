package subpub

type Consumer interface {
	Consume() ([]byte, error)
}

type Producer interface {
	Produce([]byte) error
}
