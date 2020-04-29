// Copyright (c) 2020 Hirotsuna Mizuno. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package randdata

import (
	"encoding/binary"
	"math/rand"
	"strings"

	"github.com/tunabay/go-infounit"
)

//
type textGenerator struct {
	width int
}

//
func newTextGenerator() *textGenerator {
	return &textGenerator{
		width: 80,
	}
}

//
var enWordLenTbl = []uint16{
	0x07AF, 0x34E0, 0x6963, 0x8F3E, 0xAAA4, 0xC01E, 0xD472, 0xE3AA,
	0xEF06, 0xF6E7, 0xFB6A, 0xFDDF, 0xFF34, 0xFFC6, 0xFFF9,
}

//
func (g *textGenerator) genWordLen(r *rand.Rand) int {
	seedBytes := make([]byte, 2)
	_, _ = r.Read(seedBytes) // always returns 2, nil
	seed := binary.BigEndian.Uint16(seedBytes)
	for i, t := range enWordLenTbl {
		if seed < t {
			return i + 1
		}
	}
	return 16
}

//
func (g *textGenerator) genLetter(r *rand.Rand) byte {
	const chars = "abcdefghijklmnopqrstuvwxyz"
	return chars[r.Intn(len(chars))]
}

var enSentLenTbl []int

func init() {
	enSentLenTbl = make([]int, 64)
	for i := 0; i < 64; i++ {
		switch {
		case i < 27:
			enSentLenTbl[i] = i + 4
		case i < 45:
			enSentLenTbl[i] = i - 17
		case i < 57:
			enSentLenTbl[i] = i - 30
		default:
			enSentLenTbl[i] = i - 38
		}
	}
}

//
func (g *textGenerator) genSentLen(r *rand.Rand) int {
	seedByte := make([]byte, 1)
	_, _ = r.Read(seedByte) // always returns 1, nil
	seed := int(seedByte[0] >> 2)
	return enSentLenTbl[seed]
}

//
func (g *textGenerator) genSentence(r *rand.Rand) []string {
	var words []string
	sentLen := g.genSentLen(r)
	for i := 0; i < sentLen; i++ {
		words = append(words, g.genWord(r))
	}
	words[0] = strings.Title(words[0])
	words[len(words)-1] += "."
	return words
}

//
func (g *textGenerator) genParagraph(r *rand.Rand) []string {
	numSentByte := make([]byte, 1)
	_, _ = r.Read(numSentByte) // always returns 1, nil
	numSent := int(numSentByte[0]>>5) + 5

	var words []string
	for i := 0; i < numSent; i++ {
		sw := g.genSentence(r)
		words = append(words, sw...)
	}

	var lines []string
	curLine := ""
	for _, w := range words {
		curW := len(curLine)
		sp := ""
		if 0 < curW {
			curW++
			sp = " "
		}
		if curW+len(w) <= g.width {
			curLine += sp + w
			continue
		}
		curLine += "\n"
		lines = append(lines, curLine)
		curLine = ""
	}
	if curLine != "" {
		curLine += "\n"
		lines = append(lines, curLine)
	}
	return lines
}

//
func (g *textGenerator) genWord(r *rand.Rand) string {
	wLen := g.genWordLen(r)
	wordBytes := make([]byte, wLen)
	for i := 0; i < wLen; i++ {
		wordBytes[i] = g.genLetter(r)
	}
	return string(wordBytes)
}

//
func (g *textGenerator) Gen(r *rand.Rand, pos, rem infounit.ByteCount) ([]byte, error) {
	lines := g.genParagraph(r)
	s := strings.Join(lines, "")
	if 0 < pos {
		s = "\n" + s
	}
	buf := ([]byte)(s)
	return buf, nil
}
