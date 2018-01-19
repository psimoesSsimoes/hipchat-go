package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/tbruyelle/hipchat-go/hipchat"
	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
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
		m := hist.Items[len(hist.Items)-2]
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

		if lastM != msg && strings.ContainsAny(msg, "Orlando") {
			msg = fmt.Sprintf("%s%s", strings.Replace(m.Message[:len(m.Message)], "\n", " - ", -1), moreString)
			fmt.Printf("%s [%s]: %s\n", from, m.Date, msg)
			lastM = msg
			resp, err := http.Get("https://news.ycombinator.com/")
			if err != nil {
				panic(err)
			}
			root, err := html.Parse(resp.Body)
			if err != nil {
				panic(err)
			}

			// define a matcher
			matcher := func(n *html.Node) bool {
				// must check for nil values
				if n.DataAtom == atom.A && n.Parent != nil && n.Parent.Parent != nil {
					return scrape.Attr(n.Parent.Parent, "class") == "athing"
				}
				return false
			}
			// grab all articles and print them
			articles := scrape.FindAll(root, matcher)
			for i, article := range articles {
				fmt.Printf("%2d %s (%s)\n", i, scrape.Text(article), scrape.Attr(article, "href"))

			}
		}

		time.Sleep(time.Second * 5)
	}

	// }
}
