# Contributing to cxusage

Thank you for your interest in contributing to cxusage! This document provides guidelines for contributing to the project.

## 🚀 Development Process

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone git@github.com:yourusername/cxusage.git
   cd cxusage
   ```
3. **Create a feature branch** from main:
   ```bash
   git checkout -b feature/amazing-feature
   ```
4. **Make your changes** and test locally
5. **Commit your changes** with clear messages
6. **Push to your fork**:
   ```bash
   git push origin feature/amazing-feature
   ```
7. **Create a Pull Request** on GitHub

## 🛠️ Development Setup

### Prerequisites
- **Go 1.21+** installed
- **Git** for version control
- **OpenAI Codex CLI** for testing (optional)

### Building
```bash
# Build the project
./scripts/build.sh

# Or manual build
go build -o cxusage ./cmd/cxusage
```

### Testing
```bash
# Run tests
go test ./...

# Test the binary
./cxusage demo
./cxusage --help
```

## 📝 Code Guidelines

### Go Style
- Follow standard **Go formatting** (`go fmt`)
- Use **clear, descriptive names** for functions and variables
- Add **comments for exported functions**
- Keep functions **focused and small**

### Project Structure
- **Commands** go in `internal/commands/`
- **Core logic** goes in appropriate packages (`internal/blocks/`, `internal/live/`, etc.)
- **Utilities** go in `internal/utils/`
- **Types** go in `internal/types/`

### Commit Messages
Use clear, descriptive commit messages:
```bash
# Good
git commit -m "Add token limit warnings to live dashboard"
git commit -m "Fix progress bar color calculation"

# Avoid
git commit -m "fix stuff"
git commit -m "update"
```

## 🎯 Types of Contributions

### 🐛 Bug Fixes
- Fix incorrect cost calculations
- Resolve display formatting issues
- Handle edge cases in data parsing

### ✨ New Features
- Additional output formats
- New aggregation methods
- Enhanced live monitoring features
- Support for other OpenAI tools

### 📚 Documentation
- Improve README examples
- Add usage guides
- Update help text
- Fix typos

### 🎨 UI/UX Improvements
- Better color schemes
- Improved table formatting
- Enhanced progress bars
- More responsive design

## 🧪 Testing Guidelines

### Manual Testing
- Test with real Codex CLI data if available
- Use `cxusage demo` to verify output formatting
- Test all commands and flags
- Verify error handling

### Test Coverage
- Add tests for new functionality
- Test edge cases (empty data, large datasets)
- Verify error conditions

## 📋 Pull Request Guidelines

### Before Submitting
- [ ] Code builds successfully
- [ ] All tests pass
- [ ] Code follows Go standards
- [ ] Documentation updated if needed
- [ ] Tested locally

### PR Description
Include:
- **What** the change does
- **Why** the change is needed
- **How** to test the change
- **Screenshots** if UI changes

### Example PR Template
```markdown
## Summary
Brief description of changes

## Changes
- Added feature X
- Fixed bug Y
- Updated documentation Z

## Testing
- [ ] Tested with `cxusage demo`
- [ ] Verified all commands work
- [ ] Checked edge cases

## Screenshots
(if applicable)
```

## 🚨 Important Notes

### Data Privacy
- Never commit actual usage data or logs
- Respect user privacy in examples
- Use synthetic data for tests

### Performance
- Consider memory usage for large datasets
- Optimize file parsing performance
- Test with realistic data sizes

### Compatibility
- Maintain backwards compatibility when possible
- Test on different operating systems
- Consider different terminal environments

## 🤔 Questions?

- **Issues**: Open a GitHub issue for bugs or feature requests
- **Discussions**: Use GitHub discussions for questions
- **Security**: Email security issues privately

## 🎉 Recognition

Contributors will be:
- Added to the acknowledgments section
- Mentioned in release notes for significant contributions
- Credited in commit messages

Thank you for helping make cxusage better! 🙏

---

**Happy coding!** 💻✨