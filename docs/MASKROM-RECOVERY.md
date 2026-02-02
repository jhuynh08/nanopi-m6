# MaskROM Recovery Procedure

Recovery procedure for NanoPi M6 when SD card boot fails or the board appears bricked.

## Overview

MaskROM mode is the lowest-level recovery mode on Rockchip SoCs. When the RK3588S cannot find a valid bootloader on any boot media, it automatically enters MaskROM mode. You can also force MaskROM mode manually using the Mask button.

**When to use MaskROM:**
- Board appears completely dead (no LED activity)
- Need to erase/reflash eMMC boot sectors
- SD card boot fails and board is unresponsive
- Testing boot media priority issues

## Prerequisites

### Install rkdeveloptool (macOS)

```bash
# Install dependencies
brew install libusb

# Clone and build rkdeveloptool
git clone https://github.com/rockchip-linux/rkdeveloptool.git
cd rkdeveloptool
autoreconf -i
./configure
make
sudo make install

# Verify installation
rkdeveloptool --version
```

**Alternative: Pre-built binary**
```bash
# If available via Homebrew (check first)
brew install rkdeveloptool
```

### Download Loader File

You need the RK3588 SPL loader for USB communication:

```bash
# Clone rkbin repository (or use existing clone)
git clone https://github.com/rockchip-linux/rkbin.git

# Loader location
ls rkbin/bin/rk35/rk3588_spl_loader_*.bin

# Current recommended: rk3588_spl_loader_v1.16.112.bin
```

## Entering MaskROM Mode

### Method 1: Mask Button (Recommended)

1. **Power off** the board completely (remove power cable)
2. **Locate the Mask button** - small button near the eMMC socket
3. **Press and hold** the Mask button
4. **Connect power** while holding the button
5. **Hold for 3 seconds** after power LED lights up
6. **Release** the Mask button

### Method 2: No Bootloader Present

If both eMMC and SD card have no valid bootloader, the board automatically enters MaskROM mode when powered on.

### Method 3: Connect USB First (Alternative)

1. **Power off** the board
2. **Connect USB-C** from board to your Mac
3. **Press and hold** Mask button
4. **Connect power** while holding Mask
5. **Release** after 3 seconds

## Verify MaskROM Connection

```bash
# List connected Rockchip devices
rkdeveloptool list
```

**Expected output:**
```
DevNo=1 Vid=0x2207,Pid=0x350b,LocationID=xxx Maskrom
```

If you see "Maskrom" in the output, the connection is successful.

**If no device found:**
- Try a different USB cable (use data cable, not charge-only)
- Try different USB port on Mac
- Ensure board is in MaskROM mode (repeat entry steps)
- Check System Information > USB to see if device appears

## Recovery Operations

### Step 1: Download Loader

Before any operation, you must download the SPL loader:

```bash
rkdeveloptool db rkbin/bin/rk35/rk3588_spl_loader_v1.16.112.bin
```

**Expected output:**
```
Downloading bootance at 0x00000000...
Download bootance success!
```

### Step 2a: Erase eMMC Boot Sectors

Use this to clear a corrupted bootloader from eMMC:

```bash
# Erase first 32MB (covers all bootloader regions)
rkdeveloptool ef rkbin/bin/rk35/rk3588_spl_loader_v1.16.112.bin
```

Or erase specific sectors:
```bash
# Erase sectors 0-65536 (32MB)
rkdeveloptool es 0 65536
```

After erasing eMMC, the board will boot from SD card (if valid bootloader present).

### Step 2b: Flash Bootloader to eMMC

If you need to write a new bootloader to eMMC:

```bash
# Write u-boot-rockchip.bin at sector 64 (0x40)
rkdeveloptool wl 0x40 u-boot-rockchip.bin
```

### Step 2c: Flash Complete Image to eMMC

For flashing a complete disk image:

```bash
# Write image starting at sector 0
rkdeveloptool wl 0 /path/to/image.img
```

### Step 3: Reset Device

```bash
rkdeveloptool rd
```

The board will reboot and attempt normal boot sequence.

## Troubleshooting

| Symptom | Likely Cause | Solution |
|---------|--------------|----------|
| `rkdeveloptool list` shows nothing | Not in MaskROM mode | Retry entering MaskROM with Mask button |
| `rkdeveloptool list` shows nothing | Bad USB cable | Use a known-good data cable |
| "Download Boot Fail!" | Wrong loader file | Use rk3588_spl_loader_*.bin specifically |
| "Creating Comm Object failed!" | USB permission issue | Try with `sudo` |
| Device listed but commands fail | USB instability | Try different USB port, shorter cable |
| Device shows "Loader" not "Maskrom" | Board is in loader mode, not MaskROM | Power cycle and retry with Mask button |

## Safety Notes

1. **Never flash to eMMC during initial bring-up** - Use SD card only until boot is proven
2. **Keep Armbian SD card as recovery** - Verified working image from Phase 1
3. **Document every flash attempt** - Track configurations in boot test checklist
4. **Power supply matters** - Use 5V 4A supply, underpowered board may behave erratically
5. **MaskROM cannot brick** - This mode is in silicon, always accessible via Mask button

## Boot Media Priority

The RK3588S boot ROM has fixed priority:
1. SPI Flash (not used on NanoPi M6)
2. eMMC
3. SD Card

**Implication:** If eMMC has a bootloader (even broken), it will try eMMC before SD card. Use MaskROM to erase eMMC boot sectors if SD card testing is blocked.

## Quick Reference

```bash
# Full recovery sequence
rkdeveloptool list                                          # Verify connection
rkdeveloptool db rkbin/bin/rk35/rk3588_spl_loader_v1.16.112.bin  # Download loader
rkdeveloptool ef rkbin/bin/rk35/rk3588_spl_loader_v1.16.112.bin  # Erase eMMC
rkdeveloptool rd                                            # Reset

# Flash new bootloader to eMMC
rkdeveloptool list
rkdeveloptool db rkbin/bin/rk35/rk3588_spl_loader_v1.16.112.bin
rkdeveloptool wl 0x40 u-boot-rockchip.bin
rkdeveloptool rd
```

## Related Documentation

- [Flash Workflow](FLASH-WORKFLOW.md) - SD card flashing procedure
- [Boot Test Checklist](BOOT-TEST-CHECKLIST.md) - Systematic boot attempt tracking

---

*Based on: Rockchip documentation, FriendlyElec wiki, Radxa wiki*
