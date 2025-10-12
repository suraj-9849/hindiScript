# HindiScript VSCode Extension

VSCode extension for HindiScript - a Hindi-based programming language.

## Installation

### 1. Install Dependencies

Navigate to the extension directory and install dependencies:

```bash
cd extension
npm install
```

### 2. Compile the Extension

```bash
npm run compile
```

### 3. Install the Extension

#### Option A: Install Locally for Development

1. Press `F5` in VSCode while in the extension directory
2. This will open a new VSCode window with the extension loaded

#### Option B: Package and Install

1. Install vsce (if you haven't already):
   ```bash
   npm install -g @vscode/vsce
   ```

2. Package the extension:
   ```bash
   vsce package
   ```

3. Install the `.vsix` file:
   - Open VSCode
   - Go to Extensions (Ctrl+Shift+X)
   - Click on the "..." menu at the top right
   - Select "Install from VSIX..."
   - Choose the generated `.vsix` file

### 4. Ensure HindiScript Interpreter is Available

The extension requires `hlang.exe` to be in the root of your workspace for error detection. Make sure your HindiScript interpreter is compiled and accessible.

## Usage

1. Open any `.hlang` file
2. You'll see syntax highlighting automatically applied
3. Errors will be underlined in real-time as you type

## Example

```hlang
// Variable declaration
ye naam = "Baburao"
ye umar = 20

// Print statements
bol("Naam:")
bol(naam)

// Function
firseKaro add(a, b): number {
    wapas bhejo a + b
}

// Loops
dohraye {
    bol("Namaste")
    roko
}
```

### Building

```bash
npm run compile
```

### Watching for Changes

```bash
npm run watch
```