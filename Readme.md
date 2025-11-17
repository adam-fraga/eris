# Eris CLI

**Eris** is a developer CLI tool that lets you interact with an LLM (like Qwen3) via streaming chat. It supports **capabilities (tools)** that the model can call dynamically, like fetching weather or creating files.

---

## Features

- Stream chat responses from the model.
- Dynamically execute capabilities defined in Go.
- Safe local file creation with `create_file`.
- Extensible: add your own capabilities and expose them to the LLM.

---

## Project Structure

```
eris/
├── capabilities/       # All custom capabilities (tools) the LLM can call
│   ├── capabilities.go
│   └── create_file.go
├── cmd/                # Cobra commands
│   └── prompt.go
├── handler/            # Command logic
│   └── promptHandler.go
├── output/             # Generated files by create_file capability
├── prompts/            # System prompts for LLM context
├── requests/           # HTTP and streaming requests to Ollama
├── go.mod
└── main.go
```

---

## Installation

```bash
git clone <repo-url>
cd eris
go build -o eris .
```

---

## Usage

```bash
./eris prompt
```

**Example interaction:**

```
Ask your question: Save this note for me: Remember to buy milk tomorrow
[Capability result]: File 'output/note.txt' created successfully.
Completed in 2.4s
```

---

## Adding Capabilities

Capabilities are **functions that the LLM can call**.

1. **Create a new Go file** in `capabilities/` (e.g., `mycap.go`).
2. Define your function signature (must match `CallCapability` handling):

```go
package capabilities

func MyCapability(arg1, arg2 string) string {
    return fmt.Sprintf("Did something with %s and %s", arg1, arg2)
}
```

3. Add it to the `Capabilities` map:

```go
var Capabilities = map[string]Capability{
    "my_capability": {
        Name:        "my_capability",
        Description: "Does something useful",
        Func:        MyCapability,
    },
}
```

4. Define the JSON schema for the LLM in `ToChatCapabilitiesInterface`:

```go
func ToChatCapabilitiesInterface() []interface{} {
    caps := []interface{}{}
    for _, cap := range Capabilities {
        caps = append(caps, ChatCapability{
            Type: "function",
            Function: FunctionDescr{
                Name:        cap.Name,
                Description: cap.Description,
                Parameters: map[string]CapabilityProp{
                    "arg1": {Type: "string", Description: "First argument"},
                    "arg2": {Type: "string", Description: "Second argument"},
                },
            },
        })
    }
    return caps
}
```

5. Update `CallCapability` to handle the new capability:

```go
case "my_capability":
    var args struct {
        Arg1 string `json:"arg1"`
        Arg2 string `json:"arg2"`
    }
    if err := json.Unmarshal(argsJson, &args); err != nil {
        return "", err
    }
    f := cap.Func.(func(string, string) string)
    return f(args.Arg1, args.Arg2), nil
```

6. The LLM can now call it dynamically in streaming responses.

---

## Notes

- **All files created** by `create_file` are stored in `output/`.
- **Streaming responses** from the model are handled line by line.
- The system prompt should describe all available capabilities to the LLM.
- You can extend capabilities for **any custom action**, like web requests, calculations, or internal tools.

---
