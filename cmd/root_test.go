package cmd_test

import (
	"bytes"

	"example.com/learn/btrees/cmd"
	"strings"
	"testing"
)

func TestCommands(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    []string
		wantErr bool
	}{
		{name: "root help", want: []string{"Available Commands:", "btree"}},
		{name: "explicit help", args: []string{"--help"}, want: []string{"btree"}},
		{name: "demo help", args: []string{"btree", "--help"}, want: []string{"btrees btree"}},
		{name: "demo", args: []string{"btree"}, want: []string{
			"Insert 6\n[10]\n  [5 6]\n  [20]",
			"Search for 17\nvisit [10 20]: follow child 1",
			"visit [12 17]: found",
			"Search for 15", "leaf reached; not found",
			"Search for 10\nvisit [10 20]: found",
		}},
		{name: "bplus help", args: []string{"bplus", "--help"}, want: []string{"btrees bplus"}},
		{name: "bplus demo", args: []string{"bplus"}, want: []string{
			"branch [10 20]", "leaf [10 12 17]",
			"Leaf chain: [5 6 7] -> [10 12 17] -> [20 30] -> nil",
			"Range [6, 20] (inclusive)", "follow next leaf", "Result: [6 7 10 12 17 20]",
		}},
		{name: "bplus unexpected argument", args: []string{"bplus", "extra"}, wantErr: true},
		{name: "bplus unknown flag", args: []string{"bplus", "--missing"}, wantErr: true},
		{name: "unknown command", args: []string{"missing"}, wantErr: true},
		{name: "unexpected argument", args: []string{"btree", "extra"}, wantErr: true},
		{name: "unknown flag", args: []string{"btree", "--missing"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := cmd.New()
			var output bytes.Buffer
			root.SetOut(&output)
			root.SetErr(&output)
			root.SetArgs(tt.args)
			err := root.Execute()
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v; wantErr = %v", err, tt.wantErr)
			}
			for _, want := range tt.want {
				if !strings.Contains(output.String(), want) {
					t.Errorf("output missing %q: %s", want, output.String())
				}
			}
			if tt.name != "demo" && tt.name != "bplus demo" && strings.Contains(output.String(), "Insert 10") {
				t.Fatal("demonstration ran unexpectedly")
			}
		})
	}
}

func TestDuckDBLearningCommands(t *testing.T) {
	examples := []struct {
		name string
		want []string
	}{
		{"art", []string{"Duplicate car inserted: false", "Node4", "Node16", "Node48", "Node256", "prefix="}},
		{"zonemap", []string{"ordered values:", "mixed values:", "skipped=2 scanned=1 values examined=4", "False positive:"}},
		{"columnscan", []string{"selection=[2 3]", "IDs: [3 4 5 6 7]; SUM(amount)=250", "Batches=3; filter values read=10; amount values read=5"}},
	}
	for _, example := range examples {
		t.Run(example.name, func(t *testing.T) {
			root := cmd.New()
			var output bytes.Buffer
			root.SetOut(&output)
			root.SetErr(&output)
			root.SetArgs([]string{example.name})
			if err := root.Execute(); err != nil {
				t.Fatal(err)
			}
			for _, want := range example.want {
				if !strings.Contains(output.String(), want) {
					t.Errorf("missing %q in %s", want, output.String())
				}
			}
			for _, args := range [][]string{{example.name, "extra"}, {example.name, "--missing"}} {
				command := cmd.New()
				command.SetOut(&output)
				command.SetErr(&output)
				command.SetArgs(args)
				if err := command.Execute(); err == nil {
					t.Errorf("accepted invalid arguments: %v", args)
				}
			}
			help := cmd.New()
			output.Reset()
			help.SetOut(&output)
			help.SetErr(&output)
			help.SetArgs([]string{example.name, "--help"})
			if err := help.Execute(); err != nil || !strings.Contains(output.String(), "btrees "+example.name) {
				t.Fatalf("incorrect command help: %v %s", err, output.String())
			}
		})
	}
}
