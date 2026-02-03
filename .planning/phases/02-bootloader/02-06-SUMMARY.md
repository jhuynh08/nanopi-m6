---
phase: 02-bootloader
plan: 06
subsystem: bootloader
tags: [u-boot, rk3588s, mainline, v2025.10, boot-failure, rkbin-blobs]

# Dependency graph
requires:
  - phase: 02-05
    provides: Root cause analysis - Collabora fork lacks M6 support
  - phase: 02-04
    provides: Armbian U-Boot configuration analysis
provides:
  - Mainline U-Boot v2025.10 build configuration
  - Native nanopi-m6-rk3588s_defconfig usage
  - Eliminated U-Boot source/version as root cause
  - Boot test #3 results (FAILED - identical to #1 and #2)
affects: [02-07-blob-versions, 02-08-emmc-boot, phase-3-blocked]

# Tech tracking
tech-stack:
  added:
    - "Mainline U-Boot v2025.10 (replaced Collabora fork)"
  patterns:
    - "GitHub releases tarball for U-Boot source"
    - "Native defconfig without patching (when available)"

key-files:
  created: []
  modified:
    - Pkgfile
    - artifacts/u-boot/prepare/pkg.yaml
    - artifacts/u-boot/nanopi-m6/pkg.yaml
    - Makefile
    - docs/BOOT-TEST-CHECKLIST.md

key-decisions:
  - "Switched from Collabora U-Boot fork (v2023.07-rc4) to mainline v2025.10"
  - "U-Boot source/version is NOT the root cause of boot failure"
  - "DDR/BL31 blob version mismatch is primary suspect (v1.16/v1.45 vs Armbian v1.18/v1.48)"
  - "SD card boot path may be the issue (M6 may require eMMC boot)"

patterns-established:
  - "Use GitHub releases URL for mainline U-Boot downloads"
  - "Prefer native defconfig over patching when available in source"

# Metrics
duration: ~20min (continuation after checkpoint)
completed: 2026-02-03
---

# Phase 02 Plan 06: Mainline U-Boot v2025.10 Summary

**Build system switched to mainline U-Boot v2025.10 with native M6 defconfig - boot test FAILED with identical symptoms to previous attempts**

## Performance

- **Duration:** ~20 min (post-checkpoint continuation)
- **Full plan duration:** ~1 hour (including Tasks 1-3 build and flash)
- **Started:** 2026-02-03T02:00:00Z (continuation at 02:36:04Z)
- **Completed:** 2026-02-03T02:50:00Z
- **Tasks:** 4/4 complete (Task 4 verification FAILED)
- **Files modified:** 5

## Accomplishments

- Successfully switched U-Boot source from Collabora fork to mainline v2025.10
- Build system now uses native `nanopi-m6-rk3588s_defconfig` (no patching needed)
- Eliminated U-Boot source/version as potential root cause
- Documented boot failure analysis with narrowed-down suspects
- Established cleaner build configuration (no sed patching)

## Task Commits

1. **Task 1: Switch U-Boot source to mainline v2025.10** - `4fc47ed`
   - Updated Pkgfile with v2025.10 reference and checksums
   - Changed prepare/pkg.yaml from Collabora GitLab to GitHub releases

2. **Task 2: Update NanoPi M6 build configuration** - `3f9a9f9`
   - Simplified nanopi-m6/pkg.yaml to use native defconfig
   - Removed all rock5a defconfig patching

3. **Task 3: Build U-Boot and flash to SD card** - `c088ccf`
   - Fixed build issues (Makefile dependency, pkg.yaml adjust)
   - Successfully built and flashed u-boot-rockchip.bin

4. **Task 4: Hardware boot verification** - (this commit)
   - FAILED - identical symptoms to Attempts #1 and #2
   - Updated BOOT-TEST-CHECKLIST.md with Attempt #3

## Boot Test Results

### Attempt #3 Configuration

| Setting | Value |
|---------|-------|
| U-Boot source | **Mainline v2025.10** (same as Armbian) |
| Defconfig | nanopi-m6-rk3588s_defconfig (native) |
| Device tree | rk3588s-nanopi-m6.dts (native) |
| DDR blob | v1.16 |
| BL31 | v1.45 |

### Observations

| Indicator | Result |
|-----------|--------|
| LED activity | No (SYS solid, LED1 off, no blink) |
| HDMI output | No signal |
| Network | No DHCP lease |
| Boot stage reached | TPL/DDR (failed before SPL) |

### Conclusion

**IDENTICAL to Attempts #1 and #2** - No improvement from switching to mainline U-Boot.

## Root Cause Analysis

### Eliminated Causes

| Hypothesis | Status | Evidence |
|------------|--------|----------|
| Wrong device tree | ELIMINATED | Attempt #2 with correct DTS still failed |
| Wrong defconfig base (rock5a) | ELIMINATED | Attempt #3 uses native M6 defconfig |
| Collabora fork lacks M6 support | ELIMINATED | Mainline v2025.10 has full support, still fails |
| U-Boot version too old | ELIMINATED | Using same version as Armbian |

