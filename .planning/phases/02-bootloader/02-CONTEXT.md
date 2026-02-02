# Phase 2: Bootloader Bring-Up - Context

**Gathered:** 2026-02-02
**Status:** Ready for planning

<domain>
## Phase Boundary

Get the NanoPi M6 to boot into U-Boot and display visual output on HDMI. This involves adapting U-Boot defconfig for RK3588S, integrating ATF and DDR initialization blobs, and establishing recovery procedures. No UART available — all verification via HDMI.

</domain>

<decisions>
## Implementation Decisions

### Debugging Approach
- No UART adapter available — HDMI output only
- Blind iteration strategy: flash, wait 2 minutes, check HDMI
- No fallback to UART acquisition — work within HDMI constraint
- If stuck with no HDMI output, Claude decides whether to try more configs or park

### Defconfig Strategy
- Claude's Discretion: Choose starting defconfig (R6C/R6S vs Rock 5A)
- Claude's Discretion: DDR blob sourcing (rkbin repo vs FriendlyElec BSP)
- Claude's Discretion: Iteration strategy if first config fails
- Claude's Discretion: New nanopi-m6_defconfig vs patching existing

### Boot Failure Recovery
- MaskROM mode: User is unfamiliar, needs full documentation
- Recovery fallback: Armbian SD card from Phase 1 is available
- eMMC policy: SD card only during Phase 2 — never touch eMMC
- Claude's Discretion: Brick recovery procedure documentation

### Verification Workflow
- Success criteria: U-Boot logo/splash visible on HDMI display
- Boot timeout: 2 minutes before declaring attempt failed
- Claude's Discretion: Test checklist format and documentation level

### Claude's Discretion
- Base defconfig selection (research best starting point)
- DDR blob source selection
- Iteration strategy when configs fail
- Config file organization (new vs patch)
- Recovery procedure documentation
- Test checklist format

</decisions>

<specifics>
## Specific Ideas

- User has built-in HDMI screen, no other display hardware
- Armbian SD card already verified working in Phase 1 — proven recovery path
- Conservative 2-minute timeout accounts for DDR training and display init
- Any U-Boot visual output (logo, splash, text) confirms success

</specifics>

<deferred>
## Deferred Ideas

None — discussion stayed within phase scope

</deferred>

---

*Phase: 02-bootloader*
*Context gathered: 2026-02-02*
