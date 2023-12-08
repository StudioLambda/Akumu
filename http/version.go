package http

type Version string

const (
	Version0_9 Version = "HTTP/0.9"
	Version1_0 Version = "HTTP/1.0"
	Version1_1 Version = "HTTP/1.1"
	Version2_0 Version = "HTTP/2.0"
	Version3_0 Version = "HTTP/3.0"
)
