# Contributing to TS2Go

Thank you for considering contributing to TS2Go! This document provides guidelines and instructions for contributing.

## 🎯 Getting Started

### Prerequisites
- Go 1.22.5 or later
- Node.js 20 or later
- npm or yarn
- Git

### Setting Up Development Environment

1. **Fork and clone the repository**
   ```bash
   git clone https://github.com/YOUR_USERNAME/ts2go.git
   cd ts2go
   ```

2. **Install dependencies**
   ```bash
   # Install Go dependencies (handled by go.work)
   go work sync
   
   # Install Node.js dependencies for TypeScript parser
   cd internal/transpiler/parser
   npm install
   cd ../../..
   ```

3. **Build the project**
   ```bash
   make build
   ```

4. **Run tests**
   ```bash
   make test
   ```

## 🌳 Branching Strategy

We follow a Git Flow branching model:

- `main` - Production-ready code, stable releases only
- `develop` - Integration branch for features, generally stable
- `feature/*` - Feature branches (e.g., `feature/arrow-functions`)
- `copilot/*` - Copilot-assisted development branches
- `bugfix/*` - Bug fix branches
- `release/*` - Release preparation branches

### Creating a Feature Branch

```bash
git checkout develop
git pull origin develop
git checkout -b feature/your-feature-name
```

## 📝 Making Changes

### Code Style

**Go Code:**
- Follow standard Go conventions (`gofmt`, `go vet`)
- Use meaningful variable and function names
- Add comments for exported functions and types
- Keep functions focused and small

**TypeScript/JavaScript Code:**
- Use ESLint configuration in `internal/transpiler/parser`
- Follow existing code style
- Add JSDoc comments for complex functions

### Testing

**Always add tests for new features:**

1. **Unit tests** - For individual functions and components
   ```go
   func TestFeature(t *testing.T) {
       // Test implementation
   }
   ```

2. **Integration tests** - For end-to-end transpilation
   - Add TypeScript fixtures in `tests/fixtures/`
   - Add expected Go output
   - Update `tests/integration_test.go`

3. **Run tests before committing**
   ```bash
   make test
   ```

### Commit Messages

Follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

**Types:**
- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

**Examples:**
```
feat(transpiler): add support for arrow functions
fix(cli): correct output directory handling
docs(readme): update installation instructions
test(control-flow): add tests for if/else statements
```

## 🔄 Pull Request Process

1. **Update your branch with latest develop**
   ```bash
   git checkout develop
   git pull origin develop
   git checkout feature/your-feature-name
   git rebase develop
   ```

2. **Ensure all tests pass**
   ```bash
   make test
   ```

3. **Push your changes**
   ```bash
   git push origin feature/your-feature-name
   ```

4. **Create a Pull Request**
   - Go to GitHub and create a PR from your branch to `develop`
   - Fill in the PR template
   - Link any related issues
   - Request review from maintainers

5. **Address review feedback**
   - Make requested changes
   - Push updates to your branch
   - Re-request review

6. **Merge**
   - Once approved, a maintainer will merge your PR
   - Your branch will be deleted automatically

## 🎯 What to Work On

### High Priority (Phase 22)
See [NEXT_PHASE_ROADMAP.md](docs/NEXT_PHASE_ROADMAP.md) for detailed priority list.

**Critical features needed:**
- Control flow statements (if/else, loops, switch)
- Modern JavaScript syntax (arrow functions, template literals)
- Async/await support
- Try/catch error handling

### Good First Issues
Look for issues labeled `good first issue` on GitHub. These are typically:
- Documentation improvements
- Small bug fixes
- Test additions
- Minor feature enhancements

## 📚 Architecture Overview

```
ts2go/
├── cmd/ts2go/              # CLI entry point
├── pkg/cli/                # CLI commands (transpile, analyze, ui)
├── internal/
│   ├── transpiler/         # Core transpilation logic
│   ├── analyzer/           # TypeScript analysis
│   ├── mapper/             # npm to Go package mapping
│   ├── optimizer/          # Code optimization
│   └── orchestrator/       # Build orchestration
├── runtime/                # Go runtime libraries (Node.js API implementations)
├── desktop-ui/             # Tauri desktop application
├── tests/                  # Integration tests
└── docs/                   # Documentation
```

### Key Files
- `internal/transpiler/transpiler.go` - Main transpiler logic
- `internal/transpiler/emitter.go` - Go code generation
- `internal/transpiler/parser/` - TypeScript AST parser (Node.js)
- `pkg/cli/transpile.go` - Transpile command implementation

## 🐛 Reporting Bugs

**Before submitting a bug report:**
1. Check existing issues to avoid duplicates
2. Try to reproduce with the latest version
3. Collect relevant information

**Bug report should include:**
- Clear title and description
- Steps to reproduce
- Expected behavior
- Actual behavior
- TypeScript input code (minimal example)
- Generated Go code (if applicable)
- Error messages
- Environment (OS, Go version, Node.js version)

## 💡 Suggesting Features

**Feature requests should include:**
- Clear title and description
- Use case and motivation
- Expected behavior
- Examples (TypeScript input and desired Go output)
- Impact on existing features
- Potential implementation approach (optional)

## 🧪 Running Specific Tests

```bash
# Run all tests
make test

# Run tests in a specific package
cd internal/transpiler && go test -v

# Run a specific test
cd tests && go test -v -run TestTranspileSimple

# Run with coverage
cd tests && go test -v -coverprofile=coverage.txt
```

## 📖 Documentation

**When adding features, update:**
- `README.md` - If it changes usage or features
- `docs/EXAMPLES.md` - Add before/after examples
- `docs/STATUS.md` - Update implementation status
- `SPEC.md` - If it affects supported TypeScript features
- `CHANGELOG.md` - Add to Unreleased section

## 🤝 Code of Conduct

### Our Pledge
We pledge to make participation in our project a harassment-free experience for everyone, regardless of age, body size, disability, ethnicity, gender identity and expression, level of experience, nationality, personal appearance, race, religion, or sexual identity and orientation.

### Our Standards
- Using welcoming and inclusive language
- Being respectful of differing viewpoints
- Gracefully accepting constructive criticism
- Focusing on what is best for the community
- Showing empathy towards other community members

## 📞 Getting Help

- **Documentation:** Check [docs/](docs/) folder
- **GitHub Issues:** For bugs and feature requests
- **GitHub Discussions:** For questions and general discussion
- **Pull Request Comments:** For code review questions

## 🏆 Recognition

Contributors are listed in our README and release notes. Thank you for helping make TS2Go better!

## 📄 License

By contributing, you agree that your contributions will be licensed under the MIT License.
