# Phase 4 Plan 1: CI Workflow Docker Hub Integration Summary

**Completed:** 2026-02-03
**Duration:** ~8 minutes
**Status:** SUCCESS

## One-liner

CI workflow extended with Docker Hub authentication and nanopi-m6 board matrix for automated overlay and raw image builds on tag push.

## What Was Built

### Changes Made

1. **Workflow-level Environment Variables**
   - Added `DOCKERHUB_USERNAME: ${{ vars.DOCKERHUB_USERNAME }}`
   - Added `DOCKERHUB_REGISTRY: docker.io`

2. **sbc-rk3588 Job Updates**
   - Added Docker Hub login step after GHCR login
   - Added overlay push step to Docker Hub after GHCR push
   - Key link: `crane push ... docker.io/123417/talos-sbc-nanopi-m6:$TAG`

3. **boards Job Updates**
   - Added nanopi-m6 to board matrix with `chipset: rk3588s` and `registry: dockerhub`
   - Split "Build installer image" into GHCR and Docker Hub variants
   - Split "Push installer image" into GHCR and Docker Hub variants
   - Split "Build flashable image" into GHCR and Docker Hub variants
   - nanopi-m6 uses `docker.io/123417/talos-sbc-nanopi-m6` as overlay source

4. **Cleanup Step**
   - Updated to logout from both ghcr.io and docker.io

### Key Artifacts

| File | Purpose |
|------|---------|
| `.github/workflows/ci.yaml` | CI workflow with Docker Hub support and nanopi-m6 board |

### Registry Configuration

| Board | Registry | Overlay Image |
|-------|----------|---------------|
| rock-5a | GHCR | ghcr.io/{owner}/talos-sbc-rk3588 |
| rock-5b | GHCR | ghcr.io/{owner}/talos-sbc-rk3588 |
| nanopi-m6 | Docker Hub | docker.io/123417/talos-sbc-nanopi-m6 |

## Decisions Made

| Decision | Rationale |
|----------|-----------|
| Conditional step splitting | Cleaner than complex bash conditionals in single step |
| Use `|| true` for logout | Prevents failure if not logged in to a registry |
| Keep v* tag trigger | Supports v1.10.6-nanopi-m6 pattern without modification |

## Deviations from Plan

None - plan executed exactly as written.

## Verification Results

| Check | Status |
|-------|--------|
| YAML syntax valid | PASS |
| Docker Hub login in sbc-rk3588 | PASS |
| Docker Hub overlay push in sbc-rk3588 | PASS |
| nanopi-m6 in board matrix | PASS |
| Docker Hub overlay reference in boards job | PASS |
| Existing boards unchanged | PASS |
| Tag trigger supports v*-nanopi-m6 | PASS |

## Commits

| Hash | Message |
|------|---------|
| 2d3148e | feat(04-01): add Docker Hub auth and overlay push to sbc-rk3588 job |
| 9e0b293 | feat(04-01): add nanopi-m6 to board matrix with Docker Hub config |
| 6408ea7 | chore(04-01): update cleanup to logout from both registries |

## User Setup Required

Before triggering the workflow, configure GitHub repository secrets and variables:

### GitHub Repository Variables
- `DOCKERHUB_USERNAME`: Set to `123417`

### GitHub Repository Secrets
- `DOCKERHUB_TOKEN`: Docker Hub access token with Read/Write permissions
  - Generate at: https://hub.docker.com/settings/security

## Next Phase Readiness

**Ready for:** Plan 04-02 (Tag push and hardware validation)

**Prerequisites satisfied:**
- CI workflow supports nanopi-m6 board
- Docker Hub authentication configured
- Overlay push to Docker Hub enabled
- Raw image generation configured

**Next steps:**
1. Configure DOCKERHUB_USERNAME variable and DOCKERHUB_TOKEN secret in GitHub
2. Push a version tag (e.g., v1.10.6-nanopi-m6) to trigger CI
3. Download raw image artifact
4. Flash to SD card and validate boot

---

*Plan: 04-01*
*Phase: 04-overlay-integration*
