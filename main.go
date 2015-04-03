package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	ct "github.com/daviddengcn/go-colortext"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	pairRegexp := regexp.MustCompile(`\w+?=(".+?"|.+?)(\s|$)`)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ": ", 2)
		rest := parts[0]
		if len(parts) > 1 {
			ct.ChangeColor(ct.White, false, ct.None, false)
			fmt.Print(parts[0])
			ct.ResetColor()
			rest = parts[1]
		}

		matches := pairRegexp.FindAllString(rest, -1)
		if len(matches) == 0 {
			ct.ChangeColor(ct.Yellow, true, ct.None, false)
			fmt.Printf(" %s\n", rest)
			continue
		}

		// fmt.Println(rest)
		keys := []string{}
		vals := map[string]string{}
		maxKeyLen := 0
		for _, m := range matches {
			parts := strings.SplitN(m, "=", 2)
			k, v := parts[0], strings.Trim(parts[1], `" `)
			if len(k) > maxKeyLen {
				maxKeyLen = len(k)
			}
			keys = append(keys, k)
			vals[k] = v
		}

		if lvl, ok := vals["lvl"]; ok {
			fmt.Printf(" [%s]", lvl)
			delete(vals, "lvl")
		}

		if event, ok := vals["event"]; ok {
			fmt.Printf(" %s", event)
			delete(vals, "event")
		}

		if msg, ok := vals["msg"]; ok {
			fmt.Printf(" %s", msg)
			delete(vals, "msg")
		}

		delete(vals, "app")

		sort.Strings(keys)

		for _, k := range keys {
			v, ok := vals[k]
			if !ok {
				continue
			}
			fmtString := fmt.Sprintf("\n\t%%-%ds  ", maxKeyLen)
			ct.ChangeColor(ct.Green, true, ct.None, false)
			fmt.Printf(fmtString, k)
			ct.ChangeColor(ct.White, true, ct.None, false)
			fmt.Print(v)
			ct.ResetColor()
		}
		fmt.Println()
	}
}
