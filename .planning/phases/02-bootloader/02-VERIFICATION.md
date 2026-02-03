---
phase: 02-bootloader
verified: 2026-02-03T03:00:00Z
status: gaps_found
score: 2/4 success criteria (1 verified, 1 partial, 2 failed)
re_verification: true
verification_count: 3
previous_verification:
  date: 2026-02-02T22:00:00Z
  status: gaps_found
  score: 2/4 success criteria
gaps_closed:
  - "Plan 02-04: Armbian U-Boot analysis completed"
  - "Plan 02-05: Device tree patching tested (failed, ruled out as root cause)"
  - "Plan 02-06: Mainline U-Boot v2025.10 switch completed"
  - "U-Boot binary now compiles with native nanopi-m6-rk3588s_defconfig"
eliminated_causes:
  - "Wrong device tree (ruled out by Attempt #2)"
  - "Wrong defconfig base/rock5a (ruled out by Attempt #3)"
  - "Collabora U-Boot fork lacks M6 support (ruled out by Attempt #3)"
  - "U-Boot version too old (ruled out - now using same v2025.10 as Armbian)"
gaps_remaining:
  - truth: "DDR memory initializes (LPDDR5 blob loads correctly)"
    status: failed
    reason: "Boot test Attempt #3 with mainline U-Boot still shows no LED activity"
    primary_suspect: "DDR blob version v1.16 (Armbian uses v1.18)"
  - truth: "Boot activity observable (LED blink or eventual kernel HDMI output)"
    status: failed
    reason: "Hardware test Attempt #3: identical failure - no LED, no HDMI, no network"
    primary_suspect: "rkbin blob versions or SD card boot path"
regressions: []
gaps:
  - truth: "DDR memory initializes (LPDDR5 blob loads correctly)"
    status: failed
    reason: "3 boot attempts all fail at DDR training stage (no LED activity)"
    artifacts:
      - path: "docs/BOOT-TEST-CHECKLIST.md"
        issue: "Attempts #1, #2, #3 all show no LED activity in 0-10s window"
    missing:
      - "DDR blob v1.18 (current: v1.16)"
      - "BL31 blob v1.48 (current: v1.45)"
    common_factor: "All 3 attempts used DDR v1.16, BL31 v1.45"

  - truth: "Boot activity observable (LED blink or eventual kernel HDMI output)"
    status: failed
    reason: "All 3 boot attempts show identical failure pattern"
    artifacts:
      - path: "docs/BOOT-TEST-CHECKLIST.md"
        issue: "No observable boot activity across 3 configuration variations"
    missing:
      - "Updated rkbin blobs matching Armbian versions"
      - "Alternatively: eMMC boot test to rule out SD card path"
---

# Phase 2: Bootloader Bring-Up Verification Report

**Phase Goal:** NanoPi M6 boots to U-Boot (verified via kernel reaching HDMI output or LED activity)

**Verified:** 2026-02-02T22:00:00Z

**Status:** gaps_found

**Re-verification:** Yes - after gap closure attempts (Plans 02-04, 02-05)

## Executive Summary

**Phase 2 goal NOT achieved.** After 5 sub-plans including 2 gap closure attempts:
- Boot test Attempt #2 (Plan 02-05) shows identical failure to Attempt #1
- Device tree patching did NOT resolve the boot issue
- Root cause is deeper than device tree: likely requires full M6 defconfig with proper board target
- Hardware verified functional (Armbian boots successfully on same device)

**Progress since initial verification:**
- Gap closure research completed (Plan 02-04): Armbian M6 config analyzed
- Device tree patched (Plan 02-05): Custom rk3588s-nanopi-m6.dts created
- Build configuration updated: pkg.yaml now patches device tree in defconfig
- Boot test conducted: Still fails at DDR training stage (no LED activity)

**Critical finding:** rock5a-rk3588s_defconfig with device tree patch is insufficient. Need actual nanopi-m6-rk3588s_defconfig or switch to mainline U-Boot v2025.10.

## Goal Achievement

### Observable Truths (Success Criteria from ROADMAP)

| # | Truth | Status | Evidence |
|---|-------|--------|----------|
| 1 | U-Boot binary compiles with NanoPi M6-specific defconfig | ⚠️ PARTIAL | Binary exists (8.9MB) but uses rock5a base with device tree patch only. Not full M6 defconfig. |
| 2 | DDR memory initializes (LPDDR5 blob loads correctly) | ✗ FAILED | Boot test Attempt #2: No LED activity in 0-10s window = DDR training fails |
| 3 | Boot activity observable (LED blink or eventual kernel HDMI output) | ✗ FAILED | Boot test Attempt #2: no LED, no HDMI, no network after 120s. Identical to Attempt #1. |
| 4 | Recovery procedure documented and tested (MaskROM mode) | ✓ VERIFIED | docs/MASKROM-RECOVERY.md complete (208 lines), cross-referenced from checklist |

