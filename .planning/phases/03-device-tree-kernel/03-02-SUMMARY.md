---
phase: 03-device-tree-kernel
plan: 02
subsystem: installer
tags: [go, rk3588, u-boot, bootloader, talos-installer]

# Dependency graph
requires:
  - phase: 02-bootloader
    provides: Vendor U-Boot binaries (idbloader.img + uboot.img) at correct flash offsets
  - phase: 03-01
    provides: Pre-extracted vendor DTB in artifacts
provides:
  - NanoPi M6 board recognition in RK3588 installer
  - Vendor U-Boot flash layout support (two files at sectors 64 and 16384)
  - DTB path resolution via rk3588s chipset
affects: [phase-04, talos-build, image-generation]

# Tech tracking
tech-stack:
  added: []
  patterns:
    - Board-specific conditional in Install function for different bootloader formats

key-files:
  created: []
  modified:
    - installers/rk3588/src/main.go
    - go.work

key-decisions:
  - "Use board name conditional rather than separate install function for vendor U-Boot"
  - "idbloader.img reuses ubootOffset constant (sector 64), new ubootImgOffset for uboot.img"

patterns-established:
  - "Board-specific bootloader format: Check board name before choosing flash layout"

# Metrics
duration: 1min
completed: 2026-02-03
---

# Phase 3 Plan 2: Update RK3588 Installer Summary

**RK3588 installer now recognizes NanoPi M6 and flashes vendor U-Boot (idbloader.img + uboot.img) at correct sectors**

## Performance

- **Duration:** 1 min 20 sec
- **Started:** 2026-02-03T06:34:47Z
- **Completed:** 2026-02-03T06:36:07Z
- **Tasks:** 3 (2 committed, 1 was verification)
- **Files modified:** 2

## Accomplishments
- Added nanopi-m6 case to ChipsetName function returning "rk3588s"
- Implemented board-specific Install logic for vendor U-Boot two-file format
- idbloader.img written at sector 64 (32KB offset)
- uboot.img written at sector 16384 (8MB offset)
- Maintained backward compatibility for Rock 5A/5B boards

## Task Commits

Each task was committed atomically:

1. **Task 1: Add NanoPi M6 to ChipsetName function** - `d0549ee` (feat)
2. **Task 2: Update Install function for vendor U-Boot flash layout** - `908d081` (feat)
3. **Task 3: Verify installer compiles** - (verified as part of Task 2, no separate commit)

## Files Created/Modified
- `installers/rk3588/src/main.go` - Added nanopi-m6 board support with vendor U-Boot flash layout
- `go.work` - Updated Go version from 1.22.3 to 1.24.0 to match module requirements

## Decisions Made
- Used board name conditional (`if options.ExtraOptions.Board == "nanopi-m6"`) rather than a separate function for vendor U-Boot handling
- Reused `ubootOffset` constant for idbloader.img since it goes to the same sector 64
- Added new `ubootImgOffset` constant for uboot.img at sector 16384

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 3 - Blocking] Updated go.work version to match module requirements**
- **Found during:** Task 3 (build verification)
- **Issue:** go.work specified go 1.22.3 but installers/rk3588/src/go.mod requires go >= 1.24.0
- **Fix:** Updated go.work from `go 1.22.3` to `go 1.24.0`
- **Files modified:** go.work
- **Verification:** `go build -o /dev/null .` succeeded
- **Committed in:** 908d081 (part of Task 2 commit)

---

**Total deviations:** 1 auto-fixed (1 blocking)
**Impact on plan:** Minor version alignment fix. No scope creep.

## Issues Encountered
None beyond the go.work version fix.

## User Setup Required
None - no external service configuration required.

## Next Phase Readiness
- Installer now handles NanoPi M6 vendor U-Boot format
- Ready for full image build testing in Plan 03-03
- DTB path will resolve to `rockchip/rk3588s-nanopi-m6.dtb`

---
*Phase: 03-device-tree-kernel*
*Completed: 2026-02-03*
