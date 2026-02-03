# Phase 4: Overlay Integration - Research

**Researched:** 2026-02-03
**Domain:** Talos Linux overlay container CI/CD, Docker Hub publishing, ARM64 image generation
**Confidence:** HIGH

## Summary

This phase focuses on producing a bootable Talos Linux raw image for the NanoPi M6 via CI/CD pipeline. The overlay artifacts (DTB, U-Boot, installer) exist from Phase 3 at committed locations. This phase packages them into an overlay container image, sets up CI/CD to build on tag push, publishes to Docker Hub, generates a raw image using the Talos imager, and validates boot to maintenance mode.

The existing CI workflow (`.github/workflows/ci.yaml`) already implements most of the required functionality for GHCR. The main changes are: (1) add Docker Hub authentication, (2) add the NanoPi M6 board to the matrix, (3) configure tag-only triggers per user decision. The Talos imager generates raw images via `docker run` with `--overlay-image`, `--overlay-name`, and `--overlay-option` flags.

**Primary recommendation:** Adapt the existing CI workflow to push overlay images to Docker Hub (`docker.io/123417/talos-sbc-nanopi-m6`) on version tag push only, then use the imager to generate and release raw images.

## Standard Stack

The established tools for Talos overlay CI/CD:

### Core
| Tool | Version | Purpose | Why Standard |
|------|---------|---------|--------------|
| docker/login-action | v3 | Docker Hub authentication | Official Docker action, supports multiple registries |
| docker/setup-qemu-action | v3 | ARM64 emulation support | Required for arm64 builds on x86 runners |
| docker/setup-buildx-action | v3 | Multi-platform builds | Enables cross-compilation via BuildKit |
| siderolabs/imager | v1.10.6+ | Raw image generation | Official Talos image generator |
| crane | latest | Container image manipulation | Push tarball images to registry |

### Supporting
| Tool | Version | Purpose | When to Use |
|------|---------|---------|-------------|
| actions/checkout | v4 | Repository checkout | Every workflow |
| actions/upload-artifact | v4 | Build artifact storage | Release raw images |
| crazy-max/ghaction-github-release | v2 | GitHub releases | Attach raw images to releases |

### Alternatives Considered
| Instead of | Could Use | Tradeoff |
|------------|-----------|----------|
| QEMU emulation | Native ARM64 runners | `ubuntu-22.04-arm` is faster but only for public repos |
| crane push | docker push | crane handles tarballs directly, no daemon needed |
| Self-hosted runners | GitHub-hosted | Existing workflow uses self-hosted for build speed |

**Docker Hub Authentication:**
```yaml
- name: Login to Docker Hub
  uses: docker/login-action@v3
  with:
    username: ${{ vars.DOCKERHUB_USERNAME }}
    password: ${{ secrets.DOCKERHUB_TOKEN }}
```

## Architecture Patterns

### Overlay Container Structure
```
/
├── artifacts/
│   └── arm64/
│       ├── dtb/rockchip/rk3588s-nanopi-m6.dtb
│       └── u-boot/nanopi-m6/
│           ├── idbloader.img
│           └── uboot.img
├── installers/
│   └── rk3588                    # Statically linked Go binary
└── profiles/
    └── board-rk3588.yaml         # Disk image configuration
```

### CI Workflow Architecture

**Existing Pattern (GHCR):**
```
[Tag Push] -> [Build Overlay] -> [Push to GHCR] -> [Build Imager] -> [Generate Raw Image] -> [Release]
```

**Required Pattern (Docker Hub):**
```
[Tag Push v*.*.* pattern] -> [Build Overlay] -> [Push to Docker Hub] -> [Generate Raw Image] -> [Release]
```

### Tag-Only Trigger Pattern
```yaml
on:
  push:
    tags:
      - 'v*-nanopi-m6'     # e.g., v1.9.0-nanopi-m6
```

### Imager Command Pattern
```bash
docker run --rm -t -v ./_out:/out -v /dev:/dev --privileged \
  ghcr.io/siderolabs/imager:v1.10.6 metal \
  --arch arm64 \
  --overlay-image=docker.io/123417/talos-sbc-nanopi-m6:v1.9.0-nanopi-m6 \
  --overlay-name=rk3588 \
  --overlay-option="board=nanopi-m6" \
  --overlay-option="chipset=rk3588s" \
  --base-installer-image=ghcr.io/siderolabs/installer:v1.10.6
```

Output: `_out/metal-arm64.raw.xz`

