package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"time"
)

const (
	stateNewLine = iota
	stateQuoteEscape
	stateInLine
	stateInQuote
	stateSpace
	stateValue
)

func main() {
	br := bufio.NewReader(os.Stdin)
	bw := bufio.NewWriter(os.Stdout)
	token := bytes.NewBuffer(nil)

	defer func() {
		bw.WriteString("\n\n")
		bw.Flush()
	}()

	go func() {
		for range time.Tick(time.Second) {
			bw.Flush()
		}
	}()

	s := stateNewLine

	var (
		c      byte
		err    error
		indent = "    "
		// didOutputHeader = false
		values     = map[string]string{}
		currentKey = bytes.NewBuffer(nil)
		keys       = []string{}
	)

	for {
		c, err = br.ReadByte()
		if err != nil {
			if err != io.EOF {
				log.Printf("error reading byte: %s", err)
			}
			return
		}

		switch c {
		case '"':
			// log.Println("- QUOTE")
			if s == stateQuoteEscape {
				s = stateInQuote
				token.WriteByte(c)
				continue
			}
			if s == stateValue {
				s = stateInQuote
				continue
			}
			if s == stateInQuote {
				values[currentKey.String()] = token.String()
				log.Printf("finished quoted value %s=%s", currentKey.String(), token.String())
				currentKey.Reset()
				token.Reset()
				s = stateInLine
				continue
			}
		case '\n':
			// log.Println("- NEWLINE")
			// Write all key-value pairs
			if currentKey.Len() != 0 {
				// log.Printf("NEWLINE values %+v", values)
				values[currentKey.String()] = token.String()
				token.Reset()
				currentKey.Reset()
			} else {
				changeColor(bw, Blue, false, None, false)
				io.Copy(bw, token)
				resetColor(bw)
				s = stateNewLine
				bw.WriteByte(c)
			}

			// Get longest value
			maxKeyLen := 0
			for k := range values {
				if len(k) > maxKeyLen {
					maxKeyLen = len(k)
				}
				keys = append(keys, k)
			}

			sort.Strings(keys)

			fmtStr := fmt.Sprintf("%%-%ds", maxKeyLen)
			// log.Printf("fmtStr %s", fmtStr)

			for _, k := range keys {
				if k == "t" {
					continue
				}

				bw.WriteByte('\n')
				bw.WriteString(indent)
				color := Green
				if k == "msg" {
					color = Yellow
				}
				changeColor(bw, color, false, None, false)
				fmt.Fprintf(bw, fmtStr, k)
				resetColor(bw)
				bw.WriteString("  ")
				bw.WriteString(values[k])
			}
			bw.WriteByte('\n')

			keys = keys[0:0]
			values = map[string]string{}
			s = stateNewLine
		case '=':
			if s == stateInQuote {
				token.WriteByte(c)
				continue
			}
			currentKey.Reset()
			io.Copy(currentKey, token)
			token.Reset()
			s = stateValue
			log.Printf("got key %s", currentKey.String())
		case ' ':
			if token.Len() == 0 {
				continue
			}

			if s == stateInQuote {
				token.WriteByte(c)
				continue
			}

			if s == stateValue {
				values[currentKey.String()] = token.String()
				// log.Printf("got unquoted value %s=%s", currentKey.String(), token.String())
				token.Reset()
				s = stateInLine
				continue
			}

			changeColor(bw, Blue, false, None, false)
			io.Copy(bw, token)
			resetColor(bw)
			s = stateNewLine
			bw.WriteByte(c)
		case '\\':
			if s == stateInQuote {
				s = stateQuoteEscape
			}
			fallthrough
		default:
			token.WriteByte(c)
		}
	}
}

// Color is the type of color to be set.
type Color int

const (
	// No change of color
	None = Color(iota)
	Black
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

func resetColor(w io.Writer) {
	fmt.Fprint(w, "\x1b[0m")
}

func changeColor(w io.Writer, fg Color, fgBright bool, bg Color, bgBright bool) {
	if fg == None && bg == None {
		return
	} // if

	s := ""
	if fg != None {
		s = fmt.Sprintf("%s%d", s, 30+(int)(fg-Black))
		if fgBright {
			s += ";1"
		} // if
	} // if

	if bg != None {
		if s != "" {
			s += ";"
		} // if
		s = fmt.Sprintf("%s%d", s, 40+(int)(bg-Black))
	} // if

	s = "\x1b[0;" + s + "m"
	fmt.Fprint(w, s)
}
