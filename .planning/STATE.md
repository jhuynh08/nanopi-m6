# Project State

## Project Reference

See: .planning/PROJECT.md (updated 2026-02-02)

**Core value:** The NanoPi M6 boots Talos Linux and successfully registers with Talos Omni to join the home cluster.
**Current focus:** Phase 2 - Bootloader Bring-up

## Current Position

Phase: 2 of 6 (Bootloader Bring-up)
Plan: 7 of 7 in current phase (BOOT FAILED - 4 attempts, SD card boot path suspected)
Status: BLOCKED - Boot test #4 failed, blob versions ELIMINATED as root cause
Last activity: 2026-02-03 - Completed 02-07-PLAN.md (rkbin blob update - FAILED)

Progress: [=========...........] 45%

## Performance Metrics

**Velocity:**
- Total plans completed: 9
- Average duration: ~31min
- Total execution time: ~5h 11min

**By Phase:**

| Phase | Plans | Total | Avg/Plan |
|-------|-------|-------|----------|
| 01-environment-setup | 3 | ~2h 50min | ~55min |
| 02-bootloader | 6 | ~2h 21min | ~23min |

**Recent Trend:**
- Last 5 plans: ~45min, ~8min, ~45min, ~1h (02-06 build+test)
- Trend: Hardware testing cycles dominate execution time

*Updated after each plan completion*

## Accumulated Context

### Decisions

Decisions are logged in PROJECT.md Key Decisions table.
Recent decisions affecting current work:

- [Init]: Use milas/talos-sbc-rk3588 as base (VALIDATED - forked successfully)
- [Init]: Use Armbian device tree as reference (VALIDATED - Armbian boots on hardware)
- [Init]: Target eMMC boot for production (pending validation)
- [01-01]: Origin remote points to jhuynh08/nanopi-m6
- [01-01]: Upstream remote points to milas/talos-sbc-rk3588 for syncing
- [01-02]: Keep upstream CI workflow using GHCR (not Docker Hub)
- [01-02]: Local build command: `make local-talos-sbc-rk3588-mainline DEST=_out`
- [01-03]: Use rdisk raw device for faster flash writes on macOS
- [01-03]: LED-based verification for non-UART debugging
- [02-01]: ~~Use nanopi-r6c-rk3588s_defconfig as base~~ SUPERSEDED by 02-06
- [02-01]: DDR blob v1.16 with LPDDR5 support, BL31 v1.45 (SUSPECT - may need update)
- [02-02]: LED-based verification primary method (U-Boot has no HDMI driver)
- [02-02]: 3-tier iteration strategy for systematic debugging
- [02-03]: nanopi-r6c/rock5a defconfig DOES NOT WORK for NanoPi M6 (VALIDATED - boot failed)
- [02-04]: Use Armbian nanopi-m6-rk3588s_defconfig directly (COMPLETED in 02-06)
- [02-05]: Device tree patching DOES NOT fix boot (same failure as Attempt #1)
- [02-05]: ~~Collabora U-Boot fork (v2023.07) may lack M6 support entirely~~ CONFIRMED, switched to mainline
- [02-06]: **Switched to mainline U-Boot v2025.10** (same as Armbian)
- [02-06]: **U-Boot source/version is NOT the root cause** (still fails with mainline)
- [02-06]: ~~Primary suspect: DDR/BL31 blob version mismatch~~ ELIMINATED by 02-07
- [02-07]: **Updated rkbin blobs to DDR v1.18, BL31 v1.48** (matches Armbian exactly)
- [02-07]: **Blob versions are NOT the root cause** (still fails with matching blobs)
- [02-07]: **Primary suspect: SD card boot path or build process difference**
- [02-07]: Build configuration now matches Armbian EXACTLY - different outcome = build/media issue

### Pending Todos

None.

### Blockers/Concerns

**ACTIVE BLOCKER - Phase 2 (Tier 3 - All Configuration Matches Armbian):**
- Boot test #4 FAILED: Mainline U-Boot v2025.10 with EXACT Armbian configuration does not boot
- Identical symptoms across ALL 4 attempts (no LED activity, no HDMI, no network)
- Failure occurs at DDR/TPL stage (before device tree parsing, before SPL)
- **Configuration now matches Armbian exactly** - issue is NOT in config

**Root Cause Analysis - Eliminated (6 hypotheses ruled out):**
- ~~Wrong device tree~~ ELIMINATED by Attempt #2
- ~~Wrong defconfig base (rock5a)~~ ELIMINATED by Attempt #3
- ~~Collabora fork lacks M6 support~~ ELIMINATED by Attempt #3
- ~~U-Boot version too old~~ ELIMINATED by Attempt #3
- ~~DDR blob version mismatch~~ ELIMINATED by Attempt #4 (v1.18 = Armbian)
- ~~BL31 blob version mismatch~~ ELIMINATED by Attempt #4 (v1.48 = Armbian)

**Root Cause Analysis - Current Suspects:**
1. **SD card boot path not supported** (HIGH PROBABILITY)
   - All 4 tests used SD card boot
   - Armbian may have tested on eMMC
   - SD card controller init may differ
   - Need to test: Armbian's binary OR eMMC boot

2. **Build process difference** (MEDIUM PROBABILITY)
   - Different toolchain produces different binary
   - Make flags or environment variables
   - Need to test Armbian's exact binary for diagnosis

3. **Partition/boot sector layout** (LOW-MEDIUM PROBABILITY)
   - U-Boot SPL offset or alignment issue
   - May need to verify exact write offset

**Next Steps - Choose one:**
1. **Option A (RECOMMENDED):** Test Armbian's exact u-boot-rockchip.bin
   - Pros: Quick diagnostic (10 min), definitive answer
   - If boots: Issue is our build process
   - If fails: Issue is SD card boot path
2. **Option B:** Test eMMC boot path via MaskROM
   - Pros: Rules out SD card boot as issue
   - Cons: More complex procedure
3. **Option C:** Acquire UART for detailed debugging
   - Pros: Definitive diagnosis
   - Cons: Hardware purchase, time delay

**CI Configuration:**
- Upstream workflow uses self-hosted runners (`pkgs` label)
- May need to modify for GitHub-hosted runners or set up self-hosted

**Mitigated:**
- Hardware verified working with Armbian baseline (01-03 complete)
- U-Boot build configuration created (02-01 complete)
- Recovery and debugging documentation created (02-02 complete)
- Build and flash workflow working (02-03 partial - build/flash OK)
- Armbian config extracted and analyzed (02-04 complete)
- Device tree patching tested and ruled out (02-05 complete - FAILED)
- Mainline U-Boot v2025.10 configured and tested (02-06 complete - FAILED)

## Session Continuity

Last session: 2026-02-03
Stopped at: Completed 02-07-PLAN.md (rkbin blob update - BOOT FAILED)
Resume file: None - awaiting next investigation decision

**Gap status:** OPEN - All configuration now matches Armbian exactly. Still fails.
**Key insight:** Same config, different outcome = build process or boot media issue
Recommended next action: Test Armbian's exact u-boot-rockchip.bin to determine if issue is build process or SD card boot path.

---
*State initialized: 2026-02-02*
*Last updated: 2026-02-03 (after 02-07 completion)*
