package ServerEndpoint

import (
	"bufio"
	"bytes"
	"github.com/kxg3030/shermie-nat-traversal/Concrate"
	"github.com/kxg3030/shermie-nat-traversal/Log"
	"net"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

type Server struct {
	port      string
	clients   *sync.Map
	browser   *sync.Map
	password  string
	clientId  int
	browserId int
}

func NewServer(port, password string) *Server {
	return &Server{
		port:     port,
		clients:  &sync.Map{},
		browser:  &sync.Map{},
		password: password,
	}
}

func (i *Server) Run() {
	listener, err := net.Listen("tcp", i.port)
	if err != nil {
		Log.Log.Println("创建tcp服务失败：" + err.Error())
		return
	}
	Log.Log.Println("服务监听地址：" + i.port)
	for {
		connect, err := listener.Accept()
		if err != nil {
			if e, ok := err.(net.Error); ok && e.Temporary() {
				Log.Log.Println("接受连接失败：" + err.Error())
				time.Sleep(time.Second / 20)
			} else {
				Log.Log.Println("接受连接失败：" + err.Error())
			}
			continue
		}
		go i.Handle(NewConnPeer(connect))
	}
}

func (i *Server) RemoveClient(peer *ConnPeer) {
	i.clients.Delete(peer.uniqueId)
}

func (i *Server) RemoveBrowser(peer *ConnPeer) {
	_ = peer.conn.Close()
	i.browser.Delete(peer.uniqueId)
}

func (i *Server) Handle(peer *ConnPeer) {
	// 默认tcp读取超时时间10秒
	_ = peer.conn.SetReadDeadline(time.Now().Add(time.Second * 10))
	// 读取第一个字节
	first, err := peer.readerWriter.Peek(1)
	if err != nil {
		Log.Log.Println("读取客户端数据错误-1：" + err.Error())
		peer.Close()
		return
	}
	// 0x1|0x2是内网客户端请求数据命令,其他的默认为外部浏览器请求
	if int(first[0]) != Concrate.ActionData && int(first[0]) != Concrate.ActionConnect {
		// 读取外部http请求数据
		i.HandleBrowserRequest(peer)
		return
	}
	defer func() {
		_ = peer.conn.Close()
	}()
	_ = peer.conn.SetReadDeadline(time.Time{})
	for {
		message, err := peer.ReadMessage()
		if err != nil {
			Log.Log.Println("读取客户端数据错误-2：" + err.Error())
			return
		}
		Log.Log.Printf("服务器收到客户端消息：%s 客户端标识：%d 数据类型：%d 发送标识：%d", message.DataString(), message.UniqueId, message.Action, message.To)
		switch message.Action {
		// 密码校验
		case Concrate.ActionConnect:
			if string(message.Data) != i.password {
				err := peer.WriteMessage([]byte("密码错误"), -1, -1)
				if err != nil {
					Log.Log.Println("服务器发送数据失败：" + err.Error())
					return
				}
			}
			// 每个内网客户端颁发独立id
			atomic.AddUint32((*uint32)(unsafe.Pointer(&i.clientId)), 1)
			peer.SetUniqueId(i.clientId)
			i.clients.Store(i.clientId, peer)
			_ = peer.WriteMessage([]byte("success"), i.clientId, -1)
			Log.Log.Println("当前客户端列表：" + strconv.Itoa(i.ClientLen()))
			break
		// 数据交换
		case Concrate.ActionData:
			// 判断客户端是否通过验证
			_, ok := i.clients.Load(int(message.UniqueId))
			if !ok {
				_ = peer.WriteMessage([]byte("非法请求"), int(message.UniqueId), -1)
				i.RemoveClient(peer)
				return
			}
			i.HandleBusinessRequest(message)
			break
		}
	}
}

// 处理外网浏览器请求
func (i *Server) HandleBrowserRequest(peer *ConnPeer) {
	var response http.Response
	var storagePeer *ConnPeer
	atomic.AddUint32((*uint32)(unsafe.Pointer(&i.browserId)), 1)
	bodyByte := make([]byte, peer.readerWriter.Reader.Buffered())
	Log.Log.Println("服务器收到浏览器消息：" + string(bodyByte))
	_, err := peer.readerWriter.Reader.Read(bodyByte)
	request, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(bodyByte)))
	if err != nil {
		response = http.Response{
			StatusCode: http.StatusBadGateway,
			Close:      true,
		}
		_ = response.Write(peer.conn)
		return
	}
	// 获取客户端uniqueId
	uniqueId := request.URL.Query().Get("unique-id")
	if uniqueId == "" {
		response = http.Response{
			StatusCode: http.StatusNotFound,
		}
		_ = response.Write(peer.conn)
		return
	}
	uniqueIdInt, _ := strconv.Atoi(uniqueId)
	val, ok := i.clients.Load(uniqueIdInt)
	if !ok {
		response = http.Response{
			StatusCode: http.StatusNotFound,
		}
		_ = response.Write(peer.conn)
		return
	}
	storagePeer = val.(*ConnPeer)
	_ = storagePeer.WriteMessage(bodyByte, storagePeer.uniqueId, i.browserId)
	i.browser.Store(i.browserId, peer)
}

// 处理内网客户端返回的数据
func (i *Server) HandleBusinessRequest(message *Concrate.Message) {
	var storagePeer *ConnPeer
	val, ok := i.browser.Load(int(message.To))
	if !ok {
		Log.Log.Println("客户端返回浏览器数据错误：浏览器连接不存在")
		return
	}
	storagePeer = val.(*ConnPeer)
	defer func() {
		i.RemoveBrowser(storagePeer)
	}()
	if storagePeer != nil {
		defer func() {
			i.RemoveBrowser(storagePeer)
		}()
		response, _ := http.ReadResponse(bufio.NewReader(bytes.NewReader(message.Data)), new(http.Request))
		if response != nil {
			_ = response.Write(storagePeer.conn)
		}
	}
}

// 内网客户端数量
func (i *Server) ClientLen() int {
	total := 0
	i.clients.Range(func(key, value interface{}) bool {
		total++
		return true
	})
	return total
}
