package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/fatih/color"
)

func main() {
	node := yaml.Node{}
	decoder := yaml.NewDecoder(os.Stdin)

	for {
		err := decoder.Decode(&node)

		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}

		var recurse func(*yaml.Node, int)
		recurse = func(node *yaml.Node, depth int) {
			fmt.Printf("%*s%s\n", depth*2, "", printNode(node))
			if node.Kind == yaml.SequenceNode || node.Kind == yaml.MappingNode || node.Kind == yaml.DocumentNode {
				if len(node.Content) > 0 {
					for i := range node.Content {
						recurse(node.Content[i], depth+1)
					}
				}
			}
		}

		recurse(&node, 0)
	}

	fmt.Println("---")
	encoder := yaml.NewEncoder(os.Stdout)
	encoder.SetIndent(2)
	err := encoder.Encode(&node)

	if err != nil {
		panic(err)
	}
}

func printNode(node *yaml.Node) string {
	buf := bytes.NewBuffer(nil)
	fmt.Fprintf(buf, "%s:%s ", color.YellowString(fmt.Sprint(node.Line)), color.YellowString(fmt.Sprint(node.Column)))
	fmt.Fprintf(buf, "kind=%s ", color.YellowString(kind(node.Kind)))
	fmt.Fprintf(buf, "HeadComment=\"%s\" ", color.RedString(strings.Join(strings.Split(node.HeadComment, "\n"), "\\n")))
	fmt.Fprintf(buf, "LineComment=\"%s\" ", color.CyanString(strings.Join(strings.Split(node.LineComment, "\n"), "\\n")))
	fmt.Fprintf(buf, "FootComment=\"%s\" ", color.MagentaString(strings.Join(strings.Split(node.FootComment, "\n"), "\\n")))
	fmt.Fprintf(buf, "Tag=\"%s\" ", color.HiGreenString(strings.Join(strings.Split(node.Tag, "\n"), "\\n")))
	fmt.Fprintf(buf, "Value=\"%s\"", color.HiYellowString(strings.Join(strings.Split(node.Value, "\n"), "\\n")))

	return buf.String()
}

func kind(k yaml.Kind) string {
	switch k {
	case yaml.DocumentNode:
		return "Document"
	case yaml.AliasNode:
		return "Alias"
	case yaml.MappingNode:
		return "Mapping"
	case yaml.ScalarNode:
		return "Scalar"
	case yaml.SequenceNode:
		return "Sequence"
	default:
		return "Unknown"
	}
}
