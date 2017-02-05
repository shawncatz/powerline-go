package powerline

import (
	"bytes"
	"fmt"
)

type Symbols struct {
	Lock          string
	Network       string
	Separator     string
	SeparatorThin string
	Ellipsis      string
}

func DefaultSymbols() Symbols {
	return Symbols{
		Lock:          "\uE0A2",
		Network:       "\uE0A2",
		Separator:     "\uE0B0",
		SeparatorThin: "\uE0B1",
		Ellipsis:      "\u2026",
	}
}

type Powerline struct {
	ShTemplate    string // still not quite get it
	ColorTemplate string // how to output color
	ShellBg       string
	Reset         string
	Symbols       Symbols
	Segments      []Segment
}

func (p *Powerline) color(prefix string, code string) string {
	return fmt.Sprintf(
		p.ShTemplate,
		fmt.Sprintf(p.ColorTemplate, prefix, code),
	)
}

func (p *Powerline) fgColor(code string) string {
	return p.color("38", code)
}

func (p *Powerline) bgColor(code string) string {
	return p.color("48", code)
}

func (p *Powerline) AppendSegment(s Segment) {
	p.Segments = append(p.Segments, s)
}

func (p *Powerline) PrintSegments() string {
	if len(p.Segments) == 0 {
		return ""
	}

	var nextBg string
	var buffer bytes.Buffer

	for i, segment := range p.Segments {
		if (i + 1) == len(p.Segments) {
			// if it is the last one, switch to shell bg
			nextBg = p.ShellBg
		} else {
			nextBg = p.Segments[i+1].Bg
		}

		// set background for segment
		buffer.WriteString(p.bgColor(segment.Bg))

		// print parts with correct foregrounds
		for j, segPart := range segment.values {
			buffer.WriteString(p.fgColor(segment.Fg))
			buffer.WriteString(fmt.Sprintf(" %s ", segPart))
			if (j + 1) == len(segment.values) {
				// last part switches background to next
				buffer.WriteString(p.fgColor(segment.Bg))
				buffer.WriteString(p.bgColor(nextBg))
				buffer.WriteString(p.Symbols.Separator)
			} else {
				// while not last part
				buffer.WriteString(p.fgColor(segment.sepFg))
				buffer.WriteString(p.Symbols.SeparatorThin)
			}
		}
	}

	buffer.WriteString(p.Reset)
	buffer.WriteString(" ")

	return buffer.String()
}

func NewPowerline(shell string, sym Symbols, segs []Segment, t Theme) Powerline {
	var p Powerline
	if shell == "zsh" {
		p = Powerline{
			ShTemplate: "%s",
			ColorTemplate: "%%{[%s;5;%sm%%}",
			Reset: "%{$reset_color%}",
		}
	} else {
		p = Powerline{
			ShTemplate:    "\\[\\e%s\\]",
			ColorTemplate: "[%s;5;%sm",
			Reset:         "\\[\\e[0m\\]",
		}
	}
	p.ShellBg = t.ShellBg
	p.Symbols = sym

	clean := []Segment{}
	for _, seg := range segs {
		if len(seg.values) > 0 {
			clean = append(clean, seg)
		}
	}
	p.Segments = clean

	return p
}
