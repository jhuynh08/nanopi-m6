# Feature Landscape: Talos Linux SBC Overlay for NanoPi M6

**Domain:** Talos Linux SBC overlay for headless Kubernetes node
**Researched:** 2026-02-02
**Overall Confidence:** HIGH (verified with official Talos documentation and community reference projects)

## Executive Summary

A working Talos Linux SBC overlay requires a precise set of components that enable the boot chain, hardware initialization, and cluster connectivity. For a headless Kubernetes node like the NanoPi M6, the feature set is narrow but each component is critical. Missing any table stakes item results in a non-booting system or inability to join the cluster. The RK3588S presents additional complexity due to incomplete mainline support, requiring Collabora's fork-based kernel and U-Boot.

---

## Table Stakes

Features that **must** work or the system fails to boot or join the cluster. These are non-negotiable.

| Feature | Why Required | Complexity | Dependencies | Notes |
|---------|--------------|------------|--------------|-------|
| **U-Boot Bootloader** | First-stage boot after BL31. Loads kernel, passes device tree. No U-Boot = no boot. | High | BL31, DDR blob, defconfig | NanoPi M6 has no upstream defconfig. Must adapt from R6C/R6S or create new. |
| **Device Tree Blob (DTB)** | Hardware initialization. Kernel cannot detect storage, network, USB without correct DTB. | Medium | Armbian DTS source | NanoPi M6 DTB exists in Armbian (`rk3588s-nanopi-m6.dts`). Must compile and include. |
| **Linux Kernel with RK3588S Support** | OS foundation. Must have drivers for storage (eMMC, NVMe) and Gigabit Ethernet. | High | Collabora kernel fork, DTB | Mainline 6.6 LTS insufficient. Need Collabora 6.9+ with RK3588 enablement. |
| **ARM Trusted Firmware (BL31)** | Secure monitor, EL3 boot stage. Required by RK3588 boot flow. | Low | None | Open source via TF-A v2.12+. Binary available from rkbin. |
| **DDR Memory Blob** | Memory training firmware. No DDR init = no boot. | Low | None | Closed-source binary from Rockchip rkbin. `rk3588_ddr_lp4_2112MHz_lp5_2400MHz_v1.19.bin`. |
| **Overlay Installer Binary** | Go binary implementing `overlay.Installer` interface. Talos imager requires this to generate bootable images. | Medium | Go, overlay/adapter package | Must handle U-Boot installation to boot media. |
| **Overlay Profile YAML** | Defines disk image parameters (size, bootloader type, format). Imager reads this to generate correct image. | Low | None | Standard format, rarely changes from template. |
| **Gigabit Ethernet Driver** | Network connectivity to cluster. No network = cannot join Omni or cluster. | Medium | Kernel config, DTB | GMAC driver for RK3588S native Ethernet. Must be built-in (not module) for boot reliability. |
| **eMMC/SD Card Storage Driver** | Boot and root filesystem storage. | Medium | Kernel config, DTB | sdhci-of-dwcmshc driver. Must work in HS400 mode for eMMC performance. |

### Table Stakes Verification Checklist

- [ ] Device powers on and U-Boot console appears
- [ ] U-Boot loads kernel and DTB
- [ ] Kernel boots without panic
- [ ] Storage devices detected (eMMC, SD)
- [ ] Network interface detected and gets IP address
- [ ] Talos enters maintenance mode (first boot)
- [ ] Machine appears in Talos Omni dashboard
- [ ] Machine successfully joins cluster as node

---

## Important (Production)

Features needed for reliable production use. System boots without them, but not suitable for cluster operation.

| Feature | Value Proposition | Complexity | Dependencies | Notes |
|---------|-------------------|------------|--------------|-------|
| **NVMe Support** | Fast local storage for persistent volumes. M.2 slot is key differentiator of NanoPi M6. | Low | Kernel config (nvme-core, nvme) | PCIe 2.1 x1 lane. Built-in to mainline kernel. |
| **USB 3.0 Support** | External storage, recovery tools. | Low | Kernel config | xhci-hcd, dwc3 drivers. Well-supported in mainline. |
| **CPU Frequency Scaling** | Thermal management, power efficiency. Without it, CPU runs at max frequency constantly. | Low | Kernel config, DTB | cpufreq-dt driver. Collabora kernel includes this. |
| **Hardware Watchdog** | System recovery from hangs. Kubernetes liveness depends on node recovery. | Low | Kernel config, DTB | rk3588-wdt driver. |
| **RTC (Real-Time Clock)** | Accurate time for certificates, logs. | Low | Kernel config | hym8563 RTC chip on NanoPi M6. Important for TLS. |
| **Thermal Management** | Prevent thermal throttling/shutdown under load. | Medium | Kernel config, DTB | thermal-tsadc, fan-gpio if heatsink has fan. |
| **Power Management (PMIC)** | Proper shutdown/reboot. Without it, ungraceful power cycling. | Medium | Kernel config, DTB | rk806 PMIC driver. |
| **I2C/GPIO Access** | Sensor integration, LED status indicators. | Low | Kernel config | Standard kernel support. SYS LED is useful for status. |

