package extract

import (
	"bytes"
	"code.google.com/p/go.net/html"
	"code.google.com/p/go.net/html/atom"
	"strings"
)

type TextExtractor struct {
	bytes.Buffer
	ignoring    atom.Atom
	inline      bool
	currentElem atom.Atom
	elemStack   []atom.Atom
}

func NewTextExtractor() *TextExtractor {
	return &TextExtractor{elemStack: make([]atom.Atom, 0)}
}

var (
	ignoreElements = map[atom.Atom]struct{}{
		atom.Audio: {}, atom.Canvas: {}, atom.Command: {}, atom.Embed: {},
		atom.Iframe: {}, atom.Map: {}, atom.Math: {}, atom.Menu: {},
		atom.Noscript: {}, atom.Object: {}, atom.Script: {}, atom.Style: {},
		atom.Svg: {}, atom.Video: {},
	}
	inlineElements = map[atom.Atom]struct{}{
		atom.A: {}, atom.Abbr: {}, atom.B: {}, atom.Bdo: {}, atom.Big: {},
		atom.Br: {}, atom.Button: {}, atom.Cite: {}, atom.Code: {},
		atom.Dfn: {}, atom.Em: {}, atom.I: {}, atom.Img: {}, atom.Input: {},
		atom.Kbd: {}, atom.Label: {}, atom.Map: {}, atom.Object: {}, atom.Q: {},
		atom.Samp: {}, atom.Script: {}, atom.Select: {}, atom.Small: {},
		atom.Span: {}, atom.Strong: {}, atom.Sub: {}, atom.Sup: {},
		atom.Textarea: {}, atom.Tt: {}, atom.Var: {},
	}
)

func (t *TextExtractor) HandleToken(token html.Token) {
	switch token.Type {
	case html.StartTagToken:
		t.push(token.DataAtom)

		switch token.DataAtom {
		// Upon hitting a opening body tag, clear the content seen so
		// far. This works around pages that don't have a body tag for
		// some reason by just grabbing all the textual content.
		case atom.Body:
			t.Reset()
		default:
			if _, ignore := ignoreElements[token.DataAtom]; ignore && t.ignoring == 0 {
				t.ignoring = token.DataAtom
			}
		}
	case html.EndTagToken:
		t.maybePop(token.DataAtom)

		switch token.DataAtom {
		case atom.Head:
			t.Reset()
		default:
			if token.DataAtom != 0 && token.DataAtom == t.ignoring {
				t.ignoring = 0
			}
		}
	case html.TextToken:
		if t.ignoring != 0 {
			return
		}
		words := strings.Fields(token.Data)
		if len(words) > 0 {
			if !t.inline && t.Len() > 0 {
				t.WriteString(" ")
			}
			t.WriteString(words[0])
			for _, word := range words[1:] {
				t.WriteString(" ")
				t.WriteString(word)
			}
		}
	}
}

func (t *TextExtractor) push(elem atom.Atom) {
	t.currentElem = elem
	t.elemStack = append(t.elemStack, elem)
	_, t.inline = inlineElements[elem]
}

// Pop the top element from the stack, but only if the closing tag
// matches. Spurious closing tags will simply be ignored when they
// don't match.
func (t *TextExtractor) maybePop(closing atom.Atom) {
	i := len(t.elemStack) - 1
	elem := t.elemStack[i]

	if elem != closing {
		return
	}
	t.elemStack = t.elemStack[:i]
	_, t.inline = inlineElements[elem]
}
