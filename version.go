package main

// shaman version information (populated by go linker)
// -ldflags="-X main.version=${tag} -X main.branch=${branch} -X main.commit=${commit}"
var (
	version string
	branch  string
	commit  string
)
