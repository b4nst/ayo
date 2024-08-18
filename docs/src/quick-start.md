# Quick Start

## Pre-requisites

You will need [ollama](https://ollama.com/download) installed and running on your machine, or an accessible server.

<div class="warning">

Tool support is a recent addition to ollama, and not all models are compatible. Please visit the [models page](https://ollama.com/search?c=tools) to select a compatible model.

</div>

## Install

Visit the [latest release page](https://github.com/b4nst/ayo/releases/latest) and download the appropriate binary for your system.
Put the binary somewhere accessible in your PATH.

## Add Tools

Put a list of tools you want ayo to be able to use in the toolbox directory.
By default ayo will look into `~/.config/ayo/toolbox`.
More about tools in the [tools section](./usage/tools.md).

## Run

```shell,lang=shell,icon=.devicon-bash-plain
ayo explain my ~/.profile file
```
