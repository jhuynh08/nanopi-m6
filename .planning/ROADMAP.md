# Roadmap: Talos Linux NanoPi M6 Image

## Overview

This roadmap delivers a bootable Talos Linux image for the NanoPi M6 (RK3588S) that joins an Omni-managed Kubernetes cluster. The journey follows hardware bring-up order: establish development environment, bring up bootloader (highest risk), enable kernel and device tree, create Talos overlay, integrate with cluster, then harden for production. Each phase builds on the previous - no phase can succeed without its dependencies working.

## Phases

**Phase Numbering:**
- Integer phases (1, 2, 3): Planned milestone work
- Decimal phases (2.1, 2.2): Urgent insertions (marked with INSERTED)

Decimal phases appear between their surrounding integers in numeric order.

- [ ] **Phase 1: Environment Setup** - Fork repo, verify build pipeline, establish UART debug access
- [ ] **Phase 2: Bootloader Bring-Up** - U-Boot defconfig, ATF, DDR blob - boot to U-Boot console
- [ ] **Phase 3: Device Tree & Kernel** - DTB compilation, kernel boot, essential driver validation
- [ ] **Phase 4: Overlay Integration** - Installer binary, profile YAML, bootable Talos image
- [ ] **Phase 5: Cluster Integration** - Boot Talos, network connectivity, Omni registration, cluster join
- [ ] **Phase 6: Production Hardening** - NVMe, thermal, watchdog, CPU frequency scaling

## Phase Details

### Phase 1: Environment Setup
**Goal**: Development environment enables building and debugging NanoPi M6 artifacts
**Depends on**: Nothing (first phase)
**Requirements**: None (foundation work enabling requirements)
**Success Criteria** (what must be TRUE):
  1. Forked repository builds successfully with bldr/kres toolchain
  2. UART console (1500000 baud) shows output from connected NanoPi M6
  3. Docker buildx produces ARM64 artifacts on development machine
  4. Build artifacts can be flashed to microSD for testing
**Plans**: TBD

Plans:
- [ ] 01-01: Fork and baseline upgrade
- [ ] 01-02: Build pipeline verification
- [ ] 01-03: UART console and flash workflow

### Phase 2: Bootloader Bring-Up
**Goal**: NanoPi M6 boots to U-Boot console via UART
**Depends on**: Phase 1
**Requirements**: BOOT-01, BOOT-02, BOOT-03
**Success Criteria** (what must be TRUE):
  1. U-Boot console prompt appears via UART after power-on
  2. DDR memory test passes (LPDDR5 initialized correctly)
  3. U-Boot can read files from microSD card
  4. Recovery procedure documented and tested (MaskROM mode)
**Plans**: TBD

Plans:
- [ ] 02-01: U-Boot defconfig adaptation
- [ ] 02-02: ATF and DDR blob integration
- [ ] 02-03: Boot validation and recovery procedures

### Phase 3: Device Tree & Kernel
**Goal**: Linux kernel boots with essential NanoPi M6 hardware functional
**Depends on**: Phase 2
**Requirements**: BOOT-04, KERN-01, KERN-02, KERN-03, KERN-04
**Success Criteria** (what must be TRUE):
  1. Linux kernel boots to console (dmesg visible via UART)
  2. Gigabit Ethernet interface appears (ip link shows eth0)
  3. eMMC storage detected and accessible (/dev/mmcblk* exists)
  4. USB host ports enumerate connected devices (lsusb works)
  5. Device tree correctly identifies board as NanoPi M6
**Plans**: TBD

Plans:
- [ ] 03-01: Device tree compilation from Armbian source
- [ ] 03-02: Kernel build with RK3588S support
- [ ] 03-03: Driver validation (Ethernet, eMMC, USB)

### Phase 4: Overlay Integration
**Goal**: Talos imager produces bootable NanoPi M6 image
**Depends on**: Phase 3
**Requirements**: OVRL-01, OVRL-02, OVRL-03, OVRL-04
**Success Criteria** (what must be TRUE):
  1. Installer binary implements overlay.Installer interface (compiles without error)
  2. Profile YAML passes Talos imager validation
  3. Overlay container image builds and pushes to registry
  4. Generated .raw.xz image boots to Talos maintenance mode
**Plans**: TBD

Plans:
- [ ] 04-01: Installer binary implementation
- [ ] 04-02: Profile YAML and overlay packaging
- [ ] 04-03: Image generation and boot validation

### Phase 5: Cluster Integration
**Goal**: NanoPi M6 joins Omni-managed Kubernetes cluster as functional node
**Depends on**: Phase 4
**Requirements**: CLST-01, CLST-02, CLST-03, CLST-04, CLST-05
**Success Criteria** (what must be TRUE):
  1. Device boots Talos Linux from eMMC without manual intervention
  2. Network connectivity established (can ping external hosts)
  3. Device appears in Talos Omni dashboard with correct name
  4. Node shows Ready status in kubectl get nodes
  5. Test pod scheduled on node runs successfully
**Plans**: TBD

Plans:
- [ ] 05-01: eMMC flash and boot configuration
- [ ] 05-02: Network and machine configuration
- [ ] 05-03: Omni registration and cluster join
- [ ] 05-04: Workload validation

### Phase 6: Production Hardening
**Goal**: Node operates reliably under production conditions
**Depends on**: Phase 5
**Requirements**: KERN-05, KERN-06, KERN-07, KERN-08
**Success Criteria** (what must be TRUE):
  1. NVMe drive detected and usable for persistent volumes
  2. CPU frequency scales dynamically under load (verify with cpufreq-info)
  3. Thermal sensors report temperature (thermal zone visible)
  4. Hardware watchdog triggers reboot on simulated hang
**Plans**: TBD

Plans:
- [ ] 06-01: NVMe enablement and testing
- [ ] 06-02: Power and thermal management
- [ ] 06-03: Watchdog configuration and validation

## Progress

**Execution Order:**
Phases execute in numeric order: 1 -> 2 -> 3 -> 4 -> 5 -> 6

| Phase | Plans Complete | Status | Completed |
|-------|----------------|--------|-----------|
| 1. Environment Setup | 0/3 | Not started | - |
| 2. Bootloader Bring-Up | 0/3 | Not started | - |
| 3. Device Tree & Kernel | 0/3 | Not started | - |
| 4. Overlay Integration | 0/3 | Not started | - |
| 5. Cluster Integration | 0/4 | Not started | - |
| 6. Production Hardening | 0/3 | Not started | - |

---
*Roadmap created: 2026-02-02*
*Last updated: 2026-02-02*
