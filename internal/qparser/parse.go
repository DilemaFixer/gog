package qparser

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/DilemaFixer/gog/internal/api"
)

func ParseQuery() (api.SearchQuery, error) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		for scanner.Scan() {
			text := scanner.Text()
			if text == "" {
				break
			}

			if text == "q" {
				return api.SearchQuery{}, fmt.Errorf("capitulation cancelled")
			}
			return parseInput(text)
		}

		if err := scanner.Err(); err != nil {
			return api.SearchQuery{}, err
		}
	}
}

func parseInput(input string) (api.SearchQuery, error) {
	input = strings.TrimSpace(input)
	query := api.SearchQuery{
		Input:  []string{},
		Output: []string{},
	}

	parts := strings.Split(input, "->")

	left := strings.TrimSpace(parts[0])
	right := ""
	if len(parts) > 1 {
		right = strings.TrimSpace(parts[1])
	}

	if left != "" {
		in := splitParams(left)
		query.Input = append(query.Input, in...)
	}

	if right != "" {
		out := splitParams(right)
		query.Output = append(query.Output, out...)
	}

	return query, nil
}

func splitParams(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}

	var params []string
	var buf strings.Builder
	depth := 0

	for _, r := range s {
		switch r {
		case '(':
			depth++
			buf.WriteRune(r)
		case ')':
			depth--
			buf.WriteRune(r)
		case ',':
			if depth == 0 {
				param := strings.TrimSpace(buf.String())
				if param != "" {
					params = append(params, param)
				}
				buf.Reset()
			} else {
				buf.WriteRune(r)
			}
		default:
			buf.WriteRune(r)
		}
	}

	if buf.Len() > 0 {
		param := strings.TrimSpace(buf.String())
		if param != "" {
			params = append(params, param)
		}
	}

	return params
}
