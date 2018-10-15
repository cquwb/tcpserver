package tcp

import (
	
)

/**
*/


type Server interface {
	NewSession() Sessioner
	Close()
}

type Sessioner interface {
	Init(*ClientSession)
	Process([]byte)
	Close()
}