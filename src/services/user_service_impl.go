package services

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/pfthink/whatsmeow/store/sqlstore"
	"github.com/pfthink/whatsmeow/types"
	"os"
	"path/filepath"
	"whatsappproxy/structs"
	"whatsappproxy/utils"
)

type UserServiceImpl struct {
	storeContainer *sqlstore.Container
}

func NewUserService(storeContainer *sqlstore.Container) UserService {
	return &UserServiceImpl{
		storeContainer: storeContainer,
	}
}

func (service UserServiceImpl) Logout(c *fiber.Ctx, request structs.UserRequest) (err error) {
	cli, _, _ := utils.InitWaClientIfNeed(request.Jid, service.storeContainer)
	utils.MustLogin(cli)

	// delete history
	files, err := filepath.Glob("./history-*")
	if err != nil {
		panic(err)
	}

	for _, f := range files {
		err = os.Remove(f)
		if err != nil {
			return err
		}
	}
	// delete qr images
	qrImages, err := filepath.Glob("./statics/images/qrcode/scan-*")
	if err != nil {
		panic(err)
	}

	for _, f := range qrImages {
		err = os.Remove(f)
		if err != nil {
			return err
		}
	}

	err = cli.Logout()
	return
}

func (service UserServiceImpl) Reconnect(c *fiber.Ctx, request structs.UserRequest) (err error) {
	cli, _, _ := utils.InitWaClientIfNeed(request.Jid, service.storeContainer)
	utils.MustLogin(cli)
	cli.Disconnect()
	return cli.Connect()
}

func (service *UserServiceImpl) UserInfo(_ *fiber.Ctx, request structs.UserRequest) (response structs.UserInfoResponse, err error) {
	cli, jid, _ := utils.InitWaClientIfNeed(request.Jid, service.storeContainer)
	utils.MustLogin(cli)
	phones := []string{"+" + jid.User}
	res, err := cli.IsOnWhatsApp(phones)
	fmt.Println(res)
	var jids []types.JID
	//jid, ok := utils.ParseJID(request.Phone)

	jids = append(jids, jid)
	resp, err := cli.GetUserInfo(jids)
	if err != nil {
		return response, err
	}

	for _, userInfo := range resp {
		var device []structs.UserInfoResponseDataDevice
		for _, j := range userInfo.Devices {
			device = append(device, structs.UserInfoResponseDataDevice{
				User:   j.User,
				Agent:  j.Agent,
				Device: utils.GetPlatformName(int(j.Device)),
				Server: j.Server,
				AD:     j.AD,
			})
		}

		data := structs.UserInfoResponseData{
			Status:    userInfo.Status,
			PictureID: userInfo.PictureID,
			Devices:   device,
		}
		if userInfo.VerifiedName != nil {
			data.VerifiedName = fmt.Sprintf("%v", *userInfo.VerifiedName)
		}
		response.Data = append(response.Data, data)
	}

	return response, nil
}

func (service *UserServiceImpl) UserAvatar(_ *fiber.Ctx, request structs.UserAvatarRequest) (response structs.UserAvatarResponse, err error) {
	cli, jid, _ := utils.InitWaClientIfNeed(request.Jid, service.storeContainer)
	utils.MustLogin(cli)

	pic, err := cli.GetProfilePictureInfo(jid, false)
	if err != nil {
		return response, err
	} else if pic == nil {
		return response, errors.New("no avatar found")
	} else {
		response.URL = pic.URL
		response.ID = pic.ID
		response.Type = pic.Type
		response.Jid = request.Jid
		return response, nil
	}
}

func (service *UserServiceImpl) OnlineStatus(_ *fiber.Ctx, request structs.UserOnlineStatusRequest) (response structs.UserOnlineStatusResponse, err error) {
	cli, _, _ := utils.InitWaClientIfNeed(request.Jid, service.storeContainer)
	utils.MustLogin(cli)
	targetJid, _ := utils.ParseJID(request.TargetJid)
	var targetJids = []string{"+" + targetJid.User}
	resp, err := cli.IsOnWhatsApp(targetJids)

	response.Jid = request.Jid
	response.TargetJid = request.TargetJid
	if resp[0].IsIn {
		response.OnlineStatus = 1
	}
	return response, nil
}
