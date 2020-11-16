package main

import (
	"fmt"
	"net"
)

type pkg struct {
	sender string
	msg    string
}
type Client struct {
	C    chan pkg
	Name string
	Addr string
}

var online map[string]Client
var rename map[string]bool
var message = make(chan pkg)

func sendmsg() {
	online = make(map[string]Client)
	for {
		msg := <-message
		for _, cli := range online {
			cli.C <- msg
		}
	}
}

func writetoclient(cli Client, conn net.Conn) {
	for mpkg := range cli.C {
		if mpkg.sender == online[cli.Addr].Name {
			continue
		}
		fmt.Println(mpkg.sender, "发言")
		conn.Write([]byte(mpkg.sender + "：" + mpkg.msg))
	}
}

func Handle(conn net.Conn) {
	defer conn.Close()
	flag := false
	cliAddr := conn.RemoteAddr().String()
	cli := Client{make(chan pkg), cliAddr, cliAddr}
	online[cliAddr] = cli
	tmp := pkg{"system", "[" + string(cliAddr) + "]login"}
	message <- tmp

	go writetoclient(cli, conn)

	go func() {
		buf := make([]byte, 2048)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				fmt.Println("conn.read err:", err)
				return
			}
			msg := string(buf[:n-1])
			tmp := pkg{online[cliAddr].Name, msg}
			if len(msg) >= 8 && msg[:7] == "rename|" {
				newname := string(buf[7 : n-2])
				if rename[newname] {
					conn.Write([]byte("改昵称已被占用！\n"))
				} else {
					rename[cli.Name] = false
					cli.Name = newname
					online[cliAddr] = cli
					conn.Write([]byte("改名成功！"))
					rename[newname] = true
				}
			} else if len(msg) == 4 && msg[:3] == "who" {
				conn.Write([]byte("当前在线："))
				for _, nowcli := range online {
					conn.Write([]byte(nowcli.Name + "\n"))
				}

			} else if len(msg) == 7 && msg[:6] == "logout" {
				message <- pkg{"system", "用户[" + cli.Name + "]logout"}
				delete(online, cliAddr)
				flag = true
			} else {
				message <- tmp
			}

		}
	}()
	for {
		if flag {
			return
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:8000")
	if err != nil {
		fmt.Println("listen err:", err)
		return
	}
	defer listener.Close()
	//转发消息的携程
	rename = make(map[string]bool)
	go sendmsg()
	//主携程，阻塞等待用户连接
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept err:", err)
			continue
		}
		go Handle(conn) //处理用户连接
	}
}
