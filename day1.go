package main

import (
	"fmt"
	
	"net"
)

func main(){
	//1 监听端口
	listener, err := net.Listen("tcp", "0.0.0.0:9091")
	if err != nil {
		fmt.Println("Listen fail, err: %v \n", err)
		return
	}
	
	//2 建立套接字连接
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Accept fail, err: %v\n", err)
			continue
		}
		
		//创建处理协程
		go func(conn net.Conn) {
			for {
				var buf [128]byte
				n , err := conn.Read(buf[:])
				if err != nil {
					fmt.Printf("Read from connect failed, err: %v\n", err)
					break
				}
				str := string(buf[:n])
				fmt.Printf("recieve from client, data : %v\n", str)
			}
		}(conn)
	}
}
