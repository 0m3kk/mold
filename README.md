# **Mold**

**A simple and powerful CLI tool for scaffolding projects from templates.**

Mold helps you generate project structures, configuration files, or any boilerplate code from reusable templates. It uses Go's built-in text/template engine for rendering and provides a straightforward command-line interface to streamline your development workflow.

## **Features**

- **Template Initialization**: Quickly create a home for your templates.
- **Template Management**: Easily list all available template sets.
- **Flexible Template Path**: Specify a custom directory for your templates using a global flag.
- **Data-Driven Rendering**: Use JSON or YAML files to provide data for your templates, ensuring a clean separation between logic and configuration.
- **Direct File Copying**: Non-template files are copied as-is, preserving your project structure perfectly.
- **Smart Suggestions**: Recommends an example data file if one is found in your template directory.

## **Installation**

To get started, clone the repository and build the binary:

```sh
# Clone the repository
git clone https://github.com/om3kk/mold
cd mold

# Build the binary
go build -o mold ./cmd/mold

# (Optional) Move the binary to a location in your PATH
# For example, on Linux/macOS:
# sudo mv mold /usr/local/bin/
```

## **Usage**

Mold is operated through a series of commands and flags.

### **Global Flag**

- `--dir`, `-t <path>`: Specifies the directory where your templates are stored. This flag works with `init`, `list`, and `apply`. It defaults to `./templates`.

### **Commands**

#### **mold init**

Initializes the templates directory.

```sh
# Create the default 'templates' directory
mold init

# Create a custom directory for templates
mold init --dir ./my-custom-templates
```

#### **mold list**

Lists all available template sets (subdirectories) within the templates directory.

```sh
# List templates from the default directory
mold list

# List templates from a custom directory
mold list -t ./my-custom-templates
```

#### **mold apply <template_path>**

Applies a template from a specific path, rendering `.tmpl` files and copying others to an output directory.

**Arguments:**

- `<template_path>`: The direct path to the template directory you want to use.

**Flags:**

- `--output`, `-o <path>`: The directory where the project will be generated. Defaults to the current directory (`.`).
- `--data-file`, `-d <path>`: **(Required)** The path to a JSON or YAML file containing data for your placeholders.

**Example:**

```sh
mold apply ./templates/go-cli -d ./project-data.yml -o ./my-new-app
```

## **Example Workflow**

Let's create a simple "go-cli" template and use it to scaffold a new project.

### **1. Project Structure**

First, set up your templates directory and a template for a Go CLI tool.

```
.
├── project-data.yml      # Data for our new project
└── templates/            # Our main templates directory
    └── go-cli/           # The 'go-cli' template set
        ├── .gitignore.tmpl # A template for the gitignore file
        ├── go.mod.tmpl     # A template for the go.mod file
        └── main.go         # A static file that will be copied directly
```

### **2. Create Template Files**

Populate the template files with placeholders.

**templates/go-cli/go.mod.tmpl:**

```go
module {{.ModuleName}}

go 1.24.3
```

**templates/go-cli/.gitignore.tmpl:**

```gitignore
# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with `go test -c`
*.test

# Output of the go coverage tool, created with `go test -coverprofile`
*.out

# Project-specific binary
{{.BinaryName}}
```

**templates/go-cli/main.go:**

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World\!")
}
```

### **3. Create a Data File**

Create a YAML file with the data for the placeholders we defined.

**project-data.yml:**

```yaml
ModuleName: "github.com/user/new-awesome-cli"
BinaryName: "new-awesome-cli"
```

### **4. Run Mold**

Now, use the apply command to generate the project.

```sh
# Apply the 'go-cli' template using our data file
# and output to a new directory called 'my-new-app'

mold apply ./templates/go-cli -d ./project-data.yml -o ./my-new-app
```

### **5. Verify the Output**

Check the newly created `my-new-app` directory.

**my-new-app/go.mod (Rendered):**

```go
module github.com/user/new-awesome-cli

go 1.24.3
```

**my-new-app/.gitignore (Rendered):**

```gitignore
# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with `go test -c`
*.test

# Output of the go coverage tool, created with `go test -coverprofile`
*.out

# Project-specific binary
new-awesome-cli
```

**my-new-app/main.go (Copied):**

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World\!")
}
```

You have successfully scaffolded a new project using Mold!

## **License**

This project is licensed under the MIT License.
