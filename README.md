# HindiScript

A programming language designed for Hindi speakers! HindiScript lets you write code using Hindi keywords and syntax, making programming more accessible and intuitive.

## HindiScript Keywords

| Hindi         | English Equivalent | Description             |
| ------------- | ----------------- | ----------------------- |
| `ye`          | var/let           | Declare a variable      |
| `bol`         | print             | Output to console       |
| `agar`        | if                | Conditional statement   |
| `ya`          | else              | Alternate condition     |
| `ya fir`      | else if           | Additional condition    |
| `firseKaro`   | function          | Define a function       |
| `jabtak`      | while             | While loop              |
| `dohraye`     | repeat            | Infinite loop           |
| `roko`        | break             | Exit loop/statement     |
| `aage badho`  | continue          | Skip to next iteration  |
| `wapas bhejo` | return            | Return from function    |

## Getting Started

### Installation

```bash
git clone https://github.com/suraj-9849/hindiLang.git
cd hindiLang
go build -o hlang ./cmd/main.go
```

### Running HindiScript Programs

```bash
./hlang.exe run main.hlang
```

### Command Line Options

```bash
./hlang.exe run <file.hlang>
./hlang.exe version
./hlang.exe help
```

---

**Inspired by:**  
- [The Programming Language Pipeline](https://www.freecodecamp.org/news/the-programming-language-pipeline-91d3f449c919/)  
- [Building Your Own Programming Language from Scratch](https://hackernoon.com/building-your-own-programming-language-from-scratch)  
- [go/token](https://pkg.go.dev/go/token)  
- [go/ast](https://pkg.go.dev/go/ast)