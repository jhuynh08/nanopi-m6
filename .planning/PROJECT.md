# Talos Linux NanoPi M6 Image

## What This Is

A custom Talos Linux overlay and bootable image for the FriendlyELEC NanoPi M6 single-board computer (RK3588S SoC). This enables the NanoPi M6 to run Talos Linux and join an existing Talos Omni-managed Kubernetes cluster.

## Core Value

The NanoPi M6 boots Talos Linux and successfully registers with Talos Omni to join the home cluster.

## Requirements

### Validated

(None yet — ship to validate)

### Active

- [ ] Fork milas/talos-sbc-rk3588 as project base
- [ ] Add NanoPi M6 U-Boot configuration (defconfig, device tree)
- [ ] Add NanoPi M6 device tree for Linux kernel
- [ ] Create NanoPi M6 installer profile
- [ ] Build flashable Talos image for NanoPi M6
- [ ] Boot Talos on NanoPi M6 hardware
- [ ] Register device with Talos Omni
- [ ] Join existing home cluster as a node

### Out of Scope

- Upstream contribution to milas/talos-sbc-rk3588 — personal use only
- Upstream contribution to siderolabs/sbc-rockchip — not pursuing official support
- Mainline kernel/U-Boot patches — using existing Collabora/Armbian work
- LCD/display support — headless server use case
- Audio support — not needed for Kubernetes node
- WiFi/Bluetooth — using wired Ethernet

## Context

### Hardware
- **Board**: FriendlyELEC NanoPi M6
- **SoC**: Rockchip RK3588S (octa-core: 4x Cortex-A76 @ 2.4GHz + 4x Cortex-A55 @ 1.8GHz)
- **RAM**: LPDDR5 (2400MHz)
- **Storage**: eMMC socket, microSD, M.2 NVMe (PCIe 2.1 x1)
- **Network**: Native Gigabit Ethernet
- **Boot**: U-Boot v2017.09 (vendor), GPT partitioning

### Existing Talos Setup
- Home cluster running Talos Linux
- Managed via Talos Omni
- NanoPi M6 will be added as an additional node

### Research Findings

**Base Project**: [milas/talos-sbc-rk3588](https://github.com/milas/talos-sbc-rk3588)
- Community project for RK3588 Talos support
- Currently supports Rock 5A/5B
- Uses Collabora's forked U-Boot and Linux kernel
- Provides overlay structure compatible with Talos imager

**Device Tree Source**: [Armbian linux-rockchip](https://github.com/armbian/linux-rockchip/blob/rk-6.1-rkr3/arch/arm64/boot/dts/rockchip/rk3588s-nanopi-m6.dts)
- Merged October 2024 (PR #258 by @efectn)
- Compatible: `"friendlyelec,nanopi-m6", "rockchip,rk3588"`
- Based on `rk3588s-nanopi-r6-common.dtsi`

**U-Boot Status**:
- No mainline defconfig for NanoPi M6
- Related boards exist: `nanopi-r6c-rk3588s_defconfig`, `nanopi-r6s-rk3588s_defconfig`
- Same RK3588S SoC as Orange Pi 5 (which is supported)
- May need to create defconfig based on R6C/R6S or use generic-rk3588

**Kernel/Firmware**:
- Collabora maintains RK3588 kernel patches
- ARM Trusted Firmware: `rk3588/bl31.elf`
- DDR blob: `rk3588_ddr_lp4_2112MHz_lp5_2400MHz_v1.19.bin` (from rkbin)

## Constraints

- **Tech stack**: Must use Talos Linux overlay system (not a custom distro)
- **Compatibility**: Must work with Talos Omni for cluster management
- **Kernel**: Requires Collabora's RK3588 patches (mainline 6.6 LTS insufficient)
- **Boot media**: Target eMMC or microSD for initial boot

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| Use milas/talos-sbc-rk3588 as base | Has working RK3588 infrastructure, Collabora kernel support | — Pending |
| Use Armbian device tree as reference | Only existing NanoPi M6 DTS, merged and tested | — Pending |
| Target eMMC boot | Production use case, more reliable than SD | — Pending |

---
*Last updated: 2026-02-02 after initialization*