### Production Readiness Checklist

- [ ] NVMe drive detected and usable as disk
- [ ] CPU scales frequency based on load
- [ ] Thermal sensors report values
- [ ] Watchdog can reboot hung system
- [ ] RTC maintains time across reboots
- [ ] Graceful shutdown works

---

## Optional (Nice to Have)

Features that enhance functionality but can be deferred or omitted for headless Kubernetes node.

| Feature | Value Proposition | Complexity | Why Defer | Notes |
|---------|-------------------|------------|-----------|-------|
| **NPU (Neural Processing Unit)** | 6 TOPS AI acceleration. Useful for edge ML workloads. | High | Requires out-of-tree kernel module (rknn). Not needed for general Kubernetes node. | Available as Talos system extension if needed later. |
| **SPI Flash Boot** | Boot from SPI NOR flash, freeing SD/eMMC. | Medium | Not essential. eMMC boot is reliable. | NanoPi M6 has 16MB SPI flash. |
| **USB-C Power Delivery** | Negotiate power from USB-C chargers. | Low | Basic USB-C power works. PD is optimization. | U-Boot 2025.01 has USB-PD handling. |
| **Wake-on-LAN** | Remote power-on capability. | Low | Nice for remote management but not critical. | May require PMIC/PHY configuration. |
| **UART Console** | Serial debug access. | Low | Useful for development, not production. | GPIO header exposes UART2. |
| **Custom Kernel Arguments** | Tune kernel parameters. | Low | Default args work for most cases. | Can be added via overlay later. |

---

## Out of Scope (Explicitly Not Building)

Features intentionally excluded from this project scope.

| Anti-Feature | Why Avoid | What to Do Instead |
|--------------|-----------|-------------------|
| **HDMI Display Output** | Headless server use case. Display adds complexity (DDC/EDID, resolution negotiation) and HDMI driver is still incomplete in mainline. | Use serial console or SSH for management. |
| **GPU Acceleration (Mali-G610)** | No display, no graphics workloads planned. Mesa/Panfrost support is incomplete. | If needed later, add as system extension. |
| **Audio Output** | Not needed for Kubernetes node. Adds unnecessary drivers. | Omit audio device tree nodes if possible. |
| **WiFi/Bluetooth** | Using wired Gigabit Ethernet. WiFi adds complexity (driver modules, firmware blobs, regulatory). | Hardware doesn't have built-in WiFi anyway (M.2 E-key slot optional). |
| **Camera (MIPI-CSI)** | No camera workloads planned. | If needed, add as separate project. |
| **Display Interface (MIPI-DSI)** | Same as HDMI - headless use case. | Omit. |
| **IR Receiver** | Remote control not needed for server. | Omit. |
| **H.264/H.265 Hardware Decode/Encode** | No media transcoding workloads planned. VPU drivers are incomplete in mainline anyway. | Use CPU encoding if somehow needed. |
| **Upstream Contribution** | Personal use only. Upstream contribution to milas or siderolabs adds maintenance burden. | Fork and maintain privately. |
| **Multi-board Support** | Building for NanoPi M6 only. | Add other boards as separate projects if needed. |

---

## Feature Dependencies

```
Boot Chain (sequential, all required):
DDR Blob → BL31 (TF-A) → U-Boot → Kernel + DTB → Talos Init

Hardware Detection (kernel requires DTB):
DTB → Ethernet Driver → Network Stack → Omni Connection
DTB → eMMC/SD Driver → Root Filesystem Mount
DTB → NVMe Driver (if using M.2 storage)

Overlay Build (build-time dependencies):
bldr + kres → Pkgfile builds → U-Boot binary + DTB files + Kernel
Installer Go binary + Profile YAML → Overlay container image
Overlay image + Talos imager → Bootable disk image

Cluster Join (runtime sequence):
Boot → Network DHCP → Omni connection → Machine config → Cluster join
```

