package openhab_test

import (
	"bytes"
	"net/http"
	"reflect"
	"testing"

	"github.com/choppedpork/ohctl/openhab"

	"github.com/go-test/deep"
	"github.com/jarcoal/httpmock"
)

func Test_Client_GetItem(t *testing.T) {

	type fields struct {
		Host string
		Port uint16
	}

	type args struct {
		itemName string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    openhab.Item
		wantErr bool
	}{
		{
			"success",
			fields{"meow", 12345},
			args{itemName: "cat"},
			openhab.Item{
				Link:       "http://meow:12345/rest/items/cat",
				State:      "16",
				Type:       "Dimmer",
				Name:       "cat",
				Label:      "cat",
				Tags:       []string{},
				GroupNames: []string{},
			},
			false,
		},
		{"item doesn't exist", fields{"meow", 12345}, args{itemName: "dog"}, openhab.Item{}, true},
		{"non-json response", fields{"meow", 12345}, args{itemName: "huh"}, openhab.Item{}, true},
		{"bad host", fields{}, args{itemName: "please"}, openhab.Item{}, true},
	}

	httpmock.Activate()

	httpmock.RegisterResponder("GET", "http://meow:12345/rest/items/cat",
		httpmock.NewStringResponder(200, `{"link":"http://meow:12345/rest/items/cat","state":"16","editable":false,"type":"Dimmer","name":"cat","label":"cat","category":"soundvolume","tags":[],"groupNames":[]}`))
	httpmock.RegisterResponder("GET", "http://meow:12345/rest/items/dog",
		httpmock.NewStringResponder(404, ""))
	httpmock.RegisterResponder("GET", "http://meow:12345/rest/items/huh",
		httpmock.NewStringResponder(200, "what is this i can't even"))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &openhab.Client{
				Host: tt.fields.Host,
				Port: tt.fields.Port,
			}
			got, err := c.GetItem(tt.args.itemName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := deep.Equal(got, tt.want); diff != nil {
				t.Errorf("Client.GetItem() = %v, want %v, difference is %v", got, tt.want, diff)
			}
		})
	}
}

func Test_Client_GetItems(t *testing.T) {

	type fields struct {
		Host string
		Port uint16
	}

	type args struct {
		itemName string
	}

	tests := []struct {
		name    string
		fields  fields
		want    []openhab.Item
		wantErr bool
	}{
		{
			"success",
			fields{"meow", 12345},
			[]openhab.Item{{
				Link:       "http://meow:12345/rest/items/cat",
				State:      "16",
				Type:       "Dimmer",
				Name:       "cat",
				Label:      "cat",
				Tags:       []string{},
				GroupNames: []string{},
			}},
			false,
		},
		{"bad response", fields{"meow", 123}, nil, true},
		{"bad host", fields{}, nil, true},
	}

	httpmock.Activate()

	httpmock.RegisterResponder("GET", "http://meow:12345/rest/items?recursive=false",
		httpmock.NewStringResponder(200, `[{"link":"http://meow:12345/rest/items/cat","state":"16","editable":false,"type":"Dimmer","name":"cat","label":"cat","category":"soundvolume","tags":[],"groupNames":[]}]`))
	httpmock.RegisterResponder("GET", "http://meow:123/rest/items?recursive=false",
		httpmock.NewStringResponder(200, "hi i am random text"))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &openhab.Client{
				Host: tt.fields.Host,
				Port: tt.fields.Port,
			}
			got, err := c.GetItems()
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := deep.Equal(got, tt.want); diff != nil {
				t.Errorf("Client.GetItem() = %v, want %v, difference is %v", got, tt.want, diff)
			}
		})
	}
}

func Test_Client_Cmd(t *testing.T) {
	type fields struct {
		Host string
		Port uint16
	}
	type args struct {
		item string
		cmd  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"success", fields{"meow", 12345}, args{"cat", "purr"}, false},
		{"invalid command", fields{"meow", 12345}, args{"cat", "bark"}, true},
		{"item doesn't exist", fields{"meow", 12345}, args{"dog", "purr"}, true},
		{"bad host", fields{}, args{"cat", "purr"}, true},
	}

	httpmock.RegisterResponder("POST", "http://meow:12345/rest/items/cat",
		func(req *http.Request) (*http.Response, error) {
			cmd := new(bytes.Buffer)
			cmd.ReadFrom(req.Body)

			if cmd.String() == "purr" {
				return httpmock.NewStringResponse(200, ""), nil
			}

			return httpmock.NewStringResponse(400, ""), nil
		})

	httpmock.RegisterResponder("POST", "http://meow:12345/rest/items/dog",
		httpmock.NewStringResponder(404, ""))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &openhab.Client{
				Host: tt.fields.Host,
				Port: tt.fields.Port,
			}
			if err := c.Cmd(tt.args.item, tt.args.cmd); (err != nil) != tt.wantErr {
				t.Errorf("Client.Cmd() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	type args struct {
		host string
		port uint16
	}
	tests := []struct {
		name string
		args args
		want *openhab.Client
	}{
		{"yay", args{"meow", 12345}, &openhab.Client{"meow", 12345}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := openhab.NewClient(tt.args.host, tt.args.port); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClient() = %v, want %v", got, tt.want)
			}
		})
	}
}
