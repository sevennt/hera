package hera

import "fmt"

// Plugin will be executed after configuration modified.
type Plugin struct {
	Key     []string
	Handler func() error
}

func (p Plugin) fire() {
	if len(p.Key) == 1 && p.Key[0] == "*" {
		if err := p.Handler(); err != nil {
			fmt.Printf("plugin fire fail, %#v, %s", p, err.Error())
		}
		return
	}
}
