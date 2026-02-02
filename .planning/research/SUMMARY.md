# Project Research Summary

**Project:** Talos Linux SBC Overlay for NanoPi M6
**Domain:** Custom embedded Linux image for ARM64 Kubernetes node
**Researched:** 2026-02-02
**Confidence:** HIGH

## Executive Summary

Building a Talos Linux overlay for the NanoPi M6 (RK3588S) requires forking the milas/talos-sbc-rk3588 project and adapting it for NanoPi-specific hardware. The recommended approach is to use Collabora's mainline-focused kernel/U-Boot forks (not vendor BSP kernels) combined with the Siderolabs build ecosystem (bldr, kres, imager). The NanoPi M6 presents specific challenges: no upstream U-Boot defconfig exists (must adapt from R6C/R6S), device trees are only available from Armbian linux-rockchip, and the RK3588S boot chain requires closed-source DDR training blobs. The critical path to a working node is: U-Boot defconfig adaptation → Device Tree compilation → Installer/Profile creation → Talos integration.

The primary risk is bootloader bring-up: incorrect U-Boot configuration can brick the board, requiring MaskROM recovery. Secondary risks include kernel module signing mismatches (modules must be built with the same signing key as kernel) and device tree incompatibilities between bootloader and kernel. Mitigation requires staged testing with UART console access and building the complete overlay (kernel + modules + U-Boot) as a single artifact to ensure signature compatibility.

For a headless Kubernetes node, the feature scope is minimal but non-negotiable: U-Boot must boot, device tree must initialize hardware (Ethernet, eMMC, USB), network must connect to cluster, and kernel must support storage drivers. Optional features like NPU acceleration, HDMI output, and GPU can be deferred indefinitely. The MVP delivers a bootable Talos image that joins an Omni-managed cluster.

## Key Findings

### Recommended Stack

Talos SBC overlay development requires the Siderolabs build ecosystem combined with hardware-specific sources. The milas/talos-sbc-rk3588 project provides a proven foundation, but must be upgraded from Talos v1.7.x to v1.10.x for better ARM64 support and active maintenance.

**Core technologies:**
- **bldr v0.3.0+**: Container-based package builder using BuildKit for reproducible overlay builds. Required for processing Pkgfile definitions.
- **kres**: Makefile and CI config generator from .kres.yaml. Eliminates manual build scaffolding and ensures consistent build structure across Siderolabs projects.
- **Talos Linux v1.10.x**: Target OS version with systemd-boot UEFI support, unified /usr filesystem, and kernel 6.12.x baseline. More stable than v1.11/v1.12 for overlays.
- **Collabora kernel 6.9+**: Mainline-focused RK3588 enablement branch with most hardware drivers merged. Avoids vendor BSP bloat while providing necessary RK3588 support.
- **Collabora U-Boot v2024.01+**: Mainline-focused bootloader with RK3588 support and USB-PD handling (v2025.01). Critical for NanoPi M6 boot chain.
- **Armbian linux-rockchip**: Source for NanoPi M6 device tree (`rk3588s-nanopi-m6.dts`, added October 2024). Only verified DTS source for this board.
- **Rockchip rkbin**: Closed-source DDR training blob (`rk3588_ddr_lp4_2112MHz_lp5_2400MHz_v1.19.bin`) and optional BL31 binary. Unavoidable dependency.
- **ARM Trusted Firmware v2.12+**: Secure boot firmware (BL31). Open-source as of TF-A v2.12, but Collabora fork recommended until mainline RK3588 support matures.

**Critical version dependencies:**
- Kernel, U-Boot, and ATF must use compatible versions from Collabora forks (not mixing mainline/vendor/Collabora)
- DDR blob must support LPDDR5 for NanoPi M6 (not just LP4)
- Go 1.24.x for installer binaries (using siderolabs/talos overlay/adapter package)

### Expected Features

For a headless Kubernetes node, the feature landscape is narrow but every component is critical. Missing any table stakes item results in a non-booting system or inability to join the cluster.

