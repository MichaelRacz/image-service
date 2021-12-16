package docker

// NOTE: Files are kept as strings to save implementation time,
// depending on the use cases, persisting the files may be the
// better option
type Dockerfile string

// NewDockerfile converts bytes to a Dockerfile
func NewDockerfile(bs []byte) Dockerfile {
	return Dockerfile(string(bs))
}

// String implementes the Stringer interface
func (d Dockerfile) String() string {
	return string(d)
}
