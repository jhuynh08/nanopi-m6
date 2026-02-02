# Phase 2: Bootloader Bring-Up - Research

**Researched:** 2026-02-02
**Domain:** U-Boot bootloader for RK3588S / NanoPi M6, ATF (BL31), DDR blobs
**Confidence:** MEDIUM (verified with official docs but HDMI output is a significant limitation)

## Summary

This phase involves bringing up U-Boot on the NanoPi M6 (RK3588S) to the point where it boots successfully and produces visible output. The research reveals a **critical constraint**: mainline U-Boot does NOT have video/HDMI driver support for RK3588. There is no VOP2 display driver in U-Boot, meaning **U-Boot itself will not display anything on HDMI** regardless of configuration. The user's constraint of "HDMI only, no UART" makes this phase significantly more challenging because boot success/failure must be inferred from downstream behavior (kernel booting) rather than U-Boot console output.

The standard approach is to adapt the existing `nanopi-r6c-rk3588s_defconfig` or use Armbian's `nanopi-m6-rk3588s_defconfig` as a base. The NanoPi M6 uses LPDDR5 memory, requiring the correct DDR blob version. For blind iteration without UART, success will only be confirmed when Linux kernel eventually boots and HDMI output becomes active (kernel 6.15+ required for mainline HDMI).

**Primary recommendation:** Use Armbian's `nanopi-m6-rk3588s_defconfig` as the starting point since it's already tested on the exact hardware. Accept that U-Boot stage will be "blind" (no HDMI output) and verify success through eventual kernel boot or LED activity patterns.

## Standard Stack

### Core Components

| Component | Version | Purpose | Why Standard |
|-----------|---------|---------|--------------|
| **U-Boot** | v2024.01+ (Collabora fork) | Bootloader | Collabora fork has RK3588 support; milas project uses this |
| **U-Boot** | v2025.01+ (mainline) | Bootloader | NanoPi R6C/R6S defconfigs in v2024.10+; Armbian uses v2025.01+ |
| **ARM Trusted Firmware (BL31)** | v1.45+ elf | Secure boot firmware | Open-source TF-A merged for RK3588; required for boot chain |
| **DDR Blob** | v1.16+ bin | LPDDR5 memory training | Closed-source, unavoidable; NanoPi M6 uses LPDDR5 |
| **rkbin** | Pinned commit | Binary blobs repository | Source for BL31 and DDR blob files |

### Blob Files (Current in Project)

| File | Path in rkbin | Notes |
|------|---------------|-------|
| **DDR Blob** | `bin/rk35/rk3588_ddr_lp4_2112MHz_lp5_2400MHz_v1.16.bin` | Supports LPDDR5 at 2400MHz |
| **BL31** | `bin/rk35/rk3588_bl31_v1.45.elf` | TF-A ARM Trusted Firmware |

### Latest Available (If Update Needed)

| Component | Latest Version | Notes |
|-----------|---------------|-------|
| DDR Blob | v1.19 (2025-03-13) | Supports RK3588S-B variants |
| BL31 | v1.51 (2025-04-25) | Added DDR debug interface |

### Defconfig Options

| Option | Source | Status |
|--------|--------|--------|
| `nanopi-m6-rk3588s_defconfig` | Armbian U-Boot patches | Tested on NanoPi M6 hardware |
| `nanopi-r6c-rk3588s_defconfig` | Mainline U-Boot v2024.10+ | Similar RK3588S board, good starting point |
| `nanopi-r6s-rk3588s_defconfig` | Mainline U-Boot v2024.10+ | Similar RK3588S board |
| `rock5a-rk3588s_defconfig` | Mainline U-Boot | Rock 5A is RK3588S, current project uses this |

**Recommendation:** Start with Armbian's `nanopi-m6-rk3588s_defconfig` (source it from Armbian build system patches), as it's specifically tested for NanoPi M6 hardware.

## Architecture Patterns

