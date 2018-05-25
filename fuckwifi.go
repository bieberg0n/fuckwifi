package main

import (
	"net"
	"fmt"
)

func logs(args ...interface{}) {
	fmt.Println(args...)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func reply(remoteConn net.Conn, serv net.PacketConn, cli net.Addr) {
	buf := make([]byte, 2048)
	logs("reply start")
	for {
		n, err := remoteConn.Read(buf)
		//check(err)
		if err != nil {
			logs(err)
			remoteConn.Close()
			return
		}
		logs(remoteConn.RemoteAddr().String(), "=>", cli.String(), "Len:", n)
		serv.WriteTo(buf[:n], cli)
	}
}

func main() {
	remoteStr := "j.bjong.me:5353"

	serv, err := net.ListenPacket("udp", "127.0.0.1:15353")
	check(err)
	defer serv.Close()

	m := map[string]net.Conn{}

	buf := make([]byte, 2048)
	//var buf []byte
	logs("start")

	for {
		n, cli, err := serv.ReadFrom(buf)
		check(err)

		var remoteConn net.Conn
		logs(m, cli)
		_, exist := m[cli.String()]
		if exist {
			//logs("存在旧连接")
			remoteConn = m[cli.String()]
		} else {
			//logs("创建新连接")
			remoteConn, err = net.Dial("udp", remoteStr)
			check(err)

			m[cli.String()] = remoteConn
			remoteConn.Write([]byte("\x1c?\x01 \x00\x01\x00\x00\x00\x00\x00\x01\x05bjong\x02me\x00\x00\x01\x00\x01\x00\x00)\x10\x00\x00\x00\x00\x00\x00\x00"))
			go reply(remoteConn, serv, cli)
		}

		logs(cli.String(), "=>", remoteStr, "Len:", n)
		remoteConn.Write(buf[:n])
	}
}
