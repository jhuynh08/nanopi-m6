# Project State

## Project Reference

See: .planning/PROJECT.md (updated 2026-02-02)

**Core value:** The NanoPi M6 boots Talos Linux and successfully registers with Talos Omni to join the home cluster.
**Current focus:** Phase 2 - Bootloader Bring-up

## Current Position

Phase: 2 of 6 (Bootloader Bring-up)
Plan: 5 of 5 in current phase (BOOT FAILED - Tier 3 decision needed)
Status: BLOCKED - Boot test failed, device tree approach insufficient
Last activity: 2026-02-03 - Completed 02-05-PLAN.md (Apply M6 Config - FAILED)

Progress: [========............] 40%

## Performance Metrics

**Velocity:**
- Total plans completed: 8
- Average duration: ~32min
- Total execution time: ~4h 31min

**By Phase:**

| Phase | Plans | Total | Avg/Plan |
|-------|-------|-------|----------|
| 01-environment-setup | 3 | ~2h 50min | ~55min |
| 02-bootloader | 5 | ~1h 41min | ~20min |

**Recent Trend:**
- Last 5 plans: ~2min, ~45min, ~8min, ~45min (02-05 build+test)
- Trend: Hardware testing adds significant time

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
- [02-01]: Use nanopi-r6c-rk3588s_defconfig as base (same RK3588S SoC)
- [02-01]: DDR blob v1.16 with LPDDR5 support, BL31 v1.45
- [02-02]: LED-based verification primary method (U-Boot has no HDMI driver)
- [02-02]: 3-tier iteration strategy for systematic debugging
- [02-03]: nanopi-r6c/rock5a defconfig DOES NOT WORK for NanoPi M6 (VALIDATED - boot failed)
- [02-03]: Need Armbian U-Boot configuration extraction for M6-specific support
- [02-04]: ~~Root cause: wrong device tree (rock5a vs nanopi-m6)~~ INVALIDATED by 02-05
- [02-04]: Use Armbian nanopi-m6-rk3588s_defconfig directly (Option A recommended)
- [02-04]: Download defconfig and DTS from Armbian during build
- [02-05]: Device tree patching DOES NOT fix boot (same failure as Attempt #1)
- [02-05]: Root cause is deeper: rock5a defconfig fundamentally incompatible
- [02-05]: Collabora U-Boot fork (v2023.07) may lack M6 support entirely
- [02-05]: DECISION NEEDED: Switch to mainline U-Boot or acquire UART

### Pending Todos

None.

### Blockers/Concerns

**ACTIVE BLOCKER - Phase 2 (Tier 3 Decision Point):**
- Boot test #2 FAILED: Device tree patching did not fix boot
- Identical symptoms to Attempt #1 (no LED activity, no HDMI, no network)
- Device tree is NOT the root cause - failure occurs before DT parsing
- rock5a defconfig has fundamental incompatibilities with NanoPi M6
- Collabora U-Boot fork (v2023.07) may lack M6 board support entirely

**Root Cause Revision:**
- ~~Wrong device tree~~ INVALIDATED - same failure with correct DT
- **Actual cause:** Deeper incompatibility in SPL/DDR/board init code
- Armbian uses mainline U-Boot v2025.10, not Collabora fork
- ~2 year version gap between Collabora fork and Armbian's U-Boot

**Decision Required - Choose one:**
1. **Option A (Recommended):** Switch to mainline U-Boot v2025.10
   - Pros: Armbian proves this works, actual M6 defconfig available
   - Cons: Breaking change to build system
2. **Option B:** Apply Armbian patches to Collabora fork
   - Pros: Keeps current base
   - Cons: Patches may not apply cleanly, significant effort
3. **Option C:** Extract Armbian U-Boot binary for direct test
   - Pros: Quick diagnostic test
   - Cons: Not a solution, just confirms diagnosis
4. **Option D:** Acquire UART for detailed debugging
   - Pros: Exact failure point visible
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

## Session Continuity

Last session: 2026-02-03
Stopped at: Completed 02-05-PLAN.md (Apply M6 Config - BOOT FAILED)
Resume file: None - awaiting architectural decision

**Gap status:** OPEN - Device tree approach failed, need different strategy.
Recommended next action: User decision on which option (A/B/C/D) to pursue.

---
*State initialized: 2026-02-02*
*Last updated: 2026-02-03*
