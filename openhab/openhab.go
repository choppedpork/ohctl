package openhab

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
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
	// Category         string      `json:"category"`

}

type openhabClient struct {
	Host string
	Port uint16
}

func NewClient(host string, port uint16) *openhabClient {
	return &openhabClient{Host: host, Port: port}
}

func (c *openhabClient) GetItems() ([]openhabItem, error) {

	resp, err := http.Get("http://" + c.Host + ":" + strconv.Itoa(int(c.Port)) + "/rest/items?recursive=false") // lol error handling

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var items []openhabItem
	err = json.Unmarshal(body, &items)

	if err != nil {
		return nil, err
	}

	return items, nil

}

func (c *openhabClient) GetItem(itemName string) (openhabItem, error) {

	resp, err := http.Get("http://" + c.Host + ":" + strconv.Itoa(int(c.Port)) + "/rest/items/" + itemName)

	if err != nil {
		return openhabItem{}, errors.New("fetching item: " + err.Error())
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return openhabItem{}, err
	}

	if resp.StatusCode == 404 {
		return openhabItem{}, errors.New("item " + itemName + " doesn't exist")
	}

	var item openhabItem
	err = json.Unmarshal(body, &item)

	if err != nil {
		return openhabItem{}, errors.New("coulnd't unmarshal: " + err.Error() + " (not openhab?)")
	}

	return item, nil

}

func (c *openhabClient) Cmd(item string, cmd string) error {

	resp, err := http.Post("http://"+c.Host+":"+strconv.Itoa(int(c.Port))+"/rest/items/"+item, "text/plain", strings.NewReader(cmd)) // more classy error handling

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case 400:
		return errors.New("error: invalid command: " + cmd)
	case 404:
		return errors.New("error: item " + item + " doesn't exist")
	}

	return nil
}