### Boot Chain Architecture

```
Power On
    |
    v
+-------------------+
| RK3588S Boot ROM  |  (Hardcoded in silicon)
+-------------------+
    | Loads from: SPI Flash -> eMMC -> SD Card (fixed order)
    v
+-------------------+
| TPL (DDR Init)    |  <- rk3588_ddr_lp4_2112MHz_lp5_2400MHz_v1.16.bin
+-------------------+
    | DDR training complete
    v
+-------------------+
| SPL (U-Boot SPL)  |  <- Part of idbloader.img
+-------------------+
    | Loads BL31 and U-Boot
    v
+-------------------+
| BL31 (TF-A)       |  <- rk3588_bl31_v1.45.elf
+-------------------+
    | Secure world setup
    v
+-------------------+
| U-Boot Proper     |  <- u-boot.itb
+-------------------+
    | NO HDMI OUTPUT (no video driver)
    | Loads kernel/DTB
    v
+-------------------+
| Linux Kernel      |  <- HDMI works here (6.15+)
+-------------------+
```

### File Layout on SD Card

```
Sector 64 (0x8000 / 32KB):
    +---------------------------+
    | idbloader.img             |  TPL + SPL combined
    +---------------------------+

Sector 16384 (0x800000 / 8MB):
    +---------------------------+
    | u-boot.itb                |  U-Boot proper + DTB + ATF
    +---------------------------+

Partition 1 (typically starts at sector 32768 / 16MB):
    +---------------------------+
    | Boot partition (FAT/ext4) |  Kernel, initramfs, etc.
    +---------------------------+
```

### Project Structure for NanoPi M6 U-Boot

```
artifacts/u-boot/
├── pkg.yaml                    # Aggregates all U-Boot builds
├── prepare/
│   └── pkg.yaml                # Downloads Collabora U-Boot source
├── rock5a/
│   └── pkg.yaml                # Existing Rock 5A build
├── rock5b/
│   └── pkg.yaml                # Existing Rock 5B build
└── nanopi-m6/                  # NEW: NanoPi M6 build
    └── pkg.yaml                # New defconfig for NanoPi M6
```

### Pattern: Adding New Board Build

**pkg.yaml for new board:**
```yaml
name: u-boot-nanopi-m6
variant: scratch
shell: /toolchain/bin/bash
dependencies:
  - stage: rkbin
  - stage: u-boot-prepare
steps:
  - env:
      SOURCE_DATE_EPOCH: {{ .BUILD_ARG_SOURCE_DATE_EPOCH }}
      ROCKCHIP_TPL: /libs/rkbin/bin/rk35/rk3588_ddr_lp4_2112MHz_lp5_2400MHz_v1.16.bin
      BL31: /libs/rkbin/bin/rk35/rk3588_bl31_v1.45.elf
    prepare:
      - |
        cd /src
        # Option 1: Use mainline defconfig as base
        make nanopi-r6c-rk3588s_defconfig
        # Option 2: Create custom defconfig with patches
        # Apply patches for NanoPi M6 specific settings
    build:
      - |
        cd /src
        make -j $(nproc) HOSTLDLIBS_mkimage="-lssl -lcrypto"
    install:
      - |
        mkdir -p /rootfs/artifacts/arm64/u-boot/nanopi-m6
        cp -v /src/u-boot-rockchip.bin /rootfs/artifacts/arm64/u-boot/nanopi-m6
finalize:
  - from: /rootfs
    to: /rootfs
```

### Anti-Patterns to Avoid

