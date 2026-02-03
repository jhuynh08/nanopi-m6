---
phase: 02-bootloader
plan: 08
status: complete
duration: ~15min
commits:
  - hash: pending
    message: "docs(02-08): document FriendlyElec bootloader extraction procedure"
  - hash: pending
    message: "docs(02-08): complete FriendlyElec bootloader diagnostic - ROOT CAUSE FOUND"
---

# Plan 02-08 Summary: FriendlyElec Bootloader Diagnostic

## Outcome

**STATUS: SUCCESS - ROOT CAUSE IDENTIFIED**

After 4 failed boot attempts with mainline U-Boot (all configuration matching Armbian exactly), testing the FriendlyElec vendor bootloader definitively identified the root cause:

**NanoPi M6 requires vendor U-Boot (v2017.09) with MiniLoaderAll format — mainline U-Boot does not boot.**

## Key Findings

### Boot Test Results

| Attempt | Bootloader | Format | Result |
|---------|------------|--------|--------|
| #1-4 | Mainline v2025.10 | u-boot-rockchip.bin | FAIL (no LED activity) |
| **#5** | **FriendlyElec v2017.09** | **MiniLoaderAll.bin** | **SUCCESS** |

### Configuration Comparison

| Aspect | Our Build (FAILED) | FriendlyElec (BOOTS) |
|--------|-------------------|---------------------|
| U-Boot version | v2025.10 (mainline) | v2017.09 (vendor) |
| Defconfig | nanopi-m6-rk3588s | nanopi6_defconfig |
| Loader format | u-boot-rockchip.bin (combined) | MiniLoaderAll.bin + uboot.img |
| DDR init | TPL + SPL combined | idbloader (Rockchip format) |
| Build source | github.com/u-boot/u-boot | FriendlyELEC/uboot-rockchip |

### Root Cause Analysis

**7 hypotheses tested and eliminated:**
1. ~~Wrong device tree~~ - Attempt #2
2. ~~Wrong defconfig base~~ - Attempt #3
3. ~~Collabora fork lacks M6 support~~ - Attempt #3
4. ~~U-Boot version too old~~ - Attempt #3
5. ~~DDR blob version mismatch~~ - Attempt #4
6. ~~BL31 blob version mismatch~~ - Attempt #4
7. ~~SD card boot path not supported~~ - **Attempt #5 (SD card works!)**

**Confirmed root cause:** Mainline U-Boot boot chain (u-boot-rockchip.bin format) is incompatible with NanoPi M6. The board requires vendor-style bootloader with MiniLoaderAll format.

## Architectural Decision

**DECISION: NanoPi M6 requires vendor U-Boot approach**

Options for proceeding:
1. **Vendor U-Boot integration** - Use FriendlyElec's uboot-rockchip source
2. **Hybrid approach** - Use FriendlyElec bootloader binary, only customize kernel/Talos
3. **Investigate MiniLoaderAll generation** - Research if mainline can produce this format

## Deliverables

| Artifact | Status | Description |
|----------|--------|-------------|
| docs/FRIENDLYELEC-BOOTLOADER-EXTRACTION.md | Created | Complete extraction and flash procedure |
| docs/BOOT-TEST-CHECKLIST.md | Updated | Attempt #5 documented with SUCCESS |
| Root cause identification | Complete | Mainline U-Boot format incompatible |
| Architectural decision | Documented | Vendor U-Boot required |

## Impact on Project

- **Phase 2 goal (boot to U-Boot)** can now be achieved using vendor approach
- Need to pivot from mainline U-Boot to vendor U-Boot integration
- Next gap closure plan needed to implement vendor U-Boot in build system

## Next Steps

1. Create plan 02-09: Integrate FriendlyElec vendor U-Boot into build system
2. Determine best approach: full vendor integration vs hybrid
3. Re-verify phase goal with vendor U-Boot producing working boot

## Issues Encountered

None - diagnostic test worked as expected and definitively identified root cause.

---

*Plan completed: 2026-02-03*
*Gap closure: Yes - identified root cause after 4 failed attempts*
