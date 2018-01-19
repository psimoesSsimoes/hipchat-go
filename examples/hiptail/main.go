package main

import (
	"flag"
	"fmt"
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
			fmt.Printf("%s [%s]: %s\n", from, m.Date, msg)

			uriString := strings.Fields(msg)[11]

			out, err := exec.Command("/usr/bin/python", "/home/psimoes/Github/hipchat-go/examples/hiptail/login.py", fmt.Sprintf("%s", uriString)).Output()

			if err != nil {
				fmt.Println("atum")
			} else {
				if fmt.Sprintf("%s", out) == "true" {
					//curl
				}
			}
			lastM = msg
		}
		time.Sleep(time.Second * 15)
	}

	// }
}