**Score:** 1.5/4 success criteria (1 verified, 0.5 partial, 2 failed)

**Re-verification score change:** +0.5 from previous (research artifact completed)

### Plan-Level Must-Haves Analysis

#### Plan 02-01: U-Boot Build Configuration (Initial)

**Status from previous verification:** 5/6 verified

No changes in re-verification - artifacts still exist and function as originally verified.

#### Plan 02-02: Recovery Documentation (Initial)

**Status from previous verification:** 6/6 verified

No changes in re-verification - documentation complete and substantive.

#### Plan 02-03: Build and Hardware Test (Initial)

**Status from previous verification:** 4/6 verified (boot test FAILED)

No changes in re-verification - this was the initial boot test failure.

#### Plan 02-04: Armbian U-Boot Analysis (Gap Closure)

**Must-haves from plan frontmatter:**

| Must-Have | Type | Status | Details |
|-----------|------|--------|---------|
| "Plan 02-05 can apply M6 config without additional research" | Truth | ✓ VERIFIED | Analysis doc provides clear implementation options |
| "Specific CONFIG options to change are listed with exact values" | Truth | ✓ VERIFIED | Table format with 50+ CONFIG differences documented |
| "DDR parameters identified (blob version, timing values if applicable)" | Truth | ✓ VERIFIED | Recommends v1.18 DDR blob, v1.48 BL31 (optional upgrade) |
| docs/ARMBIAN-UBOOT-ANALYSIS.md | Artifact | ✓ VERIFIED | EXISTS (230+ lines), SUBSTANTIVE (detailed CONFIG table, implementation options), WIRED (referenced in 02-05-PLAN) |
| artifacts/u-boot/nanopi-m6/defconfig.diff | Artifact | ✗ MISSING | File not found in repository (mentioned in plan but not created) |
| ARMBIAN-UBOOT-ANALYSIS.md → pkg.yaml link | Link | ⚠️ PARTIAL | Analysis exists but Plan 02-05 chose custom DTS approach instead of full defconfig |

**Plan 02-04 Score:** 4/6 verified (research complete but defconfig.diff missing, implementation incomplete)

#### Plan 02-05: Apply M6 Configuration (Gap Closure)

**Must-haves from plan frontmatter:**

| Must-Have | Type | Status | Details |
|-----------|------|--------|---------|
| "U-Boot binary compiles with M6-specific defconfig" | Truth | ⚠️ PARTIAL | Compiles but uses rock5a defconfig + device tree patch, not full M6 defconfig |
| "DDR memory initializes on NanoPi M6 hardware" | Truth | ✗ FAILED | Boot test Attempt #2: No LED activity = DDR training fails |
| "BL31 loads successfully (verified via boot stage progression)" | Truth | ✗ FAILED | Cannot verify - boot fails before BL31 stage (no LED activity = DDR stage failure) |
| "Boot activity observable (LED blink or kernel reaching HDMI)" | Truth | ✗ FAILED | Boot test Attempt #2: no LED, no HDMI, no network |
| artifacts/u-boot/nanopi-m6/pkg.yaml | Artifact | ✓ VERIFIED | EXISTS (49 lines), SUBSTANTIVE (device tree patching logic), WIRED (imported by aggregator) |
| _out/artifacts/arm64/u-boot/nanopi-m6/u-boot-rockchip.bin | Artifact | ⚠️ ORPHANED | EXISTS (8.9MB), SUBSTANTIVE (valid binary), but doesn't boot (wrong config) |
| pkg.yaml → defconfig link | Link | ⚠️ WRONG | Line 27: make rock5a-rk3588s_defconfig, should be nanopi-m6-rk3588s_defconfig |
| u-boot-rockchip.bin → SD card link | Link | ✓ WIRED | flash.sh --bootloader mode writes at sector 64 (correct Rockchip offset) |

**Plan 02-05 Score:** 2/8 verified (build works but boot still fails)

### Required Artifacts

