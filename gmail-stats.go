package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/codegangsta/cli"
)

const ATOM_FEED = "https://mail.google.com/a/gmail.com/feed/atom"

type Message struct {
	Subject string `xml:"title"`
	Summary string `xml:"summary"`
	From    string `xml:"author>email"`
}

type Stats struct {
	Count    int       `xml:"fullcount"`
	Messages []Message `xml:"entry"`
}

func main() {
	app := cli.NewApp()
	app.Name = "gmailstats"
	app.Usage = `Shows your gmail mailbox unread messages`

	app.Flags = []cli.Flag{
		cli.BoolFlag{Name: "count, c", Usage: "Return only count of unread messages."},
		cli.BoolFlag{Name: "daemon, d", Usage: "Run as daemon, uses notify-send event to display email status."},
		cli.IntFlag{Name: "interval, i", Value: 60, Usage: "Daemon check interval, defaults to 60 seconds."},
		cli.StringFlag{Name: "username, u", Value: "", Usage: "Gmail account username."},
		cli.StringFlag{Name: "password, p", Value: "", Usage: "Gmail account password."},
	}

	app.Action = func(c *cli.Context) {
		username := c.String("username")
		password := c.String("password")
		if len(username) == 0 {
			println("An account username must be provided")
			os.Exit(1)
		}

		if len(password) == 0 {
			println("An account password must be provided")
			os.Exit(1)
		}

		// if daemon mode, run forever
		if c.Bool("daemon") {
			lastUnreadCount := 0
			interval := time.Duration(c.Int("interval"))
			for {
				time.Sleep(interval * time.Second)
				stats, err := unread(username, password)
				if err != nil {
					continue // continue if there was an error
				}
				if stats.Count > 0 && lastUnreadCount != stats.Count {
					msg := fmt.Sprintf("%d unread email(s) for %s", stats.Count, username)
					exec.Command("notify-send", "Gmail stats", msg, "--icon=email").Run()
				}
				lastUnreadCount = stats.Count
			}
		}

		// non daemon mode
		stats, err := unread(username, password)
		if err != nil {
			fmt.Printf("Encountered an error: %s\n", err)
			os.Exit(1)
		}

		if c.Bool("count") {
			fmt.Printf("%d\n", stats.Count)
		} else {
			for _, m := range stats.Messages {
				fmt.Printf("%s;%s;%s\n", m.Subject, m.From, m.Summary)
			}
		}
	}

	app.Run(os.Args)
}

// unread - reads statistics from mailbox and returns Stats struct
func unread(usr, psw string) (s Stats, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", ATOM_FEED, nil)
	if err != nil {
		return
	}
	req.SetBasicAuth(usr, psw)
	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return s, fmt.Errorf(res.Status)
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	err = xml.Unmarshal(bytes, &s)
	return
}
