package export

import "fmt"

type command struct {
	base   string
	params []string
}

func newCommand() *command {
	return &command{
		base:   "clickhouse-client",
		params: []string{},
	}
}

func (c *command) appendParam(name, value string) {
	if c == nil {
		return
	}
	c.params = append(c.params, fmt.Sprintf("--%s", name), value)
}

func (c *command) getBase() string {
	return c.base
}

func (c *command) getParams() []string {
	return c.params
}
