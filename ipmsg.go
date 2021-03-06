package ipmsg

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"time"
)

type IPMSG struct {
	ClientData ClientData
	Conn       *net.UDPConn
	Conf       *IPMSGConfig
	Handlers   []*EventHandler
	PacketNum  int
}

type IPMSGConfig struct {
	NickName  string
	GroupName string
	UserName  string
	HostName  string
	Port      int
	Local     string
}

const (
	DefaultPort int = 2425
	Buflen      int = 65535
)

func NewIPMSGConf() *IPMSGConfig {
	conf := &IPMSGConfig{
		Port: DefaultPort,
	}
	return conf
}

func NewIPMSG(conf *IPMSGConfig) (*IPMSG, error) {
	ipmsg := &IPMSG{
		PacketNum: 0,
	}
	ipmsg.Conf = conf
	// UDP server
	service := fmt.Sprintf("%v:%d", conf.Local, conf.Port)
	//fmt.Println("service =", service)
	udpAddr, err := net.ResolveUDPAddr("udp", service)
	if err != nil {
		return ipmsg, err
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return ipmsg, err
	}
	ipmsg.Conn = conn
	return ipmsg, err
}

func (ipmsg *IPMSG) Close() error {
	conn := ipmsg.Conn
	if conn == nil {
		err := errors.New("Conn is not defined")
		return err
	}
	err := conn.Close()
	return err
}

func (ipmsg *IPMSG) BuildData(addr *net.UDPAddr, msg string, cmd Command) *ClientData {
	conf := ipmsg.Conf
	clientdata := NewClientData("", addr)
	clientdata.Version = 1
	clientdata.PacketNum = ipmsg.GetNewPacketNum()
	clientdata.User = conf.UserName
	clientdata.Host = conf.HostName
	clientdata.Command = cmd
	clientdata.Option = msg
	return clientdata
}

func (ipmsg *IPMSG) SendMSG(addr *net.UDPAddr, msg string, cmd Command) error {
	clientdata := ipmsg.BuildData(addr, msg, cmd)
	conn := ipmsg.Conn
	_, err := conn.WriteToUDP([]byte(clientdata.String()), addr)
	if err != nil {
		return err
	}
	return nil
}

func (ipmsg *IPMSG) RecvMSG() (*ClientData, error) {
	var buf [Buflen]byte
	conn := ipmsg.Conn
	_, addr, err := conn.ReadFromUDP(buf[0:])
	if err != nil {
		return nil, err
	}
	trimmed := bytes.Trim(buf[:], "\x00")
	clientdata := NewClientData(string(trimmed[:]), addr)

	handlers := ipmsg.Handlers
	for _, v := range handlers {
		err := v.Run(clientdata, ipmsg)
		if err != nil {
			return clientdata, err
		}
	}
	return clientdata, nil
}

// convert net.Addr to net.UDPAddr
func (ipmsg *IPMSG) UDPAddr() (*net.UDPAddr, error) {
	conn := ipmsg.Conn
	if conn == nil {
		err := errors.New("Conn is not defined")
		return nil, err
	}
	addr := conn.LocalAddr()
	network := addr.Network()
	str := addr.String()
	//fmt.Println("str =", str)
	udpAddr, err := net.ResolveUDPAddr(network, str)
	return udpAddr, err
}

func (ipmsg *IPMSG) AddEventHandler(ev *EventHandler) {
	sl := ipmsg.Handlers
	sl = append(sl, ev)
	ipmsg.Handlers = sl
}

func (ipmsg *IPMSG) GetNewPacketNum() int {
	ipmsg.PacketNum++
	return int(time.Now().Unix()) + ipmsg.PacketNum
}

func (ipmsg *IPMSG) Myinfo() string {
	conf := ipmsg.Conf
	return fmt.Sprintf("%v\x00%v\x00", conf.NickName, conf.GroupName)
}
