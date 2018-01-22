package main

import (
	"flag"
	"fmt"
	"net/url"
	"os/exec"
	"strings"
	"time"

	"github.com/tbruyelle/hipchat-go/hipchat"
)

const (
	maxMsgLen  = 128
	moreString = " [MORE]"
)

var (
	token  = flag.String("token", "nTx2TZMf5HMxB2Oi2WPmEdsWM2cDY2hLxc8Uolf4", "The HipChat AuthToken")
	roomId = flag.String("room", "962414", "The HipChat room id")
)

func main() {
	flag.Parse()
	if *token == "" || *roomId == "" {
		flag.PrintDefaults()
		return
	}
	c := hipchat.NewClient(*token)
	hist, resp, err := c.Room.History(*roomId, &hipchat.HistoryOptions{})
	if err != nil {
		fmt.Printf("Error during room history req %q\n", err)
		fmt.Printf("Server returns %+v\n", resp)
		return
	}
	lastM := ""
	for {
		m := hist.Items[len(hist.Items)-1]
		from := ""
		switch m.From.(type) {
		case string:
			from = m.From.(string)
		case map[string]interface{}:
			f := m.From.(map[string]interface{})
			from = f["name"].(string)
		}
		msg := m.Message
		msg = fmt.Sprintf("%s%s", strings.Replace(m.Message[:len(m.Message)], "\n", " - ", -1), moreString)
		if lastM != msg {
			if len(strings.Fields(msg)) > 11 {
				uriString := strings.Fields(msg)[11]
				if isValidUrl(uriString) {
					fmt.Println("isValid")
					exec.Command("/usr/bin/python", "/home/psimoes/Github/hipchat-go/examples/hiptail/login.py", fmt.Sprintf("%s", uriString)).Output()
				}
				fmt.Println("Ping")

			}
		}
		lastM = msg
		time.Sleep(time.Second * 5)
	}

	// }
}
func isValidUrl(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	} else {
		return true
	}
}
