---
phase: 01-environment-setup
plan: 03
subsystem: infra
tags: [flash, sd-card, hardware-verification, armbian, macos]

# Dependency graph
requires:
  - phase: 01-02
    provides: built artifacts to flash
provides:
  - SD card flash script with safety checks
  - Hardware verification workflow documentation
  - Confirmed NanoPi M6 hardware boots (Armbian baseline)
affects: [02-bootloader, 03-talos-integration]

# Tech tracking
tech-stack:
  added: []
  patterns: [flash-and-verify workflow, LED-based debugging]

key-files:
  created:
    - hack/flash.sh
    - docs/FLASH-WORKFLOW.md
  modified: []

key-decisions:
  - "Use rdisk raw device for faster flash writes on macOS"
  - "Refuse to write to internal disks as safety measure"
  - "Document LED-based verification for non-UART debugging"

patterns-established:
  - "Flash workflow: use hack/flash.sh for all SD card writes"
  - "Verification without UART: LED timing, HDMI output, network ping"

# Metrics
duration: ~15min
completed: 2026-02-02
---

# Phase 1 Plan 3: SD Card Flash Workflow Summary

**macOS flash script with safety checks and LED-based hardware verification workflow for NanoPi M6**

## Performance

- **Duration:** ~15 min
- **Started:** 2026-02-02
- **Completed:** 2026-02-02
- **Tasks:** 3
- **Files created:** 2

## Accomplishments

- Created safe SD card flash script with internal disk protection
- Documented non-UART verification methods (LED, HDMI, network)
- Confirmed NanoPi M6 hardware boots successfully with Armbian baseline
- Development environment ready for custom Talos image testing

## Task Commits

Each task was committed atomically:

1. **Task 1: Create SD card flash script** - `0ea3ae0` (feat)
2. **Task 2: Create verification workflow documentation** - `648b79e` (docs)
3. **Task 3: Verify hardware baseline with Armbian** - User verified (checkpoint)

**Plan metadata:** (pending)

## Files Created/Modified

- `hack/flash.sh` - SD card flash script with safety checks, confirmation prompts
- `docs/FLASH-WORKFLOW.md` - Complete flash and verification documentation

## Decisions Made

- **rdisk for performance:** Using raw device `/dev/rdiskN` instead of `/dev/diskN` for significantly faster write speeds on macOS
- **Safety checks:** Script refuses to write to internal disks, requires explicit "yes" confirmation
- **Non-UART debugging:** Documented LED timing patterns as primary early-boot diagnostic since UART not initially available

## Deviations from Plan

None - plan executed exactly as written.

## Issues Encountered

None - both tasks completed successfully and hardware verification passed.

## User Setup Required

None - no external service configuration required.

## Next Phase Readiness

- Hardware verified: NanoPi M6 boots Armbian successfully
- Flash workflow established: can test any built images
- Ready to proceed with Phase 2 (Bootloader Development)
- Blockers: None

---
*Phase: 01-environment-setup*
*Completed: 2026-02-02*
