# CircleCI Pipeline Setup Guide

This guide explains how to set up and test the CircleCI pipeline for the Go GraphQL Example project.

## Prerequisites

1. **CircleCI Account**: Sign up at [circleci.com](https://circleci.com)
2. **GitHub Repository**: Your project should be hosted on GitHub
3. **Slack Workspace**: For notifications (optional but recommended)

## Setup Steps

### 1. Connect Repository to CircleCI

1. Go to [CircleCI Dashboard](https://app.circleci.com/)
2. Click "Add Projects"
3. Find your repository: `captain-corgi/go-graphql-example`
4. Click "Set Up Project"
5. Choose "Fastest" setup method
6. CircleCI will automatically detect the `.circleci/config.yml` file

### 2. Configure Environment Variables

In your CircleCI project settings, add the following environment variables:

#### Required Variables

- `SLACK_WEBHOOK_URL`: Your Slack webhook URL for notifications

#### Optional Variables

- `GITHUB_TOKEN`: For enhanced GitHub integration
- `COVERALLS_REPO_TOKEN`: If using Coveralls for coverage reporting

### 3. Slack Integration Setup

1. Go to your Slack workspace
2. Navigate to Apps → Incoming Webhooks
3. Create a new webhook for your desired channel
4. Copy the webhook URL
5. Add it as `SLACK_WEBHOOK_URL` in CircleCI environment variables

## Pipeline Configuration

The pipeline includes the following workflows:

### 1. PR to Main Branch
- **Trigger**: Pull requests targeting the `main` branch
- **Jobs**: Lint → Test
- **Purpose**: Ensure code quality before merging

### 2. Cursor Review Tag
- **Trigger**: When a commit is tagged with `cursor-review`
- **Jobs**: Lint → Test
- **Purpose**: Manual trigger for code review

### 3. PR with Cursor Review Tag
- **Trigger**: Pull requests with `cursor-review` tag
- **Jobs**: Lint → Test
- **Purpose**: Special review process

## Testing the Pipeline

### Method 1: Create a Test Pull Request

1. Create a new branch:
   ```bash
   git checkout -b test-circleci-pipeline
   ```

2. Make a small change (e.g., add a comment to a file)

3. Commit and push:
   ```bash
   git add .
   git commit -m "test: trigger CircleCI pipeline"
   git push origin test-circleci-pipeline
   ```

4. Create a pull request to `main` branch

5. Check CircleCI dashboard for pipeline execution

### Method 2: Test with Cursor Review Tag

1. Make a commit on any branch:
   ```bash
   git commit -m "test: trigger cursor-review pipeline"
   ```

2. Tag the commit:
   ```bash
   git tag cursor-review
   git push origin cursor-review
   ```

3. Check CircleCI dashboard for pipeline execution

### Method 3: Local Testing

Test the linting locally before pushing:

```bash
# Test linting
make lint

# Test the full build process
make build

# Run tests
make test
```

## Pipeline Jobs

### Lint Job
- **Purpose**: Code quality checks
- **Commands**: `make lint`
- **Docker Image**: `cimg/go:1.24.3`
- **Notifications**: Slack notification on completion

### Test Job
- **Purpose**: Run tests and generate coverage
- **Commands**: `make test`, `make coverage`
- **Artifacts**: Coverage reports stored
- **Notifications**: Slack notification on completion

## Monitoring and Debugging

### CircleCI Dashboard
- View pipeline status at: `https://app.circleci.com/pipelines/github/captain-corgi/go-graphql-example`
- Check job logs for detailed execution information
- View artifacts (coverage reports, test results)

### Slack Notifications
- Success: Green notification with ✅ emoji
- Failure: Red notification with ❌ emoji
- Includes: Project name, branch, commit hash, author

### Common Issues

1. **Pipeline not triggering**:
   - Check branch filters in workflow configuration
   - Ensure `.circleci/config.yml` is in the correct location
   - Verify repository is connected to CircleCI

2. **Linting failures**:
   - Run `make lint` locally to reproduce
   - Check Go version compatibility
   - Verify all dependencies are installed

3. **Slack notifications not working**:
   - Verify `SLACK_WEBHOOK_URL` environment variable
   - Test webhook URL manually
   - Check Slack channel permissions

## Customization

### Adding New Jobs
1. Define job in `jobs` section
2. Add to appropriate workflow
3. Configure filters and dependencies

### Modifying Triggers
1. Update workflow filters
2. Add new workflow conditions
3. Test with different branch/tag combinations

### Environment-Specific Configuration
1. Add environment variables in CircleCI
2. Use conditional logic in pipeline
3. Configure different behaviors per environment

## Best Practices

1. **Always test locally first**: Run `make lint` and `make test` before pushing
2. **Use meaningful commit messages**: Helps with debugging and tracking
3. **Monitor pipeline performance**: Check execution times and optimize if needed
4. **Keep configurations simple**: Avoid overly complex workflow conditions
5. **Document changes**: Update this guide when modifying pipeline configuration

## Troubleshooting

### Pipeline Validation
```bash
# Validate CircleCI config locally (requires CircleCI CLI)
circleci config validate .circleci/config.yml
```

### Local Docker Testing
```bash
# Test with same Docker image as CircleCI
docker run --rm -v $(pwd):/project -w /project cimg/go:1.24.3 make lint
```

### Debug Mode
Enable debug mode in CircleCI by adding to environment variables:
- `DEBUG`: `true`
- `CIRCLE_DEBUG`: `true`

## Support

- **CircleCI Documentation**: [circleci.com/docs](https://circleci.com/docs)
- **Project Issues**: [GitHub Issues](https://github.com/captain-corgi/go-graphql-example/issues)
- **Slack Integration**: [Slack API Documentation](https://api.slack.com/messaging/webhooks)