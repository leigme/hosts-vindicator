package cmd

import (
	_ "embed"
	"testing"
)

//go:embed template/hosts.tpl
var hosts string

func Test_writeFileByTemplate(t *testing.T) {
	type args struct {
		filePath string
		tpl      string
		data     map[string]any
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test1",
			args: args{
				filePath: "test",
				tpl:      hosts,
				data: map[string]any{
					"Headers":  []string{"1", "2"},
					"StartTag": "Start",
					"Contents": []string{"3", "4"},
					"EndTag":   "End",
					"Footers":  []string{"5", "6"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writeFileByTemplate(tt.args.filePath, tt.args.tpl, tt.args.data)
		})
	}
}
