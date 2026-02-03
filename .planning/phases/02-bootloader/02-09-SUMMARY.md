---
phase: 02-bootloader
plan: 09
subsystem: bootloader
tags: [u-boot, rockchip, rk3588s, vendor-uboot, friendlyelec, idbloader]

# Dependency graph
requires:
  - phase: 02-08
    provides: "Root cause identification (vendor U-Boot required)"
  - phase: 01-03
    provides: "Hardware verification and flash workflow"
provides:
  - "Vendor U-Boot build integration in Talos bldr"
  - "Pre-extracted FriendlyELEC bootloader binaries"
  - "Boot test Attempt #6 SUCCESS (full boot to login screen)"
  - "Phase 2 goal achieved: NanoPi M6 boots to U-Boot"
affects:
  - phase-03-device-tree-kernel
  - phase-04-talos-integration

# Tech tracking
tech-stack:
  added:
    - "FriendlyELEC uboot-rockchip v2017.09"
  patterns:
    - "Vendor U-Boot integration via pre-extracted binaries"
    - "idbloader.img + uboot.img flash layout (sectors 64 and 16384)"

key-files:
  created:
    - "artifacts/u-boot/prepare-vendor/pkg.yaml"
  modified:
    - "Pkgfile"
    - "artifacts/u-boot/nanopi-m6/pkg.yaml"
    - "docs/BOOT-TEST-CHECKLIST.md"

key-decisions:
  - "Use pre-extracted vendor binaries (GCC 6.x incompatibility with bldr toolchain)"
  - "Flash layout: idbloader.img at sector 64, uboot.img at sector 16384"
  - "Vendor U-Boot v2017.09 required for NanoPi M6 (mainline format incompatible)"

patterns-established:
  - "Vendor U-Boot boards may require pre-extracted binaries"
  - "idbloader + uboot.img format for Rockchip RK3588S vendor boot chain"

# Metrics
duration: ~45min (including checkpoint and continuation)
completed: 2026-02-03
---

# Phase 02 Plan 09: Vendor U-Boot Integration Summary

**FriendlyELEC vendor U-Boot integrated into Talos build system with pre-extracted binaries, validated by boot test Attempt #6 (full boot to login screen)**

## Performance

- **Duration:** ~45 min (including checkpoint and continuation)
- **Started:** 2026-02-03T04:30:00Z
- **Completed:** 2026-02-03T05:40:00Z
- **Tasks:** 3 (2 automated + 1 human verification)
- **Files modified:** 4

## Accomplishments

- Integrated FriendlyELEC vendor U-Boot source into Pkgfile
- Created prepare-vendor pkg.yaml for vendor U-Boot source preparation
- Updated nanopi-m6 pkg.yaml to use pre-extracted vendor binaries
- Boot test Attempt #6: SUCCESS - full boot to login screen
- **Phase 2 goal achieved: NanoPi M6 boots to U-Boot (and beyond!)**

## Task Commits

Each task was committed atomically:

1. **Task 1: Add FriendlyELEC U-Boot source to Pkgfile** - `ccc8c4b` (feat)
2. **Task 2: Update nanopi-m6 pkg.yaml for vendor binaries** - `40bb33f` (feat)
3. **Task 3: Hardware boot test Attempt #6** - Human verification (SUCCESS)

**Plan metadata:** Committed with this summary

## Files Created/Modified

- `Pkgfile` - Added friendlyelec_uboot_ref variable and pre-extracted binary path
- `artifacts/u-boot/prepare-vendor/pkg.yaml` - Created vendor U-Boot source preparation stage
- `artifacts/u-boot/nanopi-m6/pkg.yaml` - Updated to copy pre-extracted binaries
- `docs/BOOT-TEST-CHECKLIST.md` - Added Attempt #6 SUCCESS entry

## Decisions Made

1. **Use pre-extracted binaries instead of building from source**
   - FriendlyELEC uboot-rockchip requires GCC 6.x
   - Talos bldr uses GCC 14.x toolchain (incompatible)
   - Pre-extracted binaries from FriendlyELEC image work correctly
   - Long-term: May need to investigate Rockchip EDK2 or cross-compilation

