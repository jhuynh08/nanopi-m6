---
phase: 02-bootloader
plan: 04
subsystem: bootloader
tags: [u-boot, rk3588s, armbian, defconfig, analysis, gap-closure]

# Dependency graph
requires:
  - phase: 02-03
    provides: Boot test results (FAILED) identifying need for M6-specific config
  - phase: 02-01
    provides: Initial pkg.yaml build configuration (used rock5a defconfig)
provides:
  - Armbian NanoPi M6 U-Boot configuration analysis
  - Key CONFIG differences documented in table format
  - Defconfig diff (rock5a vs nanopi-m6)
  - Implementation recommendation for Plan 02-05
affects: [02-05-u-boot-apply-config]

# Tech tracking
tech-stack:
  added: []
  patterns:
    - "Download defconfig from Armbian during build"
    - "Use mainline U-Boot v2025.10 with Armbian patches"

key-files:
  created:
    - docs/ARMBIAN-UBOOT-ANALYSIS.md
    - artifacts/u-boot/nanopi-m6/defconfig.diff
  modified: []

key-decisions:
  - "Use Armbian nanopi-m6-rk3588s_defconfig directly (Option A recommended)"
  - "Root cause of boot failure: wrong device tree (rock5a vs nanopi-m6)"
  - "Download configs from Armbian during build (vs storing locally)"
  - "Device tree rk3588s-nanopi-m6.dts also required from Armbian"

patterns-established:
  - "Config extraction: Use GitHub API via gh CLI to fetch Armbian configs"
  - "Gap closure: Research external working configs when mainline fails"

# Metrics
duration: 8min
completed: 2026-02-03
---

# Phase 02 Plan 04: Armbian U-Boot Analysis Summary

**Extracted NanoPi M6 defconfig and device tree from Armbian, identifying wrong device tree as root cause of boot failure**

## Performance

- **Duration:** 8 min
- **Started:** 2026-02-03T00:46:42Z
- **Completed:** 2026-02-03T00:55:00Z
- **Tasks:** 2/2 complete
- **Files created:** 2

## Accomplishments

- Located NanoPi M6 U-Boot configuration in Armbian build repository
- Extracted and analyzed complete defconfig with 90+ CONFIG options
- Documented key differences from rock5a-rk3588s_defconfig in table format
- Identified root cause: CONFIG_DEFAULT_DEVICE_TREE points to wrong DTB
- Provided actionable implementation plan for Plan 02-05

## Task Commits

Each task was committed atomically:

1. **Task 1: Locate NanoPi M6 Config in Armbian** - (research only, no commit)
2. **Task 2: Extract and Analyze M6 Configuration** - `0569171` (docs)

**Plan metadata:** (included in this commit)

## Files Created

- `docs/ARMBIAN-UBOOT-ANALYSIS.md` - Complete analysis with CONFIG diff table, DDR requirements, implementation options
- `artifacts/u-boot/nanopi-m6/defconfig.diff` - Side-by-side comparison of rock5a vs M6 defconfig

## Decisions Made

1. **Recommended Option A: Download from Armbian**
   - Fetch `nanopi-m6-rk3588s_defconfig` during build
   - Fetch `rk3588s-nanopi-m6.dts` during build
   - Always uses latest tested configuration
   - Rationale: Armbian actively maintains and tests M6 support

2. **Root Cause Identification**
   - `CONFIG_DEFAULT_DEVICE_TREE="rockchip/rk3588s-rock-5a"` was wrong
   - Must be `"rockchip/rk3588s-nanopi-m6"` for correct GPIO/pinmux/PMIC init
   - Wrong DTB = wrong LED pins, wrong PMIC sequence, no boot

3. **DDR Blob Recommendation**
   - Current v1.16 is acceptable (Armbian uses v1.18)
   - BL31 v1.45 is acceptable (Armbian uses v1.48)
   - Version difference not the root cause

## Deviations from Plan

None - plan executed exactly as written.

## Issues Encountered

None - Armbian repository structure was well-organized and configs were easy to locate.

## Key Findings

### CONFIG Differences (Critical)

| Option | rock5a | M6 | Impact |
|--------|--------|-----|--------|
| DEFAULT_DEVICE_TREE | rk3588s-rock-5a | rk3588s-nanopi-m6 | **Root cause** |
| TARGET_ROCK5A_RK3588 | y | (not set) | Wrong board init |
| TARGET_EVB_RK3588 | (not set) | y | Generic RK3588 target |
| DEFAULT_FDT_FILE | rock-5a.dtb | nanopi-m6.dtb | **Root cause** |

### Device Tree Key Sections

- LEDs: GPIO1_A4 (sys_led), GPIO1_A6 (user_led) - completely different from rock5a
- PMIC: RK806 on SPI2 with M6-specific regulator config
- Ethernet: RTL8211F on GMAC1 with tx_delay 0x42

## User Setup Required

None - no external service configuration required.

## Next Phase Readiness

**READY for Plan 02-05 (Apply Configuration)**

Phase 02-05 can now:
1. Download defconfig from Armbian URL in ANALYSIS.md
2. Download device tree from Armbian URL in ANALYSIS.md
3. Update pkg.yaml per recommended implementation
4. Rebuild U-Boot with correct M6 configuration
5. Re-test on hardware

**Blockers:** None

**Files to reference:**
- `docs/ARMBIAN-UBOOT-ANALYSIS.md` - Full implementation guide
- `artifacts/u-boot/nanopi-m6/defconfig.diff` - CONFIG comparison

---
*Phase: 02-bootloader*
*Completed: 2026-02-03*
*Status: SUCCESS - Gap closure plan complete*
