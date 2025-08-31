# Security Policy

## Supported Versions

We actively maintain and provide security updates for the following versions:

| Version | Supported          |
| ------- | ------------------ |
| 0.1.x   | :white_check_mark: |
| < 0.1.0 | :x:                |

## Reporting a Vulnerability

We take security vulnerabilities seriously. If you discover a security vulnerability, please follow these steps:

### 1. **DO NOT** create a public GitHub issue
Security vulnerabilities should be reported privately to prevent exploitation.

### 2. Report the vulnerability
Send an email to the maintainers with the following information:
- **Subject**: `[SECURITY] Vulnerability in workerpool package`
- **Description**: Detailed description of the vulnerability
- **Steps to reproduce**: Clear steps to reproduce the issue
- **Impact**: Potential impact of the vulnerability
- **Suggested fix**: If you have a suggested fix (optional)

### 3. Response timeline
- **Initial response**: Within 48 hours
- **Status update**: Within 7 days
- **Fix timeline**: Depends on severity and complexity

### 4. Disclosure
- Vulnerabilities will be disclosed through GitHub Security Advisories
- CVEs will be requested for significant vulnerabilities
- Security releases will be tagged and documented in CHANGELOG.md

## Security Best Practices

### For Users
- Always use the latest stable version
- Run security scans: `make security`
- Keep dependencies updated: `make deps-update`
- Review code changes before deployment

### For Contributors
- Never commit sensitive information
- Run security scans before submitting PRs
- Follow secure coding practices
- Use the provided Makefile targets for validation

## Security Tools

The project includes several security tools:

### Automated Security Scanning
```bash
# Run security scan
make security

# Install gosec if not available
go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
```

### Pre-commit Security Checks
The project includes pre-commit hooks that run security checks automatically:
- gosec security scanner
- Private key detection
- Debug statement detection
- Merge conflict detection

### CI/CD Security
GitHub Actions automatically runs:
- Security scans on every PR
- Dependency vulnerability checks
- Code quality and security validation

## Known Security Considerations

### Worker Pool Specific
- **Resource Exhaustion**: The worker pool includes buffer overflow protection
- **Context Cancellation**: Proper context management prevents resource leaks
- **Concurrent Access**: Mutex protection prevents race conditions
- **Error Handling**: Comprehensive error propagation prevents silent failures

### General Go Security
- **Dependency Management**: Uses Go modules with checksum verification
- **Input Validation**: Generic types provide compile-time safety
- **Memory Management**: Proper cleanup prevents memory leaks

## Security Contacts

- **Primary**: Repository maintainers
- **Backup**: GitHub Security Team
- **Emergency**: Create a private security advisory

## Security Updates

Security updates are released as:
- **Patch releases** (0.1.x) for security fixes
- **Security advisories** for detailed vulnerability information
- **CVE assignments** for significant vulnerabilities

## Reporting Security Issues

**Email**: [Maintainer email - to be added]
**GitHub**: [Create a private security advisory](https://github.com/go-foundations/workerpool/security/advisories/new)

---

Thank you for helping keep the Worker Pool package secure! 🔒
