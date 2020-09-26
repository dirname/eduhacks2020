package render

import (
	"fmt"
	"testing"
)

func TestGetLayer(t *testing.T) {
	type args struct {
		t       int
		icon    int
		title   string
		content string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "layer", args: args{t: 0, icon: Smile, title: "test", content: "test"}, want: fmt.Sprintf(layerTemp, 0, "test", "test", Smile)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetLayer(tt.args.t, tt.args.icon, tt.args.title, tt.args.content); got != tt.want {
				t.Errorf("GetLayer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMsg(t *testing.T) {
	type args struct {
		content string
		Sec     int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "layerMsg", args: args{
			content: "test",
			Sec:     0,
		}, want: fmt.Sprintf(layerMsg, "test", 0)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetMsg(tt.args.content, tt.args.Sec); got != tt.want {
				t.Errorf("GetMsg() = %v, want %v", got, tt.want)
			}
		})
	}
}