- **Expecting HDMI output from U-Boot:** Mainline U-Boot has NO video driver for RK3588. Don't waste time trying to enable splash screens.
- **Using vendor U-Boot (v2017.09):** FriendlyElec's vendor U-Boot is ancient and incompatible with Talos overlay system.
- **Mixing blob versions:** DDR blob and BL31 must be compatible versions. Use tested combinations.
- **Flashing to eMMC during bring-up:** SD card only until boot is proven. eMMC is harder to recover.

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| DDR initialization | Custom timing code | rkbin DDR blob | Closed-source, reverse engineering not feasible |
| Secure boot chain | Custom TF-A port | rkbin BL31 or TF-A mainline | Complex, already done |
| Boot image packing | Manual dd offsets | `make u-boot-rockchip.bin` | U-Boot build creates combined image |
| Defconfig from scratch | Manual Kconfig | Base on existing R6C/R6S | Too many options, easy to break boot |
| HDMI output in U-Boot | Custom video driver | Accept no output | VOP2 driver not ported to U-Boot |

**Key insight:** The RK3588 boot chain requires specific binary blobs in specific formats. The U-Boot build system handles combining TPL, SPL, BL31, and U-Boot proper into the correct images. Don't try to manually construct boot images.

## Common Pitfalls

### Pitfall 1: Expecting HDMI Output from U-Boot

**What goes wrong:** User expects to see U-Boot logo/text on HDMI display
**Why it happens:** Mainline U-Boot does NOT have a VOP2 video driver for RK3588. There is no framebuffer, no splash screen capability.
**How to avoid:** Accept that U-Boot stage is "blind". Success is verified by:
1. LED activity (if board has programmable LEDs)
2. Kernel eventually booting and showing HDMI output
3. Network ping after expected boot time
**Warning signs:** Searching for CONFIG_VIDEO_ROCKCHIP for RK3588 - it doesn't exist in mainline

### Pitfall 2: Wrong DDR Blob for LPDDR5

**What goes wrong:** Board fails to boot at all - no LED activity, appears completely dead
**Why it happens:** NanoPi M6 uses LPDDR5 memory. Using a DDR blob that only supports LPDDR4 will fail
**How to avoid:** Use blob that explicitly mentions LP5: `rk3588_ddr_lp4_2112MHz_lp5_2400MHz_v1.16.bin`
**Warning signs:** Filename without "lp5" in it

### Pitfall 3: Wrong Defconfig Selection

**What goes wrong:** U-Boot compiles but fails to boot on NanoPi M6
**Why it happens:** Different RK3588S boards have different:
- Pin muxing configurations
- Memory configurations
- Peripheral connections
**How to avoid:**
- First choice: Use Armbian's `nanopi-m6-rk3588s_defconfig`
- Second choice: Use `nanopi-r6c-rk3588s_defconfig` (same SoC variant)
- Avoid: `rock5a-rk3588s_defconfig` unless testing generic config
**Warning signs:** Using generic-rk3588_defconfig or configs for different SoC variant (RK3588 vs RK3588S)

### Pitfall 4: Boot Media Priority Confusion

**What goes wrong:** Flash new image to SD, board still boots old image from eMMC
**Why it happens:** RK3588 boot ROM has fixed priority: SPI Flash -> eMMC -> SD Card
**How to avoid:**
- During Phase 2: Only use SD card, ensure eMMC has no bootloader (or erase boot sectors)
- Check boot source with MaskROM mode if confused
**Warning signs:** Changes to SD card image have no effect on boot behavior

### Pitfall 5: Incorrect Flash Offset

**What goes wrong:** Bootloader image written but board enters MaskROM mode
**Why it happens:** Rockchip expects bootloader at sector 64 (32KB offset), not sector 0
**How to avoid:** Always use: `dd if=u-boot-rockchip.bin of=/dev/sdX seek=64`
**Warning signs:** Using seek=0 or no seek parameter

### Pitfall 6: BL31/DDR Version Mismatch

**What goes wrong:** Boot hangs after DDR init, before U-Boot proper
**Why it happens:** BL31 and DDR blob have compatibility requirements
**How to avoid:** Use tested combinations from rkbin release notes:
- DDR v1.16 works with BL31 v1.45
- Check rkbin/doc/release/RK3588_EN.md for compatibility matrix
**Warning signs:** Mixing blobs from different rkbin commits

