package search_engine

import (
	"github.com/DilemaFixer/gog/internal/api"
	"github.com/DilemaFixer/gog/internal/finder"
	"github.com/DilemaFixer/gog/internal/printer"
	"github.com/DilemaFixer/gog/internal/qparser"
)

type SearchEngine struct {
	logger api.Logger
}

func NewSearchEngine(logger api.Logger) *SearchEngine {
	return &SearchEngine{
		logger: logger,
	}
}

func (eng *SearchEngine) StartCyclicExecution(root string) error {
	f, err := finder.NewFinder(root, eng.logger)
	if err != nil {
		return err
	}

	for {
		query, err := qparser.ParseQuery()
		if err != nil {
			return err
		}

		result, err := f.Search(query)
		if err != nil {
			return err
		}

		printer.PrintFileSearchResults(result)
	}
}
