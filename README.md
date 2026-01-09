# Transformer

Execute workflows by chaining functions together. Define steps in YAML that call built-in functions or custom Go code.

## Installation

```bash
go install github.com/nodyhub/transformer/cmd/transformer@latest
```

Or build from source:

```bash
make build
```

## Usage

Define your steps in a YAML file and execute with:

```bash
transformer run workflow.yml [args...]
```

See [example.yml](example.yml) for a complete example.

## Built-in Functions

- `builtin/echo` - Print messages to stdout
- `builtin/shell` - Execute shell commands  
- `builtin/file/read` - Read file contents
- `builtin/file/write` - Write content to files

## Module Hooks

Transformer supports module hooks to install system-level dependencies and modules required for your workflows. Module hooks are one-time setup commands that install system utilities and libraries before your workflow runs.

### Managing Modules

List available modules:
```bash
transformer modules list
```

Install a module or modules:
```bash
transformer modules install <module>
transformer modules install <module1>,<module2>
transformer modules install all
```

### Built-in Module Hooks

The `builtin/shell` module registers hooks that install common system dependencies like curl, wget, git, jq, and other utilities required for shell scripts and system integrations.

## Step Structure

```yaml
- id: step_id              # optional: identifier for referencing outputs
  name: Human description  # optional: displayed in logs
  uses: builtin/function   # required: function to execute
  with:                    # required: function parameters
    key: value
```

## Variable Substitution

Steps can reference:
- Previous step outputs: `${{ outputs.step_id.property }}`
- CLI arguments: `${{ args.0 }}` or `${{ args }}`
- Environment variables: `${{ env.HOME }}`

## Documentation

- [Built-in Functions](./builtin/)
- [Transformer Engine](./transformer.go)
