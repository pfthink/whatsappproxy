package services

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/pfthink/whatsmeow"
	"github.com/pfthink/whatsmeow/store/sqlstore"
	"github.com/pfthink/whatsmeow/types"
	"whatsappproxy/structs"
	"whatsappproxy/utils"
)

type UserServiceImpl struct {
	WaCli          *whatsmeow.Client
	storeContainer *sqlstore.Container
}

func NewUserService(storeContainer *sqlstore.Container) UserService {
	return &UserServiceImpl{
		storeContainer: storeContainer,
	}
}

func (service *UserServiceImpl) UserInfo(_ *fiber.Ctx, request structs.UserInfoRequest) (response structs.UserInfoResponse, err error) {
	cliMap := utils.CliMap
	jid, ok := utils.ParseJID(request.Phone)
	cli, exists := cliMap[jid.User]
	if !exists {
		client := utils.InitWaCLIByJidUser(jid.User, service.storeContainer)
		cliMap[jid.User] = client
		cli = client
	}

	utils.MustLogin(cli)

	phones := []string{"+" + jid.User}
	res, err := cli.IsOnWhatsApp(phones)
	fmt.Println(res)
	var jids []types.JID
	//jid, ok := utils.ParseJID(request.Phone)
	if !ok {
		return response, errors.New("invalid JID " + request.Phone)
	}

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
	jid, ok := utils.ParseJID(request.Phone)
	cliMap := utils.CliMap
	cli, exists := cliMap[jid.User]
	if !exists {
		client := utils.InitWaCLIByJidUser(jid.User, service.storeContainer)
		cliMap[jid.User] = client
		cli = client
	}
	utils.MustLogin(cli)

	/*if service.WaCli == nil {
		return response, errors.New("wa cli nil cok")
	}*/

	//jid, ok := utils.ParseJID(request.Phone)
	if !ok {
		return response, errors.New("invalid JID " + request.Phone)
	}
	pic, err := cli.GetProfilePictureInfo(jid, false)
	if err != nil {
		return response, err
	} else if pic == nil {
		return response, errors.New("no avatar found")
	} else {
		response.URL = pic.URL
		response.ID = pic.ID
		response.Type = pic.Type

		return response, nil
	}
}

func (service UserServiceImpl) UserMyListGroups(_ *fiber.Ctx) (response structs.UserMyListGroupsResponse, err error) {
	utils.MustLogin(service.WaCli)

	groups, err := service.WaCli.GetJoinedGroups()
	if err != nil {
		return
	}
	fmt.Printf("%+v\n", groups)
	if groups != nil {
		for _, group := range groups {
			response.Data = append(response.Data, *group)
		}
	}
	return response, nil
}

func (service UserServiceImpl) UserMyPrivacySetting(_ *fiber.Ctx) (response structs.UserMyPrivacySettingResponse, err error) {
	utils.MustLogin(service.WaCli)

	resp, err := service.WaCli.TryFetchPrivacySettings(false)
	if err != nil {
		return
	}

	response.GroupAdd = string(resp.GroupAdd)
	response.Status = string(resp.Status)
	response.ReadReceipts = string(resp.ReadReceipts)
	response.Profile = string(resp.Profile)
	return response, nil
}
