package note

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	extast "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/text"
)

// Daily notes must be named YYYY-MM-DD.md.
const dailyNoteDateFormat = "2006-01-02"

// Note is a parsed Markdown note.
type Note struct {
	Date           time.Time
	Filepath       string
	Notename       string
	Sections       []Header
	Raw            string
	Tasks          []string
	CompletedTasks []string
}

// HeaderType is a Markdown heading level.
type HeaderType int

const (
	H1 HeaderType = 1
	H2 HeaderType = 2
	H3 HeaderType = 3
	H4 HeaderType = 4
	H5 HeaderType = 5
	H6 HeaderType = 6
)

// Header is a heading found in a note.
type Header struct {
	Name string
	Type HeaderType
}

var markdown = goldmark.New(
	goldmark.WithExtensions(extension.GFM),
)

// ParseNote reads a Markdown note and extracts its date, headings, and tasks.
// path is the note's filepath. Notename is the base name without the .md extension.
func ParseNote(r io.Reader, path string) (Note, error) {
	var n Note

	data, err := io.ReadAll(r)
	if err != nil {
		return n, fmt.Errorf("read markdown: %w", err)
	}

	n.Raw = string(data)
	n.Filepath = path
	n.Notename = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))

	n.Date, _ = time.Parse(dailyNoteDateFormat, n.Notename)

	// Parse Markdown into AST.
	doc := markdown.Parser().Parse(textReader(data))

	// Walk through the AST.
	err = ast.Walk(doc, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		switch v := node.(type) {
		case *ast.Heading:
			n.Sections = append(n.Sections, Header{
				Name: inlineText(v, data),
				Type: HeaderType(v.Level),
			})

		case *ast.ListItem:
			found, checked := parseTask(v)

			if !found {
				break
			}

			task := taskText(v, data)

			if checked {
				n.CompletedTasks = append(n.CompletedTasks, task)
			} else {
				n.Tasks = append(n.Tasks, task)
			}
		}

		return ast.WalkContinue, nil
	})

	if err != nil {
		return n, fmt.Errorf("parse markdown: %w", err)
	}

	return n, nil
}

// parseTask checks whether a list item is a task and whether it is completed.
func parseTask(item ast.Node) (found bool, checked bool) {
	var walk func(ast.Node)

	walk = func(n ast.Node) {
		if found {
			return
		}

		if checkbox, ok := n.(*extast.TaskCheckBox); ok {
			found = true
			checked = checkbox.IsChecked
			return
		}

		for child := n.FirstChild(); child != nil; child = child.NextSibling() {
			walk(child)
		}
	}

	walk(item)

	return
}

func taskText(item ast.Node, source []byte) string {
	var parts []string

	for child := item.FirstChild(); child != nil; child = child.NextSibling() {
		switch child.(type) {
		case *ast.TextBlock, *ast.Paragraph:
		default:
			continue
		}

		text := inlineText(child, source)
		if text != "" {
			parts = append(parts, text)
		}
	}

	return strings.Join(parts, " ")
}

func inlineText(node ast.Node, source []byte) string {
	var b strings.Builder

	var walk func(ast.Node)

	walk = func(n ast.Node) {
		switch v := n.(type) {
		case *ast.Text:
			b.Write(v.Segment.Value(source))

			if v.SoftLineBreak() || v.HardLineBreak() {
				b.WriteByte(' ')
			}

		case *ast.String:
			b.Write(v.Value)

		case *ast.CodeSpan:
			// Code spans hold their text in child nodes. Writing both the
			// span and those children would duplicate the code.
			b.Write(v.Text(source))
			return

		case *extast.TaskCheckBox:
			return
		}

		for child := n.FirstChild(); child != nil; child = child.NextSibling() {
			walk(child)
		}
	}

	walk(node)

	return strings.TrimSpace(b.String())
}

// textReader creates a Goldmark text reader.
func textReader(source []byte) text.Reader {
	return text.NewReader(source)
}
