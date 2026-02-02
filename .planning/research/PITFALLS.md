# Domain Pitfalls

**Domain:** Talos Linux Custom SBC Image (NanoPi M6 / RK3588S)
**Researched:** 2026-02-02
**Confidence:** MEDIUM (WebSearch verified with multiple sources, GitHub issues, and official documentation)

## Critical Pitfalls

Mistakes that cause rewrites, bricked boards, or major project setbacks.

### Pitfall 1: Kernel Module Signing Key Mismatch

**What goes wrong:** Kernel modules built separately from the Talos kernel fail to load at runtime with cryptographic verification errors.

**Why it happens:** Talos requires kernel modules to be signed with a key that's only available during the Talos kernel build process. Using pre-built modules from siderolabs/extensions or building modules against a different kernel results in signature mismatches.

**Consequences:**
- Modules silently fail to load
- Critical drivers (storage, network) unavailable
- Board appears to boot but functionality is broken
- Hours of debugging "why doesn't my NVMe work"

**Warning signs:**
- `dmesg` shows module loading failures with signature errors
- Mixing extensions from different repos (e.g., siderolabs extensions with milas kernel)
- Building kernel and extensions in separate CI jobs

**Prevention:**
- Build kernel and all kernel modules in the same build pipeline with the same signing key
- Use the exact extensions from the milas/talos-sbc-rk3588 repository (or fork)
- Never mix siderolabs kernel modules with custom RK3588 kernel builds
- Verify module signatures before flashing: check that `*.ko` files have proper ELF signatures

**Recovery if it happens:**
- Rebuild the entire image with kernel + modules in one pass
- If using milas/talos-sbc-rk3588 as base, ensure you're pulling extensions from the same release tag

**Phase mapping:** Phase 1 (Environment Setup) - establish build pipeline that enforces this constraint from the start

---

### Pitfall 2: Missing or Incorrect U-Boot Defconfig

**What goes wrong:** The NanoPi M6 has no mainline U-Boot defconfig. Using an incorrect defconfig results in boot failures at the SPL/TPL stage.

**Why it happens:** NanoPi M6 is relatively new; only NanoPi R6C/R6S have upstream defconfigs. Pin configurations, DRAM timing, and peripheral initialization differ between boards even when using the same RK3588S SoC.

**Consequences:**
- Board doesn't boot past BROM (no serial output after initial boot ROM)
- DDR initialization fails silently
- Board appears dead
- eMMC may be left in inconsistent state

**Warning signs:**
- Using `generic-rk3588_defconfig` without modifications
- No UART output after flashing
- Boot ROM falls back to MaskROM mode

**Prevention:**
- Start with `nanopi-r6c-rk3588s_defconfig` or `nanopi-r6s-rk3588s_defconfig` as base
- Cross-reference with Armbian's U-Boot config for NanoPi M6
- ALWAYS have UART serial console connected for initial bring-up
- Test U-Boot independently before integrating with Talos
- Use Collabora's U-Boot fork (not mainline) until mainline has complete RK3588 support

**Recovery if it happens:**
- Use MaskROM mode to reflash (hold recovery button, or short eMMC clock to ground)
- Boot from SD card to recover eMMC
- Use `rkdeveloptool` to flash via USB

**Phase mapping:** Phase 2 (Bootloader) - dedicated phase for U-Boot bring-up before kernel work

---

### Pitfall 3: Device Tree Mismatch Between U-Boot and Kernel

**What goes wrong:** U-Boot loads one device tree, kernel expects another, resulting in kernel panic or missing peripherals.

**Why it happens:** U-Boot and Linux kernel can have different device tree sources. Rockchip's DTB structure evolved differently in vendor trees vs. mainline vs. Collabora's fork.

**Consequences:**
- Kernel panic at boot: "Unable to mount root fs"
- Ethernet, eMMC, or SD not detected
- Random MAC addresses on each boot
- Peripherals work in U-Boot but not in Linux (or vice versa)

**Warning signs:**
- Using DTB from different source than kernel (e.g., Armbian DTB with Collabora kernel)
- Device tree compatible strings don't match
- Kernel config doesn't enable drivers for devices in DTB
- U-Boot shows correct hardware but Linux doesn't

**Prevention:**
- Use device tree from the SAME kernel source tree you're building
- Copy Armbian's NanoPi M6 DTS into Collabora's kernel tree and adapt
- Verify compatible strings: `"friendlyelec,nanopi-m6", "rockchip,rk3588s"`
- Test DTB loading explicitly: `dtc -I dtb -O dts` to decompile and inspect
- Ensure kernel config enables all drivers referenced in DTB

