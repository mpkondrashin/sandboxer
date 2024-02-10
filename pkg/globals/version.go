package globals

/*
Folloing value is defined during linking in makefile as follows:

VERSION := $(shell git describe --tags --abbrev=0)
BUILD_OPTS := -ldflags "-X 'github.com/mpkondrashin/tunneleffect/internal/version.MajorMinorRevision=$(VERSION)'"
*/
var (
	Version = "X.X.X"
	Build   = "0"
)
