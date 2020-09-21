package websocket

import (
	"context"
	"eduhacks2020/Go/pkg/setting"
	"eduhacks2020/Go/protocol/pb"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"sync"
)

func grpcConn(addr string) *grpc.ClientConn {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Errorf("did not connect: %v", err)
	}
	return conn
}

func SendRpc2Client(addr string, messageId, sendUserID, clientId string, code int, message string, data *string) {
	conn := grpcConn(addr)
	defer conn.Close()

	log.WithFields(log.Fields{
		"host":     setting.GlobalSetting.LocalHost,
		"port":     setting.CommonSetting.Port,
		"add":      addr,
		"clientId": clientId,
		"msg":      data,
	}).Info("发送到服务器")

	c := pb.NewCommonServiceClient(conn)
	_, err := c.Send2Client(context.Background(), &pb.Send2ClientReq{
		MessageID:  messageId,
		SendUserID: sendUserID,
		ClientID:   clientId,
		Code:       int32(code),
		Message:    message,
		Data:       *data,
	})
	if err != nil {
		log.Errorf("failed to call: %v", err)
	}
}

func CloseRpcClient(addr string, clientId, systemId string) {
	conn := grpcConn(addr)
	defer conn.Close()

	log.WithFields(log.Fields{
		"host":     setting.GlobalSetting.LocalHost,
		"port":     setting.CommonSetting.Port,
		"add":      addr,
		"clientId": clientId,
	}).Info("发送关闭连接到服务器")

	c := pb.NewCommonServiceClient(conn)
	_, err := c.CloseClient(context.Background(), &pb.CloseClientReq{
		SystemID: systemId,
		ClientID: clientId,
	})
	if err != nil {
		log.Errorf("failed to call: %v", err)
	}
}

//绑定分组
func SendRpcBindGroup(addr string, systemId string, groupName string, clientId string, userId string, extend string) {
	conn := grpcConn(addr)
	defer conn.Close()

	c := pb.NewCommonServiceClient(conn)
	_, err := c.BindGroup(context.Background(), &pb.BindGroupReq{
		SystemID:  systemId,
		GroupName: groupName,
		ClientID:  clientId,
		UserID:    userId,
		Extend:    extend,
	})
	if err != nil {
		log.Errorf("failed to call: %v", err)
	}
}

//发送分组消息
func SendGroupBroadcast(systemId string, messageId, sendUserID, groupName string, code int, message string, data *string) {
	setting.GlobalSetting.ServerListLock.Lock()
	defer setting.GlobalSetting.ServerListLock.Unlock()
	for _, addr := range setting.GlobalSetting.ServerList {
		conn := grpcConn(addr)
		defer conn.Close()

		c := pb.NewCommonServiceClient(conn)
		_, err := c.Send2Group(context.Background(), &pb.Send2GroupReq{
			SystemID:   systemId,
			MessageID:  messageId,
			SendUserID: sendUserID,
			GroupName:  groupName,
			Code:       int32(code),
			Message:    message,
			Data:       *data,
		})
		if err != nil {
			log.Errorf("failed to call: %v", err)
		}
	}
}

//发送系统信息
func SendSystemBroadcast(systemId string, messageId, sendUserID string, code int, message string, data *string) {
	setting.GlobalSetting.ServerListLock.Lock()
	defer setting.GlobalSetting.ServerListLock.Unlock()
	for _, addr := range setting.GlobalSetting.ServerList {
		conn := grpcConn(addr)
		defer conn.Close()

		c := pb.NewCommonServiceClient(conn)
		_, err := c.Send2System(context.Background(), &pb.Send2SystemReq{
			SystemID:   systemId,
			MessageID:  messageId,
			SendUserID: sendUserID,
			Code:       int32(code),
			Message:    message,
			Data:       *data,
		})
		if err != nil {
			log.Errorf("failed to call: %v", err)
		}
	}
}

func GetOnlineListBroadcast(systemId *string, groupName *string) (clientIdList []string) {
	setting.GlobalSetting.ServerListLock.Lock()
	defer setting.GlobalSetting.ServerListLock.Unlock()

	serverCount := len(setting.GlobalSetting.ServerList)

	onlineListChan := make(chan []string, serverCount)
	var wg sync.WaitGroup

	wg.Add(serverCount)
	for _, addr := range setting.GlobalSetting.ServerList {
		go func(addr string) {
			conn := grpcConn(addr)
			defer conn.Close()
			c := pb.NewCommonServiceClient(conn)
			response, err := c.GetGroupClients(context.Background(), &pb.GetGroupClientsReq{
				SystemID:  *systemId,
				GroupName: *groupName,
			})
			if err != nil {
				log.Errorf("failed to call: %v", err)
			} else {
				onlineListChan <- response.List
			}
			wg.Done()

		}(addr)
	}

	wg.Wait()

	for i := 1; i <= serverCount; i++ {
		list, ok := <-onlineListChan
		if ok {
			clientIdList = append(clientIdList, list...)
		} else {
			return
		}
	}
	close(onlineListChan)

	return
}
