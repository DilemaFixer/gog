# GOG - Go Function Signature Finder

A command-line tool that searches Go functions by their signatures (input and output parameters).

## What It Does

GOG scans your Go project for functions matching specific type signatures. Instead of searching by function name, you search by what parameters a function accepts and returns. It displays matching functions with their line numbers, full declarations, and documentation comments.

## Installation

Make sure you have Go 1.25.1 or higher installed, then build the project:

```bash
go build -o gog ./cmd/main.go
```

## Usage

Run the binary with optional flags:

```bash
./gog [flags]
```

### Flags

- `-root <path>` - Path to search in (default: current directory)
- `-log` - Enable debug logging output

### Query Syntax

Enter function signatures in an interactive prompt using the following format:

```
input_params -> output_params
```

### Search Examples

**Search for functions with specific inputs only:**
```
int,string
```
Finds functions that take exactly an `int` and a `string` as parameters and return nothing.

**Search for functions with inputs and outputs:**
```
int,string -> bool
```
Finds functions that take `int` and `string`, and return a `bool`.

**Search for functions by return types only:**
```
-> bool,error
```
Finds functions that take no parameters and return a `bool` and an `error`.

### Parameter Matching

Parameters can be matched by:
- Type name: `string`, `int`, `error`, etc.
- Parameter name: `ctx`, `reader`, etc.
- Combined: `name:type`

Supports complex types:
- Pointers: `*Writer`
- Slices: `[]byte`
- Maps: `map[string]int`

## Example Session

```bash
$ ./gog -root ./myproject
int -> error
Searching for functions...

handler.go (./myproject/handlers.go)
├── Line 12: func (h *Handler) Process(count int) error
   │ // Process validates and handles the request
└── Line 28: func Validate(n int) error
   │ // Validate checks if integer is valid

-> bool,error
result.go (./myproject/result.go)
├── Line 5: func GetResult() (bool, error)
   │ // GetResult retrieves the cached result

q
```

Press `q` to quit the interactive session.

## Exit

Type `q` at any prompt to exit the application.

---

**Repository:** github.com/DilemaFixer/gog
