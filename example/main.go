package main

import (
	"fmt"
	"net"
	"os"

	"github.com/Songmu/prompter"
	"github.com/k0kubun/pp"
	goipmsg "github.com/masanorih/go-ipmsg"
)

var commands = []string{"help", "quit", "join", "list"}

func main() {
	conf := goipmsg.NewIPMSGConf()
	ipmsg, err := goipmsg.NewIPMSG(conf)
	if err != nil {
		panic(err)
	}
	ev := goipmsg.NewEventHandler()
	ev.Regist(goipmsg.BR_ENTRY, RECEIVE_BR_ENTRY)
	ipmsg.AddEventHandler(ev)

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
			recv <- cd.Option
		}
	}()

	<-quit
	os.Exit(1)
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
	fmt.Println("list up known users")
}

// Join sends BR_ENTRY packet to the broadcast address
func Join(ipmsg *goipmsg.IPMSG) {
	addr := brAddr(ipmsg)
	err := ipmsg.SendMSG(addr, ipmsg.Myinfo(), goipmsg.BR_ENTRY)
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
