package mediastreamvalidator

import (
	"encoding/json"
	"fmt"
	"github.com/jkittell/array"
	"testing"
	"time"
)

var host = "127.0.0.1"

func TestA(t *testing.T) {
	url := "https://cph-p2p-msl.akamaized.net/hls/live/2000341/test/master.m3u8"

	info := Post(host, url)

	for range time.Tick(10 * time.Second) {
		i := Get(host, info.Id)
		if i.Status == "completed" {
			break
		}
	}

	info = Get(host, info.Id)
	data, err := json.MarshalIndent(info, "", "    ")
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Println(string(data))
}

func TestB(t *testing.T) {
	ids := array.New[string]()
	url := "https://cph-p2p-msl.akamaized.net/hls/live/2000341/test/master.m3u8"

	time.Sleep(2 * time.Second)

	for i := 0; i < 2; i++ {
		info := Post(host, url)
		ids.Push(info.Id)
	}

	var done bool
	for range time.Tick(10 * time.Second) {
		for i := 0; i < 2; i++ {
			id := ids.Lookup(i)
			info := Get(host, id)
			fmt.Println(info.Status)
			if info.Status == "completed" {
				done = true
			}
		}
		if done {
			break
		}
	}

	all := GetAll(host)
	data, err := json.MarshalIndent(all, "", "    ")
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Println(string(data))
}