## Code Examples

### Flashing U-Boot to SD Card

```bash
# Source: Rockchip documentation, verified for RK3588
# IMPORTANT: seek=64 is sectors (512 bytes each) = 32KB offset

# Identify SD card (VERIFY THIS IS CORRECT!)
diskutil list  # macOS
# Look for external disk matching SD card size

# Unmount (not eject!)
diskutil unmountDisk /dev/disk2

# Write bootloader at correct offset
sudo dd if=u-boot-rockchip.bin of=/dev/rdisk2 seek=64 bs=512 status=progress

# Sync and safely eject
sync
diskutil eject /dev/disk2
```

### Building U-Boot with Correct Environment

```bash
# Source: U-Boot Rockchip documentation
# Environment variables must be set before make

export ROCKCHIP_TPL=/path/to/rkbin/bin/rk35/rk3588_ddr_lp4_2112MHz_lp5_2400MHz_v1.16.bin
export BL31=/path/to/rkbin/bin/rk35/rk3588_bl31_v1.45.elf

# Configure for target board
make nanopi-r6c-rk3588s_defconfig
# Or with custom defconfig:
# make nanopi-m6-rk3588s_defconfig

# Build
make -j$(nproc) CROSS_COMPILE=aarch64-linux-gnu-

# Output: u-boot-rockchip.bin (combined image ready to flash)
```

### MaskROM Recovery Procedure

```bash
# Source: Rockchip documentation, Radxa wiki
# Use when board won't boot and needs reflash via USB

# 1. Enter MaskROM mode:
#    - Power off the board
#    - Press and hold the "Mask" button (near eMMC)
#    - Connect USB-C to PC
#    - Wait for LED to light, release after 3 seconds

# 2. Verify connection
rkdeveloptool list
# Should show: DevNo=1 Vid=0x2207,Pid=0x350b,LocationID=xxx Maskrom

# 3. Download loader (required before any operations)
rkdeveloptool db rk3588_spl_loader.bin

# 4. Flash bootloader to eMMC (if needed)
rkdeveloptool wl 0x40 u-boot-rockchip.bin  # 0x40 = sector 64

# 5. Reset device
rkdeveloptool rd
```

### Verification Without UART

```bash
# Source: Practical experience documented in Phase 1 research

# Blind boot verification checklist:

# 1. Prepare timing
START_TIME=$(date +%s)

# 2. Insert SD card, power on
# Watch for any LED activity in first 10 seconds
# - LED blinks: SPL/U-Boot is running (good sign)
# - No LED activity after 10s: Likely DDR/SPL failure

# 3. Wait for kernel boot (up to 2 minutes)
# HDMI will only work once Linux kernel loads (6.15+)
# If display shows anything: SUCCESS
sleep 120

# 4. Check network (if DHCP configured)
# Look for new device in router DHCP leases
# Or scan network:
nmap -sP 192.168.1.0/24

# 5. Decision tree:
# - HDMI shows kernel/Talos: SUCCESS
# - HDMI blank but network responds: Partial success (display config issue)
# - No HDMI, no network, no LED: Boot failure - try different config
```

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| FriendlyElec vendor U-Boot v2017.09 | Collabora/Mainline U-Boot 2024.01+ | 2024 | Modern boot features, mainline compatibility |
| Closed-source BL31 from vendor | Open-source TF-A | TF-A v2.12 | Fully auditable secure boot (except DDR) |
| No mainline NanoPi support | R6C/R6S in mainline v2024.10 | Oct 2024 | Direct upstream support, easier updates |
| No M6 defconfig | Armbian patches for M6 | 2024-2025 | Working defconfig available |

**Deprecated/outdated:**
- **FriendlyElec vendor U-Boot (v2017.09):** Ancient, missing modern features, incompatible with Talos
- **generic-rk3588_defconfig:** Too generic, use board-specific configs
- **Manual boot image construction:** U-Boot build system handles this now

