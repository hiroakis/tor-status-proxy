package tor

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Status struct {
	exitNodeList NodeList
	allNodeList  NodeList
}

type NodeList struct {
	rawData      []byte
	lastModified time.Time
}

var urls = map[string]string{
	"all":  "https://torstatus.blutmagie.de/ip_list_all.php/Tor_ip_list_ALL.csv",
	"exit": "https://torstatus.blutmagie.de/ip_list_exit.php/Tor_ip_list_EXIT.csv",
}

func get(u string) ([]byte, error) {
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b := &bytes.Buffer{}
	io.Copy(b, resp.Body)
	return b.Bytes(), nil

}

func updateStatus(status *Status) {
	bCh := make(chan []byte)
	errCh := make(chan error)

	for k, v := range urls {
		go func(w string) {
			b, err := get(w)
			if err != nil {
				errCh <- err
			} else {
				bCh <- b
			}
		}(v)

		select {
		case b := <-bCh:
			switch k {
			case "all":
				status.allNodeList.rawData = b
				status.allNodeList.lastModified = time.Now()
			case "exit":
				status.exitNodeList.rawData = b
				status.exitNodeList.lastModified = time.Now()
			}
		case err := <-errCh:
			fmt.Println(err)
		}
	}
}

func NewStatus(interval int) *Status {
	var status Status
	updateStatus(&status)

	pollingFunc := func(sec int) {
		ticker := time.NewTicker(time.Duration(sec) * time.Second)
		for range ticker.C {
			updateStatus(&status)
		}
	}
	go pollingFunc(interval)
	return &status
}

func (self *Status) RawExitNodes() []byte {
	return self.exitNodeList.rawData
}

func (self *Status) ExitNodes() []string {
	return strings.Split(string(self.exitNodeList.rawData), "\n")
}

func (self *Status) ExitNodeLastModified() time.Time {
	return self.exitNodeList.lastModified
}

func (self *Status) IsExitNode(ip string) bool {
	for _, v := range self.ExitNodes() {
		if v == ip {
			return true
		}
	}
	return false
}

func (self *Status) RawAllNodes() []byte {
	return self.allNodeList.rawData
}

func (self *Status) AllNodes() []string {
	return strings.Split(string(self.allNodeList.rawData), "\n")
}

func (self *Status) AllNodeLastModified() time.Time {
	return self.allNodeList.lastModified
}

func (self *Status) IsTorNode(ip string) bool {
	for _, v := range self.AllNodes() {
		if v == ip {
			return true
		}
	}
	return false
}
