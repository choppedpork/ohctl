package openhab

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// Item is a representation of an Openhab item
type Item struct {
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

// Client represents a client for a Openhab server
type Client struct {
	Host string
	Port uint16
}

// NewClient creates a new Openhab client instance
func NewClient(host string, port uint16) *Client {
	return &Client{Host: host, Port: port}
}

// GetItems retrieves a list of all items
func (c *Client) GetItems() ([]Item, error) {

	resp, err := http.Get("http://" + c.Host + ":" + strconv.Itoa(int(c.Port)) + "/rest/items?recursive=false")

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var items []Item
	err = json.Unmarshal(body, &items)

	if err != nil {
		return nil, err
	}

	return items, nil

}

// GetItem retrieves an item by name
func (c *Client) GetItem(itemName string) (Item, error) {

	resp, err := http.Get("http://" + c.Host + ":" + strconv.Itoa(int(c.Port)) + "/rest/items/" + itemName)

	if err != nil {
		return Item{}, errors.New("fetching item: " + err.Error())
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return Item{}, err
	}

	if resp.StatusCode == 404 {
		return Item{}, errors.New("item " + itemName + " doesn't exist")
	}

	var item Item
	err = json.Unmarshal(body, &item)

	if err != nil {
		return Item{}, errors.New("coulnd't unmarshal: " + err.Error() + " (not openhab?)")
	}

	return item, nil

}

// Cmd executes a command on an item
func (c *Client) Cmd(item string, cmd string) error {

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
