# 🧰 Tools

Ayo can use all tools from its toolbox to take actions on your local system.

## Toolbox

The toolbox is a directory where ayo looks for tools to use.
By default, ayo will look into `~/.config/ayo/toolbox`.
You can change the toolbox location by setting the `AYO_TOOLBOX` environment variable.

## Tool

A tool is a JSON file that describes how ayo can use a command on your system.

```json,lang=json,icon=.devicon-json-plain,filepath=~/.config/ayo/tools/cat.json
{
  "type": "function",
  "function": {
    "name": "cat",
    "description": "Get the content of a file",
    "parameters": {
      "type": "object",
      "properties": {
        "file": {
          "type": "string",
          "description": "The file to read"
        }
      },
      "required": [
        "file"
      ]
    }
  },
  "cmd": "cat",
  "args": [
    "{{.file}}"
  ]
}
```

Breaking down the JSON file:
- `type` is the type of tool. Currently, only `function` is supported.
- `function` describes the function which will be sent to the model.
    - `name` is the name of the function. It **must** be unique in ayo's toolbox.
    - `description` is a short description of the function.
    - `parameters` describes the parameters of the function.
        - `type` is the type of the parameters. Currently, only `object` is supported.
        - `properties` describes the properties of the object.
        - `required` is an array of required properties.
- `cmd` is the command to run. It will be executed using the `os/exec` package.
- `args` is an array of arguments to pass to the command. The arguments can be templated using Go's [text/template](https://pkg.go.dev/text/template) package. The template will be executed with the function's parameters values returned by the model.
