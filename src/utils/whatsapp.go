package utils

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pfthink/agollo"
	"github.com/pfthink/whatsmeow"
	"github.com/pfthink/whatsmeow/appstate"
	waProto "github.com/pfthink/whatsmeow/binary/proto"
	"github.com/pfthink/whatsmeow/store"
	"github.com/pfthink/whatsmeow/store/sqlstore"
	"github.com/pfthink/whatsmeow/types"
	"github.com/pfthink/whatsmeow/types/events"
	waLog "github.com/pfthink/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
	"mime"
	"os"
	"strings"
	"sync/atomic"
	"time"
	"whatsappproxy/config"
	"whatsappproxy/rabbitmq"
)

var (
	cli           *whatsmeow.Client
	log           waLog.Logger
	historySyncID int32
	startupTime   = time.Now().Unix()
	CliMap        = make(map[string]*whatsmeow.Client)
	ApolloClient  agollo.Client
	Bucket        *oss.Bucket
)

func GetPlatformName(deviceID int) string {
	switch deviceID {
	case 2:
		return "UNKNOWN"
	case 3:
		return "CHROME"
	case 4:
		return "FIREFOX"
	case 5:
		return "IE"
	case 6:
		return "OPERA"
	case 7:
		return "SAFARI"
	case 8:
		return "EDGE"
	case 9:
		return "DESKTOP"
	case 10:
		return "IPAD"
	case 11:
		return "ANDROID_TABLET"
	case 12:
		return "OHANA"
	default:
		return "UNKNOWN"
	}
}

func ParseJID(arg string) (types.JID, bool) {
	if arg[0] == '+' {
		arg = arg[1:]
	}
	if !strings.ContainsRune(arg, '@') {
		return types.NewJID(arg, types.DefaultUserServer), true
	} else {
		recipient, err := types.ParseJID(arg)
		if err != nil {
			_ = fmt.Errorf("Invalid JID %s: %v", arg, err)
			return recipient, false
		} else if recipient.User == "" {
			_ = fmt.Errorf("Invalid JID %s: no server specified", arg)
			return recipient, false
		}
		return recipient, true
	}
}

func InitWaDB() *sqlstore.Container {
	// Running Whatsapp
	log = waLog.Stdout("Main", config.WhatsappLogLevel, true)
	dbLog := waLog.Stdout("Database", config.WhatsappLogLevel, true)
	//storeContainer, err := sqlstore.New("sqlite3", fmt.Sprintf("file:%s?_foreign_keys=off", config.DBName), dbLog)
	storeContainer, err := sqlstore.New("mysql", "root:123456@tcp(127.0.0.1:3306)/whatsapp", dbLog)

	//var dbDialect = flag.String("db-dialect", "mysql", "Database dialect (sqlite3 or postgres)")

	//var dbAddress = flag.String("db-address", "file:mdtest.db?_foreign_keys=on", "Database address")
	//var dbAddress = flag.String("db-address", "root:123456@tcp(127.0.0.1:3306)/whatsapp", "Database address")

	if err != nil {
		log.Errorf("Failed to connect to database: %v", err)
		panic(err)

	}
	return storeContainer
}

func InitWaCLI(storeContainer *sqlstore.Container) *whatsmeow.Client {
	device, err := storeContainer.GetFirstDevice()
	if err != nil {
		log.Errorf("Failed to get device: %v", err)
		panic(err)
	}

	store.DeviceProps.PlatformType = waProto.DeviceProps_CHROME.Enum()
	store.DeviceProps.Os = proto.String("AldinoKemal")
	cli = whatsmeow.NewClient(device, waLog.Stdout("Client", config.WhatsappLogLevel, true))
	cli.AddEventHandler(handler)

	return cli
}

func InitWaCLIByJidUser(jidUser string, storeContainer *sqlstore.Container) *whatsmeow.Client {
	device, err := storeContainer.GetDeviceByJidUser(jidUser)
	if err != nil {
		log.Errorf("Failed to get device: %v", err)
		panic(err)
	}

	store.DeviceProps.PlatformType = waProto.DeviceProps_CHROME.Enum()
	//store.DeviceProps.Os = proto.String("AldinoKemal")
	cli = whatsmeow.NewClient(device, waLog.Stdout("Client", config.WhatsappLogLevel, true))
	cli.AddEventHandler(handler)

	return cli
}

func NewWaCLI(storeContainer *sqlstore.Container) *whatsmeow.Client {
	device, err := storeContainer.GenerateDevice()
	if err != nil {
		log.Errorf("Failed to get device: %v", err)
		panic(err)
	}
	store.DeviceProps.PlatformType = waProto.DeviceProps_CHROME.Enum()
	//store.DeviceProps.Os = proto.String("AldinoKemal")
	cli = whatsmeow.NewClient(device, waLog.Stdout("Client", config.WhatsappLogLevel, true))
	cli.AddEventHandler(handler)

	return cli
}

func MustLogin(waCli *whatsmeow.Client) {
	if waCli == nil {
		panic(AuthError{Message: "wa cli nil cok"})
	} else if !waCli.IsConnected() {
		panic(AuthError{Message: "you are not connect to whatsapp server, please reconnect"})
	}
}

