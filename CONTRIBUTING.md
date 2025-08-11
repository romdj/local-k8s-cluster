# Contributing to K3s Cluster Manager

Thank you for your interest in contributing! This document provides guidelines for contributing to the project.

## Development Setup

1. **Clone the repository**:
   ```bash
   git clone https://github.com/your-username/local-k8s-cluster.git
   cd local-k8s-cluster
   ```

2. **Run the setup script**:
   ```bash
   ./setup-dev.sh
   ```
   
   This will install all necessary development tools including:
   - Pre-commit hooks
   - Go linting tools
   - Conventional commits helper

## Development Workflow

### Making Changes

1. **Create a branch** from `main`:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes** following the coding standards

3. **Test your changes**:
   ```bash
   cd local-k8s-cluster-go
   make test
   make lint
   ```

4. **Commit using conventional commits**:
   ```bash
   git cz  # Use commitizen for guided commit messages
   # or manually follow the format: type(scope): description
   ```

### Commit Message Format

We use [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

**Types:**
- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation only changes
- `style`: Changes that do not affect the meaning of the code
- `refactor`: A code change that neither fixes a bug nor adds a feature
- `perf`: A code change that improves performance
- `test`: Adding missing tests
- `build`: Changes that affect the build system
- `ci`: Changes to CI configuration
- `chore`: Other changes that don't modify src or test files

**Examples:**
```
feat(cli): add cluster status command
fix(k8s): resolve connection timeout issue
docs: update installation guide
test(apps): add deployment manager tests
```

### Pre-commit Hooks

Pre-commit hooks run automatically and will:
- Format Go code with `go fmt` and `goimports`
- Run `go vet` and `golangci-lint`
- Run quick tests (`go test -short`)
- Lint YAML and Markdown files
- Check for secrets and large files
- Validate conventional commit format

If hooks fail, fix the issues and commit again.

## Code Standards

### Go Code Style

- Follow standard Go conventions
- Use `gofmt` and `goimports` (handled by pre-commit)
- Write meaningful comments for exported functions
- Keep functions small and focused
- Use meaningful variable names

### Testing

- Write unit tests for new functionality
- Maintain or improve test coverage
- Use table-driven tests where appropriate
- Mock external dependencies

### Documentation

- Update README.md for user-facing changes
- Add inline comments for complex logic
- Update API documentation
- Include examples where helpful

## Pull Request Process

1. **Ensure CI passes**: All tests, linting, and security checks must pass

2. **Update documentation**: Include any necessary documentation updates

3. **Write descriptive PR description**:
   - What changes were made
   - Why they were made
   - How to test them

4. **Link related issues**: Reference any GitHub issues

5. **Request review**: Tag relevant maintainers

## Release Process

Releases are automated using semantic-release:

1. **Merge to main**: All changes go through pull requests to main
2. **Automatic versioning**: Based on conventional commit types
3. **Changelog generation**: Automatically generated from commits
4. **Binary releases**: Multi-platform binaries built and released
5. **Docker images**: Container images published to GitHub Container Registry

### Version Bumps

- `feat:` → Minor version bump (1.0.0 → 1.1.0)
- `fix:`, `perf:` → Patch version bump (1.0.0 → 1.0.1)
- `BREAKING CHANGE:` → Major version bump (1.0.0 → 2.0.0)

## Project Structure

```
local-k8s-cluster/
├── .github/workflows/     # GitHub Actions workflows
├── local-k8s-cluster-go/  # Go application source
│   ├── cmd/               # Application entrypoint
│   ├── internal/          # Internal packages
│   └── pkg/               # Public packages
├── docs/                  # Documentation
└── scripts/               # Utility scripts
```

## Getting Help

- **Issues**: Report bugs or request features via GitHub Issues
- **Discussions**: Use GitHub Discussions for questions
- **Documentation**: Check existing docs in the repository

## Code of Conduct

- Be respectful and inclusive
- Focus on constructive feedback
- Help others learn and grow
- Follow the Golden Rule

Thank you for contributing! 🚀