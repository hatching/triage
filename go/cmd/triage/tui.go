// Copyright (C) 2019-2020 Hatching B.V.
// All rights reserved.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type selectOptions interface {
	Len() int
	Emphasized(i int) bool
	Display(i int) string
}

func promptSelectOptions(options selectOptions, validate func([]int) bool) []int {
selectLoop:
	for sc := bufio.NewScanner(os.Stdin); ; {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Make your selection by entering the numbers as listed below separated by spaces and finish with enter.\n")

		for i := 0; i < options.Len(); i++ {
			em := " "
			if options.Emphasized(i) {
				em = ">"
			}
			fmt.Fprintf(os.Stderr, " %s% 2d %q\n", em, i, options.Display(i))
		}

		fmt.Fprintf(os.Stderr, "> ")
		sc.Scan()

		selectionStr := sc.Text()
		selection := []int{}
	parseSelectionLoop:
		for _, s := range strings.Split(selectionStr, " ") {
			if s == "" {
				continue
			}
			i, err := strconv.Atoi(s)
			if err != nil {
				fmt.Fprintf(os.Stderr, "bad input: %q\n", s)
				continue selectLoop
			}
			if i < 0 || options.Len() <= i {
				fmt.Fprintf(os.Stderr, "out of range: %d\n", i)
				continue selectLoop
			}
			for _, j := range selection {
				if i == j {
					continue parseSelectionLoop
				}
			}
			selection = append(selection, i)
		}

		if !validate(selection) {
			continue
		}

		if len(selection) > 0 {
			fmt.Fprintf(os.Stderr, "you selected:\n")
			for _, i := range selection {
				fmt.Fprintf(os.Stderr, "  %q\n", options.Display(i))
			}
		}
		return selection
	}
}

type arrayFlags []string

func (i *arrayFlags) String() string {
	return strings.Join(*i, ",")
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}
