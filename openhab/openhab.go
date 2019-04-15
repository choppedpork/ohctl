package openhab

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type openhabItem struct {
	// editable: false
	// groupNames: []
	// label: Mode
	// link: http://openhab/rest/items/nest_mode
	// name: nest_mode
	// state: HEAT
	// stateDescription:
	//   options:
	//   - label: "off"
	// 	   value: "OFF"
	//   - label: eco
	//     value: ECO
	//   - label: heating
	//     value: HEAT
	//   - label: cooling
	// 	   value: COOL
	//   - label: heat/cool
	// 	   value: HEAT_COOL
	//   pattern: '%s'
	//   readOnly: false
	// tags: []
	// type: String
	Editable         bool        `json:"editable"`
	GroupNames       []string    `json:"groupNames"`
	Label            string      `json:"label"`
	Link             string      `json:"link"`
	Name             string      `json:"name"`
	State            string      `json:"state"`
	StateDescription interface{} `json:"stateDescription"`
	Tags             []string    `json:"tags"`
	Type             string      `json:"type"`
}

type openhabClient struct {
}

func NewClient() *openhabClient {
	return &openhabClient{}
}

func (c *openhabClient) GetItems() []openhabItem {

	resp, _ := http.Get("http://openhab/rest/items?recursive=false") // lol error handling

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body) // more classy error handling

	var items []openhabItem
	_ = json.Unmarshal(body, &items)

	return items

}

func (c *openhabClient) GetItem(itemName string) openhabItem {

	resp, _ := http.Get("http://openhab/rest/items/" + itemName) // lol error handling

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body) // more classy error handling

	var item openhabItem
	_ = json.Unmarshal(body, &item)

	return item

}

func (c *openhabClient) Cmd(item string, cmd string) {

	resp, _ := http.Post("http://openhab/rest/items/"+item, "text/plain", strings.NewReader(cmd)) // more classy error handling

	switch resp.StatusCode {
	case 200:
		fmt.Println("yay")
	case 400:
		fmt.Printf("error: invalid command: %s\n", cmd)
	case 404:
		fmt.Printf("error: item %s doesn't exist\n", item)
	}
}
