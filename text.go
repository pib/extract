package extract

import (
	"bytes"
	"code.google.com/p/go.net/html"
	"code.google.com/p/go.net/html/atom"
	"strings"
	"unicode"
)

type TextExtractor struct {
	bytes.Buffer
	ignoring    atom.Atom
	ignoreDepth int
	inline      bool
	currentElem atom.Atom
	elemStack   []atom.Atom
}

func NewTextExtractor() *TextExtractor {
	return &TextExtractor{elemStack: make([]atom.Atom, 0)}
}

var (
	ignoreIds = map[string]struct{}{
		"disqus_thread": {},
	}
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
	voidElements = map[atom.Atom]struct{}{
		atom.Area: {}, atom.Base: {}, atom.Br: {}, atom.Col: {},
		atom.Command: {}, atom.Embed: {}, atom.Hr: {}, atom.Img: {},
		atom.Input: {}, atom.Keygen: {}, atom.Link: {}, atom.Meta: {},
		atom.Param: {}, atom.Source: {}, atom.Track: {}, atom.Wbr: {},
	}
)

func (t *TextExtractor) String() string {
	return string(bytes.TrimSpace(t.Buffer.Bytes()))
}

func (t *TextExtractor) HandleToken(token html.Token) {
	switch token.Type {
	case html.SelfClosingTagToken:
		t.maybeSpace() // Self-closing tags take up some amount of space, so implicit whitespace
		switch token.DataAtom {
		case atom.Img, atom.Area:
			if alt, exists := Attr(token, "alt"); exists {
				t.writeSpaceCompressed(alt)
			}
		}
	case html.StartTagToken:
		// In html5, "self-closing" tags no longer require the "/",
		// but this parser doesn't seem to know that, so we have to
		// check for those elements ourselves.
		if _, ok := voidElements[token.DataAtom]; ok {
			t.maybeSpace() // Self-closing tags take up some amount of space, so implicit whitespace
			switch token.DataAtom {
			case atom.Img, atom.Area:
				if alt, exists := Attr(token, "alt"); exists {
					t.writeSpaceCompressed(alt)
				}
			}
			return
		}

		t.push(token.DataAtom)
		if !t.inline {
			t.maybeSpace()
		}

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
			if id, exists := Attr(token, "id"); exists {
				if _, ignore := ignoreIds[id]; ignore && t.ignoring == 0 {
					t.ignoreDepth = 1

					t.ignoring = token.DataAtom
				}
			}
		}
	case html.EndTagToken:
		t.maybePop(token.DataAtom)
		if !t.inline {
			t.maybeSpace()
		}

		switch token.DataAtom {
		case atom.Head:
			t.Reset()
		default:
			if token.DataAtom != 0 && token.DataAtom == t.ignoring && t.ignoreDepth == 0 {
				t.ignoring = 0
			}
		}
	case html.TextToken:
		t.writeSpaceCompressed(token.Data)
	}
}

func (t *TextExtractor) writeSpaceCompressed(s string) {
	if t.ignoring != 0 {
		return
	}

	if len(s) > 0 && unicode.IsSpace(rune(s[0])) {
		t.maybeSpace()
	}

	words := strings.Fields(s)
	if len(words) > 0 {
		t.WriteString(words[0])
		for _, word := range words[1:] {
			t.WriteString(" ")
			t.WriteString(word)
		}
	}
	if len(s) > 0 && unicode.IsSpace(rune(s[len(s)-1])) {
		t.maybeSpace()
	}
}

func (t *TextExtractor) maybeSpace() {
	if t.Len() > 0 && !bytes.HasSuffix(t.Bytes(), []byte(" ")) {
		t.WriteString(" ")
	}
}

func (t *TextExtractor) push(elem atom.Atom) {
	t.currentElem = elem
	t.elemStack = append(t.elemStack, elem)
	_, t.inline = inlineElements[elem]
	if t.ignoreDepth != 0 {
		t.ignoreDepth += 1
	}
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
	if t.ignoreDepth != 0 {
		t.ignoreDepth -= 1
	}

}