**Must have (table stakes):**
- **U-Boot bootloader** — First-stage boot after BL31. No U-Boot = no boot. NanoPi M6 lacks upstream defconfig (HIGH complexity).
- **Device Tree Blob** — Hardware initialization. Kernel cannot detect storage/network/USB without correct DTB (MEDIUM complexity, source from Armbian).
- **Linux kernel with RK3588S drivers** — OS foundation. Must have eMMC (sdhci-of-dwcmshc) and Gigabit Ethernet (GMAC) drivers built-in (HIGH complexity if kernel customization needed).
- **BL31 and DDR blob** — Secure monitor and memory training. Required by RK3588 boot flow (LOW complexity, binary drop-in).
- **Overlay installer binary** — Go binary implementing overlay.Installer interface. Talos imager cannot generate bootable images without it (MEDIUM complexity).
- **Overlay profile YAML** — Defines disk image parameters (size, bootloader, format). Standard template (LOW complexity).
- **Gigabit Ethernet connectivity** — Network path to cluster/Omni. No network = cannot join cluster (MEDIUM complexity, DTB-dependent).
- **eMMC/SD storage driver** — Boot and root filesystem. Must support HS400 mode for performance (MEDIUM complexity).

**Should have (production readiness):**
- NVMe support (M.2 slot is key differentiator) — Fast local storage for persistent volumes
- CPU frequency scaling (cpufreq-dt) — Thermal management, prevents running at max frequency constantly
- Hardware watchdog (rk3588-wdt) — System recovery from hangs
- RTC (hym8563) — Accurate time for TLS certificates
- Thermal management (thermal-tsadc) — Prevent throttling/shutdown under load
- PMIC (rk806) — Proper shutdown/reboot handling

**Defer (v2+):**
- NPU (6 TOPS AI acceleration) — Requires out-of-tree rknn kernel module, not needed for general Kubernetes
- HDMI/GPU/Audio — Headless use case, display adds complexity
- WiFi/Bluetooth — Using wired Ethernet, hardware doesn't have built-in WiFi
- VPU hardware encode/decode — No media workloads planned, drivers incomplete in mainline

### Architecture Approach

Talos SBC overlays follow a standardized container-based structure derived from siderolabs/sbc-template. The architecture separates artifacts (firmware, bootloader, device trees), installers (Go binaries), and profiles (YAML configuration) into distinct components with clear boundaries.

**Major components:**

1. **Artifacts (`/artifacts/`)** — Contains firmware, bootloader, and device tree binaries. Organized hierarchically: SoC-level (ATF, rkbin) shared across RK3588 boards, board-level (U-Boot) specific to NanoPi M6. Built via pkg.yaml definitions and deployed by installer.

2. **Installers (`/installers/`)** — Go binaries implementing overlay.Installer interface. Must provide `GetOptions()` (returns kernel args) and `Install()` (deploys artifacts to disk). Statically linked for scratch containers, invoked by Talos imager during image generation.

3. **Profiles (`/profiles/`)** — YAML definitions specifying disk image parameters (arch, platform, diskSize, bootloader type, compression format). Consumed by Talos imager, paired with installers.

4. **Build system** — Orchestrated by bldr via Docker buildx. Pkgfile defines version variables and upstream sources, pkg.yaml per component defines build steps, kres generates Makefile from .kres.yaml configuration.

**Data flow:**
```
Pkgfile vars → bldr build → artifacts + installers + profiles →
Container image (OCI) → Talos imager → Bootable disk image (.raw.xz)
```

**Build order (dependencies):**
1. base (toolchain) → 2. ATF/rkbin (parallel) → 3. U-Boot (depends on ATF/rkbin) → 4. installer binary (parallel with U-Boot) → 5. profiles → 6. final overlay image

**Key patterns:**
- One pkg.yaml per board component with explicit dependencies
- Installers implement standard CLI contract (accept/output YAML, `install` and `get-options` subcommands)
- Profile YAML must match installer expectations (bootloader type, disk format)
- All kernel modules must be built with same signing key as kernel (critical for runtime loading)

### Critical Pitfalls

**Top 5 pitfalls to avoid:**

1. **Kernel module signing key mismatch** — Modules built separately from Talos kernel fail to load with cryptographic verification errors. Prevention: Build kernel and all modules in same pipeline with same signing key. Never mix siderolabs extensions with custom RK3588 kernel. (Phase 1 risk)

