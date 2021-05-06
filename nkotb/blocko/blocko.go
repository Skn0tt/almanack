package blocko

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"strings"

	"github.com/carlmjohnson/flagext"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const AppName = "NKOTB"

func CLI(args []string) error {
	var app appEnv
	err := app.ParseArgs(args)
	if err != nil {
		return err
	}
	if err = app.Exec(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
	return err
}

func (app *appEnv) ParseArgs(args []string) error {
	fl := flag.NewFlagSet(AppName, flag.ContinueOnError)
	fl.Usage = func() {
		version := "(unknown)"
		if i, ok := debug.ReadBuildInfo(); ok {
			version = i.Main.Version
		}

		fmt.Fprintf(fl.Output(), `NKOTB %s - extract blocks of Markdownish content from HTML
		
Usage:

	nkotb [options] <src>

If not set, src is stdin.

Options:
`, version)
		fl.PrintDefaults()
		fmt.Fprintln(fl.Output())
	}
	src := flagext.FileOrURL(flagext.StdIO, nil)
	app.src = src
	if err := fl.Parse(args); err != nil {
		return err
	}
	if err := flagext.ParseEnv(fl, AppName); err != nil {
		return err
	}
	if fl.NArg() > 0 {
		if err := src.Set(fl.Arg(0)); err != nil {
			return err
		}
	}
	return nil
}

type appEnv struct {
	src io.ReadCloser
}

func (app *appEnv) Exec() (err error) {
	defer app.src.Close()
	buf := bufio.NewReader(app.src)
	doc, err := html.Parse(buf)
	if err != nil {
		return err
	}
	bNode := findNode(doc, func(n *html.Node) *html.Node {
		if n.DataAtom == atom.Body {
			return n
		}
		return nil
	})
	if bNode == nil {
		return fmt.Errorf("could not find body")
	}
	visitAll(bNode, func(n *html.Node) {
		if n.Type != html.TextNode {
			return
		}
		n.Data = strings.ReplaceAll(n.Data, "\n", " ")
		n.Data = strings.ReplaceAll(n.Data, "\r", " ")
	})

	return outputBlocks(bNode, 0)
}

func outputBlocks(bNode *html.Node, depth int) (err error) {
	var (
		wbuf    strings.Builder
		needsNL = false
	)
loop:
	for p := bNode.FirstChild; p != nil; p = p.NextSibling {
		if needsNL {
			fmt.Print("\n\n")
		}
		if depth > 0 {
			if p == bNode.FirstChild {
				fmt.Print("- ")
			} else {
				fmt.Print(strings.Repeat("    ", depth))
			}
		}
		if !blockElements[p.DataAtom] {
			if err = html.Render(&wbuf, p); err != nil {
				return err
			}
			needsNL = output(&wbuf)
			continue
		}
		if isEmpty(p) {
			fmt.Print("")
			needsNL = false
			continue
		}
		switch p.DataAtom {
		case atom.H1:
			wbuf.WriteString("# ")
		case atom.H2:
			wbuf.WriteString("## ")
		case atom.H3:
			wbuf.WriteString("### ")
		case atom.H4:
			wbuf.WriteString("#### ")
		case atom.H5:
			wbuf.WriteString("##### ")
		case atom.H6:
			wbuf.WriteString("###### ")
		case atom.Ul, atom.Ol:
			for c := p.FirstChild; c != nil; c = c.NextSibling {
				if err = outputBlocks(c, depth+1); err != nil {
					return err
				}
				fmt.Print("\n")
			}
			continue loop
		}
		for c := p.FirstChild; c != nil; c = c.NextSibling {
			if err = html.Render(&wbuf, c); err != nil {
				return err
			}
		}
		needsNL = output(&wbuf)
	}
	return nil
}

func output(wbuf *strings.Builder) bool {
	s := wbuf.String()
	s = strings.TrimSpace(s)
	wbuf.Reset()
	fmt.Print(s)
	return s != ""
}

func findNode(n *html.Node, callback func(*html.Node) *html.Node) *html.Node {
	if r := callback(n); r != nil {
		return r
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if r := findNode(c, callback); r != nil {
			return r
		}
	}
	return nil
}

func visitAll(n *html.Node, callback func(*html.Node)) {
	findNode(n, func(n *html.Node) *html.Node {
		callback(n)
		return nil
	})
}

var blockElements = map[atom.Atom]bool{
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

var stylisticElements = map[atom.Atom]bool{
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

func isEmpty(n *html.Node) bool {
	root := n
	n = findNode(n, func(n *html.Node) *html.Node {
		if n == root {
			return nil
		}
		switch n.Type {
		case html.TextNode:
			s := strings.ReplaceAll(n.Data, "\n", " ")
			s = strings.TrimSpace(s)
			if s == "" {
				return nil
			}
		case html.ElementNode:
			if stylisticElements[n.DataAtom] {
				return nil
			}
		}
		return n
	})
	return n == nil
}
