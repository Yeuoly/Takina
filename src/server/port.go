package server

import (
	"errors"

	"github.com/Yeuoly/Takina/src/helper"
)

func (c *Takina) requestAvailablePort(callbacks ...[]func()) (int, error) {
	c.requestPortMutex.Lock()
	defer c.requestPortMutex.Unlock()

	if c.PortPool.Len() == 0 {
		// release callback
		if len(callbacks) > 0 {
			for _, cb := range callbacks[0] {
				cb()
			}
		}

		return 0, errors.New("no available port")
	}

	e := c.PortPool.Front()
	c.PortPool.Remove(e)

	port := e.Value.(int)

	if !helper.TestPortAvailable(port) {
		helper.Warn("[Takina] port %d is not available, try next", port)
		if len(callbacks) > 0 {
			callbacks[0] = append(callbacks[0], func() {
				c.PortPool.PushBack(port)
			})

			return c.requestAvailablePort(callbacks[0])
		} else {
			callbacks = append(callbacks, []func(){
				func() {
					c.PortPool.PushBack(port)
				},
			})
			return c.requestAvailablePort(callbacks[0])
		}
	}

	// release callback
	if len(callbacks) > 0 {
		for _, cb := range callbacks[0] {
			cb()
		}
	}

	return e.Value.(int), nil
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