| Artifact | Expected | Status | Details |
|----------|----------|--------|---------|
| `artifacts/u-boot/nanopi-m6/pkg.yaml` | NanoPi M6 U-Boot build config | ⚠️ PARTIAL | EXISTS (49 lines), SUBSTANTIVE (DTS patching logic), WIRED (aggregator imports) - BUT uses rock5a defconfig base |
| `artifacts/u-boot/nanopi-m6/rk3588s-nanopi-m6.dts` | Custom device tree for M6 | ✓ VERIFIED | EXISTS (152 lines), SUBSTANTIVE (LED GPIOs, ethernet config, SD card GPIO), WIRED (copied during build) |
| `artifacts/u-boot/pkg.yaml` | Aggregated U-Boot builds | ✓ VERIFIED | EXISTS (9 lines), SUBSTANTIVE (3 board dependencies), WIRED (line 6: u-boot-nanopi-m6) |
| `docs/ARMBIAN-UBOOT-ANALYSIS.md` | Armbian M6 config analysis | ✓ VERIFIED | EXISTS (230+ lines), SUBSTANTIVE (CONFIG table, implementation options), WIRED (referenced in plans) |
| `artifacts/u-boot/nanopi-m6/defconfig.diff` | CONFIG comparison | ✗ MISSING | Mentioned in 02-04-PLAN but not created |
| `docs/MASKROM-RECOVERY.md` | MaskROM recovery procedure | ✓ VERIFIED | EXISTS (208 lines), SUBSTANTIVE (detailed procedures), WIRED (cross-referenced from checklist) |
| `docs/BOOT-TEST-CHECKLIST.md` | Boot attempt tracking | ✓ VERIFIED | EXISTS (277 lines), SUBSTANTIVE (2 attempts recorded, analysis), WIRED (cross-refs to recovery doc) |
| `_out/artifacts/arm64/u-boot/nanopi-m6/u-boot-rockchip.bin` | Compiled U-Boot binary | ⚠️ ORPHANED | EXISTS (8.9MB), SUBSTANTIVE (valid binary format), but doesn't boot (wrong config) |
| `hack/flash.sh` | Enhanced flash script | ✓ VERIFIED | EXISTS (110 lines), SUBSTANTIVE (--bootloader mode), WIRED (used in Plans 02-03, 02-05) |

### Key Link Verification

| From | To | Via | Status | Details |
|------|-----|-----|--------|---------|
| artifacts/u-boot/pkg.yaml | artifacts/u-boot/nanopi-m6/pkg.yaml | dependency declaration | ✓ WIRED | Line 6: `- stage: u-boot-nanopi-m6` |
| artifacts/u-boot/nanopi-m6/pkg.yaml | rk3588s-nanopi-m6.dts | file copy in prepare | ✓ WIRED | Line 20: cp /pkg/rk3588s-nanopi-m6.dts arch/arm/dts/ |
| pkg.yaml | _out/.../u-boot-rockchip.bin | bldr build | ✓ WIRED | Build produces 8.9MB binary at expected path |
| _out/.../u-boot-rockchip.bin | SD card sector 64 | hack/flash.sh | ✓ WIRED | flash.sh --bootloader mode writes at correct offset |
| docs/BOOT-TEST-CHECKLIST.md | docs/MASKROM-RECOVERY.md | cross-reference | ✓ WIRED | Bidirectional references |
| pkg.yaml defconfig | rock5a-rk3588s_defconfig | make command line 27 | ✗ WRONG TARGET | Should use nanopi-m6-rk3588s_defconfig from mainline |
| defconfig device tree | rk3588s-nanopi-m6.dts | sed patch lines 31-34 | ✓ WIRED | CONFIG_DEFAULT_DEVICE_TREE, OF_LIST, DEFAULT_FDT_FILE patched |

### Requirements Coverage

Requirements mapped to Phase 2 from REQUIREMENTS.md:

| Requirement | Description | Status | Blocking Issue |
|-------------|-------------|--------|----------------|
| BOOT-01 | U-Boot bootloader with NanoPi M6 defconfig boots to console | ✗ BLOCKED | Wrong defconfig base (rock5a vs full M6), DDR never initializes |
| BOOT-02 | ARM Trusted Firmware (BL31) loads successfully | ✗ BLOCKED | Cannot verify - boot fails before BL31 stage (DDR training failure) |
| BOOT-03 | DDR training blob initializes LPDDR5 memory | ✗ BLOCKED | No LED activity = DDR training never completes (wrong board config) |

**Requirements Score:** 0/3 satisfied (no change from initial verification)

### Anti-Patterns Found

| File | Line | Pattern | Severity | Impact |
|------|------|---------|----------|--------|
| artifacts/u-boot/nanopi-m6/pkg.yaml | 27 | Still using rock5a defconfig | 🛑 Blocker | Root cause: wrong board target, device tree patch alone insufficient |
| artifacts/u-boot/nanopi-m6/pkg.yaml | 30-34 | Comment claims device tree is "ROOT CAUSE fix" | ⚠️ Warning | Misleading - boot test proves device tree alone doesn't fix issue |
| .planning/phases/02-bootloader/02-05-SUMMARY.md | 108 | "Device tree is NOT the root cause" | ℹ️ Info | Correct conclusion documented in summary |