### Critical Path

The minimum viable path to a working node:

1. **BL31 + DDR blob** - From rkbin (binary, no build needed)
2. **U-Boot** - Needs defconfig adaptation/creation (HIGH effort)
3. **Device Tree** - Extract from Armbian, compile (MEDIUM effort)
4. **Kernel** - Use Collabora fork as-is (LOW effort - already built by milas project)
5. **Installer** - Adapt from Rock 5A/5B installer (MEDIUM effort)
6. **Profile** - Adapt from existing profiles (LOW effort)

### Dependency Risk Assessment

| Component | Risk | Mitigation |
|-----------|------|------------|
| U-Boot defconfig | HIGH - May need iterative debugging | Start with nanopi-r6c-rk3588s_defconfig, test incrementally |
| Device Tree | MEDIUM - Armbian DTS may have Armbian-specific bits | Compare with Rock 5A DTS structure, test hardware detection |
| Ethernet | MEDIUM - GMAC needs correct DTB pinmux | Verify against working Armbian image |
| eMMC | LOW - Standard driver, well-tested on RK3588 | Should work with correct DTB |

---

## MVP Recommendation

For MVP (Minimum Viable Product), prioritize **table stakes only**:

### MVP Feature Set

1. **U-Boot with NanoPi M6 support** - Boots the system
2. **Device Tree Blob** - Hardware initialization
3. **Kernel** - Use existing Collabora fork from milas project
4. **BL31 + DDR blob** - Boot firmware (binary drop-in)
5. **Gigabit Ethernet** - Cluster connectivity
6. **eMMC storage** - Boot and root filesystem
7. **Overlay installer and profile** - Talos imager integration

### Defer to Post-MVP

| Feature | Reason to Defer |
|---------|-----------------|
| NVMe support | Works with standard kernel, test after basic boot works |
| CPU frequency scaling | Optimization, not critical for function |
| Thermal management | Add after observing thermal behavior under load |
| Hardware watchdog | Production hardening, not MVP |
| NPU | Specialized workload, not general Kubernetes |

### MVP Validation Criteria

1. Flash image to eMMC or SD card
2. Power on NanoPi M6
3. Observe boot messages (UART or network discovery)
4. Machine appears in Talos Omni
5. Apply machine config
6. Node joins cluster
7. Run test workload (e.g., nginx pod)

---

## Complexity Summary

| Effort Level | Components |
|--------------|------------|
| **High** | U-Boot defconfig creation, Linux kernel customization (if needed) |
| **Medium** | Device tree compilation, Overlay installer adaptation, Ethernet driver validation |
| **Low** | BL31/DDR blob (binary drop-in), Profile YAML, NVMe/USB (kernel built-in), CPU scaling, Watchdog |

---

## Sources

### Official Documentation (HIGH confidence)
- [Talos Linux Overlays](https://docs.siderolabs.com/talos/v1.11/build-and-extend-talos/custom-images-and-development/overlays) - Overlay structure and requirements
- [siderolabs/overlays Repository](https://github.com/siderolabs/overlays) - Official overlay reference
- [siderolabs/sbc-rockchip](https://github.com/siderolabs/sbc-rockchip) - Official Rockchip overlay (23 boards supported)
- [Talos Linux v1.10 Release](https://github.com/siderolabs/talos/discussions/10842) - ARM64 boot changes

### Community Projects (HIGH confidence)
- [milas/talos-sbc-rk3588](https://github.com/milas/talos-sbc-rk3588) - Reference RK3588 overlay project
- [Armbian linux-rockchip](https://github.com/armbian/linux-rockchip) - NanoPi M6 device tree source

### Hardware Enablement (HIGH confidence)
- [Collabora RK3588 Progress](https://www.collabora.com/news-and-blog/news-and-events/rockchip-rk3588-upstream-support-progress-future-plans.html) - Mainline kernel status
- [Almost Open-Source Boot Chain for RK3588](https://www.collabora.com/news-and-blog/blog/2024/02/21/almost-a-fully-open-source-boot-chain-for-rockchips-rk3588/) - BL31/U-Boot status
- [FriendlyELEC NanoPi M6 Wiki](https://wiki.friendlyelec.com/wiki/index.php/NanoPi_M6) - Hardware specifications

### Talos Omni (HIGH confidence)
- [Getting Started with Omni](https://docs.siderolabs.com/omni/getting-started/getting-started) - Cluster join requirements
