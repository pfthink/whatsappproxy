package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	fiberutils "github.com/gofiber/fiber/v2/utils"
	"github.com/pfthink/whatsmeow"
	"github.com/pfthink/whatsmeow/store/sqlstore"
	"github.com/skip2/go-qrcode"
	"os"
	"path/filepath"
	"time"
	"whatsappproxy/config"
	"whatsappproxy/structs"
	"whatsappproxy/utils"
)

type AppServiceImpl struct {
	WaCli          *whatsmeow.Client
	storeContainer *sqlstore.Container
}

func NewAppService(storeContainer *sqlstore.Container) AppService {
	return &AppServiceImpl{
		storeContainer: storeContainer,
	}
}

func (service AppServiceImpl) Login(c *fiber.Ctx) (response structs.LoginResponse, err error) {
	jidUser := c.Query("jid")
	service.WaCli = utils.InitWaCLIByJidUser(jidUser, service.storeContainer)

	if service.WaCli == nil {
		return response, errors.New("wa cli nil cok")
	}

	// Disconnect for reconnecting
	service.WaCli.Disconnect()

	chImage := make(chan string)

	ch, err := service.WaCli.GetQRChannel(context.Background())
	if err != nil {
		// This error means that we're already logged in, so ignore it.
		if errors.Is(err, whatsmeow.ErrQRStoreContainsID) {
			_ = service.WaCli.Connect() // just connect to websocket
			if service.WaCli.IsLoggedIn() {
				return response, errors.New("you already logged in :)")
			}
			return response, errors.New("your session have been saved, please wait to connect 2 second and refresh again")
		} else {
			return response, errors.New("Error when GetQRChannel:" + err.Error())
		}
	} else {
		go func() {
			for evt := range ch {
				response.Code = evt.Code
				response.Duration = evt.Timeout / time.Second / 2
				if evt.Event == "code" {
					fmt.Println("code:", evt.Code)
					qrPath := fmt.Sprintf("%s/scan-qr-%s.png", config.PathQrCode, fiberutils.UUIDv4())
					err = qrcode.WriteFile(evt.Code, qrcode.Medium, 512, qrPath)
					if err != nil {
						fmt.Println("error when write qrImage file", err.Error())
					}
					go func() {
						time.Sleep(response.Duration * time.Second)
						err := os.Remove(qrPath)
						if err != nil {
							fmt.Println("Failed to remove qrPath " + qrPath)
						}
					}()
					chImage <- qrPath
				} else {
					fmt.Printf("QR channel result: %s", evt.Event)
				}
			}
		}()
	}

	err = service.WaCli.Connect()
	if err != nil {
		return response, errors.New("Failed to connect bro " + err.Error())
	}
	response.ImagePath = <-chImage

	return response, nil
}

func (service AppServiceImpl) Logout(c *fiber.Ctx) (err error) {
	jidUser := c.Query("jid")
	service.WaCli = utils.InitWaCLIByJidUser(jidUser, service.storeContainer)
	if service.WaCli == nil {
		return errors.New("wa cli nil cok")
	}

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

	err = service.WaCli.Logout()
	return
}

func (service AppServiceImpl) Reconnect(c *fiber.Ctx) (err error) {
	jidUser := c.Query("jid")
	service.WaCli = utils.InitWaCLIByJidUser(jidUser, service.storeContainer)
	if service.WaCli == nil {
		return errors.New("wa cli nil cok")
	}
	service.WaCli.Disconnect()
	return service.WaCli.Connect()
}
