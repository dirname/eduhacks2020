package websocket

import (
	"eduhacks2020/Go/database"
	"eduhacks2020/Go/define/retcode"
	"eduhacks2020/Go/pkg/setting"
	"eduhacks2020/Go/utils"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

// ClientManager 连接管理
type ClientManager struct {
	ClientIDMap     map[string]*Client // 全部的连接
	ClientIDMapLock sync.RWMutex       // 读写锁

	Connect    chan *Client // 连接处理
	DisConnect chan *Client // 断开连接处理

	GroupLock sync.RWMutex
	Groups    map[string][]string

	SystemClientsLock sync.RWMutex
	SystemClients     map[string][]string
}

// NewClientManager 新的实例
func NewClientManager() (clientManager *ClientManager) {
	clientManager = &ClientManager{
		ClientIDMap:   make(map[string]*Client),
		Connect:       make(chan *Client, 10000),
		DisConnect:    make(chan *Client, 10000),
		Groups:        make(map[string][]string, 100),
		SystemClients: make(map[string][]string, 100),
	}

	return
}

// Start 管道处理程序
func (manager *ClientManager) Start() {
	for {
		select {
		case client := <-manager.Connect:
			// 建立连接事件
			manager.EventConnect(client)
		case conn := <-manager.DisConnect:
			// 断开连接事件
			manager.EventDisconnect(conn)
		}
	}
}

// EventConnect 建立连接事件
func (manager *ClientManager) EventConnect(client *Client) {
	manager.AddClient(client)
	client.Manager.Online(&database.ClientDevice{
		ClientID:    client.ClientID,
		SystemID:    client.SystemID,
		ConnectTime: client.ConnectTime,
		IsDeleted:   client.IsDeleted,
		UserName:    client.UserName,
		NickName:    client.NickName,
		UserRole:    client.UserRole,
	})
	log.WithFields(log.Fields{
		"host":     setting.GlobalSetting.LocalHost,
		"port":     setting.CommonSetting.Port,
		"clientId": client.ClientID,
		"systemId": client.SystemID,
		"counts":   Manager.Count(),
	}).Info("client connected")
}

// EventDisconnect 断开连接时间
func (manager *ClientManager) EventDisconnect(client *Client) {
	//关闭连接
	_ = client.Socket.Close()
	manager.DelClient(client)

	mJSON, _ := json.Marshal(map[string]string{
		"clientId": client.ClientID,
		"userId":   client.UserID,
		"extend":   client.Extend,
	})
	data := string(mJSON)
	sendUserID := ""

	//发送下线通知
	if len(client.GroupList) > 0 {
		for _, groupName := range client.GroupList {
			SendMessage2Group(client.SystemID, sendUserID, groupName, retcode.OfflineMessageCode, "客户端下线", &data)
		}
	}

	log.WithFields(log.Fields{
		"host":     setting.GlobalSetting.LocalHost,
		"port":     setting.CommonSetting.Port,
		"clientId": client.ClientID,
		"counts":   Manager.Count(),
		"seconds":  uint64(time.Now().Unix()) - client.ConnectTime,
	}).Info("client disconnected")
	client.Manager.Offline(client.ClientID)
	//标记销毁
	client.IsDeleted = true
	client = nil
}

// AddClient 添加客户端
func (manager *ClientManager) AddClient(client *Client) {
	manager.ClientIDMapLock.Lock()
	defer manager.ClientIDMapLock.Unlock()

	manager.ClientIDMap[client.ClientID] = client
}

// AllClient 获取所有的客户端
func (manager *ClientManager) AllClient() map[string]*Client {
	manager.ClientIDMapLock.RLock()
	defer manager.ClientIDMapLock.RUnlock()

	return manager.ClientIDMap
}

// Count 客户端数量
func (manager *ClientManager) Count() int {
	manager.ClientIDMapLock.RLock()
	defer manager.ClientIDMapLock.RUnlock()
	return len(manager.ClientIDMap)
}

// DelClient 删除客户端
func (manager *ClientManager) DelClient(client *Client) {
	manager.delClientIDMap(client.ClientID)

	//删除所在的分组
	if len(client.GroupList) > 0 {
		for _, groupName := range client.GroupList {
			manager.delGroupClient(utils.GenGroupKey(client.SystemID, groupName), client.ClientID)
		}
	}

	// 删除系统里的客户端
	manager.delSystemClient(client)
}

// 删除clientIdMap
func (manager *ClientManager) delClientIDMap(clientID string) {
	manager.ClientIDMapLock.Lock()
	defer manager.ClientIDMapLock.Unlock()

	delete(manager.ClientIDMap, clientID)
}

// GetByClientID 通过clientId获取
func (manager *ClientManager) GetByClientID(clientID string) (*Client, error) {
	manager.ClientIDMapLock.RLock()
	defer manager.ClientIDMapLock.RUnlock()

	client, ok := manager.ClientIDMap[clientID]
	if !ok {
		return nil, errors.New("客户端不存在")
	}
	return client, nil
}

// SendMessage2LocalGroup 发送到本机分组
func (manager *ClientManager) SendMessage2LocalGroup(systemID, messageID, sendUserID, groupName string, code int, msg string, data *string) {
	if len(groupName) > 0 {
		clientIds := manager.GetGroupClientList(utils.GenGroupKey(systemID, groupName))
		if len(clientIds) > 0 {
			for _, clientID := range clientIds {
				if _, err := Manager.GetByClientID(clientID); err == nil {
					//添加到本地
					SendMessage2LocalClient(messageID, clientID, sendUserID, code, msg, data)
				} else {
					//删除分组
					manager.delGroupClient(utils.GenGroupKey(systemID, groupName), clientID)
				}
			}
		}
	}
}

// SendMessage2LocalSystem 发送给指定业务系统
func (manager *ClientManager) SendMessage2LocalSystem(systemID, messageID string, sendUserID string, code int, msg string, data *string) {
	if len(systemID) > 0 {
		clientIds := Manager.GetSystemClientList(systemID)
		if len(clientIds) > 0 {
			for _, clientID := range clientIds {
				SendMessage2LocalClient(messageID, clientID, sendUserID, code, msg, data)
			}
		}
	}
}

// AddClient2LocalGroup 添加到本地分组
func (manager *ClientManager) AddClient2LocalGroup(groupName string, client *Client, userID string, extend string) {
	//标记当前客户端的userId
	client.UserID = userID
	client.Extend = extend

	//判断之前是否有添加过
	for _, groupValue := range client.GroupList {
		if groupValue == groupName {
			return
		}
	}

	// 为属性添加分组信息
	groupKey := utils.GenGroupKey(client.SystemID, groupName)

	manager.addClient2Group(groupKey, client)

	client.GroupList = append(client.GroupList, groupName)

	mJSON, _ := json.Marshal(map[string]string{
		"clientId": client.ClientID,
		"userID":   client.UserID,
		"extend":   client.Extend,
	})
	data := string(mJSON)
	sendUserID := ""

	//发送系统通知
	SendMessage2Group(client.SystemID, sendUserID, groupName, retcode.OnlineMessageCode, "客户端上线", &data)
}

// addClient2Group 添加到本地分组
func (manager *ClientManager) addClient2Group(groupKey string, client *Client) {
	manager.GroupLock.Lock()
	defer manager.GroupLock.Unlock()
	manager.Groups[groupKey] = append(manager.Groups[groupKey], client.ClientID)
}

// 删除分组里的客户端
func (manager *ClientManager) delGroupClient(groupKey string, clientID string) {
	manager.GroupLock.Lock()
	defer manager.GroupLock.Unlock()

	for index, groupClientID := range manager.Groups[groupKey] {
		if groupClientID == clientID {
			manager.Groups[groupKey] = append(manager.Groups[groupKey][:index], manager.Groups[groupKey][index+1:]...)
		}
	}
}

// GetGroupClientList 获取本地分组的成员
func (manager *ClientManager) GetGroupClientList(groupKey string) []string {
	manager.GroupLock.RLock()
	defer manager.GroupLock.RUnlock()
	return manager.Groups[groupKey]
}

// AddClient2SystemClient 添加到系统客户端列表
func (manager *ClientManager) AddClient2SystemClient(systemID string, client *Client) {
	manager.SystemClientsLock.Lock()
	defer manager.SystemClientsLock.Unlock()
	manager.SystemClients[systemID] = append(manager.SystemClients[systemID], client.ClientID)
}

// 删除系统里的客户端
func (manager *ClientManager) delSystemClient(client *Client) {
	manager.SystemClientsLock.Lock()
	defer manager.SystemClientsLock.Unlock()

	for index, clientID := range manager.SystemClients[client.SystemID] {
		if clientID == client.ClientID {
			manager.SystemClients[client.SystemID] = append(manager.SystemClients[client.SystemID][:index], manager.SystemClients[client.SystemID][index+1:]...)
		}
	}
}

// GetSystemClientList 获取指定系统的客户端列表
func (manager *ClientManager) GetSystemClientList(systemID string) []string {
	manager.SystemClientsLock.RLock()
	defer manager.SystemClientsLock.RUnlock()
	return manager.SystemClients[systemID]
}
