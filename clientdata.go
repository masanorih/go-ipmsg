package ipmsg

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

type ClientData struct {
	Version   int
	PacketNum int
	User      string
	Host      string
	Command   Message
	Option    string
	Nick      string
	Group     string
	PeerAddr  string
	PeerPort  string
	ListAddr  string
	Time      time.Time
	PubKey    string
	Encrypt   bool
	Attach    bool
}

func NewClientData(msg string, addr *net.Addr) *ClientData {
	clientdata := &ClientData{}
	if msg != "" {
		clientdata.Parse(msg)
	}
	return clientdata
}

func (c *ClientData) Parse(msg string) {
	s := strings.SplitN(msg, ":", 6)
	//fmt.Println(s)
	c.Version, _ = strconv.Atoi(s[0])
	c.PacketNum, _ = strconv.Atoi(s[1])
	c.User = s[2]
	c.Host = s[3]
	cmd, _ := strconv.Atoi(s[4])
	c.Command = Message(cmd)
	c.Option = s[5]
	c.Time = time.Now()
	//pp.Println(c)
	c.UpdateNick()
}

func (c *ClientData) UpdateNick() {
	msg := c.Command
	mode := msg.Mode()
	if mode == BR_ENTRY || mode == ANSENTRY {
		if strings.Contains(c.Option, "\x00") {
			s := strings.SplitN(c.Option, "\x00", 2)
			c.Nick = s[0]
			c.Group = strings.Trim(s[1], "\x00")
		}
		if msg.Get(ENCRYPT) {
			c.Encrypt = true
		}
	}
}

func (c ClientData) NickName() string {
	nick := "noname"
	if c.Nick != "" {
		nick = c.Nick
	} else if c.User != "" {
		nick = c.User
	}
	group := "nogroup"
	if c.Group != "" {
		group = c.Group
	} else if c.Host != "" {
		group = c.Host
	}
	return fmt.Sprintf("%s@%s", nick, group)
}

func (c ClientData) Key() string {
	return fmt.Sprintf("%s@%s:%s", c.User, c.PeerAddr, c.PeerPort)
}
