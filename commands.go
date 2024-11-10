package main

import "errors"

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmds[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	if cCmd, ok := c.cmds[cmd.name]; ok {
		return cCmd(s, cmd)
	}
	return errors.New("command not available, register command")
}