## Open Questions

### 1. HDMI Output Without UART

- **What we know:** Mainline U-Boot has NO video driver for RK3588. HDMI will not display anything until Linux kernel boots.
- **What's unclear:** Whether there's any visual indication of boot progress before kernel
- **Recommendation:** Accept blind U-Boot stage. Verify via:
  - LED activity (board-specific, need to document)
  - Kernel HDMI output (requires kernel 6.15+)
  - Network connectivity after expected boot time
  - Fallback to Armbian SD card to verify hardware works

### 2. Armbian Defconfig Sourcing

- **What we know:** Armbian has `nanopi-m6-rk3588s_defconfig` that works
- **What's unclear:** Exact patch/defconfig content, how to integrate with milas build system
- **Recommendation:**
  1. First try: Extract defconfig from Armbian build system patches
  2. Fallback: Use `nanopi-r6c-rk3588s_defconfig` from mainline
  3. Document any modifications needed

### 3. LED Behavior on NanoPi M6

- **What we know:** Board has LEDs that can indicate status
- **What's unclear:** Which LEDs are controllable from U-Boot/SPL, what patterns indicate success/failure
- **Recommendation:** Document observed LED behavior during first boot attempt

### 4. Iteration Strategy for Boot Failures

- **What we know:** Without UART, debugging is limited
- **What's unclear:** How many config variations to try before escalating
- **Recommendation:**
  1. Armbian defconfig (3 attempts with different blob versions)
  2. R6C defconfig (2 attempts)
  3. If still failing: Consider if UART acquisition is necessary
  4. Park phase if 5+ attempts fail with no progress indicators

## Sources

