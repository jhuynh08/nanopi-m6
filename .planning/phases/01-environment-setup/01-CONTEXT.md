# Phase 1: Environment Setup - Context

**Gathered:** 2026-02-02
**Status:** Ready for planning

<domain>
## Phase Boundary

Establish the development environment, build pipeline, and testing workflow for NanoPi M6 hardware bring-up. This phase delivers the ability to build, flash, and verify artifacts — not the artifacts themselves.

</domain>

<decisions>
## Implementation Decisions

### Build workflow
- Local builds on Apple Silicon Mac + GitHub Actions for CI
- Push container images to Docker Hub
- Local builds for iteration, CI for releases and validation

### Flash & test cycle
- Single NanoPi M6 board, sequential testing only
- SD card swap workflow (unplug/replug each time)
- Simple recovery: if boot fails, reflash SD card
- No MaskROM/USB flashing in this phase (SD-only)

### Debug setup
- **No UART adapter initially** — higher risk for Phase 2
- Verify boot progress via: HDMI output, network ping, LED behavior
- Built-in 2.1" LCD available once display driver loads
- Fallback: order UART adapter if stuck in bootloader work

### Artifact organization
- Build outputs to `_out/` directory (gitignored)
- Milestone releases via GitHub Releases with downloadable assets
- Version scheme: Talos-based (e.g., v1.9.0-nanopi-m6.1)

### Claude's Discretion
- Build caching strategy (Docker layers vs persistent volume)
- Which generated files to commit (source configs vs binaries)
- Whether to document failed attempts or only successes

</decisions>

<specifics>
## Specific Ideas

- User has FriendlyElec 2.1" LCD touchscreen attached — can use for visual feedback once kernel boots
- Workflow assumes physical SD card swapping, no always-connected reader

</specifics>

<deferred>
## Deferred Ideas

None — discussion stayed within phase scope

</deferred>

---

*Phase: 01-environment-setup*
*Context gathered: 2026-02-02*
