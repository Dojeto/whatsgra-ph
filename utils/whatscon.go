package utils

import (
	ctx "context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal"
	"go.mau.fi/whatsmeow"
	wp "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

var (
	UrlToLink = make(map[string]string)
)

const (
	ContentTypeBinary = "application/octet-stream"
	ContentTypeForm   = "application/x-www-form-urlencoded"
	ContentTypeJSON   = "application/json"
	ContentTypeHTML   = "text/html; charset=utf-8"
	ContentTypeText   = "text/plain; charset=utf-8"
)

var client *whatsmeow.Client

func ConnectToWP() {
	dbLog := waLog.Stdout("Database", "INFO", true)
	container, err := sqlstore.New("sqlite3", "file:whatgpt.db?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
	}
	clientLog := waLog.Stdout("Client", "INFO", true)
	client = whatsmeow.NewClient(deviceStore, clientLog)

	client.AddEventHandler(func(evt interface{}) {
		if evt, ok := evt.(*events.Message); ok {
			if evt.Message.DocumentMessage != nil {
				pdfData, _ := client.Download(evt.Message.GetDocumentMessage())

				link, err := Upload(pdfData, evt.Message.GetDocumentMessage().FileName)
				if err != nil {
					fmt.Println("Error opening file:", err)
					return
				}
				id := uuid.New()
				UrlToLink[id.String()[0:8]] = *link
				client.SendMessage(ctx.Background(), evt.Info.Sender, &wp.Message{
					Conversation: proto.String(id.String()[0:8]),
				})
			}
		}
	})
	if client.Store.ID == nil {
		qrChan, _ := client.GetQRChannel(ctx.Background())
		// Connect to WhatsApp
		err = client.Connect()
		if err != nil {
			panic(err)
		}
		// Print the QR code to the console
		for evt := range qrChan {
			if evt.Event == "code" {
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		err = client.Connect()
		if err != nil {
			panic(err)
		}
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Disconnect()
}
