package websocket

import (
	"eduhacks2020/Go/pkg/setting"
	"eduhacks2020/Go/utils"
	"github.com/gorilla/websocket"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAddClient(t *testing.T) {
	clientID := "clientID"
	systemID := "publishSystem"
	var manager = NewClientManager() // 管理者
	conn := &websocket.Conn{}
	clientSocket := NewClient(clientID, systemID, conn)

	manager.AddClient(clientSocket)

	Convey("测试添加客户端", t, func() {
		Convey("长度是否够", func() {
			So(len(manager.ClientIDMap), ShouldEqual, 1)
		})

		Convey("clientId是否存在", func() {
			_, ok := manager.ClientIDMap[clientID]
			So(ok, ShouldBeTrue)
		})
	})
}

func TestDelClient(t *testing.T) {
	clientID := "clientID"
	systemID := "publishSystem"
	var manager = NewClientManager() // 管理者
	conn := &websocket.Conn{}
	clientSocket := NewClient(clientID, systemID, conn)
	manager.AddClient(clientSocket)

	manager.DelClient(clientSocket)

	Convey("测试删除客户端", t, func() {
		Convey("长度是否够", func() {
			So(len(manager.ClientIDMap), ShouldEqual, 0)
		})

		Convey("clientId是否存在", func() {
			_, ok := manager.ClientIDMap[clientID]
			So(ok, ShouldBeFalse)
		})
	})
}

func TestCount(t *testing.T) {
	clientID := "clientID"
	systemID := "publishSystem"
	var manager = NewClientManager() // 管理者
	conn := &websocket.Conn{}
	clientSocket := NewClient(clientID, systemID, conn)

	Convey("测试获取客户端数量", t, func() {
		Convey("添加一个客户端后", func() {
			manager.AddClient(clientSocket)
			So(manager.Count(), ShouldEqual, 1)
		})

		Convey("删除一个客户端后", func() {
			manager.DelClient(clientSocket)
			So(manager.Count(), ShouldEqual, 0)
		})

		Convey("再添加两个客户端后", func() {
			manager.AddClient(clientSocket)
			manager.AddClient(clientSocket)
			So(manager.Count(), ShouldEqual, 1)
		})
	})
}

func TestGetByClientId(t *testing.T) {
	clientId := "clientId"
	systemId := "publishSystem"
	var manager = NewClientManager() // 管理者
	conn := &websocket.Conn{}
	clientSocket := NewClient(clientId, systemId, conn)

	Convey("测试通过clientId获取客户端", t, func() {
		Convey("获取一个存在的clientId", func() {
			manager.AddClient(clientSocket)
			_, err := manager.GetByClientID(clientId)
			So(err, ShouldBeNil)
		})

		Convey("获取一个不存在的clientId", func() {
			_, err := manager.GetByClientID("notExistId")
			So(err, ShouldNotBeNil)
		})
	})
}

func TestAddClient2LocalGroup(t *testing.T) {
	setting.DefaultSetting()
	clientID := "clientID"
	systemID := "publishSystem"
	userID := "userID"
	var manager = NewClientManager() // 管理者
	conn := &websocket.Conn{}
	clientSocket := NewClient(clientID, systemID, conn)
	manager.AddClient(clientSocket)
	groupName := "testGroup"

	Convey("测试添加分组", t, func() {
		Convey("添加一个客户端到分组", func() {
			manager.AddClient2LocalGroup(groupName, clientSocket, userID, "")
			So(len(manager.Groups[utils.GenGroupKey(systemID, groupName)]), ShouldEqual, 1)
		})

		Convey("再添加一个客户端到分组", func() {
			manager.AddClient2LocalGroup(groupName, clientSocket, userID, "")
			So(len(manager.Groups[utils.GenGroupKey(systemID, groupName)]), ShouldEqual, 1)
		})
	})
}

func TestGetGroupClientList(t *testing.T) {
	clientID := "clientID"
	systemID := "publishSystem"
	userID := "userID"
	var manager = NewClientManager() // 管理者
	conn := &websocket.Conn{}
	clientSocket := NewClient(clientID, systemID, conn)
	manager.AddClient(clientSocket)
	groupName := "testGroup"

	Convey("测试添加分组", t, func() {
		Convey("获取一个存在的分组", func() {
			manager.AddClient2LocalGroup(groupName, clientSocket, userID, "")
			clientIds := manager.GetGroupClientList(utils.GenGroupKey(systemID, groupName))
			So(len(clientIds), ShouldEqual, 1)
		})

		Convey("获取一个不存在的clientId", func() {
			clientIds := manager.GetGroupClientList("notExistId")
			So(len(clientIds), ShouldEqual, 0)
		})
	})
}
