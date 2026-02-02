---
phase: 02-bootloader
plan: 02
subsystem: infra
tags: [maskrom, recovery, rkdeveloptool, rk3588, debugging]

# Dependency graph
requires:
  - phase: 01-environment-setup
    provides: Armbian baseline verification, flash workflow
provides:
  - MaskROM recovery procedure for RK3588 boards
  - Boot test checklist template for non-UART debugging
  - LED-based verification methodology
affects: [02-03, 02-04, any future bootloader debugging]

# Tech tracking
tech-stack:
  added: [rkdeveloptool]
  patterns: [tiered-iteration-debugging, led-based-verification]

key-files:
  created:
    - docs/MASKROM-RECOVERY.md
    - docs/BOOT-TEST-CHECKLIST.md
  modified: []

key-decisions:
  - "Focus on LED activity as primary boot indicator (U-Boot has no HDMI)"
  - "3-tier iteration strategy for systematic debugging"
  - "Cross-reference recovery docs for failure handling"

patterns-established:
  - "Observation timeline template for non-UART debugging"
  - "MaskROM as last-resort recovery, always accessible"

# Metrics
duration: 2min
completed: 2026-02-02
---

# Phase 2 Plan 2: Recovery and Iteration Documentation Summary

**MaskROM recovery procedure and boot test checklist for systematic non-UART debugging on NanoPi M6**

## Performance

- **Duration:** 2 min
- **Started:** 2026-02-02T23:06:18Z
- **Completed:** 2026-02-02T23:08:21Z
- **Tasks:** 2
- **Files created:** 2

## Accomplishments

- Complete MaskROM recovery documentation with macOS-specific rkdeveloptool installation
- Boot test checklist template with tiered iteration strategy (Tier 1/2/3)
- Observation timeline for tracking LED/HDMI/network indicators at 0-10s, 10-30s, 30-60s, 60-120s, 120s+ intervals
- Documented critical constraint: U-Boot has no VOP2 driver, no HDMI output during bootloader stage

## Task Commits

Each task was committed atomically:

1. **Task 1: Create MaskROM Recovery Documentation** - `e840557` (docs)
2. **Task 2: Create Boot Test Checklist Template** - `72ead62` (docs)

## Files Created

- `docs/MASKROM-RECOVERY.md` - Complete MaskROM mode entry and recovery operations for NanoPi M6
- `docs/BOOT-TEST-CHECKLIST.md` - Template for tracking boot attempts with observation timeline

## Decisions Made

- **LED-based verification as primary method:** U-Boot has no HDMI driver for RK3588, so boot stage verification relies on LED activity patterns and eventual kernel boot confirmation
- **Tiered iteration strategy:** Tier 1 (3 quick attempts), Tier 2 (2 config analysis attempts), Tier 3 (decision point to park or acquire UART)
- **Bidirectional cross-references:** Both documents link to each other for easy navigation during debugging sessions

## Deviations from Plan

None - plan executed exactly as written.

## Issues Encountered

None

## User Setup Required

None - documentation-only plan. rkdeveloptool installation documented for when MaskROM recovery is needed.

## Next Phase Readiness

- Recovery procedures documented and ready for bootloader bring-up testing
- Boot test checklist ready for systematic iteration during U-Boot configuration
- Critical constraint documented: expect no HDMI output until Linux kernel boots (6.15+)

---
*Phase: 02-bootloader*
*Completed: 2026-02-02*
