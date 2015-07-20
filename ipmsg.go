package ipmsg

import (
	"fmt"
	"net"
	"os"
)

type IPMSG struct {
	ClientData ClientData
	Conn       *net.UDPConn
}

type IPMSGConfig struct {
	NickName  string
	GroupName string
	UserName  string
	HostName  string
}

const (
	DefaultPort int = 2425
)

func NewIPMSG(conf *IPMSGConfig) *IPMSG {
	ipmsg := &IPMSG{}
	// UDP server
	service := fmt.Sprintf(":%s", DefaultPort)
	udpAddr, err := net.ResolveUDPAddr("udp4", service)
	checkError(err)
	conn, err := net.ListenUDP("udp", udpAddr)
	checkError(err)
	ipmsg.Conn = conn
	return ipmsg
	//for {
	//	handleClient(conn)
	//}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal Error ", err.Error())
		os.Exit(1)
	}
}

/*
	func handleClient(conn *net.UDPConn) {
		var buf [512]byte
		_, addr, err := conn.ReadFromUDP(buf[0:])
		if err != nil {
			return
		}
		daytime := time.Now().String()
		conn.WriteToUDP([]byte(daytime), addr)
	}

	// UDP client
    service := os.Args[1]
    udpAddr, err := net.ResolveUDPAddr("udp4", service)
    checkError(err)
    conn, err := net.DialUDP("udp", nil, udpAddr)
    checkError(err)
    _, err = conn.Write([]byte("anything"))
    checkError(err)
    var buf [512]byte
    n, err := conn.Read(buf[0:])
    checkError(err)
    fmt.Println(string(buf[0:n]))
    os.Exit(0)
*/
