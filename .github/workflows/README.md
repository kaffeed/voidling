# GitHub Actions Workflows

This directory contains CI/CD workflows for automated testing and deployment.

## Workflows

### CI Workflow (`ci.yml`)

**Triggers:**
- Push to `main` branch
- Pull requests to `main` branch

**Jobs:**

1. **Lint** - Code quality checks
   - Format checking (`make fmt`)
   - Static analysis (`make vet`)
   - Linting with golangci-lint

2. **Test** - Test execution
   - Unit tests (`make test`)
   - Race detector tests (`make test-race`)

3. **Build** - Multi-platform builds
   - Linux (amd64) with CGO for SQLite
   - Windows (amd64)
   - macOS (amd64, arm64)
   - Artifacts uploaded for 7 days

4. **SQLC** - Database query generation validation
   - Verifies generated code matches queries

5. **Migrations** - Database migration testing
   - Tests migrations can be applied successfully

**Duration**: ~5-10 minutes

### CD Workflow (`cd.yml`)

**Triggers:**
- Version tags (`v*`, e.g., `v1.0.0`)
- Manual workflow dispatch

**Jobs:**

1. **Build** - Production binary
   - Linux binary with CGO enabled
   - Optimized with `-ldflags "-s -w"`
   - Uploaded as artifact

2. **Deploy** - Server deployment
   - Stops systemd service
   - Backs up current binary
   - Copies new binary via SCP
   - Starts service
   - Verifies service health
   - Automatic rollback on failure

3. **Create Release** - GitHub release
   - Creates release from tag
   - Attaches binary
   - Auto-generates release notes

**Duration**: ~3-5 minutes

**Environment**: Uses `production` environment for additional security

## Required Secrets

Configure in: Repository Settings → Secrets and variables → Actions

| Secret | Description | Example |
|--------|-------------|---------|
| `SSH_HOST` | Ubuntu server IP/hostname | `123.45.67.89` |
| `SSH_USERNAME` | SSH username | `voidling` |
| `SSH_KEY` | SSH private key (full content) | `-----BEGIN OPENSSH...` |
| `SSH_PORT` | SSH port | `22` |
| `DEPLOY_PATH` | Deployment directory | `/opt/voidling` |

See [`../implement/SECRETS.md`](../implement/SECRETS.md) for detailed setup.

## Workflow Status

Check status at: https://github.com/kaffeed/voidling/actions

**Badges:**
```markdown
[![CI](https://github.com/kaffeed/voidling/actions/workflows/ci.yml/badge.svg)](https://github.com/kaffeed/voidling/actions/workflows/ci.yml)
[![CD](https://github.com/kaffeed/voidling/actions/workflows/cd.yml/badge.svg)](https://github.com/kaffeed/voidling/actions/workflows/cd.yml)
```

## Usage Examples

### Running CI locally

```bash
# Format code
make fmt

# Run static analysis
make vet

# Run tests
make test
make test-race

# Run all checks
make check

# Build for all platforms
make build-all
```

### Triggering Deployment

```bash
# Create and push version tag
git tag v1.0.0
git push origin v1.0.0

# Or manually via GitHub Actions UI:
# Actions → CD → Run workflow → Select branch/tag
```

### Viewing Logs

**GitHub Actions:**
- Actions tab → Select workflow run → Select job → View logs

**Server Logs:**
```bash
# Real-time logs
sudo journalctl -u voidling -f

# Last 50 lines
sudo journalctl -u voidling -n 50

# Since specific time
sudo journalctl -u voidling --since "1 hour ago"
```

## Troubleshooting

### CI Failures

**Format check fails:**
```bash
make fmt
git add -A
git commit -m "Format code"
```

**Vet fails:**
- Fix reported issues
- Run `make vet` locally to verify

**Tests fail:**
- Run `make test` locally to reproduce
- Fix failing tests
- Ensure all tests pass before pushing

**Build fails:**
- Check Go version (must be 1.24.1)
- Verify all dependencies in go.mod
- Test build locally: `make build-linux`

### CD Failures

**SSH connection fails:**
- Verify `SSH_HOST`, `SSH_USERNAME`, `SSH_PORT` are correct
- Test SSH manually: `ssh -i key user@host`
- Check server firewall allows SSH

**Permission denied:**
- Verify SSH public key is in `~/.ssh/authorized_keys`
- Check file permissions (600 for authorized_keys)
- Ensure private key in `SSH_KEY` secret matches public key

**Service won't start:**
- Check logs: `sudo journalctl -u voidling -n 50`
- Verify `.env` file exists and is readable
- Check binary has execute permissions
- Verify database path is writable

**Deployment succeeds but bot offline:**
- Check Discord token in `.env`
- Verify network connectivity from server
- Check service logs for errors

### Rollback Procedure

**Automatic rollback:**
- CD workflow automatically rolls back if service fails to start
- Previous binary is restored from `.backup` file

**Manual rollback:**
```bash
cd /opt/voidling
sudo systemctl stop voidling
mv voidling voidling.failed
mv voidling.backup voidling
sudo systemctl start voidling
```

**Rollback to specific version:**
```bash
# Re-run deployment for previous version
# Actions → CD → Run workflow → Select previous tag
```

## Security Considerations

**Secrets:**
- Never commit secrets to repository
- Rotate SSH keys periodically
- Use separate keys for CI/CD (not personal keys)
- Store Discord token in server `.env`, not GitHub secrets

**Permissions:**
- CI workflows have `contents: read` (read-only)
- CD workflow has `contents: write` (for creating releases)
- Use dedicated deployment user with limited sudo access
- Service runs as non-root user with restricted permissions

**Best Practices:**
- Enable branch protection on `main`
- Require PR reviews before merging
- Require status checks to pass
- Use signed commits (recommended)
- Tag releases from main branch only

## Performance

**Cache Usage:**
- Go modules cached across runs
- Build cache improves subsequent builds
- Cache key based on go.mod/go.sum

**Optimization:**
- Parallel job execution where possible
- Artifacts retained for 7 days (CI) or 90 days (releases)
- Minimal image layers
- Stripped binaries (-s -w flags)

## Maintenance

**Regular Tasks:**
- Update Go version when new releases available
- Update action versions (check for @v5, @v2, etc.)
- Review and update golangci-lint version
- Rotate SSH keys quarterly
- Review workflow efficiency monthly

**Monitoring:**
- Watch for failed workflow runs
- Check deployment logs regularly
- Monitor server resource usage
- Track deployment frequency

## Support

**Issues:**
- GitHub Issues: https://github.com/kaffeed/voidling/issues
- Check existing issues before creating new ones

**Documentation:**
- Setup: [`../implement/SECRETS.md`](../implement/SECRETS.md)
- Development: [`../README.md`](../README.md)
- Deployment: This file
