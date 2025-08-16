# Security Policy

## Supported Versions

We actively support the following versions with security updates:

| Version | Supported          | Status |
| ------- | ------------------ | ------ |
| 1.0.x   | :white_check_mark: | Current stable release |
| < 1.0   | :x:                | Pre-release versions |

## Reporting a Vulnerability

We take security vulnerabilities seriously. If you discover a security vulnerability in this project, please follow these steps:

### 1. **Do Not** Create a Public Issue

Please do not report security vulnerabilities through public GitHub issues, discussions, or pull requests.

### 2. Report Privately

Send an email to the project maintainers with the following information:

- **Subject**: `[SECURITY] Vulnerability Report for go-graphql-example`
- **Description**: A clear description of the vulnerability
- **Steps to Reproduce**: Detailed steps to reproduce the issue
- **Impact**: Assessment of the potential impact
- **Suggested Fix**: If you have suggestions for fixing the issue

### 3. Response Timeline

- **Initial Response**: We will acknowledge receipt within 48 hours
- **Investigation**: We will investigate and provide an initial assessment within 5 business days
- **Resolution**: We aim to resolve critical vulnerabilities within 30 days
- **Disclosure**: We will coordinate responsible disclosure with you

### 4. What to Expect

**If the vulnerability is accepted:**

- We will work with you to understand and resolve the issue
- We will create a security advisory
- We will release a patch as soon as possible
- We will credit you in the security advisory (unless you prefer to remain anonymous)

**If the vulnerability is declined:**

- We will provide a detailed explanation of why we don't consider it a security issue
- We may suggest alternative ways to address your concerns

## Security Best Practices

When using this project, please follow these security best practices:

### Database Security

- Use strong, unique passwords for database connections
- Enable SSL/TLS for database connections in production
- Regularly update database credentials
- Implement proper database access controls

### Configuration Security

- Never commit sensitive configuration values to version control
- Use environment variables for sensitive data
- Regularly rotate API keys and secrets
- Implement proper access controls for configuration files

### Docker Security

- Regularly update base images
- Scan images for vulnerabilities
- Use non-root users in containers (already implemented)
- Implement proper network segmentation

### GraphQL Security

- Implement query complexity analysis
- Set appropriate query depth limits
- Use proper authentication and authorization
- Validate and sanitize all inputs

### General Security

- Keep dependencies up to date
- Regularly audit dependencies for vulnerabilities
- Implement proper logging and monitoring
- Use HTTPS in production environments

## Security Features

This project includes several security features:

- **Non-root Docker containers**: Containers run as non-privileged users
- **Input validation**: GraphQL inputs are validated and sanitized
- **Structured logging**: Security events are properly logged
- **Health checks**: Built-in monitoring for security-relevant components
- **Dependency management**: Regular dependency updates and vulnerability scanning

## Security Audits

We recommend regular security audits for production deployments:

- **Code Review**: Regular security-focused code reviews
- **Dependency Scanning**: Automated scanning for vulnerable dependencies
- **Container Scanning**: Regular scanning of Docker images
- **Penetration Testing**: Periodic security testing of deployed applications

## Contact

For security-related questions or concerns, please contact the project maintainers.

---

**Note**: This security policy is subject to change. Please check back regularly for updates.