2. **Missing or incorrect U-Boot defconfig** — NanoPi M6 has no mainline defconfig. Wrong config causes boot failure at SPL/TPL stage. Prevention: Start with nanopi-r6c-rk3588s_defconfig, cross-reference Armbian's U-Boot config, ALWAYS have UART console connected during bring-up. (Phase 2 risk)

3. **Device tree mismatch between U-Boot and kernel** — Different DTS sources cause kernel panic or missing peripherals. Prevention: Use device tree from SAME kernel source tree, copy Armbian NanoPi M6 DTS into Collabora kernel and adapt, verify compatible strings match. (Phase 3 risk)

4. **Boot media priority confusion** — RK3588 boot order (SPI → eMMC → SD) is fixed. Flashing wrong media or wrong offset causes "successful" flash with no behavior change. Prevention: Verify boot media with lsblk before flashing, zero old boot media, document flash procedures. Use sector 64 offset (32KB) for Rockchip. (Phase 2 risk)

5. **Closed-source binary version mismatches** — Incompatible TPL/BL31/DDR blob versions cause DDR init failures or memory corruption. Prevention: Use EXACT versions from milas project or document tested combination. Pin rkbin commit hash. Ensure DDR blob supports LPDDR5 for NanoPi M6. (Phase 2 risk)

**Additional moderate pitfalls:**
- Random MAC address on every boot (breaks DHCP, Omni registration)
- Talos version upgrade path violations (only adjacent minor versions supported)
- Install disk path hardcoding (eMMC vs SD vs NVMe)
- Network configuration lost after install (must be in machine config, not just installer)
- Kernel config missing required drivers (compare with Armbian config)

## Implications for Roadmap

Based on research, the project should follow hardware bring-up order with clear phase boundaries to enable incremental validation and avoid catastrophic failures.

### Phase 1: Environment Setup and Base Fork
**Rationale:** Foundation must be solid before hardware work. Kernel module signing is established here and cannot be retrofitted later.

**Delivers:**
- Forked milas/talos-sbc-rk3588 repository
- Upgraded to Talos v1.10.x baseline
- Build pipeline verified (bldr, kres, Docker buildx)
- UART console configured and tested (1500000 baud, UART2)
- ARM64 build environment (native ARM64 or x86_64 with QEMU)

**Addresses pitfalls:**
- Kernel module signing key mismatch (establish single build pipeline)
- UART console not configured (set up debug access before hardware work)

**Research needs:** Standard patterns, skip research-phase

### Phase 2: Bootloader Bring-Up
**Rationale:** Boot chain is the highest risk area. Must validate U-Boot independently before kernel work. Binary version mismatches are catastrophic.

**Delivers:**
- NanoPi M6 U-Boot defconfig (adapted from R6C/R6S)
- Tested U-Boot build with ATF + DDR blob
- Verified boot to U-Boot console via UART
- Documented flash procedures (media, offset, recovery)

**Uses stack:**
- Collabora U-Boot v2024.01+
- ARM Trusted Firmware v2.12+
- Rockchip rkbin DDR blob v1.19+ (LPDDR5 support)

**Addresses pitfalls:**
- Missing/incorrect U-Boot defconfig (dedicated validation phase)
- Boot media priority confusion (establish flash procedures)
- Binary version mismatches (verify combination works)

**Research needs:** MEDIUM - May need research-phase for U-Boot defconfig adaptation strategies

### Phase 3: Device Tree and Kernel
**Rationale:** After U-Boot is stable, bring up kernel with hardware detection. Device tree must be validated separately before kernel integration.

**Delivers:**
- NanoPi M6 device tree compiled from Armbian source
- Kernel build with correct DTB
- Hardware detection validated (eMMC, Ethernet, USB, NVMe)
- Kernel config verified against Armbian baseline

**Uses stack:**
- Collabora kernel 6.9.x
- Armbian `rk3588s-nanopi-m6.dts`
- Kernel config with required drivers (GMAC, sdhci-of-dwcmshc, nvme)

**Implements architecture:**
- artifacts/dtb/ component
- artifacts/nanopi-m6/u-boot/ with DTB integration