**Recovery if it happens:**
- Boot with known-good DTB from Armbian to identify working configuration
- Compare DTB nodes between working and failing configurations
- Add `earlycon=uart8250,mmio32,0xfeb50000` to kernel cmdline for early debug output

**Phase mapping:** Phase 3 (Device Tree) - dedicated phase after bootloader is working

---

### Pitfall 4: Boot Media Priority and Flash Location Confusion

**What goes wrong:** Flashing to wrong offset or wrong media; RK3588's fixed boot order (SPI -> eMMC -> SD) causes unexpected boot behavior.

**Why it happens:** RK3588 boot ROM has hardcoded boot order that cannot be changed. Multiple boot media (SPI flash, eMMC, SD card) can all contain bootloaders, leading to confusion about which one is active.

**Consequences:**
- Flashing SD card but board boots from old eMMC image
- "Successfully" flashing but no change in behavior
- Corrupting existing working bootloader
- Brick requiring MaskROM recovery

**Warning signs:**
- Flashing to `/dev/mmcblk1` when board boots from `/dev/mmcblk0`
- Not zeroing existing boot media before testing new image
- SPI flash contains old U-Boot that takes priority

**Prevention:**
- ALWAYS identify boot media with `lsblk` before flashing
- Know the mapping: eMMC is typically `mmcblk0`, SD is `mmcblk1` (but verify!)
- Zero out SPI flash if not using it: prevents priority conflicts
- Document your boot media configuration in project README
- Use correct offset for Rockchip: sector 64 (0x8000 bytes / 32KB)

**Recovery if it happens:**
- Use rk2aw loader to reverse boot priority (SD first)
- Enter MaskROM mode to reflash via USB
- Boot from SD card and zero eMMC boot sectors

**Phase mapping:** Phase 2 (Bootloader) - establish clear flash procedures early

---

### Pitfall 5: Closed-Source Binary Version Mismatches (TPL/BL31/DDR)

**What goes wrong:** Mixing incompatible versions of TPL (DRAM init), BL31 (TF-A), and U-Boot results in boot failures or subtle hardware issues.

**Why it happens:** RK3588 requires closed-source DDR training blob. Different versions have different DRAM timing parameters. TF-A (BL31) and U-Boot must be compatible versions.

**Consequences:**
- DDR initialization fails (no output, board appears dead)
- Memory corruption / instability under load
- CPU frequency/voltage issues
- Cryptographic features disabled (wrong TF-A version)

