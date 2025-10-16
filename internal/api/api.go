package api

type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
)

type (
	Logger interface {
		Log(LogLevel, string, ...interface{})
	}

	SearchQuery struct {
		Input  []string
		Output []string
	}

	FileSearchResult struct {
		Filename string
		Filepath string
		Results  []SearchResult
	}

	SearchResult struct {
		Line                uint
		FuncDeclarationLine string
		Coments             []string
	}
)
