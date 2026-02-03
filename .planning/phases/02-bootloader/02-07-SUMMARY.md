---
phase: 02-bootloader
plan: 07
subsystem: bootloader
tags: [u-boot, rk3588s, rkbin, ddr-v1.18, bl31-v1.48, boot-failure, gap-closure]

# Dependency graph
requires:
  - phase: 02-06
    provides: Root cause analysis - DDR/BL31 blob versions suspected
provides:
  - rkbin blob versions updated to match Armbian (DDR v1.18, BL31 v1.48)
  - Eliminated blob versions as root cause
  - Boot test #4 results (FAILED - identical to #1-3)
  - Build configuration now matches Armbian EXACTLY
affects: [02-08-emmc-boot, 02-09-armbian-binary-test, phase-3-blocked]

# Tech tracking
tech-stack:
  added: []
  patterns:
    - "rkbin commit pinning for blob version control"

key-files:
  created: []
  modified:
    - Pkgfile
    - artifacts/u-boot/nanopi-m6/pkg.yaml
    - artifacts/u-boot/rock5a/pkg.yaml
    - artifacts/u-boot/rock5b/pkg.yaml
    - docs/BOOT-TEST-CHECKLIST.md

key-decisions:
  - "DDR/BL31 blob versions are NOT the root cause of boot failure"
  - "Build configuration now matches Armbian exactly (U-Boot v2025.10, DDR v1.18, BL31 v1.48)"
  - "SD card boot path or build process difference is primary suspect"
  - "Recommend testing Armbian's exact binary as diagnostic"

patterns-established:
  - "rkbin commit hash determines blob versions"
  - "All RK3588 boards in project share same rkbin commit"

# Metrics
duration: ~25min (continuation after checkpoint)
completed: 2026-02-03
---

# Phase 02 Plan 07: rkbin Blob Version Update Summary

**Updated DDR to v1.18 and BL31 to v1.48 matching Armbian exactly - boot test FAILED with identical symptoms, ELIMINATING blob versions as root cause**

## Performance

- **Duration:** ~25 min (post-checkpoint continuation)
- **Full plan duration:** ~45 min (including Tasks 1-3 build)
- **Started:** 2026-02-03T03:00:00Z (continuation at 04:01:42Z)
- **Completed:** 2026-02-03T04:30:00Z
- **Tasks:** 4/4 complete (Task 4 verification FAILED)
- **Files modified:** 5

## Accomplishments

- Updated rkbin to commit 0f8ac860 (same as Armbian) with DDR v1.18 and BL31 v1.48
- Build configuration now matches Armbian in ALL aspects
- ELIMINATED blob versions as potential root cause
- Narrowed down suspects to: SD card boot path or build process difference
- Documented comprehensive 4-attempt failure analysis

## Task Commits

1. **Task 1: Update rkbin commit reference** - `6011674`
   - Updated Pkgfile with new rkbin commit hash and checksums

2. **Task 2: Update blob version references in pkg.yaml** - `6dc19e5`
   - Changed DDR from v1.16 to v1.18
   - Changed BL31 from v1.45 to v1.48

3. **Task 3: Build U-Boot with updated blobs** - `2cad20d`
   - Updated rock5a/rock5b blob paths for build compatibility
   - Successfully built with v1.18/v1.48 blobs