### Primary (HIGH confidence)
- [U-Boot Rockchip Documentation](https://docs.u-boot.org/en/latest/board/rockchip/rockchip.html) - Official defconfig list, build instructions
- [rkbin Release Notes (RK3588_EN.md)](https://github.com/rockchip-linux/rkbin/blob/master/doc/release/RK3588_EN.md) - DDR/BL31 version compatibility
- [FriendlyElec NanoPi M6 Wiki](https://wiki.friendlyelec.com/wiki/index.php/NanoPi_M6) - UART settings, MaskROM procedure
- [Armbian Build PR #7652](https://github.com/armbian/build/pull/7652) - NanoPi M6 SPI flash support, defconfig names
- [Armbian Build PR #7341](https://github.com/armbian/build/pull/7341) - Initial NanoPi M6 support

### Secondary (MEDIUM confidence)
- [Collabora RK3588 Upstream Status](https://www.collabora.com/news-and-blog/news-and-events/rockchip-rk3588-upstream-support-progress-future-plans.html) - Boot chain progress, HDMI status
- [CNX Software RK3588 Mainline Status 2025](https://www.cnx-software.com/2024/12/21/rockchip-rk3588-mainline-linux-support-current-status-and-future-work-for-2025/) - Overall RK3588 support timeline
- [Radxa Rock5 USB Install Guide](https://wiki.radxa.com/Rock5/install/usb-install-emmc) - MaskROM recovery procedures
- [U-Boot NanoPi R6C Patch](https://www.mail-archive.com/u-boot@lists.denx.de/msg514307.html) - Mainline R6C support details

### Tertiary (LOW confidence - verify during implementation)
- [ROCKNIX rk3588-uboot](https://github.com/ROCKNIX/rk3588-uboot) - Alternative U-Boot fork, potential video patches
- [edk2-rk3588](https://github.com/edk2-porting/edk2-rk3588) - UEFI alternative (not compatible with Talos but good reference)

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH - Official U-Boot docs, rkbin repository
- Architecture patterns: HIGH - Documented boot chain, verified flash offsets
- Pitfalls: MEDIUM - Some NanoPi M6 specific issues need validation
- HDMI limitation: HIGH - Verified no VOP2 driver in mainline U-Boot defconfigs

**Critical constraint confirmed:** Mainline U-Boot has NO HDMI/video output support for RK3588. User's "HDMI only" constraint means Phase 2 success can only be verified indirectly through kernel boot or LED activity.

**Research date:** 2026-02-02
**Valid until:** 2026-03-02 (30 days - U-Boot/rkbin move slowly)

---

## Recommendations for Claude's Discretion Items

### Base Defconfig Selection

**Recommendation:** Use `nanopi-r6c-rk3588s_defconfig` from mainline U-Boot as starting point

**Rationale:**
1. NanoPi R6C uses same RK3588S SoC as NanoPi M6
2. Defconfig is in mainline U-Boot v2024.10+
3. Armbian's `nanopi-m6-rk3588s_defconfig` would require extracting patches from Armbian build system
4. R6C is simpler to integrate with existing milas build system
5. Collabora U-Boot fork (used by milas) likely has R6C defconfig

**If R6C fails:**
1. Try `nanopi-r6s-rk3588s_defconfig` (also RK3588S)
2. Investigate extracting Armbian's M6 defconfig patches
3. Create minimal delta from R6C for M6-specific differences

### DDR Blob Source Selection

**Recommendation:** Use rkbin repository (current project approach)

**Rationale:**
1. Project already pins rkbin commit: `a2a0b89b6c8c612dca5ed9ed8a68db8a07f68bc0`
2. Current DDR blob `v1.16` supports LPDDR5
3. FriendlyElec BSP would require additional integration work
4. rkbin has documented release notes and compatibility info

**Blob versions to use:**
- DDR: `rk3588_ddr_lp4_2112MHz_lp5_2400MHz_v1.16.bin` (current)
- BL31: `rk3588_bl31_v1.45.elf` (current)

### Iteration Strategy When Configs Fail

**Recommendation:** 3-tier approach with clear exit criteria

**Tier 1: Quick iterations (3 attempts max)**
1. R6C defconfig with current blobs
2. R6C defconfig with updated blobs (v1.19 DDR, v1.51 BL31)
3. R6S defconfig as alternative

**Tier 2: Configuration investigation (2 attempts max)**
4. Extract Armbian M6 defconfig, analyze differences
5. Apply Armbian patches to R6C defconfig

**Tier 3: Decision point**
- If 5 attempts fail with no LED/network activity: Consider parking phase
- If LED activity observed but no kernel: Continue with kernel debugging
- Document all attempts with exact configs used

### Config File Organization

**Recommendation:** Create new directory, don't patch existing

**Structure:**
```
artifacts/u-boot/
├── nanopi-m6/
│   └── pkg.yaml      # New file, mirrors rock5a structure
```

**Rationale:**
1. Cleaner separation of concerns
2. Existing rock5a/rock5b structure is clean template
3. Easy to maintain and update independently
4. Consistent with project patterns

### Recovery Procedure Documentation

**Recommendation:** Create comprehensive MaskROM guide in docs/

**Content:**
1. Prerequisites (rkdeveloptool installation)
2. Entering MaskROM mode (step-by-step with timing)
3. Verifying connection
4. Flashing recovery image
5. Returning to normal boot
6. Troubleshooting common issues

### Test Checklist Format

**Recommendation:** Simple markdown checklist with timing

```markdown
## Boot Attempt #N - [Date] [Time]

### Configuration
- Defconfig: [name]
- DDR blob: [version]
- BL31: [version]

### Observation Timeline
- [ ] 0-10s: LED activity observed? [Y/N]
- [ ] 10-30s: Any LED patterns? [describe]
- [ ] 30-120s: HDMI output? [Y/N/partial]
- [ ] 120s+: Network ping response? [Y/N]

### Result
- [ ] SUCCESS: Kernel booted, HDMI works
- [ ] PARTIAL: Some activity but no full boot
- [ ] FAILURE: No activity observed

### Notes
[Observations, next steps]
```
