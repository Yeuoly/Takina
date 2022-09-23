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
		if (i.Laddr == laddr && i.Lport == lport) || (i.Raddr == raddr && i.Rport == rport) {
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

func (f *FrpAdmin) StopProxy(note *FrpcNote, laddr string, lport int) error {
	note.mtx.RLock()
	for _, i := range note.CurrentProxy {
		if i.Laddr == laddr && i.Lport == lport {
			note.mtx.RUnlock()
			note.mtx.Lock()
			delete(note.CurrentProxy, i.Id)
			note.mtx.Unlock()
			err := f.Reload(note)
			if err != nil {
				return err
			}
			return nil
		}
	}
	note.mtx.RUnlock()
	return errors.New("proxy not found")
}

func (f *FrpAdmin) Reload(note *FrpcNote) error {
	//generate config file content
	content := note.OriginalConfig
	for _, i := range note.CurrentProxy {
		content += i.GenerateConfigContent()
	}

	//write config file
	_, err := PutFrpConfig(note, content)
	if err != nil {
		return err
	}

	//wait 1s to ensure config file is written
	time.Sleep(time.Second)

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

func AutoLaunchProxy(laddr string, lport int, protocol string) (string, int, error) {
	raddr := ""
	rport := 0
	idx := rand.Int31() % int32(len(globalConfig.ClientNotes))
	raddr = globalConfig.ClientNotes[idx].RAddress
	findport := false
	for !findport {
		rport = rand.Intn(40000-25590) + 25590
		err := LaunchProxy(laddr, lport, raddr, rport, protocol)
		if err == nil {
			findport = true
		}
		if err != nil && (err.Error() != "proxy already exists" || err.Error() != "unavailable proxy") {
			return "", 0, err
		}
	}
	return raddr, rport, nil
}

func LaunchProxy(laddr string, lport int, raddr string, rport int, protocol string) error {
	idx := rand.Int31() % int32(len(globalConfig.ClientNotes))
	note := &globalConfig.ClientNotes[idx]
	return defaultAdmin.LaunchProxy(note, laddr, lport, raddr, rport, protocol)
}

func StopProxy(laddr string, lport int) error {
	//for each to find proxy
	for i := range globalConfig.ClientNotes {
		err := defaultAdmin.StopProxy(&globalConfig.ClientNotes[i], laddr, lport)
		if err == nil {
			return nil
		}
		if err.Error() != "proxy not found" {
			return err
		}
	}
	return errors.New("proxy not found")
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
