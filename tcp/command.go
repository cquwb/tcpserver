package tcp

import (
	"fmt"
)

/**
*/

type Services func(Sessioner, []byte) bool 

type Command struct {
	cmds map[uint32]Services
}

func NewCommand() *Command {
	return &Command{
		cmds:make(map[uint32]Services),
	}
}

func (this *Command) Register(id uint32, f Services) bool {
	if _,exists := this.cmds[id];exists {
		fmt.Printf("registe services fail %d", id)
		return false
	}
	this.cmds[id] = f
	return true
}

func (this *Command) Dispatch(sess Sessioner, id uint32, data []byte) bool {
	if f,exists := this.cmds[id];exists {
		return f(sess, data)
	}
	fmt.Printf("registe services fail %d", id)
	return false
}