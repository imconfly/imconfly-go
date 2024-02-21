package constants

// See:
// man sysexits.h (on GNU-based systems)
// man sysexits (on BSD systems)

const (
	ExUsage    = 65 // command line usage error
	ExConfig   = 78 // configuration error
	ExSoftware = 70 // internal software error
)