### Remaining Suspects (Ordered by Likelihood)

1. **DDR blob version mismatch** (HIGH PROBABILITY)
   - We use: v1.16
   - Armbian uses: v1.18
   - DDR training may fail with older blob on M6's specific LPDDR5 chip
   - DDR init is exactly where boot appears to stall

2. **BL31 blob version mismatch** (MEDIUM PROBABILITY)
   - We use: v1.45
   - Armbian uses: v1.48
   - ATF may have M6-specific fixes in newer version

3. **SD card boot not supported** (MEDIUM PROBABILITY)
   - M6 may only boot from eMMC by default
   - Armbian may work because it was tested on eMMC
   - SD boot may require specific configuration

4. **rkbin repository version** (LOW PROBABILITY)
   - Different rkbin commit may have interdependent blob versions
   - But unlikely to be sole cause

5. **Hardware issue** (LOW PROBABILITY)
   - SD card slot hardware different from eMMC path
   - Unlikely since Armbian boots from same SD card

## Files Created/Modified

- `Pkgfile` - Updated uboot_version to v2025.10, new checksums
- `artifacts/u-boot/prepare/pkg.yaml` - GitHub releases URL instead of Collabora GitLab
- `artifacts/u-boot/nanopi-m6/pkg.yaml` - Native defconfig usage, removed patching
- `Makefile` - Added u-boot-nanopi-m6 dependency
- `docs/BOOT-TEST-CHECKLIST.md` - Attempt #3 results

## Decisions Made

1. **Switched to mainline U-Boot v2025.10** - Breaking change from Collabora fork
   - Rationale: Armbian proves mainline works; Collabora fork likely lacks M6 support
   - Impact: Cleaner build, native M6 support

2. **U-Boot source is NOT the root cause** - New investigation needed
   - Rationale: Identical failure with proven-working U-Boot version
   - Impact: Must investigate rkbin blob versions or boot media

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 3 - Blocking] Makefile missing u-boot-nanopi-m6 dependency**
- **Found during:** Task 3 (build attempt)
- **Issue:** `make local-talos-sbc-rk3588-mainline` didn't include nanopi-m6
- **Fix:** Added `u-boot-nanopi-m6` to Makefile targets
- **Committed in:** c088ccf

**2. [Rule 3 - Blocking] pkg.yaml platform specifier issue**
- **Found during:** Task 3 (build attempt)
- **Issue:** Build phase needed adjustment for mainline source structure
- **Fix:** Updated build section in pkg.yaml
- **Committed in:** c088ccf

---

**Total deviations:** 2 auto-fixed (2 blocking)
**Impact on plan:** Both were necessary build fixes. No scope creep.

## Issues Encountered

1. **Build initially failed** - Fixed via deviations above
2. **Boot test failed** - Primary plan objective not achieved

## User Setup Required

None - but hardware testing required UART for definitive diagnosis.

## Next Phase Readiness

**NOT READY for Phase 3**

Phase 2 goal (working U-Boot) not achieved after 3 boot attempts.

### Recommended Next Steps

1. **Update rkbin blob versions** (Highest priority)
   - Create plan 02-07 to update DDR to v1.18 and BL31 to v1.48
   - Match Armbian's exact rkbin versions
   - Low effort, directly addresses #1 suspect

2. **Test eMMC boot path** (If blob update fails)
   - Flash U-Boot to eMMC via MaskROM mode
   - Rule out SD card boot as the issue

3. **Extract Armbian's u-boot-rockchip.bin** (Quick diagnostic)
   - If Armbian's exact binary boots, issue is in our blob/build
   - If it also fails, issue is boot media

4. **Acquire UART adapter** (If all else fails)
   - TTL-USB adapter for serial console
   - Definitive diagnosis of exact failure point

### Tier Status

Per iteration strategy in BOOT-TEST-CHECKLIST.md:
- Tier 1: Exhausted (3 attempts, no LED activity)
- Tier 2: Completed (analysis done, defconfig/DTS approaches exhausted)
- **Tier 3: Continued investigation - blob versions are new hypothesis**

## Summary Table

| Attempt | Date | Defconfig | DDR | BL31 | Result |
|---------|------|-----------|-----|------|--------|
| 1 | 2026-02-02 | nanopi-r6c-rk3588s | v1.16 | v1.45 | FAIL |
| 2 | 2026-02-03 | rock5a + M6 DTS patch | v1.16 | v1.45 | FAIL |
| 3 | 2026-02-03 | nanopi-m6-rk3588s (mainline v2025.10) | v1.16 | v1.45 | FAIL |

**Common factor across all failures:** DDR v1.16, BL31 v1.45

---
*Phase: 02-bootloader*
*Completed: 2026-02-03*
*Status: PARTIAL SUCCESS - Build works, hardware verification FAILED*
