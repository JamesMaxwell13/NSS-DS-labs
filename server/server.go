package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
)

type Config struct {
	Protocol string
	Adress   string
}

type Server interface {
	ConfigurateServer(conn net.Conn) (net.Conn, error)
	GetConfig(addr net.Addr) Config
}

type TcpServer struct {
	config Config
}

func (ts *TcpServer) GetConfig(addr net.Addr) Config {
	if addr != nil {
		ts.config = Config{"adress", addr.String()}
	} else {
		ts.config = Config{"adress", "0.0.0.0:8080"}
	}
	return ts.config
}

func (ts *TcpServer) ConfigurateServer(conn net.Conn) (net.Conn, error) {
	tcpConn, ok := conn.(*net.TCPConn)
	if !ok {
		return nil, errors.New("this connection isn't TCP")
	}
	err := tcpConn.SetKeepAlive(true)
	if err != nil {
		return nil, err
	}
	err = tcpConn.SetKeepAlivePeriod(30 * time.Second)
	if err != nil {
		return nil, err
	}
	return tcpConn, nil
}

func RunServer(s Server) error {
	config := s.GetConfig(nil)
	var conn net.Conn
	l, err := net.Listen(config.Protocol, config.Adress)
	if err != nil {
		return err
	}
	defer func(l net.Listener) {
		err := l.Close()
		if err != nil {
			fmt.Println("Error server shutdown...")
		}
	}(l)
	for {
		if conn == nil {
			conn, err = l.Accept()
			if err != nil {
				return err
			}
			conn, err = s.ConfigurateServer(conn)
			if err != nil {
				return err
			}
		}
		err = RequestHandler(conn)
		if err != nil {
			return err
		}
	}
}

func ServerCLI(conn net.Conn, data string) error {
	command := ParseCommands(data)
	switch command {
	case "ech":
		{
			text := strings.Replace(data, "echo ", "", -1)
			text = strings.TrimLeft(text, " ")
			_, err := conn.Write([]byte(text))
			if err != nil {
				return err
			}
			return nil
		}
	case "tim":
		{
			text := "server time: " + time.Now().String()
			_, err := conn.Write([]byte(text))
			if err != nil {
				return err
			}
			return nil
		}
	case "dwn":
		{
			return nil
		}
	case "upl":
		{
			return nil
		}
	case "cls":
		{
			fmt.Println("closing connection with current client...")
			err := conn.Close()
			if err != nil {
				return err
			}
			return nil
		}
	default:
		_, err := conn.Write([]byte("wrong command..."))
		if err != nil {
			return err
		} else {
			return nil
		}
	}
}

func UploadFile(conn net.Conn) error {
	return nil
}

func DownloadFile(conn net.Conn) error {
	return nil
}

func ParseCommands(data string) string {
	lowStr := strings.ToLower(data)
	switch {
	case strings.HasPrefix(lowStr, "echo"):
		return "ech"
	case strings.HasPrefix(lowStr, "time"):
		return "tim"
	case strings.HasPrefix(lowStr, "close"):
		return "cls"
	case strings.HasPrefix(lowStr, "download"):
		return "dwn"
	case strings.HasPrefix(lowStr, "upload"):
		return "upl"
	default:
		return ""
	}
}

func RequestHandler(conn net.Conn) error {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("Error connection shutdown...")
		}
	}(conn)
	connReader := bufio.NewReader(conn)
	for conn != nil {
		data, err := connReader.ReadString('\n')
		if err != nil {
			return err
		}
		fmt.Printf("message from %s: %s\n\n", conn.RemoteAddr().String(), data)
		err = ServerCLI(conn, data)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	RunServer(&TcpServer{})
}
