package main

import (
	"fmt"
	"net"
	"os"
	"os/user"

	"github.com/Songmu/prompter"
	"github.com/k0kubun/pp"
	goipmsg "github.com/masanorih/go-ipmsg"
)

var commands = []string{"help", "quit", "join", "list"}

func main() {
	ipmsg := setup()

	input := make(chan string)
	next := make(chan string)
	quit := make(chan string)
	recv := make(chan string)
	// main loop or receive event from channel
	go func() {
		for {
			var str string
			select {
			case str = <-input:
				SwitchInput(ipmsg, str, quit)
				next <- ""
			case str = <-recv:
				pp.Println("recv=", str)
			}
		}
	}()
	// get command from stdin
	go func() {
		for {
			input <- (&prompter.Prompter{
				Choices:    commands,
				Default:    "list",
				Message:    "ipmsg> ",
				IgnoreCase: true,
			}).Prompt()
			<-next
		}
	}()
	// recv message from socket
	go func() {
		for {
			cd, err := ipmsg.RecvMSG()
			if err != nil {
				panic(err)
			}
			//pp.Println("recv=", cd.String())
			recv <- cd.Option
		}
	}()

	<-quit
	os.Exit(1)
}

func setup() *goipmsg.IPMSG {
	conf := goipmsg.NewIPMSGConf()
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	conf.UserName = user.Username
	conf.NickName = user.Username
	conf.GroupName = user.Gid
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	conf.HostName = hostname

	ipmsg, err := goipmsg.NewIPMSG(conf)
	if err != nil {
		panic(err)
	}

	ev := goipmsg.NewEventHandler()
	// those are defined at handler.go
	ev.Regist(goipmsg.BR_ENTRY, RECEIVE_BR_ENTRY)
	ev.Regist(goipmsg.ANSENTRY, RECEIVE_ANSENTRY)
	ipmsg.AddEventHandler(ev)

	return ipmsg
}

// Switchinput dispatches actions via input
func SwitchInput(ipmsg *goipmsg.IPMSG, input string, quit chan string) {
	switch input {
	case "help":
		fmt.Println("help usage here")
	case "quit":
		fmt.Println("quitting...")
		ipmsg.Close()
		quit <- "quitting"
	case "join":
		Join(ipmsg)
	case "list":
		ListUp(ipmsg)
	}
}

// Listup printout known users
func ListUp(ipmsg *goipmsg.IPMSG) {
	for k, _ := range Userlist {
		//fmt.Printf("%v=%v\n", k, v.String())
		fmt.Println(k)
	}
}

// Join sends BR_ENTRY packet to the broadcast address
func Join(ipmsg *goipmsg.IPMSG) {
	addr := brAddr(ipmsg)
	cmd := goipmsg.BR_ENTRY
	cmd.SetOpt(goipmsg.BROADCAST)
	err := ipmsg.SendMSG(addr, ipmsg.Myinfo(), cmd)
	if err != nil {
		panic(err)
	}
	fmt.Println("sent BR_ENTRY")
}

// brAddr retrieves broadcast address
func brAddr(ipmsg *goipmsg.IPMSG) *net.UDPAddr {
	//broadcast := net.IPv4(203, 181, 79, 127) //net.IP
	broadcast := net.IPv4(255, 255, 255, 255) //net.IP
	port := ipmsg.Conf.Port
	str := fmt.Sprintf("%v:%d", broadcast.String(), port)
	//fmt.Println("str=", str)
	udpAddr, err := net.ResolveUDPAddr("udp4", str)
	if err != nil {
		panic(err)
	}
	return udpAddr
}
