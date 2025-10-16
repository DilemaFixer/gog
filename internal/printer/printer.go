package printer

import (
	"fmt"

	"github.com/DilemaFixer/gog/internal/api"
)

func PrintFileSearchResults(files []api.FileSearchResult) {
	if len(files) == 0 {
		fmt.Println("No results found.")
		return
	}

	for _, f := range files {
		fmt.Printf("%s (%s)\n", f.Filename, f.Filepath)

		if len(f.Results) == 0 {
			fmt.Println("   └── No matches")
			continue
		}

		for i, r := range f.Results {
			prefix := "├──"
			if i == len(f.Results)-1 {
				prefix = "└──"
			}

			// жирное выделение имени функции
			fmt.Printf("   %s Line %d: \033[1m%s\033[0m\n", prefix, r.Line, r.FuncDeclarationLine)

			if len(r.Coments) > 0 {
				for _, c := range r.Coments {
					fmt.Printf("       │ %s\n", c)
				}
			} else {
				fmt.Println("       │ (no comments)")
			}
		}
		fmt.Println()
	}
}
