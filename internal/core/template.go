package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"text/template/parse"
)

// IdentifyPlaceholders parses a template file and returns a slice of unique placeholder keys.
// It inspects the AST of the template to find all field nodes, e.g., {{.Name}}, {{.Version}}.
func IdentifyPlaceholders(templatePath string) ([]string, error) {
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("could not read template file '%s': %w", templatePath, err)
	}

	// Create a new template and parse the content.
	tmpl, err := template.New(filepath.Base(templatePath)).Parse(string(content))
	if err != nil {
		return nil, fmt.Errorf("could not parse template '%s': %w", templatePath, err)
	}

	// Use a map to store unique placeholder names.
	placeholders := make(map[string]struct{})

	// Walk the parse tree of the template to find all fields.
	if tmpl.Tree != nil && tmpl.Tree.Root != nil {
		walk(tmpl.Tree.Root, placeholders)
	}

	// Convert the map keys to a slice.
	keys := make([]string, 0, len(placeholders))
	for k := range placeholders {
		keys = append(keys, k)
	}

	return keys, nil
}

// walk is a recursive function to traverse the template's abstract syntax tree (AST).
//
//nolint:gocognit // acceptance
func walk(node parse.Node, placeholders map[string]struct{}) {
	if node.Type() == parse.NodeAction {
		// An ActionNode is a template action, like {{.Field}}.
		// We need to look inside its pipeline.
		action := node.(*parse.ActionNode) //nolint:errcheck // it is predictable type
		if action.Pipe != nil {
			for _, cmd := range action.Pipe.Cmds {
				for _, arg := range cmd.Args {
					if fieldNode, ok := arg.(*parse.FieldNode); ok {
						// A FieldNode represents a field access, e.g., .Name
						// The Ident slice holds the parts of the field.
						// We join them with dots for nested fields.
						fieldName := strings.Join(fieldNode.Ident, ".")
						placeholders[fieldName] = struct{}{}
					}
				}
			}
		}
	}

	// Recursively walk through list nodes.
	if list, ok := node.(*parse.ListNode); ok {
		for _, n := range list.Nodes {
			walk(n, placeholders)
		}
	}
	// Add more checks for other node types if needed (e.g., if, range)
	// For range nodes
	if rangeNode, ok := node.(*parse.RangeNode); ok {
		walk(rangeNode.List, placeholders)
		if rangeNode.ElseList != nil {
			walk(rangeNode.ElseList, placeholders)
		}
	}
	// For if nodes
	if ifNode, ok := node.(*parse.IfNode); ok {
		walk(ifNode.List, placeholders)
		if ifNode.ElseList != nil {
			walk(ifNode.ElseList, placeholders)
		}
	}
}