### Anti-Patterns to Avoid
- **Building raw images locally**: CI/CD is the decided path (user decision)
- **Using GITHUB_TOKEN for tag push triggers**: Tags pushed with GITHUB_TOKEN won't trigger workflows
- **Hardcoding Talos version**: Use env variable that can be updated (TALOS_VERSION)
- **Skipping privileged mode for imager**: Device access required for raw image generation

## Don't Hand-Roll

Problems that have existing solutions in the Talos ecosystem:

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| Container registry auth | Custom auth scripts | docker/login-action@v3 | Handles logout, token refresh |
| Multi-arch builds | Manual docker build | docker/build-push-action@v6 | BuildKit, caching, attestations |
| Raw image generation | Custom dd scripts | siderolabs/imager container | Handles partitioning, overlay installation |
| Image tagging | Manual tag logic | `git describe --tags` | Consistent with project pattern |
| Release artifacts | Custom upload scripts | ghaction-github-release | Handles checksums, draft releases |

**Key insight:** The existing CI workflow already solves 80% of the problem. The adaptations are additive, not replacement.

## Common Pitfalls

### Pitfall 1: GITHUB_TOKEN Cannot Trigger Workflows
**What goes wrong:** Creating a tag with GITHUB_TOKEN does not trigger tag-based workflows
**Why it happens:** Security design - prevents workflow loops
**How to avoid:** Tags must be pushed manually or with a PAT (not GITHUB_TOKEN)
**Warning signs:** Workflow never runs after automated tag push

### Pitfall 2: Docker Hub Rate Limits for Anonymous Pulls
**What goes wrong:** Builds fail with "toomanyrequests" error
**Why it happens:** Anonymous pulls limited to 100/6h per IP
**How to avoid:** Always authenticate, even for pulls; use `--scope` parameter
**Warning signs:** Intermittent failures in CI, especially on shared runners

### Pitfall 3: Wrong Overlay Option Names
**What goes wrong:** Imager fails with "unknown board" or similar
**Why it happens:** `--overlay-option` keys must match installer's expected options
**How to avoid:** Use exact keys: `board=nanopi-m6`, `chipset=rk3588s`
**Warning signs:** Imager produces default/wrong DTB, installer errors

### Pitfall 4: Missing Privileged Mode for Imager
**What goes wrong:** Raw image generation fails with permission errors
**Why it happens:** Imager needs device access for disk image creation
**How to avoid:** Include `--privileged -v /dev:/dev` in docker run
**Warning signs:** "permission denied" errors, incomplete images

### Pitfall 5: Talos Version Mismatch
**What goes wrong:** Overlay incompatible with imager
**Why it happens:** Overlay API changes between Talos versions
**How to avoid:** Pin TALOS_VERSION env var, update overlay when Talos updates
**Warning signs:** "unknown profile", interface errors during image generation

### Pitfall 6: Self-Hosted Runner Availability
**What goes wrong:** Jobs queue indefinitely
**Why it happens:** Existing workflow requires `[self-hosted, pkgs]` labels
**How to avoid:** Either set up self-hosted runner or adapt to GitHub-hosted with QEMU
**Warning signs:** Jobs stay in "queued" state, timeout after max wait

## Code Examples

Verified patterns from official sources:

### Docker Hub Login for Workflow
```yaml
# Source: https://github.com/docker/login-action
- name: Login to Docker Hub
  if: github.event_name != 'pull_request'
  uses: docker/login-action@v3
  with:
    username: ${{ vars.DOCKERHUB_USERNAME }}
    password: ${{ secrets.DOCKERHUB_TOKEN }}
```

### Tag-Only Workflow Trigger
```yaml
# Source: https://docs.github.com/actions/using-workflows/workflow-syntax-for-github-actions
on:
  push:
    tags:
      - 'v*-nanopi-m6'
```

### Multi-Platform Build with Docker Hub Push
```yaml
# Source: https://docs.docker.com/build/ci/github-actions/multi-platform/
- name: Set up QEMU
  uses: docker/setup-qemu-action@v3

- name: Set up Docker Buildx
  uses: docker/setup-buildx-action@v3

- name: Build and push
  uses: docker/build-push-action@v6
  with:
    platforms: linux/arm64
    push: ${{ github.event_name != 'pull_request' }}
    tags: docker.io/123417/talos-sbc-nanopi-m6:${{ github.ref_name }}
```