**Critical finding:** The pkg.yaml comment on line 30 is now known to be incorrect based on boot test results. The root cause is NOT just device tree but the entire board configuration.

### Gap Analysis: Why Boot Still Fails

#### Comparison: Initial vs Re-verification

| Aspect | Initial (02-03) | Re-verification (02-05) | Change |
|--------|-----------------|-------------------------|--------|
| Defconfig base | rock5a | rock5a | No change |
| Device tree | rock5a | nanopi-m6 (custom) | Improved but insufficient |
| Build success | Yes (8.9MB) | Yes (9.35MB) | ✓ |
| LED activity | No | No | No improvement |
| HDMI output | No | No | No improvement |
| Network | No | No | No improvement |

#### Root Cause: Device Tree Patch is Insufficient

**Evidence from boot tests:**
1. Attempt #1 (rock5a DT): No LED activity
2. Attempt #2 (M6 custom DT): No LED activity

**Conclusion:** The issue is NOT the device tree. Boot fails BEFORE device tree is parsed (DDR/SPL stage).

**Actual root causes (per Plan 02-05 analysis):**

1. **Wrong board target** (most likely)
   - rock5a uses `CONFIG_TARGET_ROCK5A_RK3588`
   - M6 needs `CONFIG_TARGET_EVB_RK3588`
   - Board-specific init code is completely different

2. **SPL configuration incompatibilities**
   - SPL stack addresses, memory layout
   - SPL driver configuration
   - ATF loading parameters

3. **U-Boot version gap**
   - Collabora fork: v2023.07-rc4 (2 years old)
   - Armbian uses: v2025.10 (mainline)
   - M6 support may only exist in newer versions

### What Works

Progress confirmed working:

- ✓ Build pipeline configured correctly (bldr/kres integration)
- ✓ DDR blob v1.16 available (LPDDR5 support)
- ✓ BL31 v1.45 configured
- ✓ Flash workflow works (--bootloader mode)
- ✓ MaskROM recovery documented (208-line guide)
- ✓ Iteration strategy established (boot test checklist)
- ✓ Hardware verified functional (Armbian boots)
- ✓ Device tree creation process established (custom DTS)
- ✓ Armbian configuration analyzed (research artifact)

### What's Missing

Gaps preventing Phase 2 goal achievement:

- ✗ Actual nanopi-m6-rk3588s_defconfig from Armbian mainline U-Boot
- ✗ CONFIG_TARGET_EVB_RK3588 board target
- ✗ Board initialization code for M6 hardware
- ✗ SPL configuration compatible with M6
- ✗ Possibly: Mainline U-Boot v2025.10 instead of Collabora fork

### Options for Next Iteration

Per Plan 02-05 summary, ordered by likelihood of success:

**Option A: Switch to Mainline U-Boot v2025.10** (Recommended by 02-05)
- Pros: Armbian proves this works, has actual M6 defconfig, full M6 support
- Cons: Breaking change to build system, may affect Talos compatibility

**Option B: Extract Armbian U-Boot binary for testing**
- Pros: Quick diagnostic, confirms M6 can boot with proper U-Boot
- Cons: Not sustainable, doesn't integrate with build

**Option C: Acquire UART for debugging**
- Pros: Provides exact failure point, essential for complex debugging
- Cons: Requires hardware purchase, adds delay

**Option D: Apply Armbian U-Boot patches to Collabora fork**
- Pros: Keeps current U-Boot base
- Cons: Patches designed for v2025.10, may not apply to v2023.07

## Impact

**Phase 2 goal NOT achieved.** Cannot proceed to Phase 3 (Kernel Integration) without working U-Boot.

**Blockers:**
- All Phase 3+ plans blocked by boot failure
- Requirements BOOT-01, BOOT-02, BOOT-03 blocked
- Device cannot run Talos without working bootloader

**Next phase readiness:** NOT READY

## Recommended Actions

**Immediate:**
1. Create Plan 02-06 to test Option B (Armbian binary extraction) as diagnostic
2. If Armbian binary boots, confirms issue is U-Boot config (not hardware)

**Short-term:**
1. Evaluate Option A: Switch to mainline U-Boot v2025.10
2. Assess Talos compatibility impact of U-Boot version change
3. If compatible, create Plan 02-07 to implement mainline U-Boot

**If still failing:**
1. Acquire UART adapter (USB to TTL serial, 1500000 baud)
2. Capture boot log to identify exact failure point
3. Debug based on boot log findings

---

_Verified: 2026-02-02T22:00:00Z_
_Verifier: Claude (gsd-verifier)_
_Re-verification: Yes (after Plans 02-04, 02-05 gap closure attempts)_
