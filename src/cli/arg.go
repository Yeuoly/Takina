package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/table"
)

var (
	mode         = flag.String("mode", "none", "start/stop/list")
	takina_port  = flag.Int("takina_port", 40002, "takina port")
	takina_token = flag.String("takina_token", "none", "takina token")
)

var (
	laddr = flag.String("laddr", "none", "local address")
	lport = flag.Int("lport", 0, "local port")
	typ   = flag.String("type", "none", "proxy type, tcp/udp")
)

func checkArgs(requires []string) {
	for _, arg := range requires {
		switch arg {
		case "laddr":
			if *laddr == "none" {
				panic("laddr is required")
			}
		case "lport":
			if *lport == 0 {
				panic("lport is required")
			}
		case "type":
			if *typ == "none" {
				panic("type is required")
			}
		}
	}
}

func Run() {
	// parse args
	flag.Parse()
	switch *mode {
	case "start":
		checkArgs([]string{"laddr", "lport", "type"})
		proxy, err := startProxy(*laddr, *lport, *typ)
		if err != nil {
			panic(err)
		}
		fmt.Printf("proxy started: %s:%d -> %s:%d\n", proxy.Laddr, proxy.Lport, proxy.Raddr, proxy.Rport)
	case "stop":
		checkArgs([]string{"laddr", "lport"})
		err := stopProxy(*laddr, *lport)
		if err != nil {
			panic(err)
		}
		fmt.Printf("proxy stopped: %s:%d\n", *laddr, *lport)
	case "list":
		proxies, err := listProxy()
		if err != nil {
			panic(err)
		}
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"#", "Laddr", "Lport", "Raddr", "Rport", "Type"})
		for i, proxy := range proxies {
			t.AppendRow(table.Row{i, proxy.Laddr, proxy.Lport, proxy.Raddr, proxy.Rport, proxy.Type})
		}
		t.Render()
	default:
		panic("invalid mode")
	}
}
