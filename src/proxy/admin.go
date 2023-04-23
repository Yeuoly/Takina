package proxy

import (
	"encoding/json"
	"errors"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type FrpAdmin struct{}

func (f *FrpAdmin) GetFrpStatus(note *FrpcNote) StatusResponse {
	var status StatusResponse
	result := defaultApi.getFrpStatus(note)
	json.Unmarshal([]byte(result), &status)
	return status
}

func (f *FrpAdmin) LaunchProxy(note *FrpcNote, laddr string, lport int, raddr string, rport int, protocol string) error {
	note.mtx.RLock()
	for _, i := range note.CurrentProxy {
		if i.Raddr == raddr && i.Rport == rport {
			note.mtx.RUnlock()
			return errors.New("proxy already exists")
		}
	}
	note.mtx.RUnlock()

	//get uuid
	id := uuid.New().String()
	proxy := Proxy{
		Id:    id,
		Laddr: laddr,
		Lport: lport,
		Raddr: raddr,
		Rport: rport,
		Type:  protocol,
	}

	//add proxy to note
	note.mtx.Lock()
	note.CurrentProxy[id] = proxy
	note.mtx.Unlock()

	err := f.Reload(note)
	if err != nil {
		note.mtx.Lock()
		delete(note.CurrentProxy, id)
		note.mtx.Unlock()
		return err
	}

	//check
	err = f.CheckAndSync(note)
	if err != nil {
		return err
	}

	return nil
}

func (f *FrpAdmin) StopProxy(note *FrpcNote, laddr string, lport int) (string, int, error) {
	note.mtx.RLock()
	for _, i := range note.CurrentProxy {
		if i.Laddr == laddr && i.Lport == lport {
			note.mtx.RUnlock()
			note.mtx.Lock()
			delete(note.CurrentProxy, i.Id)
			note.mtx.Unlock()
			err := f.Reload(note)
			if err != nil {
				return "", 0, err
			}
			return i.Raddr, i.Rport, nil
		}
	}
	note.mtx.RUnlock()
	return "", 0, errors.New("proxy not found")
}

func (f *FrpAdmin) Reload(note *FrpcNote) error {
	//generate config file content
	content := note.OriginalConfig
	for _, i := range note.CurrentProxy {
		content += GenerateConfigContent(&i)
	}

	//write config file
	_, err := PutFrpConfig(note, content)
	if err != nil {
		return err
	}

	//time.Sleep(time.Millisecond * 50)

	//reload
	_, err = ReloadFrp(note)
	if err != nil {
		return err
	}
	return nil
}

func (f *FrpAdmin) CheckAndSync(note *FrpcNote) (err error) {
	//cycle check until proxy status is not wait start
	started := false
	note.mtx.Lock()
	for !started {
		result := GetFrpStatus(note)
		note.CurrentProxy = make(map[string]Proxy)
		for _, i := range result.Tcp {
			if i.Status == "wait start" {
				started = false
				break
			}
			started = true
			if i.Status == "start error" {
				err = errors.New("unavailable proxy")
				continue
			}
			note.CurrentProxy[i.Name] = Proxy{
				Id:    i.Name,
				Laddr: i.LocalAddress(),
				Lport: i.LocalPort(),
				Raddr: i.RemoteAddress(),
				Rport: i.RemotePort(),
				Type:  i.Type,
			}
		}
	}
	note.mtx.Unlock()
	f.Reload(note)

	return err
}

var defaultAdmin = FrpAdmin{}

func init() {
	rand.Seed(time.Now().Unix())
}

func AutoLaunchProxy(laddr string, lport int, protocol string, raddr string, rport int) (string, int, error) {
	if len(globalConfig.ClientNotes) == 0 {
		return "", 0, errors.New("no client note")
	}
	// find raddr's note
	idx := -1
	for i, v := range globalConfig.ClientNotes {
		if v.RAddress == raddr {
			idx = i
			break
		}
	}

	if idx == -1 {
		return "", 0, errors.New("no client note")
	}

	err := LaunchProxy(&globalConfig.ClientNotes[idx], laddr, lport, raddr, rport, protocol)
	if err != nil {
		return "", 0, err
	}
	return raddr, rport, nil
}

func LaunchProxy(node *FrpcNote, laddr string, lport int, raddr string, rport int, protocol string) error {
	return defaultAdmin.LaunchProxy(node, laddr, lport, raddr, rport, protocol)
}

func StopProxy(laddr string, lport int) (string, int, error) {
	//for each to find proxy
	for i := range globalConfig.ClientNotes {
		raddr, rport, err := defaultAdmin.StopProxy(&globalConfig.ClientNotes[i], laddr, lport)
		if err == nil {
			return raddr, rport, nil
		}
		if err.Error() != "proxy not found" {
			return "", 0, err
		}
	}
	return "", 0, errors.New("proxy not found")
}

func GetProxies() []Proxy {
	var proxies []Proxy
	for i := range globalConfig.ClientNotes {
		result := GetFrpStatus(&globalConfig.ClientNotes[i])
		for _, j := range result.Tcp {
			proxy := Proxy{
				Id:    j.Name,
				Laddr: j.LocalAddress(),
				Lport: j.LocalPort(),
				Raddr: j.RemoteAddress(),
				Rport: j.RemotePort(),
				Type:  j.Type,
			}
			proxies = append(proxies, proxy)
		}
		//sync
		globalConfig.ClientNotes[i].mtx.Lock()
		globalConfig.ClientNotes[i].CurrentProxy = make(map[string]Proxy)
		for _, j := range result.Tcp {
			globalConfig.ClientNotes[i].CurrentProxy[j.Name] = Proxy{
				Id:    j.Name,
				Laddr: j.LocalAddress(),
				Lport: j.LocalPort(),
				Raddr: j.RemoteAddress(),
				Rport: j.RemotePort(),
				Type:  j.Type,
			}
		}
		globalConfig.ClientNotes[i].mtx.Unlock()
	}
	return proxies
}
