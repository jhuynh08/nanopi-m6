---
phase: 02-bootloader
plan: 05
subsystem: bootloader
tags: [u-boot, rk3588s, device-tree, boot-failure, gap-closure]

# Dependency graph
requires:
  - phase: 02-04
    provides: Armbian U-Boot analysis and M6 defconfig extraction
provides:
  - Custom NanoPi M6 device tree (Collabora U-Boot compatible)
  - Updated pkg.yaml with M6 DTS patching
  - Boot test results (FAILED - identical to Attempt #1)
  - Root cause revision (device tree is NOT the issue)
affects: [02-06-further-investigation]

# Tech tracking
tech-stack:
  added: []
  patterns:
    - "Custom DTS creation for U-Boot fork compatibility"
    - "Defconfig patching via sed for device tree substitution"

key-files:
  created:
    - artifacts/u-boot/nanopi-m6/rk3588s-nanopi-m6.dts
  modified:
    - artifacts/u-boot/nanopi-m6/pkg.yaml
    - docs/BOOT-TEST-CHECKLIST.md

key-decisions:
  - "Device tree was NOT the root cause of boot failure"
  - "rock5a defconfig has fundamental incompatibilities beyond device tree"
  - "Collabora U-Boot fork (v2023.07) may lack M6 support entirely"
  - "Need different approach: mainline U-Boot or Armbian patches"

patterns-established:
  - "Minimal DTS creation for older U-Boot compatibility"
  - "CONFIG_OF_LIST must be updated alongside CONFIG_DEFAULT_DEVICE_TREE"

# Metrics
duration: 45min
completed: 2026-02-03
---

# Phase 02 Plan 05: Apply M6 Configuration Summary

**Boot test FAILED - Device tree patching did not resolve the issue. Root cause is deeper than expected.**

## Performance

- **Duration:** ~45 min (including build and hardware test)
- **Started:** 2026-02-03T00:52:43Z
- **Completed:** 2026-02-03T01:40:00Z
- **Tasks:** 3/3 complete (but hardware test failed)
- **Files created:** 1 (rk3588s-nanopi-m6.dts)
- **Files modified:** 2 (pkg.yaml, BOOT-TEST-CHECKLIST.md)

## Accomplishments

- Created minimal NanoPi M6 device tree compatible with Collabora U-Boot fork
- Successfully built U-Boot (9.35MB binary) with M6 device tree
- Updated pkg.yaml with proper defconfig patching (DEFAULT_DEVICE_TREE, OF_LIST, DEFAULT_FDT_FILE)
- Flashed and tested on hardware
- Documented failure and revised root cause analysis

## Task Commits

1. **Tasks 1-2: Update config and build** - `32b6cb9`
   - Created rk3588s-nanopi-m6.dts
   - Updated pkg.yaml with DTS patching
   - Build verified successful

2. **Task 3: Hardware test** - (this commit)
   - FAILED - identical symptoms to Attempt #1
   - Updated BOOT-TEST-CHECKLIST.md

## Boot Test Results

### Attempt #2 Configuration

| Setting | Value |
|---------|-------|
| Base defconfig | rock5a-rk3588s_defconfig |
| Device tree | rk3588s-nanopi-m6.dts (custom, minimal) |
| DDR blob | v1.16 |
| BL31 | v1.45 |
| Binary size | 9.35MB |

### Observations

| Indicator | Result |
|-----------|--------|
| LED activity | No (power only, no blink) |
| HDMI output | No signal |
| Network | No DHCP lease |
| Boot stage reached | TPL/DDR (failed) |

### Conclusion

**IDENTICAL to Attempt #1** - No improvement from device tree patching.

## Root Cause Analysis Revision

### What We Learned

1. **Device tree is NOT the root cause**
   - Custom M6 DTS with correct LED GPIOs (GPIO1_A4/A6), ethernet tx_delay (0x42), and SD card GPIO did not fix boot
   - Same early-stage failure as with rock5a device tree
   - Failure occurs BEFORE device tree is parsed (DDR/SPL stage)

2. **rock5a defconfig has deeper incompatibilities**
   - `CONFIG_TARGET_ROCK5A_RK3588` may trigger Rock5A-specific initialization
   - SPL configuration may differ significantly
   - Board-specific code paths in U-Boot source may be incompatible

3. **U-Boot version gap is significant**
   - Collabora fork: v2023.07-rc4
   - Armbian uses: v2025.10 (mainline)
   - ~2 years of development between versions
   - M6 support may only exist in newer versions

### Potential Root Causes (Ordered by Likelihood)

1. **Wrong board target** - `CONFIG_TARGET_ROCK5A_RK3588` vs `CONFIG_TARGET_EVB_RK3588`
   - Armbian uses EVB target for M6
   - Board-specific init code may be completely different

2. **SPL configuration incompatibilities**
   - SPL stack addresses, memory layout
   - SPL driver configuration
   - ATF loading parameters

3. **DDR timing in rkbin blobs**
   - DDR blob v1.16 may not have M6-specific timing
   - Armbian uses v1.18
   - But: blob versions unlikely to be board-specific

4. **PMIC initialization**
   - M6 uses RK806 PMIC on SPI2
   - Rock5A may have different PMIC or different bus
   - But: PMIC init happens after DDR

5. **U-Boot source code gaps**
   - Collabora fork may lack M6 board support entirely
   - Board detection/identification code may be missing

## Options for Next Iteration

### Option A: Switch to Mainline U-Boot v2025.10 (Recommended)

**Pros:**
- Armbian proves this works on M6
- Uses actual nanopi-m6-rk3588s_defconfig
- Full M6 support in DTS and board code

**Cons:**
- Breaking change to build system
- May require updating Pkgfile and deps
- Collabora fork was chosen for a reason (Talos compatibility?)

### Option B: Apply Armbian U-Boot Patches to Collabora Fork

**Pros:**
- Keeps current U-Boot base
- May work if patches are backward-compatible

**Cons:**
- Patches designed for v2025.10, may not apply to v2023.07
- Significant effort to adapt patches
- May introduce regressions

### Option C: Extract Working Binary from Armbian

**Pros:**
- Quick test to verify M6 boot works with proper U-Boot
- No build changes required
- Confirms whether issue is U-Boot config or something else

**Cons:**
- Not a sustainable solution
- Doesn't integrate with Talos build
- Just a diagnostic step

### Option D: Acquire UART for Debugging

**Pros:**
- Provides exact failure point in boot log
- Can diagnose DDR training, SPL, or ATF issues
- Essential for any complex debugging

**Cons:**
- Requires hardware purchase
- Adds time before next iteration

## Deviations from Plan

### Unplanned Work

1. **[Rule 3 - Blocking] Armbian DTS incompatibility with Collabora U-Boot**
   - Armbian DTS references nodes not in Collabora's rk3588s.dtsi
   - Had to create minimal custom DTS stripped of incompatible references

2. **[Rule 3 - Blocking] OF_LIST not auto-updating**
   - Changing CONFIG_DEFAULT_DEVICE_TREE didn't update OF_LIST
   - Binman failed with FDT list mismatch
   - Had to add explicit sed for OF_LIST

3. **[Rule 3 - Blocking] GNU Make version incompatibility**
   - macOS ships Make 3.81, project requires Make 4.x
   - Used gmake (Homebrew) instead

## User Setup Required

None - but hardware testing required for verification.

## Next Phase Readiness

**NOT READY for Phase 3**

Phase 2 goal (working U-Boot) not achieved. Boot failure persists.

### Recommended Path Forward

1. **Immediate:** Test Option C (Armbian binary) to confirm M6 can boot with proper U-Boot
2. **Short-term:** Evaluate switching to mainline U-Boot v2025.10 (Option A)
3. **If still failing:** Acquire UART adapter for detailed boot log analysis

### Tier Status

Per iteration strategy in BOOT-TEST-CHECKLIST.md:
- Tier 1: Exhausted (2 attempts, no LED activity)
- Tier 2: Partially completed (analysis done, but patch approach failed)
- **Tier 3: Decision point reached**

Decision needed: Continue with Option A/B/C, or acquire UART for debugging.

---
*Phase: 02-bootloader*
*Completed: 2026-02-03*
*Status: FAILED - Requires architectural decision for next approach*
