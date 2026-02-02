# Requirements: Talos Linux NanoPi M6 Image

**Defined:** 2026-02-02
**Core Value:** The NanoPi M6 boots Talos Linux and successfully registers with Talos Omni to join the home cluster.

## v1 Requirements

Requirements for initial release. Each maps to roadmap phases.

### Boot Chain

- [ ] **BOOT-01**: U-Boot bootloader with NanoPi M6 defconfig boots to console
- [ ] **BOOT-02**: ARM Trusted Firmware (BL31) loads successfully
- [ ] **BOOT-03**: DDR training blob initializes LPDDR5 memory
- [ ] **BOOT-04**: Device tree blob correctly describes NanoPi M6 hardware

### Kernel & Drivers

- [ ] **KERN-01**: Linux kernel boots with RK3588S support
- [ ] **KERN-02**: Gigabit Ethernet driver (GMAC) enables network connectivity
- [ ] **KERN-03**: eMMC storage driver (sdhci-of-dwcmshc) enables boot media access
- [ ] **KERN-04**: USB host driver enables peripheral connectivity
- [ ] **KERN-05**: NVMe driver enables M.2 slot for storage
- [ ] **KERN-06**: CPU frequency scaling (cpufreq-dt) enables dynamic frequency
- [ ] **KERN-07**: Thermal management (thermal-tsadc) prevents overheating
- [ ] **KERN-08**: Hardware watchdog (rk3588-wdt) enables crash recovery

### Talos Overlay

- [ ] **OVRL-01**: Go installer binary implements overlay.Installer interface
- [ ] **OVRL-02**: Board profile YAML defines disk image parameters
- [ ] **OVRL-03**: Overlay builds successfully with bldr/kres toolchain
- [ ] **OVRL-04**: Talos imager generates bootable .raw.xz image

### Cluster Integration

- [ ] **CLST-01**: Device boots Talos Linux from eMMC
- [ ] **CLST-02**: Network connectivity established (Gigabit Ethernet)
- [ ] **CLST-03**: Device appears in Talos Omni dashboard
- [ ] **CLST-04**: Device successfully joins home cluster as node
- [ ] **CLST-05**: Workloads can be scheduled on the node

## v2 Requirements

Deferred to future release. Tracked but not in current roadmap.

### Extended Hardware

- **EXT-01**: RTC support (hym8563) for accurate time without network
- **EXT-02**: PMIC support (rk806) for clean shutdown/reboot
- **EXT-03**: GPIO access for external sensors/actuators

### Maintenance

- **MAINT-01**: Automated CI/CD builds on new Talos releases
- **MAINT-02**: Upstream contribution to milas/talos-sbc-rk3588

## Out of Scope

Explicitly excluded. Documented to prevent scope creep.

| Feature | Reason |
|---------|--------|
| HDMI/Display output | Headless server use case, adds complexity |
| GPU acceleration | Not needed for Kubernetes node workloads |
| Audio support | Not needed for server use case |
| WiFi/Bluetooth | Using wired Gigabit Ethernet |
| NPU (AI accelerator) | Requires out-of-tree rknn kernel module, not needed |
| VPU (video encode/decode) | No media workloads planned, drivers incomplete |
| Upstream mainline patches | Personal use only, not pursuing kernel/U-Boot upstreaming |

## Traceability

Which phases cover which requirements. Updated during roadmap creation.

| Requirement | Phase | Status |
|-------------|-------|--------|
| BOOT-01 | Phase 2 | Pending |
| BOOT-02 | Phase 2 | Pending |
| BOOT-03 | Phase 2 | Pending |
| BOOT-04 | Phase 3 | Pending |
| KERN-01 | Phase 3 | Pending |
| KERN-02 | Phase 3 | Pending |
| KERN-03 | Phase 3 | Pending |
| KERN-04 | Phase 3 | Pending |
| KERN-05 | Phase 6 | Pending |
| KERN-06 | Phase 6 | Pending |
| KERN-07 | Phase 6 | Pending |
| KERN-08 | Phase 6 | Pending |
| OVRL-01 | Phase 4 | Pending |
| OVRL-02 | Phase 4 | Pending |
| OVRL-03 | Phase 4 | Pending |
| OVRL-04 | Phase 4 | Pending |
| CLST-01 | Phase 5 | Pending |
| CLST-02 | Phase 5 | Pending |
| CLST-03 | Phase 5 | Pending |
| CLST-04 | Phase 5 | Pending |
| CLST-05 | Phase 5 | Pending |

**Coverage:**
- v1 requirements: 21 total
- Mapped to phases: 21
- Unmapped: 0

**Phase 1 Note:** Environment Setup has no requirements - it delivers the foundation (fork, build pipeline, UART access) that enables all subsequent requirement work.

---
*Requirements defined: 2026-02-02*
*Last updated: 2026-02-02 after roadmap creation*
