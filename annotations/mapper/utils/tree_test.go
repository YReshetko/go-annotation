package utils_test

import (
	"bytes"
	"fmt"
	"github.com/YReshetko/go-annotation/annotations/mapper/utils"
	"strings"
	"testing"
)

func TestTree(t *testing.T) {
	tests := []struct {
		keys  []string
		value string
	}{
		{value: "outside1"},
		{
			keys: []string{
				"a", "a.b", "a.b.d", "a.b.d.k", "a.b.d.k.x",
			},
			value: "first",
		}, {
			keys: []string{
				"a", "a.b", "a.b.d", "a.b.d.p",
			},
			value: "second",
		}, {
			keys: []string{
				"a", "a.e", "a.e.l", "a.e.l.y",
			},
			value: "third",
		}, {
			keys: []string{
				"a", "a.f", "a.f.n",
			},
			value: "forth",
		},
		{value: "outside2"},
		{
			keys: []string{
				"a", "a.f", "a.f.n",
			},
			value: "fifth",
		}, {
			keys: []string{
				"a", "a.f", "a.f.m", "a.f.m.z",
			},
			value: "sixth",
		},
		{
			keys: []string{
				"A", "a.b", "a.b.d", "a.b.d.k", "a.b.d.k.x",
			},
			value: "first",
		}, {
			keys: []string{
				"A", "a.b", "a.b.d", "a.b.d.p",
			},
			value: "second",
		}, {
			keys: []string{
				"a", "a.e", "a.e.l", "a.e.l.y",
			},
			value: "third",
		},
		{value: "outside3"},
		{
			keys: []string{
				"A", "a.f", "a.f.n",
			},
			value: "forth",
		}, {
			keys: []string{
				"A", "a.f", "a.f.n",
			},
			value: "fifth",
		}, {
			keys: []string{
				"A", "a.f", "a.f.m", "a.f.m.z",
			},
			value: "sixth",
		},
	}

	n := utils.NewNode[string, string]()
	for _, test := range tests {
		n.Add(test.keys, test.value)
	}
	n.Optimize()

	buffer := bytes.NewBufferString("")
	tab := &tabs{}
	n.Execute(preProcessing(buffer, tab), postProcessing(buffer, tab))

	fmt.Println(buffer.String())

}

type tabs struct {
	t []string
}

func preProcessing(buffer *bytes.Buffer, t *tabs) func([]string, []string) {
	return func(k []string, value []string) {
		var keys []string
		for _, s := range k {
			keys = append(keys, s+" != nil")
		}
		ifLine := strings.Join(keys, " && ")
		if len(k) != 0 {
			buffer.Write([]byte(strings.Join(t.t, "") + "if " + ifLine + " {\n"))
			t.t = append(t.t, "\t")
		}
		if len(value) != 0 {
			buffer.Write([]byte(strings.Join(t.t, "") + strings.Join(value, "\n"+strings.Join(t.t, "")) + "\n"))
		}
	}
}

func postProcessing(buffer *bytes.Buffer, t *tabs) func([]string, []string) {
	return func(k []string, value []string) {
		if len(t.t) > 0 {
			t.t = t.t[1:]
		}
		if len(k) > 0 {
			buffer.Write([]byte(strings.Join(t.t, "") + "}\n"))
		}
	}
}
