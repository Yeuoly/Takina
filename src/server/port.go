package server

import (
	"errors"

	"github.com/Yeuoly/Takina/src/helper"
)

func (c *Takina) requestAvailablePort() (int, error) {
	c.requestPortMutex.Lock()
	defer c.requestPortMutex.Unlock()

	if c.PortPool.Len() == 0 {
		// release callback
		return 0, errors.New("no available port")
	}

	first_port := 0

	for {
		e := c.PortPool.Front()
		c.PortPool.Remove(e)

		port := e.Value.(int)
		if first_port == 0 {
			first_port = port
		} else if port == first_port {
			c.PortPool.PushBack(port)
			return 0, errors.New("no available port")
		}

		// check if port is available
		if helper.TestPortAvailable(port) {
			return port, nil
		}

		// if not available, release it
		c.PortPool.PushBack(port)
	}
}

func (c *Takina) releasePort(port int) error {
	c.requestPortMutex.Lock()
	defer c.requestPortMutex.Unlock()

	// check if port exists
	exists := false
	for e := c.PortPool.Front(); e != nil; e = e.Next() {
		if e.Value.(int) == port {
			exists = true
			break
		}
	}

	if exists {
		return errors.New("port already exists")
	}

	// check if port is in range
	inRange := false
	for _, p := range c.RealPortRange {
		if p == port {
			inRange = true
			break
		}
	}

	if !inRange {
		return errors.New("port is not in range")
	}

	c.PortPool.PushBack(port)

	return nil
}
