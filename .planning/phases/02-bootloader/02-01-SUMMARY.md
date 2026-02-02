---
phase: 02-bootloader
plan: 01
subsystem: bootloader
tags: [u-boot, rk3588s, build-config]

dependency-graph:
  requires:
    - 01-environment-setup
  provides:
    - nanopi-m6-uboot-config
  affects:
    - 02-02 (U-Boot build validation)
    - 02-03 (defconfig customization)

tech-stack:
  added: []
  patterns:
    - bldr/kres build system pkg.yaml pattern

key-files:
  created:
    - artifacts/u-boot/nanopi-m6/pkg.yaml
  modified:
    - artifacts/u-boot/pkg.yaml

decisions:
  - id: uboot-defconfig
    choice: nanopi-r6c-rk3588s_defconfig
    reason: Same RK3588S SoC as NanoPi M6, verified compatible per research

metrics:
  duration: ~1 minute
  completed: 2026-02-02
---

# Phase 02 Plan 01: U-Boot Build Configuration Summary

U-Boot pkg.yaml created for NanoPi M6 using nanopi-r6c-rk3588s_defconfig as base, with DDR v1.16 blob (LPDDR5 support) and BL31 v1.45.

## What Was Done

### Task 1: Create NanoPi M6 U-Boot pkg.yaml

Created `artifacts/u-boot/nanopi-m6/pkg.yaml` following the rock5a/rock5b pattern from the upstream repository.

**Key configuration choices:**
- **defconfig:** `nanopi-r6c-rk3588s_defconfig` - Uses the NanoPi R6C config since it shares the same RK3588S SoC as the M6
- **DDR blob:** `rk3588_ddr_lp4_2112MHz_lp5_2400MHz_v1.16.bin` - Supports both LPDDR4 and LPDDR5 (NanoPi M6 uses LPDDR5)
- **BL31 blob:** `rk3588_bl31_v1.45.elf` - ARM Trusted Firmware compatible with DDR v1.16
- **CFLAGS:** Added debug logging flags (same as rock5a/rock5b) for boot troubleshooting

**Commit:** 52ae3d6

### Task 2: Update U-Boot Aggregator

Updated `artifacts/u-boot/pkg.yaml` to include `u-boot-nanopi-m6` in the dependencies list. This registers the NanoPi M6 as a build target in the bldr/kres build system.

**Commit:** da5137d

## Verification Results

| Check | Result |
|-------|--------|
| Directory exists: `artifacts/u-boot/nanopi-m6/` | PASS |
| pkg.yaml valid YAML syntax | PASS |
| Uses `nanopi-r6c-rk3588s_defconfig` | PASS |
| DDR blob includes `lp5` for LPDDR5 support | PASS |
| Aggregator includes `u-boot-nanopi-m6` | PASS |

## Deviations from Plan

None - plan executed exactly as written.

## Decisions Made

| Decision | Choice | Rationale |
|----------|--------|-----------|
| U-Boot defconfig | nanopi-r6c-rk3588s_defconfig | Same RK3588S SoC as NanoPi M6; research validated this approach |
| DDR blob version | v1.16 | Supports LPDDR5 (required for NanoPi M6 hardware) |
| BL31 version | v1.45 | Compatible with DDR v1.16, matches rock5a/rock5b |
| CFLAGS | Debug logging enabled | Consistent with other boards, aids boot debugging |

## Files Changed

### Created
- `artifacts/u-boot/nanopi-m6/pkg.yaml` - Build configuration for NanoPi M6 U-Boot

### Modified
- `artifacts/u-boot/pkg.yaml` - Added nanopi-m6 to aggregator dependencies

## Next Phase Readiness

**Ready for:** 02-02-PLAN.md (Build U-Boot)

**Prerequisites met:**
- Build configuration exists
- Aggregator includes NanoPi M6 target
- DDR and BL31 blob paths configured

**Note:** The defconfig `nanopi-r6c-rk3588s_defconfig` may need customization in 02-03 if the initial build doesn't boot correctly on NanoPi M6 hardware. Key areas to watch:
- UART console configuration (may differ from R6C)
- Display/HDMI output (M6 has different display connectors)
- GPIO/LED assignments
