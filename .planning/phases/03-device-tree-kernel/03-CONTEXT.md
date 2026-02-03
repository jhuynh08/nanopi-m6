# Phase 3: Device Tree & Kernel - Context

**Gathered:** 2026-02-02
**Status:** Ready for planning

<domain>
## Phase Boundary

Linux kernel boots with essential NanoPi M6 hardware functional. This includes getting the device tree compiled, kernel booting to console, and validating critical drivers (Ethernet, USB, NVMe). **Important: This device has no eMMC — storage strategy is SD boot + NVMe root.**

</domain>

<decisions>
## Implementation Decisions

### Device Tree Source
- Use pre-extracted FriendlyELEC vendor DTB (same approach as vendor U-Boot)
- Extracted from: rk3588-sd-ubuntu-noble-minimal-6.1-arm64-20251222.img
- Single DTB variant, not multiple configurations
- Note: FriendlyELEC's DTS is incompatible with Talos/mainline kernel (uses vendor bindings), so we use the compiled DTB directly

### Kernel Source
- Use Talos-provided kernel (milas/talos-sbc-rk3588 base)
- Consistent with Talos overlay pattern

### Driver Priorities
- **Priority 1:** Ethernet — network is essential for Talos cluster join
- **Priority 2:** USB host — for peripherals and debugging
- **Priority 3:** NVMe — **moved from Phase 6** because device has no eMMC; NVMe is required for root filesystem
- eMMC validation removed from scope (device doesn't have eMMC)

### Hardware Configuration
- No eMMC on this NanoPi M6 unit
- Boot strategy: SD card boot + NVMe root filesystem
- NVMe driver must be validated in Phase 3 (not deferred to Phase 6)

### Validation Approach
- Primary debug output: HDMI console
- NVMe validation: Device detection (shows in lsblk/dmesg as /dev/nvme0n1)

### Claude's Discretion
- Device tree customization approach (minimal vs curated for Talos)
- DTS file location in repo structure
- Kernel config approach (existing RK3588 config vs customization)
- Kernel module strategy (built-in vs modular)
- Kernel version (follow upstream Talos repo)
- Network validation method (link LED, DHCP, ping)
- Test image strategy for driver validation

</decisions>

<specifics>
## Specific Ideas

- Device has no eMMC — this is a hardware constraint, not a choice
- NVMe becomes critical path since it's needed for root filesystem
- Vendor DTS + vendor U-Boot should provide consistent boot chain

</specifics>

<deferred>
## Deferred Ideas

None — discussion stayed within phase scope

**Note:** eMMC-related success criteria in ROADMAP.md need to be updated to reflect NVMe instead.

</deferred>

---

*Phase: 03-device-tree-kernel*
*Context gathered: 2026-02-02*
