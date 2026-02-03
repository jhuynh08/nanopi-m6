# Phase 4: Overlay Integration - Context

**Gathered:** 2026-02-03
**Status:** Ready for planning

<domain>
## Phase Boundary

Produce a bootable Talos Linux raw image for the NanoPi M6 via CI/CD pipeline. The overlay artifacts (DTB, U-Boot, installer) exist from Phase 3. This phase packages them correctly, sets up CI/CD to build the image, publishes to registry, and validates boot to maintenance mode with hardware driver checks.

</domain>

<decisions>
## Implementation Decisions

### Image Generation Method
- CI/CD only — no local Talos imager setup
- Trigger: Git tag only (e.g., v1.9.0-nanopi-m6) — prevents accidental builds
- Full CI setup in this phase — configure workflows, secrets, produce real image
- CI workflow approach: Claude's discretion to adapt existing or create new

### Registry and Naming
- Registry: Docker Hub (not GHCR)
- Namespace: `123417`
- Image name: `talos-sbc-nanopi-m6`
- Full path: `docker.io/123417/talos-sbc-nanopi-m6`
- Tagging: Match Talos version (e.g., `v1.9.0-nanopi-m6`)

### Boot Validation Scope
- Success criteria: Boot to Talos maintenance mode
- Observation method: HDMI monitor
- No formal boot sequence documentation needed
- Debug approach: Fix and re-tag (clean version history)

### Deferred Hardware Checks (from Phase 3)
- Validate full Phase 3 list: Ethernet, NVMe, USB, board identification
- Blocking criteria: Only Ethernet and NVMe failures block Phase 4
- USB issues noted but don't block completion
- Ethernet test: Interface appears (`ip link` shows eth0)
- NVMe test: Device detected (`/dev/nvme0n1` exists)
- Full connectivity and mount tests deferred to Phase 5

### Claude's Discretion
- CI workflow: Adapt existing upstream workflows vs create new
- Exact workflow file structure and job names
- Secret naming conventions for Docker Hub credentials
- Order of validation checks during boot testing

</decisions>

<specifics>
## Specific Ideas

- Use Talos version in image tag to show compatibility (v1.9.0-nanopi-m6 pattern)
- Docker Hub credentials need to be set up as GitHub secrets
- Existing Phase 3 outputs should be verified to include all needed artifacts

</specifics>

<deferred>
## Deferred Ideas

None — discussion stayed within phase scope

</deferred>

---

*Phase: 04-overlay-integration*
*Context gathered: 2026-02-03*
