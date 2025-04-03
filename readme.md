# Bootstraper CLI (bt)

[![Go Report Card](https://goreportcard.com/badge/github.com/sharik709/bootstraper)](https://goreportcard.com/report/github.com/sharik709/bootstraper)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Version](https://img.shields.io/badge/version-0.2.0-blue.svg)](https://github.com/sharik709/bootstraper/releases)
[![npm](https://img.shields.io/npm/v/bootstraper-cli)](https://www.npmjs.com/package/bootstraper-cli)

Bootstraper (`bt`) is a unified CLI tool that simplifies project initialization across multiple frameworks, languages, and platforms. Create new projects with your favorite frameworks using a single command-line interface.

## Features

- **Universal Interface**: One command to bootstrap projects with any supported framework
- **Extensible**: Easily add new framework providers via JSON configuration
- **Smart Defaults**: Sensible defaults with customizable options
- **Dependency Checking**: Automatic verification that required tools are installed

## Installation

### Using npm (Recommended)

```bash
npm install -g bootstraper-cli
```

### Using Go

```bash
go install github.com/sharik709/bootstraper@latest
```

### From Binary Releases

Download pre-built binaries from the [releases page](https://github.com/sharik709/bootstraper/releases).

### From Source

```bash
git clone https://github.com/sharik709/bootstraper.git
cd bootstraper
make build
make install
```

## Usage

### Create a New Project

```bash
# Basic usage
bt new [framework] [project-name]

# Examples
bt new next my-nextjs-app
bt new vue my-vue-app
bt new go myproject --module=github.com/username/myproject
```

### With Framework-specific Options

```bash
# Next.js with TypeScript and Tailwind
bt new next my-app --typescript=true --tailwind=true

# Vue with Router and Pinia
bt new vue my-app --typescript=true --router=true --pinia=true

# Laravel with specific version
bt new laravel my-app --version=10.0
```

### List Available Frameworks

```bash
bt list
```

## Supported Frameworks

Bootstraper includes support for many popular frameworks:

- **Frontend**: Next.js, Vue, Angular, Svelte, React, Astro, SolidJS, Nuxt
- **Backend**: Express, Laravel, Django, Spring Boot, Rails, NestJS, FastAPI
- **Mobile**: Flutter, React Native
- **Languages**: Go, Rust, Python

Run `bt list` to see all available frameworks and their descriptions.

## Extending Bootstraper

Bootstraper uses a JSON-based provider registry that makes it easy to add new frameworks without changing the code.

To add a custom framework, modify the `providers/registry.json` file following this structure:

```json
{
  "name": "your-framework",
  "description": "Description of your framework",
  "command": "installation-command",
  "args": ["command", "args", "{project-name}"],
  "dependencies": ["required-commands"],
  "options": {
    "option1": "Description of option1",
    "option2": "Description of option2"
  }
}
```

## Publishing to npm

If you're forking this project and want to publish your own version to npm:

1. Update the package name in `package.json`
2. Build the binary: `make build`
3. Publish to npm: `npm publish`

You'll need to have an npm account and be logged in via `npm login`.

## Contributing

Contributions are welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for details.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
