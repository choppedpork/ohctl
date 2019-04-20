package openhab

import (
	"testing"

	"github.com/go-test/deep"
	"github.com/jarcoal/httpmock"
)

func Test_openhabClient_GetItem(t *testing.T) {

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
		want    openhabItem
		wantErr bool
	}{
		{
			"good cat",
			fields{Host: "meow", Port: 12345},
			args{itemName: "cat"},
			openhabItem{
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
		{"not a cat", fields{Host: "meow", Port: 12345}, args{itemName: "dog"}, openhabItem{}, true},
		{"what cat", fields{Host: "meow", Port: 12345}, args{itemName: "huh"}, openhabItem{}, true},
		{"bad host", fields{}, args{itemName: "please"}, openhabItem{}, true},
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
			c := &openhabClient{
				Host: tt.fields.Host,
				Port: tt.fields.Port,
			}
			got, err := c.GetItem(tt.args.itemName)
			if (err != nil) != tt.wantErr {
				t.Errorf("openhabClient.GetItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := deep.Equal(got, tt.want); diff != nil {
				t.Errorf("openhabClient.GetItem() = %v, want %v, difference is %v", got, tt.want, diff)
			}
		})
	}
}

func Test_openhabClient_GetItems(t *testing.T) {

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
		want    []openhabItem
		wantErr bool
	}{
		{
			"good cat",
			fields{Host: "meow", Port: 12345},
			[]openhabItem{{
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
		{"bad host", fields{}, nil, true},
	}

	httpmock.Activate()

	httpmock.RegisterResponder("GET", "http://meow:12345/rest/items?recursive=false",
		httpmock.NewStringResponder(200, `[{"link":"http://meow:12345/rest/items/cat","state":"16","editable":false,"type":"Dimmer","name":"cat","label":"cat","category":"soundvolume","tags":[],"groupNames":[]}]`))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &openhabClient{
				Host: tt.fields.Host,
				Port: tt.fields.Port,
			}
			got, err := c.GetItems()
			if (err != nil) != tt.wantErr {
				t.Errorf("openhabClient.GetItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := deep.Equal(got, tt.want); diff != nil {
				t.Errorf("openhabClient.GetItem() = %v, want %v, difference is %v", got, tt.want, diff)
			}
		})
	}
}

func Test_openhabClient_Cmd(t *testing.T) {
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
		// TODO: Add test cases.
	}

	httpmock.RegisterResponder("GET", "http://meow:12345/rest/items/cat",
		httpmock.NewStringResponder(200, `{"link":"http://meow:12345/rest/items/cat","state":"16","editable":false,"type":"Dimmer","name":"cat","label":"cat","category":"soundvolume","tags":[],"groupNames":[]}`))
	httpmock.RegisterResponder("GET", "http://meow:12345/rest/items/dog",
		httpmock.NewStringResponder(404, ""))
	httpmock.RegisterResponder("GET", "http://meow:12345/rest/items/huh",
		httpmock.NewStringResponder(200, "what is this i can't even"))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &openhabClient{
				Host: tt.fields.Host,
				Port: tt.fields.Port,
			}
			if err := c.Cmd(tt.args.item, tt.args.cmd); (err != nil) != tt.wantErr {
				t.Errorf("openhabClient.Cmd() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