**Addresses pitfalls:**
- Device tree mismatch (use same DTS source for U-Boot and kernel)
- Kernel config missing drivers (validate against Armbian)

**Research needs:** LOW - Device tree compilation is standard pattern

### Phase 4: Overlay Integration
**Rationale:** With working hardware, create Talos overlay components. Installer and profile are templates with board-specific customization.

**Delivers:**
- installers/nanopi-m6/ Go binary
- profiles/nanopi-m6.yaml configuration
- Bootable Talos image (.raw.xz)
- Verified boot to Talos maintenance mode

**Implements architecture:**
- Installer implementing overlay.Installer interface
- Profile defining disk image parameters
- Final overlay OCI image

**Addresses pitfalls:**
- Install disk path hardcoding (configure explicitly in profile)
- Image size (verify fits on target eMMC)

**Research needs:** LOW - Template-based, established patterns

### Phase 5: Network and Cluster Integration
**Rationale:** After bootable image exists, establish cluster connectivity. Network configuration and MAC address handling are critical for Omni integration.

**Delivers:**
- Network configuration in machine config
- Static or persistent MAC address
- Talos machine appears in Omni
- Node successfully joins cluster
- Test workload validated

**Addresses pitfalls:**
- Random MAC address (configure static MAC)
- Network config lost after install (include in machine config)
- Talos upgrade path (document version requirements)

**Research needs:** LOW - Standard Talos configuration

### Phase 6: Production Hardening (Optional)
**Rationale:** After MVP cluster node works, add production features. Can be deferred if basic node is sufficient.

**Delivers:**
- CPU frequency scaling
- Hardware watchdog
- Thermal management
- NVMe performance tuning

**Research needs:** LOW - Kernel config changes

### Phase Ordering Rationale

- **Hardware-first approach:** Boot chain → Device tree → Kernel follows physical boot sequence. Each phase depends on previous working.
- **Fail-fast validation:** U-Boot phase catches defconfig issues before kernel complexity. Device tree phase isolates DTB problems before Talos integration.
- **Minimize rework:** Building kernel+modules together in Phase 1 avoids signing key issues. Single build pipeline prevents later module loading failures.
- **UART debugging throughout:** Serial console access established in Phase 1 enables debugging in Phases 2-4 without network dependency.

### Research Flags

**Phases needing deeper research:**
- **Phase 2 (Bootloader):** U-Boot defconfig adaptation for NanoPi M6 may require iterative debugging. Consider research-phase if R6C/R6S configs are insufficient. Research topics: RK3588S pin configurations, DRAM timing differences, peripheral initialization.

**Phases with standard patterns:**
- **Phase 1 (Environment):** Fork and build pipeline setup follows standard git/Docker patterns
- **Phase 3 (Device Tree):** DTS compilation and kernel integration is well-documented
- **Phase 4 (Overlay):** Installer/profile creation follows siderolabs templates
- **Phase 5 (Integration):** Talos machine configuration is standard

## Confidence Assessment

| Area | Confidence | Notes |
|------|------------|-------|
| Stack | HIGH | Verified with official Siderolabs repos, active milas community project, Collabora hardware enablement docs. All sources cross-referenced and consistent. |
| Features | HIGH | Table stakes validated against Talos overlay requirements and RK3588 hardware capabilities. MVP scope is well-defined (headless K8s node). |
| Architecture | HIGH | Official siderolabs/sbc-template and sbc-rockchip provide reference implementations. milas project confirms patterns work for RK3588. |
| Pitfalls | MEDIUM-HIGH | Top 5 pitfalls verified from GitHub issues and community reports. Phase-specific warnings derived from boot chain dependencies. Some hardware-specific issues (MAC address, UART baud rate) need validation on actual NanoPi M6 hardware. |

**Overall confidence:** HIGH

The research provides a clear technical path based on proven reference implementations (milas/talos-sbc-rk3588) and official Siderolabs tooling. The main uncertainties are hardware-specific (U-Boot defconfig, device tree validation) which require iterative testing but have documented recovery procedures (MaskROM, serial console).

### Gaps to Address