4. **Task 4: Hardware boot verification (Attempt #4)** - `49579a8`
   - FAILED - identical symptoms to Attempts #1, #2, and #3
   - Updated BOOT-TEST-CHECKLIST.md with Attempt #4

## Boot Test Results

### Attempt #4 Configuration

| Setting | Value |
|---------|-------|
| U-Boot source | Mainline v2025.10 (same as Armbian) |
| Defconfig | nanopi-m6-rk3588s_defconfig (same as Armbian) |
| DDR blob | **v1.18** (same as Armbian) |
| BL31 | **v1.48** (same as Armbian) |
| rkbin commit | 0f8ac860 (same as Armbian) |

### Observations

| Indicator | Result |
|-----------|--------|
| LED activity | No (SYS solid ON, LED1 off, no blink) |
| HDMI output | No signal |
| Network | No DHCP lease |
| Boot stage reached | TPL/DDR (failed before SPL) |

### Conclusion

**IDENTICAL to Attempts #1, #2, and #3** - No improvement from updating blob versions.

## Root Cause Analysis

### Eliminated Causes (After 4 Attempts)

| Hypothesis | Status | Evidence |
|------------|--------|----------|
| Wrong device tree | ELIMINATED | Attempt #2 with correct DTS still failed |
| Wrong defconfig base (rock5a) | ELIMINATED | Attempt #3 uses native M6 defconfig |
| Collabora fork lacks M6 support | ELIMINATED | Mainline v2025.10 has full support, still fails |
| U-Boot version too old | ELIMINATED | Using same version as Armbian |
| DDR blob version mismatch | ELIMINATED | Attempt #4 uses v1.18 (same as Armbian) |
| BL31 blob version mismatch | ELIMINATED | Attempt #4 uses v1.48 (same as Armbian) |

### Remaining Suspects (After Systematic Elimination)

1. **SD card boot path not supported** (HIGH PROBABILITY)
   - All tests used SD card boot
   - Armbian may have been tested on eMMC
   - SD card controller init may differ
   - NanoPi M6 may require specific SD card boot configuration

2. **Build process difference** (MEDIUM PROBABILITY)
   - Armbian build system may produce different binary
   - Cross-compilation toolchain differences
   - Make flags or environment variables
   - Need to test Armbian's exact binary for diagnosis

3. **Partition/boot sector layout** (LOW-MEDIUM PROBABILITY)
   - U-Boot SPL offset or alignment issue
   - Our SD card uses Armbian's partitions, but may need verification

4. **Hardware-specific issue** (LOW PROBABILITY)
   - This specific NanoPi M6 unit
   - Unlikely since Armbian SD card boots successfully on same hardware

## Configuration Comparison

| Component | Our Build (02-07) | Armbian | Match? |
|-----------|-------------------|---------|--------|
| U-Boot version | v2025.10 | v2025.10 | YES |
| Defconfig | nanopi-m6-rk3588s | nanopi-m6-rk3588s | YES |
| DDR blob | v1.18 | v1.18 | YES |
| BL31 blob | v1.48 | v1.48 | YES |
| rkbin commit | 0f8ac860 | 0f8ac860 | YES |
| Boot result | FAIL | SUCCESS | NO |

**All configuration matches, but different boot outcome = build process or boot media issue**

## Files Created/Modified

- `Pkgfile` - Updated rkbin_ref to 0f8ac860 with new checksums
- `artifacts/u-boot/nanopi-m6/pkg.yaml` - DDR v1.18, BL31 v1.48 blob paths
- `artifacts/u-boot/rock5a/pkg.yaml` - Updated blob paths for consistency
- `artifacts/u-boot/rock5b/pkg.yaml` - Updated blob paths for consistency
- `docs/BOOT-TEST-CHECKLIST.md` - Attempt #4 results and analysis

## Decisions Made

1. **Blob versions are NOT the root cause** - Major finding
   - Rationale: Identical failure with exact Armbian blob versions
   - Impact: Must investigate build process or boot media next

2. **SD card boot path is primary suspect** - New hypothesis
   - Rationale: Only untested variable after matching all configuration
   - Impact: Need eMMC boot test or Armbian binary test

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 3 - Blocking] Rock5a/Rock5b pkg.yaml blob paths**
- **Found during:** Task 3 (build attempt)
- **Issue:** Other board configs still referenced v1.16/v1.45 blobs
- **Fix:** Updated blob paths in rock5a and rock5b pkg.yaml files
- **Committed in:** 2cad20d

---

**Total deviations:** 1 auto-fixed (1 blocking)
**Impact on plan:** Necessary for successful build. No scope creep.

## Issues Encountered

1. **Boot test failed** - Primary objective not achieved
2. **Blob version hypothesis disproven** - But this is valuable information

## User Setup Required

None for build. Hardware testing may require:
- eMMC flashing capability (via MaskROM mode)
- UART adapter for definitive boot log analysis

## Next Phase Readiness

**NOT READY for Phase 3**

Phase 2 goal (working U-Boot) not achieved after 4 boot attempts.

### Recommended Next Steps (Prioritized)

1. **Test Armbian's exact u-boot-rockchip.bin** (HIGHEST PRIORITY - Diagnostic)
   - Extract bootloader from working Armbian image
   - Flash ONLY the bootloader to our SD card
   - If boots: Issue is in our build process/toolchain
   - If fails: Issue is in SD card layout or hardware path
   - Estimated time: 10 minutes

2. **Test eMMC boot path** (If Armbian binary fails on SD)
   - Use MaskROM mode to flash U-Boot to eMMC
   - Boot from eMMC instead of SD card
   - Rules out SD card boot path as issue

3. **Binary comparison** (If Armbian binary works)
   - Compare file sizes and structure
   - Identify build process differences
   - May need to match Armbian's exact toolchain

4. **Acquire UART adapter** (If all else fails)
   - TTL-USB adapter for serial console
   - Definitive diagnosis of exact failure point

### Summary Table

| Attempt | Date | Defconfig | DDR | BL31 | Result |
|---------|------|-----------|-----|------|--------|
| 1 | 2026-02-02 | nanopi-r6c-rk3588s | v1.16 | v1.45 | FAIL |
| 2 | 2026-02-03 | rock5a + M6 DTS patch | v1.16 | v1.45 | FAIL |
| 3 | 2026-02-03 | nanopi-m6-rk3588s (mainline v2025.10) | v1.16 | v1.45 | FAIL |
| 4 | 2026-02-03 | nanopi-m6-rk3588s (mainline v2025.10) | v1.18 | v1.48 | FAIL |

**After 4 attempts:** All configuration matches Armbian, failure persists = build process or boot media issue

---
*Phase: 02-bootloader*
*Completed: 2026-02-03*
*Status: GAP CLOSURE FAILED - Build works, blob versions updated, hardware verification FAILED*
