package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"io"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
)

// initialize global variables
var apiID int
var apiHash string
var sessionName string
var target int64

type Secrets struct {
	ApiID       int    `json:"api_id"`
	ApiHash     string `json:"api_hash"`
	SessionName string `json:"session_name"`
	Target      int64  `json:"target"`
	PhoneNumber string `json:"phone_number"`
}

var secrets Secrets

func initClient() {
	// Read from secrets.json in project root
	data, err := os.ReadFile("secrets.json")
	if err != nil {
		log.Fatalf("Error reading secrets.json: %v", err)
		return
	}

	err = json.Unmarshal(data, &secrets)
	if err != nil {
		log.Fatalf("Error parsing secrets.json: %v", err)
		return
	}

	apiID = secrets.ApiID
	apiHash = secrets.ApiHash
	sessionName = secrets.SessionName
	target = secrets.Target
}

func main() {
	initClient()

	// Create session storage using FileStorage
	sessionStorage := &session.FileStorage{
		Path: filepath.Join(".", sessionName+".json"),
	}

	client := telegram.NewClient(apiID, apiHash, telegram.Options{
		SessionStorage: sessionStorage,
	})
	ctx := context.Background()

	err := client.Run(ctx, func(ctx context.Context) error {
		// Check if we're already authenticated
		authStatus, err := client.Auth().Status(ctx)
		if err != nil {
			return fmt.Errorf("failed to get auth status: %v", err)
		}

		// If not authenticated, perform authentication
		if !authStatus.Authorized {
			fmt.Println("Not authenticated. Starting authentication...")

			if err := client.Auth().IfNecessary(ctx, TerminalAuth(secrets.PhoneNumber)); err != nil {
				return fmt.Errorf("auth failed: %v", err)
			}
			fmt.Println("Authentication successful!")
		} else {
			fmt.Println("Already authenticated!")
		}

		// Now that we're authenticated, we read the dialogs.
		readDialogs(client, ctx)

		return nil
	})

	if err != nil {
		log.Fatalf("Error running client: %v", err)
	}
}

func readLastLine(filename string) (string, error) {
	// Open fd
	fd, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer fd.Close()

	stat, err := fd.Stat()
	if err != nil {
		return "", err
	}
	size := stat.Size()

	if size == 0 {
		return "", nil // empty file!
	}

	buf := make([]byte, 1)
	var line []byte

	for offset := int64(1); offset <= size; offset++ {
		// Seek until the end of the file
		_, err = fd.Seek(-offset, io.SeekEnd)
		if err != nil {
			return "", err
		}

		// Read into the buffer created earlier
		_, err = fd.Read(buf)
		if err != nil {
			return "", err
		}

		lastByte := []byte{buf[0]}

		if string(lastByte) == "\n" && offset != 1 {
			// we want to stop only if we found a newline AND we're not at the first iteration.
			// otherwise we wouldn't have actually read the last line, only the last char!
			break
		}
		// build the line backwards, since I'm reading backwards...
		line = append(lastByte, line...)
	}

	return string(line), nil

}

const layout = "2006-01-02T15:04:05"

func shouldRaiseAlert() (string, error) {
	// See if the alarm log already showed an alert today
	lastLine, err := readLastLine("logs/alarms.log")
	if err != nil {
		log.Fatalf("reading from logs failed %v", err)
		return "not-send", nil
	}

	lastLine = strings.TrimSpace(lastLine) // removes \n
	lastAlarmDate, err := time.Parse(layout, lastLine)
	if err != nil {
		log.Fatalf("Last line of alarms.log file did not have the expected format. Last line: %v", lastLine)
		return "not-send", nil
	}

	now := time.Now()

	now = now.Truncate(24 * time.Hour)
	lastAlarmDate = lastAlarmDate.Truncate(24 * time.Hour)

	if lastAlarmDate.Before(now) {
		return "send", nil
	} else {
		return "not-send", nil
	}

}

func sendAlert() {
	shouldSend, err := shouldRaiseAlert()
	if err != nil {
		log.Fatalf("Could not read the last line from the logs/alarms.log directory")
	}
	cmd := exec.Command("bash", "src/send_alert.sh", shouldSend)
	cmd.Stdout = os.Stdout

	err = cmd.Run()
	if err != nil {
		log.Fatalf("sending alert failed to run, %v", err)
	}
}

func checkIfDialogIsTargetAndUnread(dialog *tg.Dialog) {
	peer := dialog.Peer
	var id int64
	switch p := peer.(type) {
	case *tg.PeerUser:
		id = p.UserID
	case *tg.PeerChat:
		id = p.ChatID
	case *tg.PeerChannel:
		id = p.ChannelID
	default:
		id = 0
		fmt.Println("Unknown peer type")
	}

	if id == target && dialog.UnreadCount > 0 {
		sendAlert()
	}

}

func readDialogs(client *telegram.Client, ctx context.Context) (bool, error) {
	api := client.API()
	dialogsResp, err := api.MessagesGetDialogs(ctx, &tg.MessagesGetDialogsRequest{
		Limit:      10,
		OffsetPeer: &tg.InputPeerEmpty{},
	})
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	// Type switch to get the concrete type
	switch d := dialogsResp.(type) {
	case *tg.MessagesDialogs:
		// This has both Dialogs and Chats/Users
		fmt.Printf("Got %d dialogs\n", len(d.Dialogs))
		for _, dialogClass := range d.Dialogs {
			// dialogClass is a DialogClass interface, need to type-switch
			if dialog, ok := dialogClass.(*tg.Dialog); ok {
				checkIfDialogIsTargetAndUnread(dialog)
			}
		}
	case *tg.MessagesDialogsSlice:
		// This is used when there are more dialogs (pagination)
		fmt.Printf("Got %d dialogs (slice, total: %d)\n", len(d.Dialogs), d.Count)
		for _, dialogClass := range d.Dialogs {
			// dialogClass is a DialogClass interface, need to type-switch
			if dialog, ok := dialogClass.(*tg.Dialog); ok {
				checkIfDialogIsTargetAndUnread(dialog)
			}
		}
	default:
		fmt.Printf("Unexpected dialog type: %T\n", d)
	}

	return false, nil
}
