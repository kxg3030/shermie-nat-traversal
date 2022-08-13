package ClientEndPoint

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/kxg3030/shermie-nat-traversal/Concrate"
	Log "github.com/kxg3030/shermie-nat-traversal/Log"
	"io"
	"net"
	"net/http"
	"time"
)

type Client struct {
	host     string
	port     string
	peer     *ConnPeer
	uniqueId int
}

func NewClient(host, port string) *Client {
	return &Client{
		host: host,
		port: port,
	}
}

func (i *Client) Run() {
	connect, err := net.Dial("tcp", fmt.Sprintf("%s:%s", i.host, i.port))
	if err != nil {
		Log.Log.Println("连接服务器失败：" + err.Error())
		return
	}
	// 发送密码验证
	i.peer = NewConnPeer(connect)
	_ = i.ConnectAuth()
	// 读取消息
	message, err := i.peer.ReadMessage()
	if err != nil {
		Log.Log.Println("服务器验证错误：" + err.Error())
		return
	}
	i.uniqueId = int(message.UniqueId)
	Log.Log.Printf("服务器返回数据：%s 标识ID：%d\n", message.DataString(), i.uniqueId)
	if string(message.Data) != "success" {
		return
	}
	_ = connect.SetDeadline(time.Time{})
	for {
		message, err = i.peer.ReadMessage()
		if err != nil {
			Log.Log.Println("读取服务器数据错误：" + err.Error())
			i.peer.Close()
			return
		}
		go func() {
			Log.Log.Println("服务器返回数据：" + message.DataString())
			var response http.Response
			conn, err := net.DialTimeout("tcp", "47.100.10.87:80", time.Second*10)
			if err != nil {
				Log.Log.Println("转发请求连接失败：" + err.Error())
				response.StatusCode = http.StatusBadGateway
				_ = i.peer.WriteMessage([]byte("HTTP/1.1 502 Bad Gateway\r\n\r\n"), Concrate.ActionData, i.uniqueId, int(message.To))
				return
			}
			_, err = io.Copy(conn, bytes.NewReader(message.Data))
			reader := bufio.NewReader(conn)
			responseBuffer := make([]byte, 10*1024)
			n, _ := reader.Read(responseBuffer)
			responseBuffer = responseBuffer[:n]
			_ = i.peer.WriteMessage(responseBuffer, Concrate.ActionData, i.uniqueId, int(message.To))
		}()
	}
}

func (i *Client) ConnectAuth() error {
	return i.peer.WriteMessage([]byte("123456"), Concrate.ActionConnect, -1, -1)
}
