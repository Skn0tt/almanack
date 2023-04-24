package xhtml

import (
	"bytes"
	"strings"
	"sync"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func ToBuffer(n *html.Node) *bytes.Buffer {
	var buf bytes.Buffer
	if err := html.Render(&buf, n); err != nil {
		panic(err)
	}
	return &buf
}

var pool sync.Pool

func poolGet() *bytes.Buffer {
	if v := pool.Get(); v != nil {
		buf := v.(*bytes.Buffer)
		buf.Reset()
		return buf
	}
	return &bytes.Buffer{}
}

func ToString(n *html.Node) string {
	buf := poolGet()
	defer pool.Put(buf)
	if err := html.Render(buf, n); err != nil {
		panic(err)
	}
	return buf.String()
}

func ContentsToString(n *html.Node) string {
	buf := poolGet()
	defer pool.Put(buf)

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if err := html.Render(buf, c); err != nil {
			panic(err)
		}
	}
	return buf.String()
}

// InnerBlocksToString is the same as ContentsToString,
// but it separates top level nodes with a line break.
func InnerBlocksToString(n *html.Node) string {
	buf := poolGet()
	defer pool.Put(buf)

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if err := html.Render(buf, c); err != nil {
			panic(err)
		}
		buf.WriteString("\n")
	}
	return buf.String()
}

func InnerText(n *html.Node) string {
	buf := poolGet()
	defer pool.Put(buf)

	VisitAll(n, func(n *html.Node) {
		if n.Type == html.TextNode {
			buf.WriteString(n.Data)
		}
	})
	return buf.String()
}

var MarkdownBlockElements = map[atom.Atom]bool{
	atom.P:  true,
	atom.H1: true,
	atom.H2: true,
	atom.H3: true,
	atom.H4: true,
	atom.H5: true,
	atom.H6: true,
	atom.Ul: true,
	atom.Ol: true,
}

var InlineElements = map[atom.Atom]bool{
	atom.A:       true,
	atom.Abbr:    true,
	atom.Acronym: true,
	atom.B:       true,
	atom.Bdi:     true,
	atom.Bdo:     true,
	atom.Big:     true,
	atom.Br:      true,
	atom.Cite:    true,
	atom.Code:    true,
	atom.Del:     true,
	atom.Dfn:     true,
	atom.Em:      true,
	atom.I:       true,
	atom.Ins:     true,
	atom.Kbd:     true,
	atom.Label:   true,
	atom.Mark:    true,
	atom.Meter:   true,
	atom.Output:  true,
	atom.Q:       true,
	atom.Ruby:    true,
	atom.S:       true,
	atom.Samp:    true,
	atom.Small:   true,
	atom.Span:    true,
	atom.Strong:  true,
	atom.Sub:     true,
	atom.Sup:     true,
	atom.U:       true,
	atom.Tt:      true,
	atom.Var:     true,
	atom.Wbr:     true,
}

func IsEmpty(n *html.Node) bool {
	root := n
	n = Find(n, func(n *html.Node) bool {
		if n == root {
			return false
		}
		switch n.Type {
		case html.TextNode:
			s := strings.ReplaceAll(n.Data, "\n", " ")
			s = strings.TrimSpace(s)
			if s == "" {
				return false
			}
		case html.ElementNode:
			if InlineElements[n.DataAtom] {
				return false
			}
		}
		return true
	})
	return n == nil
}

func ReplaceWith(old, new *html.Node) {
	old.Parent.InsertBefore(new, old)
	old.Parent.RemoveChild(old)
}

func AdoptChildren(dst, src *html.Node) {
	if dst == src {
		return
	}
	for src.FirstChild != nil {
		c := src.FirstChild
		src.RemoveChild(c)
		dst.AppendChild(c)
	}
}

func Attr(n *html.Node, name string) string {
	for _, attr := range n.Attr {
		if attr.Key == name {
			return attr.Val
		}
	}
	return ""
}
