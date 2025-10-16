package finder

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/DilemaFixer/gog/internal/api"
)

type (
	Finder struct {
		root   string
		logger api.Logger
	}
)

func NewFinder(root string, logger api.Logger) (*Finder, error) {
	if err := validatePath(root); err != nil {
		return nil, err
	}

	return &Finder{
		root:   root,
		logger: logger,
	}, nil
}

func validatePath(path string) error {
	exist, err := exists(path)
	if err != nil {
		return fmt.Errorf("geting error when try check is root folder exist : %v\n", err)
	}

	if !exist {
		return fmt.Errorf("path %s doesn't exist\n", path)
	}
	return nil
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func (finder *Finder) Search(query api.SearchQuery) ([]api.FileSearchResult, error) {
	files := find(finder.root, ".go")
	if !isGoFilesExist(files) {
		return nil, fmt.Errorf("in folder %s not exist any '.go' files\n", finder.root)
	}

	result := make([]api.FileSearchResult, 0)
	for _, filePath := range files {
		finder.logger.Log(api.LogLevelDebug, "start working on %s", filePath)

		searchResult, err := finder.searchInFile(query, filePath)
		if err != nil {
			return nil, fmt.Errorf("error parsing %s : %v", filePath, err)
		}

		if len(searchResult) != 0 {
			result = append(result, api.FileSearchResult{Results: searchResult, Filepath: filePath, Filename: filepath.Base(filePath)})
		}
	}

	return result, nil
}

func find(root, ext string) []string {
	var a []string
	filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if filepath.Ext(d.Name()) == ext {
			a = append(a, s)
		}
		return nil
	})
	return a
}

func isGoFilesExist(files []string) bool {
	if len(files) == 0 {
		return false
	}
	return true
}

func (finder *Finder) searchInFile(query api.SearchQuery, filePath string) ([]api.SearchResult, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, nil, parser.AllErrors|parser.ParseComments)
	if err != nil {
		return nil, err
	}
	finder.logger.Log(api.LogLevelDebug, "parsed %s", filePath)
	return finder.searchInDeclarations(file, fset, query)
}

func (finder *Finder) searchInDeclarations(file *ast.File, fset *token.FileSet, query api.SearchQuery) ([]api.SearchResult, error) {
	result := make([]api.SearchResult, 0)

	for _, decl := range file.Decls {
		if funcDecl, isFuncDecl := isFuncDeclaration(decl); isFuncDecl {
			if !isMatch(query, funcDecl.Type) {
				continue
			}

			sresult, err := funcDeclarationToSearchResult(funcDecl, fset)
			if err != nil {
				return nil, err
			}
			result = append(result, sresult)
		}
	}

	return result, nil
}

func isFuncDeclaration(decl ast.Decl) (*ast.FuncDecl, bool) {
	fn, ok := decl.(*ast.FuncDecl)
	return fn, ok
}

func isMatch(expected api.SearchQuery, fn *ast.FuncType) bool {
	if fn.Params == nil || len(expected.Input) != len(fn.Params.List) {
		return false
	}
	for i, field := range fn.Params.List {
		if !matchField(expected.Input[i], field) {
			return false
		}
	}

	if fn.Results == nil && len(expected.Output) > 0 {
		return false
	}
	results := 0
	if fn.Results != nil {
		results = len(fn.Results.List)
	}
	if len(expected.Output) != results {
		return false
	}
	if fn.Results != nil {
		for i, field := range fn.Results.List {
			if !matchField(expected.Output[i], field) {
				return false
			}
		}
	}

	return true
}

func matchField(expected string, field *ast.Field) bool {
	if expected == "" || field == nil {
		return false
	}

	name := ""
	if len(field.Names) > 0 && field.Names[0] != nil {
		name = field.Names[0].Name
	}

	typ := ""
	if field.Type != nil {
		typ = typeToString(field.Type)
	}

	switch expected {
	case name, typ, name + ":" + typ:
		return true
	default:
		return false
	}
}

func typeToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + typeToString(t.X)
	case *ast.SelectorExpr:
		return typeToString(t.X) + "." + t.Sel.Name
	case *ast.ArrayType:
		return "[]" + typeToString(t.Elt)
	case *ast.MapType:
		return "map[" + typeToString(t.Key) + "]" + typeToString(t.Value)
	case *ast.FuncType:
		return "func"
	default:
		return ""
	}
}

func funcDeclarationToSearchResult(fn *ast.FuncDecl, fset *token.FileSet) (api.SearchResult, error) {
	if fn == nil {
		return api.SearchResult{}, fmt.Errorf("converting func declaration ast node to search result error: ast node is nil")
	}

	res := api.SearchResult{}

	if fset != nil && fn.Pos().IsValid() {
		pos := fset.Position(fn.Pos())
		res.Line = uint(pos.Line)
	}

	var declBuf strings.Builder
	start := fn.Pos()
	end := fn.Body.Lbrace - 1
	srcBytes, err := os.ReadFile(fset.File(start).Name())
	if err != nil {
		return res, fmt.Errorf("reading source file: %w", err)
	}
	if int(end) > len(srcBytes) {
		end = token.Pos(len(srcBytes))
	}
	declBuf.Write(srcBytes[fset.File(start).Offset(start):fset.File(start).Offset(end)])
	res.FuncDeclarationLine = strings.TrimSpace(declBuf.String())

	if fn.Doc != nil {
		for _, c := range fn.Doc.List {
			res.Coments = append(res.Coments, strings.TrimSpace(c.Text))
		}
	}

	return res, nil
}