**Warning signs:**
- Mixing rkbin files from different dates/versions
- Using mainline TF-A (doesn't support RK3588 yet)
- DDR blob version doesn't match your RAM type (LP4 vs LP5)
- Build warnings about missing BL31 or TPL

**Prevention:**
- Use EXACT versions from milas/talos-sbc-rk3588 or document your tested combination
- NanoPi M6 uses LPDDR5 - ensure DDR blob supports LP5 (e.g., `rk3588_ddr_lp4_2112MHz_lp5_2736MHz_v1.08.bin`)
- Pin rkbin commit hash in your build scripts
- Use Collabora's TF-A fork until mainline merges RK3588 support
- Set environment variables explicitly: `export ROCKCHIP_TPL=... export BL31=...`

**Recovery if it happens:**
- Identify which binary is wrong from UART output
- Return to known-good combination from reference project
- Use MaskROM mode if board won't boot

**Phase mapping:** Phase 2 (Bootloader) - validate binary combination before kernel work

---

## Moderate Pitfalls

Mistakes that cause delays, rework, or technical debt.

### Pitfall 6: Random MAC Address on Every Boot

**What goes wrong:** Ethernet interface gets a different MAC address each boot, breaking DHCP reservations, Kubernetes node identity, and Omni registration.

**Why it happens:** RK3588 boards may not have MAC address burned into efuse. U-Boot reads invalid MAC (00:00:00:00:00:00) and generates random one. Different images may store MAC differently.

**Consequences:**
- Node gets different IP on each reboot
- Talos Omni sees multiple registrations for same physical board
- Kubernetes node identity issues
- Network automation breaks

**Warning signs:**
- `dmesg` shows "setting random MAC address"
- MAC address changes between reboots
- DHCP server shows multiple leases for same hostname

**Prevention:**
- Check if your NanoPi M6 has MAC in efuse (FriendlyELEC may have programmed it)
- Configure static MAC in Talos machine config if not
- Use `ethaddr` in U-Boot environment
- Document MAC address handling in your build process

**Recovery if it happens:**
- Set static MAC address in Talos machine configuration
- Use `machine.network.interfaces[].permanentAddr` in machine config
- For multiple boards, maintain MAC address inventory

**Phase mapping:** Phase 4 (Kernel/Image) - address during Talos image configuration

---

### Pitfall 7: Talos Version Upgrade Path Violations

**What goes wrong:** Attempting to upgrade Talos across multiple minor versions fails; overlay/kernel incompatibilities cause boot loops.

**Why it happens:** Talos only supports upgrades between adjacent minor versions. Custom SBC overlays must be rebuilt for each Talos version. Kernel module ABIs change between versions.

**Consequences:**
- Upgrade fails, node stuck in boot loop
- Automatic rollback works but node is on old version
- Cluster has mixed Talos versions
- Omni shows unhealthy node state

**Warning signs:**
- Skipping Talos versions (e.g., 1.7 -> 1.9 directly)
- Using overlay built for different Talos version
- Not testing upgrade path before production

**Prevention:**
- Track upstream Talos releases; rebuild overlay for each minor version
- Test upgrade path: fresh install -> upgrade -> verify
- Document which overlay version matches which Talos version
- Use Talos 1.7+ overlay system (not pre-1.7 image factory approach)
- Keep `--talos-version` flag consistent across all generated configs

**Recovery if it happens:**
- Node will auto-rollback to previous working image
- Rebuild overlay for target Talos version
- Upgrade incrementally through each minor version

**Phase mapping:** Phase 5 (Integration) - establish upgrade testing before Omni registration

---

### Pitfall 8: Install Disk Path Hardcoding

**What goes wrong:** Talos installer targets wrong disk; installation goes to SD card instead of eMMC or vice versa.

**Why it happens:** Default Talos config assumes `/dev/sda`. ARM SBCs use `/dev/mmcblk0` (eMMC) or `/dev/mmcblk1` (SD). NVMe is `/dev/nvme0n1`. Path depends on kernel probe order.

**Consequences:**
- Installation appears to succeed but boots old image
- Data written to wrong storage device
- Confusion about which media contains the OS

**Warning signs:**
- Using default `controlplane.yaml` without disk customization
- Installing to production board without verifying disk path
- Different behavior between development and production hardware

**Prevention:**
- Always specify `machine.install.disk` explicitly in machine config
- Verify disk paths on actual hardware with `lsblk` before generating config
- For NanoPi M6: eMMC is likely `/dev/mmcblk0`, SD is `/dev/mmcblk1`
- Test installation path with `talosctl disks --insecure`

**Recovery if it happens:**
- Edit machine config to specify correct disk
- Re-apply configuration with correct path
- Boot from correct media and fix persistent config

**Phase mapping:** Phase 4 (Kernel/Image) - verify disk paths during image testing

---

### Pitfall 9: Network Configuration Lost After Install

**What goes wrong:** Node boots into Talos but is unreachable; network configuration from installer doesn't persist.

**Why it happens:** Talos's kernel cmdline `ip=` parameter is only for installer boot. Permanent network config must be in machine configuration. Missing network config = no connectivity after reboot.

**Consequences:**
- Node boots but cannot reach API
- Cannot apply configuration
- Must re-flash or use serial console
- Omni cannot communicate with node

**Warning signs:**
- Relying on DHCP without verifying it works post-install
- Not including `machine.network` in configuration
- Testing only with installer boot, not installed system

**Prevention:**
- ALWAYS include network configuration in machine config
- Use static IP or ensure DHCP reservation exists
- Test network connectivity after install AND after reboot
- Include `machine.network.interfaces[]` even for DHCP

**Recovery if it happens:**
- Connect via serial console
- Apply network configuration via `talosctl apply-config --insecure --nodes <serial-console>`
- Re-flash with correct machine config embedded

**Phase mapping:** Phase 5 (Integration) - critical for Omni connectivity

---

### Pitfall 10: Kernel Config Missing Required Drivers

**What goes wrong:** Hardware works in Armbian but not in Talos; kernel config doesn't enable RK3588-specific drivers.

**Why it happens:** milas/talos-sbc-rk3588 kernel config was manually merged from Armbian with "probably some mistakes." Some drivers may be missing or misconfigured for NanoPi M6.

**Consequences:**
- Ethernet not working (Realtek RTL8125 driver)
- eMMC not detected
- USB ports non-functional
- Thermal management disabled

**Warning signs:**
- `lspci`/`lsusb` doesn't show expected devices
- `dmesg` shows "no driver found" messages
- Hardware works in Armbian but not Talos

**Prevention:**
- Compare kernel config with working Armbian config
- Ensure these are enabled for NanoPi M6:
  - `CONFIG_ROCKCHIP_DW_HDMI` (if needed)
  - `CONFIG_MMC_SDHCI_OF_DWCMSHC` (eMMC)
  - `CONFIG_STMMAC_ETH` (built-in Ethernet)
  - `CONFIG_R8169` or `CONFIG_R8125` (if Realtek NIC)
- Test each peripheral after kernel build
- Review milas kernel config against device tree requirements

**Recovery if it happens:**
- Identify missing driver from `dmesg` or device tree
- Add driver to kernel config
- Rebuild kernel and overlay

**Phase mapping:** Phase 3 (Device Tree) and Phase 4 (Kernel) - verify after DTB is working

---

## Minor Pitfalls

Mistakes that cause annoyance but are recoverable.

### Pitfall 11: UART Console Not Configured

**What goes wrong:** No serial output during boot; unable to debug early boot failures.

**Why it happens:** Wrong UART port configured, baud rate mismatch, or console not enabled in kernel cmdline.

**Prevention:**
- NanoPi M6 uses UART2 at 1500000 baud (verify with FriendlyELEC documentation)
- Add `console=ttyS2,1500000n8` to kernel cmdline
- Test serial cable before flashing custom images

**Phase mapping:** Phase 1 (Environment) - set up before any board work

---

### Pitfall 12: Image Size Too Large for Boot Media

**What goes wrong:** Built image doesn't fit on target SD card or eMMC.

**Why it happens:** Debug symbols, extra modules, or uncompressed kernel bloat image size.

**Prevention:**
- Strip debug symbols in release builds
- Use kernel compression (gzip or zstd)
- Test on target storage media before finalizing
- NanoPi M6 has 32GB eMMC option - verify your target

**Phase mapping:** Phase 4 (Kernel/Image) - check during image build

---

### Pitfall 13: Forgetting to Update Overlay for New Boards

**What goes wrong:** Adding NanoPi M6 support but overlay installer doesn't recognize the new board.

**Why it happens:** Overlay system uses installer profiles to match board. New boards need profile entries.

**Prevention:**
- Add NanoPi M6 to installer profile list
- Ensure `get-options` returns valid board identifier
- Test overlay detection: `talosctl get platformstatus`

**Phase mapping:** Phase 4 (Kernel/Image) - part of installer profile work

---

## Phase-Specific Warnings

| Phase | Likely Pitfall | Mitigation |
|-------|---------------|------------|
| 1. Environment Setup | UART not configured | Set up serial console FIRST |
| 2. Bootloader | Defconfig wrong, binary version mismatch | Use reference project binaries, test U-Boot standalone |
| 3. Device Tree | DTB source mismatch | Use DTB from same kernel tree |
| 4. Kernel/Image | Module signing, missing drivers | Build kernel + modules together, compare with Armbian config |
| 5. Integration | Network config lost, MAC randomization | Test full boot cycle, configure static MAC/IP |
| 6. Omni Registration | Upgrade path issues | Test upgrade before production |

---

## Sources

### Official Documentation
- [Talos Overlays Documentation](https://www.talos.dev/v1.10/advanced/overlays/)
- [Talos Adding Kernel Modules](https://www.talos.dev/v1.11/advanced/kernel-module/)
- [U-Boot Rockchip Documentation](https://docs.u-boot.org/en/latest/board/rockchip/rockchip.html)
- [Rockchip Boot Options](https://opensource.rock-chips.com/wiki_Boot_option)

### GitHub Issues & Discussions
- [milas/talos-sbc-rk3588 Random MAC Address Issue](https://github.com/milas/talos-sbc-rk3588/issues/2)
- [Community Managed SBCs - siderolabs/talos#8065](https://github.com/siderolabs/talos/issues/8065)
- [RK3588 MAC Address Persistence - radxa/u-boot#34](https://github.com/radxa/u-boot/issues/34)
- [Talos SBC Upgrade Issues](https://github.com/siderolabs/talos/discussions/10447)

### Community Resources
- [milas/talos-sbc-rk3588 GitHub Repository](https://github.com/milas/talos-sbc-rk3588)
- [Armbian NanoPi M6 Support](https://github.com/armbian/build/pull/7763)
- [Collabora RK3588 Mainline Status](https://www.cnx-software.com/2024/12/21/rockchip-rk3588-mainline-linux-support-current-status-and-future-work-for-2025/)
- [RK3588 Cluster Boot Process Guide](https://soliddowant.github.io/2024/01/23/rk3588-cluster-4)
- [rk2aw Boot Order Solution](https://xnux.eu/rk2aw/)

### Confidence Notes
- Pitfalls 1-5: HIGH confidence - documented in multiple sources
- Pitfalls 6-10: MEDIUM confidence - derived from issue reports and documentation
- Pitfalls 11-13: MEDIUM confidence - common patterns, verify against your hardware