func handler(rawEvt interface{}) {
	switch evt := rawEvt.(type) {
	case *events.AppStateSyncComplete:
		if len(cli.Store.PushName) > 0 && evt.Name == appstate.WAPatchCriticalBlock {
			err := cli.SendPresence(types.PresenceAvailable)
			if err != nil {
				log.Warnf("Failed to send available presence: %v", err)
			} else {
				log.Infof("Marked self as available")
			}
		}
	case *events.Connected, *events.PushNameSetting:
		if len(cli.Store.PushName) == 0 {
			return
		}
		// Send presence available when connecting and when the pushname is changed.
		// This makes sure that outgoing messages always have the right pushname.
		err := cli.SendPresence(types.PresenceAvailable)
		if err != nil {
			log.Warnf("Failed to send available presence: %v", err)
		} else {
			log.Infof("Marked self as available")
		}
	case *events.StreamReplaced:
		os.Exit(0)
	case *events.Message:
		metaParts := []string{fmt.Sprintf("pushname: %s", evt.Info.PushName), fmt.Sprintf("timestamp: %s", evt.Info.Timestamp)}
		if evt.Info.Type != "" {
			metaParts = append(metaParts, fmt.Sprintf("type: %s", evt.Info.Type))
		}
		if evt.Info.Category != "" {
			metaParts = append(metaParts, fmt.Sprintf("category: %s", evt.Info.Category))
		}
		if evt.IsViewOnce {
			metaParts = append(metaParts, "view once")
		}
		if evt.IsViewOnce {
			metaParts = append(metaParts, "ephemeral")
		}

		//log.Infof("Received message %s from %s (%s): %+v", evt.Info.ID, evt.Info.SourceString(), strings.Join(metaParts, ", "), evt.Message)
		fmt.Printf("Received message %s from %s (%s): %+v", evt.Info.ID, evt.Info.SourceString(), strings.Join(metaParts, ", "), evt.Message)

		m := make(map[string]interface{})
		/*
			m["messageId"] = evt.Info.ID
			m["fromJid"] = fmt.Sprintf("%s@%s", evt.Info.MessageSource.Sender.User, evt.Info.MessageSource.Sender.Server)
			m["toJid"] = fmt.Sprintf("%s@%s", cli.Store.ID.User, cli.Store.ID.Server)
			m["conversation"] = evt.Message.Conversation
			m["chat"] = fmt.Sprintf("%s@%s", evt.Info.MessageSource.Chat.User, evt.Info.MessageSource.Chat.Server)*/
		msg, err := json.Marshal(evt)
		json.Unmarshal(msg, &m)
		m["jid"] = fmt.Sprintf("%s@%s", cli.Store.ID.User, cli.Store.ID.Server)
		msgByte, err := json.Marshal(&m)
		if err != nil {
			log.Errorf("Failed to parsing message %v", err)
			return
		}
		rabbitmq.SendBossImMsg(msgByte)
		img := evt.Message.GetImageMessage()
		if img != nil {
			data, err := cli.Download(img)
			if err != nil {
				log.Errorf("Failed to download image: %v", err)
				return
			}
			exts, _ := mime.ExtensionsByType(img.GetMimetype())
			path := fmt.Sprintf("%s%s", evt.Info.ID, exts[0])
			err = os.WriteFile(path, data, 0600)
			if err != nil {
				log.Errorf("Failed to save image: %v", err)
				return
			}
			log.Infof("Saved image in message to %s", path)
		}

		if config.WhatsappAutoReplyMessage != "" {
			_, _ = cli.SendMessage(evt.Info.Sender, "", &waProto.Message{Conversation: proto.String(config.WhatsappAutoReplyMessage)})
		}
	case *events.Receipt:
		if evt.Type == events.ReceiptTypeRead || evt.Type == events.ReceiptTypeReadSelf {
			log.Infof("%v was read by %s at %s", evt.MessageIDs, evt.SourceString(), evt.Timestamp)
		} else if evt.Type == events.ReceiptTypeDelivered {
			log.Infof("%s was delivered to %s at %s", evt.MessageIDs[0], evt.SourceString(), evt.Timestamp)
		}
	case *events.Presence:
		if evt.Unavailable {
			if evt.LastSeen.IsZero() {
				log.Infof("%s is now offline", evt.From)
			} else {
				log.Infof("%s is now offline (last seen: %s)", evt.From, evt.LastSeen)
			}
		} else {
			log.Infof("%s is now online", evt.From)
		}
	case *events.HistorySync:
		id := atomic.AddInt32(&historySyncID, 1)
		jid := fmt.Sprintf("%s@%s", cli.Store.ID.User, cli.Store.ID.Server)
		fileName := fmt.Sprintf("history-%s-%d.json", jid, id)
		file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			log.Errorf("Failed to open file to write history sync: %v", err)
			return
		}
		enc := json.NewEncoder(file)
		enc.SetIndent("", "  ")
		err = enc.Encode(evt.Data)
		m := make(map[string]interface{})
		msg, err := json.Marshal(evt.Data)
		json.Unmarshal(msg, &m)
		m["jid"] = fmt.Sprintf("%s@%s", cli.Store.ID.User, cli.Store.ID.Server)
		msgByte, err := json.Marshal(&m)
		rabbitmq.SendBossImMsg(msgByte)
		if err != nil {
			log.Errorf("Failed to write history sync: %v", err)
			return
		}
		log.Infof("Wrote history sync to %s", fileName)
		_ = file.Close()
	case *events.AppState:
		log.Debugf("App state event: %+v / %+v", evt.Index, evt.SyncActionValue)
	case *events.PairSuccess:
		//noise := base64.StdEncoding.EncodeToString(cli.Store.NoiseKey.Pub[:])
		//fmt.Printf("PairSuccess......,noise:%s", evt.ID, noise)
	}

}