### Talos Imager Raw Image Generation
```yaml
# Source: https://docs.siderolabs.com/talos/v1.9/platform-specific-installations/boot-assets/
- name: Generate raw image
  run: |
    docker run --rm -t -v $PWD/_out:/out -v /dev:/dev --privileged \
      ghcr.io/siderolabs/imager:${{ env.TALOS_VERSION }} metal \
      --arch arm64 \
      --overlay-image=docker.io/123417/talos-sbc-nanopi-m6:${{ github.ref_name }} \
      --overlay-name=rk3588 \
      --overlay-option="board=nanopi-m6" \
      --overlay-option="chipset=rk3588s" \
      --base-installer-image=ghcr.io/siderolabs/installer:${{ env.TALOS_VERSION }}
```

### Crane Push Tarball to Registry
```yaml
# Source: Existing ci.yaml pattern
- name: Push image
  run: |
    echo "${{ secrets.DOCKERHUB_TOKEN }}" | crane auth login docker.io --username "${{ vars.DOCKERHUB_USERNAME }}" --password-stdin
    crane push _out/overlay.tar docker.io/123417/talos-sbc-nanopi-m6:${{ github.ref_name }}
```

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| Board-specific imager flags | `--overlay-image` + `--overlay-option` | Talos v1.7+ | Unified overlay system |
| `--board` flag | `--overlay-name` + options | Talos v1.7 | Deprecated board flag |
| Manual DTB/U-Boot copy | Overlay container artifacts | Talos v1.7+ | Automated installation |
| Local raw image builds | CI/CD with imager | Current best practice | Reproducibility |

**Current Talos Version:** v1.10.6 (per existing ci.yaml)

**Deprecated/outdated:**
- `--board` flag: Replaced by `--overlay-name` + `--overlay-option`
- Board names as image platform: Use `metal` with overlays instead

## Hardware Validation Commands

Commands for validating Phase 3 deferred hardware checks in maintenance mode:

### Ethernet Validation
```bash
# Interface should appear (blocking if fails)
ip link | grep eth0
```

### NVMe Validation
```bash
# Device should be detected (blocking if fails)
ls -la /dev/nvme0n1
```

### USB Validation
```bash
# Enumerate devices (non-blocking if fails)
lsusb
```

### Board Identification
```bash
# Check device tree
cat /proc/device-tree/compatible
```

## Open Questions

Things that couldn't be fully resolved:

1. **Self-hosted runner availability**
   - What we know: Existing workflow uses `[self-hosted, pkgs]` runner labels
   - What's unclear: Whether self-hosted runners are available for this fork
   - Recommendation: Plan for GitHub-hosted fallback with QEMU if needed

2. **Imager rootless support**
   - What we know: Recent Talos docs mention imager supports rootless
   - What's unclear: Whether this eliminates need for `--privileged`
   - Recommendation: Keep `--privileged` for compatibility, test without later

3. **Realtek firmware extension requirement**
   - What we know: Existing workflow includes REALKTEK_FIRMWARE_EXTENSION_IMAGE
   - What's unclear: Whether NanoPi M6 requires this extension
   - Recommendation: Include initially, can remove if testing shows not needed

## Sources

### Primary (HIGH confidence)
- [docker/login-action](https://github.com/docker/login-action) - Docker Hub authentication patterns
- [docker/build-push-action](https://github.com/docker/build-push-action) - Multi-platform build patterns
- [Talos Overlays Documentation](https://docs.siderolabs.com/talos/v1.8/build-and-extend-talos/custom-images-and-development/overlays) - Overlay structure
- [siderolabs/overlays](https://github.com/siderolabs/overlays) - Official overlay catalog
- [siderolabs/sbc-rockchip](https://github.com/siderolabs/sbc-rockchip) - Reference RK3588 overlay implementation
- Existing `.github/workflows/ci.yaml` - Working GHCR workflow to adapt

### Secondary (MEDIUM confidence)
- [GitHub Actions Workflow Syntax](https://docs.github.com/actions/using-workflows/workflow-syntax-for-github-actions) - Tag trigger patterns
- [Docker Multi-Platform Builds](https://docs.docker.com/build/ci/github-actions/multi-platform/) - QEMU/Buildx setup
- [GitHub ARM64 Runners](https://github.blog/changelog/2025-01-16-linux-arm64-hosted-runners-now-available-for-free-in-public-repositories-public-preview/) - Native ARM64 option

### Tertiary (LOW confidence)
- Web search results on QEMU performance - General guidance, verify with testing

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH - Official Docker actions, documented Talos patterns
- Architecture: HIGH - Based on existing working ci.yaml and official docs
- Pitfalls: HIGH - Documented issues, verified against official sources

**Research date:** 2026-02-03
**Valid until:** 30 days (stable tooling, watch for Talos version updates)
