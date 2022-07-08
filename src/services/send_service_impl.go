package services

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/nacos-group/nacos-sdk-go/common/logger"
	"github.com/pfthink/whatsmeow"
	waProto "github.com/pfthink/whatsmeow/binary/proto"
	"github.com/pfthink/whatsmeow/store/sqlstore"
	"google.golang.org/protobuf/proto"
	"whatsappproxy/structs"
	"whatsappproxy/utils"
)

type SendServiceImpl struct {
	storeContainer *sqlstore.Container
}

func NewSendService(storeContainer *sqlstore.Container) SendService {
	return &SendServiceImpl{
		storeContainer: storeContainer,
	}
}

func (service SendServiceImpl) SendText(_ *fiber.Ctx, request structs.SendMessageRequest) (response structs.SendMessageResponse, err error) {
	cli, _, _ := utils.InitWaClientIfNeed(request.FromJid, service.storeContainer)
	utils.MustLogin(cli)
	recipient, ok := utils.ParseJID(request.ToJid)
	if !ok {
		return response, errors.New("invalid JID " + request.ToJid)
	}
	msg := &waProto.Message{Conversation: proto.String(request.Message)}
	msgId := whatsmeow.GenerateMessageID()
	ts, err := cli.SendMessage(recipient, msgId, msg)
	if err != nil {
		return response, err
	} else {
		response.MessageId = msgId
		logger.Infof("Message sent to %s from %s ,server timestamp: %s,msgId:%s", request.ToJid, request.FromJid, ts, msgId)
	}
	return response, nil
}
