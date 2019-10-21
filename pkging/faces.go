package pkging

type Adder interface {
	Add(files ...File) error
}
