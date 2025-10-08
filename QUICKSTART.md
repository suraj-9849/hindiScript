# HindiScript

## Installation

1. Make sure Go is installed
2. Clone or download this repository
3. Open Command Prompt or PowerShell in the project folder
4. Build the executable:
   ```
   go build -o hlang.exe ./cmd/main.go
   ```
5. Add `hlang.exe` to your PATH or use it directly with `.\hlang.exe`

## Your First Program

1. Create a file called `helloworld.hlang`

2. Write this code:

```hlang
ye name = "baburao"
bol("Namaste " + name + "!")

ye numbers = 1
jabtak numbers <= 5 {
    bol(numbers)
    numbers = numbers + 1
}
```

3. Run it:

```bash
./hlang.exe run helloworld.hlang
```

## Getting Help

```bash
./hlang.exe help      # Show help
./hlang.exe version   # Show version
```