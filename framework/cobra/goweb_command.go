package cobra

import "github.com/wxsatellite/goweb/framework"

func (c *Command) SetContainer(container framework.Container) {
	c.container = container
}

func (c *Command) Container() framework.Container {
	return c.container
}
