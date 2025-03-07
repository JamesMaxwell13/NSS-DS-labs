package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

type Config struct {
	Protocol string
	Adress   string
}

type Client interface {
	GetConfig(addr net.Addr) Config
}

type TcpClient struct {
	config Config
}

func (tc *TcpClient) GetConfig(addr net.Addr) Config {
	if addr != nil {
		tc.config = Config{Protocol: "adress", Adress: addr.String()}
	} else {
		tc.config = Config{Protocol: "adress", Adress: "0.0.0.0:8080"}
	}
	return tc.config
}

func RunClient(c Client) error {
	config := c.GetConfig(nil)
	conn, err := net.Dial(config.Protocol, config.Adress)
	if err != nil {
		return err
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("Error client shutdown...")
		}
	}(conn)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		message := scanner.Text()
		_, err := conn.Write([]byte(message + "\n"))
		if err != nil {
			return err
		}
		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			return err
		}
		fmt.Print("Server response: " + response)
	}
	return nil
}

func ClientCLI(conn net.Conn, data string) error {
	lowStr := strings.ToLower(data)
	switch {
	case strings.HasPrefix(lowStr, "upload"):
		{

			_, err :=
			if err != nil {
				return err
			}
			return nil
		}
	case strings.HasPrefix(lowStr, "download"):
		{
			_, err :=
			if err != nil {
				return err
			}
			return nil
		}
	default:
		return nil
	}
}

func main() {
	RunClient(&TcpClient{})
}