**Hardware validation gaps:**
- **U-Boot defconfig differences:** R6C/R6S configs may not perfectly match M6 hardware (different PMU, power rails, GPIO assignments). Mitigation: Compare with Armbian's M6 U-Boot config during Phase 2, test incrementally with UART access.
- **MAC address persistence:** Unknown if NanoPi M6 has MAC in efuse or requires software configuration. Mitigation: Test during Phase 5, document static MAC configuration in machine config.
- **eMMC performance tuning:** HS400 mode may require specific device tree properties. Mitigation: Benchmark storage I/O during Phase 3, compare with Armbian DTB if performance issues.

**Version upgrade path:**
- Talos v1.10.x is recommended but milas project is on v1.7.x. Upgrade compatibility needs validation. Mitigation: Test v1.10 overlay build during Phase 1, document any breaking changes from milas baseline.

**Kernel config validation:**
- milas acknowledges kernel config "probably has mistakes" from Armbian merge. Unknown which drivers are missing for NanoPi M6. Mitigation: Systematic peripheral testing during Phase 3, compare `dmesg` output with working Armbian system.

**Build environment:**
- Research assumes ARM64 native build environment or x86_64 with QEMU. Build times and cross-compilation issues not quantified. Mitigation: Establish build environment in Phase 1, document actual build times and resource requirements.

## Sources

### Primary (HIGH confidence)
- [Talos Linux Overlays Documentation](https://docs.siderolabs.com/talos/v1.10/build-and-extend-talos/custom-images-and-development/overlays) — Overlay architecture, build system requirements
- [siderolabs/overlays Repository](https://github.com/siderolabs/overlays) — Reference overlay implementations
- [siderolabs/sbc-rockchip](https://github.com/siderolabs/sbc-rockchip) — Official Rockchip overlay (23 boards)
- [siderolabs/sbc-template](https://github.com/siderolabs/sbc-template) — Overlay project template
- [siderolabs/bldr Repository](https://github.com/siderolabs/bldr) — Build system documentation
- [siderolabs/kres Repository](https://github.com/siderolabs/kres) — Build config generator
- [milas/talos-sbc-rk3588 Repository](https://github.com/milas/talos-sbc-rk3588) — Working RK3588 overlay (Rock 5A/5B), base for NanoPi M6
- [Armbian linux-rockchip Repository](https://github.com/armbian/linux-rockchip) — NanoPi M6 device tree source (rk-6.1-rkr3 branch)
- [Collabora RK3588 Upstream Status](https://www.collabora.com/news-and-blog/news-and-events/rockchip-rk3588-upstream-support-progress-future-plans.html) — Mainline kernel/U-Boot status
- [Collabora RK3588 Boot Chain](https://www.collabora.com/news-and-blog/blog/2024/02/21/almost-a-fully-open-source-boot-chain-for-rockchips-rk3588/) — BL31/U-Boot integration

### Secondary (MEDIUM confidence)
- [FriendlyELEC NanoPi M6 Wiki](https://wiki.friendlyelec.com/wiki/index.php/NanoPi_M6) — Hardware specifications, GPIO pinout
- [Armbian NanoPi M6 Support](https://github.com/armbian/build/pull/7763) — Community enablement efforts
- [milas/talos-sbc-rk3588 Issues](https://github.com/milas/talos-sbc-rk3588/issues) — Random MAC address, kernel config issues
- [CNX Software: RK3588 Mainline 2025](https://www.cnx-software.com/2024/12/21/rockchip-rk3588-mainline-linux-support-current-status-and-future-work-for-2025/) — Hardware enablement roadmap
- [U-Boot Rockchip Documentation](https://docs.u-boot.org/en/latest/board/rockchip/rockchip.html) — Boot process documentation
- [Rockchip Boot Options](https://opensource.rock-chips.com/wiki_Boot_option) — Boot ROM behavior, recovery procedures

### Tertiary (LOW confidence - verify during implementation)
- [RK3588 Cluster Boot Process](https://soliddowant.github.io/2024/01/23/rk3588-cluster-4) — Community boot chain guide
- [rk2aw Boot Order Tool](https://xnux.eu/rk2aw/) — Potential solution for boot priority issues

---
*Research completed: 2026-02-02*
*Ready for roadmap: yes*
