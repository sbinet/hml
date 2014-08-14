package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

var (
	g_ifname = flag.String("i", "", "path to input file (or STDIN if empty)")
)

func main() {
	flag.Parse()

	var err error
	var r io.Reader
	if *g_ifname == "" {
		if flag.NArg() <= 0 {
			r = os.Stdin
		} else {
			f, err := os.Open(flag.Arg(0))
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()
			r = f
		}
	} else {
		f, err := os.Open(*g_ifname)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		r = f
	}

	rows := make([]Row, 0, 1024)
	dec := NewDecoder(r)
	for {
		i := len(rows)
		rows = append(rows, Row{})
		err = dec.Decode(&rows[i])
		if err != nil {
			rows = rows[:i]
			break
		}
	}

	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
	err = nil

	for _, row := range rows {
		fmt.Fprintf(os.Stdout, "%#v\n", row)
	}

	fmt.Fprintf(os.Stdout, "#rows: %d\n", len(rows))
}
