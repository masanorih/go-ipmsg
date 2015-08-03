package main

import (
	"fmt"
	"net"
	"os"
	"os/user"
	"strconv"

	"github.com/Songmu/prompter"
	"github.com/k0kubun/pp"
	goipmsg "github.com/masanorih/go-ipmsg"
)

var commands = []string{"help", "quit", "join", "list", "send", "read"}

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
				Message:    "ipmsg>",
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
	ev.Debug = false
	// those are defined at handler.go
	ev.Regist(goipmsg.BR_ENTRY, RECEIVE_BR_ENTRY)
	ev.Regist(goipmsg.ANSENTRY, RECEIVE_ANSENTRY)
	ev.Regist(goipmsg.SENDMSG, RECEIVE_SENDMSG)
	ipmsg.AddEventHandler(ev)

	return ipmsg
}

// Switchinput dispatches actions via input
func SwitchInput(ipmsg *goipmsg.IPMSG, input string, quit chan string) {
	switch input {
	case "help":
		fmt.Println("usage:")
		fmt.Println("\thelp: show this help.")
		fmt.Println("\tquit: quit this programme.")
		fmt.Println("\tjoin: let others know U.")
		fmt.Println("\tlist: list up users U know.")
		fmt.Println("\tsend: send message to user.")
		fmt.Println("\tread: read received message.")
	case "quit":
		fmt.Println("quitting...")
		ipmsg.Close()
		quit <- "quitting"
	case "join":
		Join(ipmsg)
	case "list":
		List(ipmsg)
	case "send":
		Send(ipmsg)
	case "read":
		Read(ipmsg)
	}
}

func Read(ipmsg *goipmsg.IPMSG) {
	if len(Messages) == 0 {
		fmt.Println("There is no message to read.")
	} else {
		for _, v := range Messages {
			fmt.Printf("From: %v\n", v.Key())
			fmt.Printf("Date: %v\n", v.Time)
			fmt.Printf("Message: %v\n\n", v.Option)
		}
		// clear all datas
		Messages = []*goipmsg.ClientData{}
	}
}

// Send get specified user from prompt and actually send it to the target
func Send(ipmsg *goipmsg.IPMSG) {
	userIdx := []string{}
	i := 0
	m := make(map[int]string)
	for k, _ := range Users {
		i++
		fmt.Printf("%d %v\n", i, k)
		userIdx = append(userIdx, strconv.Itoa(i))
		m[i] = k
	}

	chosen := (&prompter.Prompter{
		Choices:    userIdx,
		UseDefault: false,
		Message:    "Choose the user to send message>",
		IgnoreCase: true,
	}).Prompt()
	i, _ = strconv.Atoi(chosen)
	key := m[i]

	promptMessage := fmt.Sprintf("Enter message(to %v)", key)
	message := prompter.Prompt(promptMessage, "")
	cd := Users[key]
	addr := cd.Addr

	cmd := goipmsg.SENDMSG
	cmd.SetOpt(goipmsg.SECRET)
	err := ipmsg.SendMSG(addr, message, cmd)
	if err != nil {
		panic(err)
	}
	fmt.Println("sent SENDMSG")
}

// List print out known users
func List(ipmsg *goipmsg.IPMSG) {
	if len(Users) == 0 {
		fmt.Println("There is no users.")
	} else {
		fmt.Println("known users below.")
		i := 0
		for k, _ := range Users {
			i++
			fmt.Printf("\t%d %v\n", i, k)
		}
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
