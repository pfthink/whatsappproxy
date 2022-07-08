package utils

import (
	"errors"
	"github.com/pfthink/whatsmeow"
	"github.com/pfthink/whatsmeow/store/sqlstore"
	"github.com/pfthink/whatsmeow/types"
)

func InitWaClientIfNeed(jid string, storeContainer *sqlstore.Container) (cli *whatsmeow.Client, typeJid types.JID, err error) {
	typeJid, ok := ParseJID(jid)
	if !ok {
		return nil, typeJid, errors.New("invalid JID: " + jid)
	}
	cliMap := CliMap
	cli, exists := cliMap[typeJid.User]
	if !exists {
		client := InitWaCLIByJidUser(typeJid.User, storeContainer)
		cliMap[typeJid.User] = client
		cli = client
	}
	if !cli.IsConnected() {
		cli.Connect()
	}

	return cli, typeJid, nil
}