2. **Flash layout validated**
   - idbloader.img: sector 64 (0x8000 = 32KB offset)
   - uboot.img: sector 16384 (0x800000 = 8MB offset)
   - This matches FriendlyELEC's expected layout

## Deviations from Plan

### Build Approach Deviation

**[Rule 3 - Blocking] Switched to pre-extracted binaries**
- **Found during:** Task 2 (Build attempt)
- **Issue:** FriendlyELEC uboot-rockchip (v2017.09) requires GCC 6.x toolchain
  - Talos bldr uses GCC 14.x which produces incompatible code
  - Build fails with syntax/compilation errors for older codebase
- **Fix:** Extract working binaries from FriendlyELEC Ubuntu image
  - Source: rk3588-sd-ubuntu-noble-minimal-6.1-arm64-20251222.img
  - Files extracted: idbloader.img (2.0MB), uboot.img (4.0MB)
  - Stored in artifacts/u-boot/friendlyelec-vendor/
- **Verification:** Boot test Attempt #6 - SUCCESS (full boot to login)
- **Committed in:** 40bb33f

---

**Total deviations:** 1 (blocking - toolchain incompatibility)
**Impact on plan:** Pre-extracted binaries achieve same goal. No scope creep. Long-term solution may be needed for customization.

## Issues Encountered

- **GCC version incompatibility:** FriendlyELEC U-Boot v2017.09 cannot build with modern GCC 14.x
  - Resolution: Use pre-extracted binaries from working FriendlyELEC image
  - This is a known issue with older Rockchip vendor U-Boot code

## User Setup Required

None - no external service configuration required.

## Boot Test Results

### Attempt #6 Summary

| Setting | Value |
|---------|-------|
| Defconfig | nanopi6_defconfig (via pre-extracted binaries) |
| U-Boot source | FriendlyELEC vendor v2017.09 |
| Format | idbloader.img + uboot.img |
| Result | **SUCCESS - Full boot to login screen** |

### Boot Attempt History (Phase 2)

| # | Bootloader | Format | Result |
|---|------------|--------|--------|
| 1-4 | Mainline v2025.10 | u-boot-rockchip.bin | FAIL |
| 5 | FriendlyELEC pre-built | MiniLoaderAll | SUCCESS |
| **6** | **Talos-integrated vendor** | **idbloader + uboot.img** | **SUCCESS** |

## Phase 2 Completion

**PHASE 2 GOAL ACHIEVED**

- **Goal:** NanoPi M6 boots to U-Boot
- **Actual:** NanoPi M6 boots to login screen (exceeds goal)
- **Method:** Vendor U-Boot (FriendlyELEC v2017.09) with pre-extracted binaries

### Summary of Phase 2

| Plan | Focus | Status | Key Outcome |
|------|-------|--------|-------------|
| 02-01 | U-Boot build setup | Complete | Initial build configuration |
| 02-02 | Boot recovery docs | Complete | LED verification strategy |
| 02-03 | Boot test #1 | Complete | FAIL - wrong defconfig |
| 02-04 | Armbian analysis | Complete | Extracted M6 configs |
| 02-05 | DTS patching | Complete | FAIL - DTS not root cause |
| 02-06 | Mainline switch | Complete | FAIL - version not root cause |
| 02-07 | Blob update | Complete | FAIL - blobs not root cause |
| 02-08 | Vendor diagnostic | Complete | SUCCESS - root cause found |
| **02-09** | **Vendor integration** | **Complete** | **SUCCESS - Phase 2 done** |

### Root Cause Confirmed

Mainline U-Boot u-boot-rockchip.bin format is incompatible with NanoPi M6.
The board requires vendor U-Boot with MiniLoaderAll format (idbloader.img + uboot.img).

## Next Phase Readiness

**Ready for Phase 3: Device Tree & Kernel**

- Bootloader is working and integrated into Talos build system
- Full boot chain validated to login screen
- Kernel boot confirmed working

**Remaining work for Talos image:**
- Phase 3: Device tree and kernel configuration for Talos
- Phase 4: Talos integration (SBC profile, omnictl)
- Phase 5: CI/CD automation
- Phase 6: Documentation and release

---
*Phase: 02-bootloader*
*Plan: 09 (Gap Closure)*
*Completed: 2026-02-03*
