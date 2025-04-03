# Contributing to Bootstraper

Thank you for considering contributing to Bootstraper! This document provides guidelines and instructions for contributing to this project.

## Code of Conduct

By participating in this project, you agree to abide by our [Code of Conduct](CODE_OF_CONDUCT.md).

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check the existing issues to avoid duplicates. When you create a bug report, include as many details as possible:

- A clear and descriptive title
- Steps to reproduce the issue
- Expected behavior
- Actual behavior
- Screenshots if applicable
- Your environment information (OS, Go version, etc.)

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion:

- Use a clear and descriptive title
- Provide a detailed description of the suggested enhancement
- Explain why this enhancement would be useful
- Include code examples or mockups if applicable

### Adding New Framework Support

One of the easiest ways to contribute is to add support for new frameworks in the `providers/registry.json` file:

1. Identify a framework that is not yet supported
2. Test the command-line installation process for that framework
3. Add a new entry to the registry.json file following the existing pattern
4. Test your implementation
5. Submit a pull request

Example entry:

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

### Pull Requests

1. Fork the repository
2. Create a new branch: `git checkout -b feature/amazing-feature`
3. Make your changes
4. Test your changes thoroughly
5. Commit your changes: `git commit -m 'Add amazing feature'`
6. Push to the branch: `git push origin feature/amazing-feature`
7. Open a pull request

## Development Setup

### Prerequisites

- Go 1.20 or higher
- Git

### Local Development

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/bootstraper.git
   cd bootstraper
   ```

2. Build the project:
   ```bash
   make build
   ```

3. Run tests:
   ```bash
   make test
   ```

## Style Guidelines

### Git Commit Messages

- Use the present tense ("Add feature" not "Added feature")
- Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
- Limit the first line to 72 characters or less
- Reference issues and pull requests liberally after the first line

### Go Code

- Follow standard Go formatting and style (run `go fmt` before committing)
- Include comments for exported functions, types, and methods
- Write tests for new functionality

## Additional Notes

- If you're making significant changes, please open an issue first to discuss the proposed changes
- Feel free to ask for help if you need guidance or have questions

Thank you for contributing to Bootstraper!
