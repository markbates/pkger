package pkger

// Open opens the named file for reading.
func Open(p string) (*File, error) {
	pt, err := Parse(p)
	if err != nil {
		return nil, err
	}
	return rootIndex.Open(pt)
}

// Create creates the named file with mode 0666 (before umask), truncating it if it already exists. If successful, methods on the returned File can be used for I/O; the associated file descriptor has mode O_RDWR. If there is an error, it will be of type *PathError.
func Create(p string) (*File, error) {
	pt, err := Parse(p)
	if err != nil {
		return nil, err
	}
	return rootIndex.Create(pt)
}
