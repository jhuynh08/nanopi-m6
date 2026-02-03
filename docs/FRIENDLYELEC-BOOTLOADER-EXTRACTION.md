# FriendlyELEC Bootloader Extraction Procedure

Extract the bootloader from FriendlyELEC's official Ubuntu image to test on our SD card.

## Purpose

After 4 failed boot attempts with our U-Boot build (all configuration matching FriendlyELEC's exactly), this diagnostic test determines whether the issue is:
- **Our build process** (if FriendlyELEC's binary boots successfully)
- **SD card boot path** (if FriendlyELEC's binary also fails)

## Source Image

| Property | Value |
|----------|-------|
| Image | rk3588-sd-ubuntu-noble-minimal-6.1-arm64-20251222.img.gz |
| Source | FriendlyELEC official release |
| OS | Ubuntu 24.04 Noble Minimal |
| Kernel | 6.1 LTS |
| Architecture | ARM64 |
| Board | NanoPi M6 (RK3588S) |
| File Size | ~1.17GB compressed |

**Download:** [FriendlyELEC Downloads](https://drive.google.com/drive/folders/1tOXhfO_PSq7rk_ZS0_lh3M_r4kHJAFfq) or [wiki.friendlyelec.com](https://wiki.friendlyelec.com/wiki/index.php/NanoPi_M6)

## Bootloader Location (Rockchip Standard)

On RK3588/RK3588S devices, the bootloader occupies:

| Component | Sector Start | Byte Offset | Size |
|-----------|--------------|-------------|------|
| idbloader.img (TPL+SPL) | 64 | 32 KB | ~512 KB |
| u-boot.itb (U-Boot + DTB) | 16384 | 8 MB | ~1-2 MB |
| **Total bootloader region** | 64 | 32 KB | ~10 MB |

The combined `u-boot-rockchip.bin` contains both idbloader and u-boot.itb.

## Extraction Procedure

### Step 1: Decompress the Image

```bash
# Decompress (keeps original .gz file intact with -k)
gunzip -k /Users/johnsonhuynh/Downloads/rk3588-sd-ubuntu-noble-minimal-6.1-arm64-20251222.img.gz

# Verify decompressed file exists
ls -lh /Users/johnsonhuynh/Downloads/rk3588-sd-ubuntu-noble-minimal-6.1-arm64-20251222.img
# Expected: ~7-8GB uncompressed
```

### Step 2: Extract Bootloader Binary

```bash
# Extract bootloader region (sectors 64-20544, covering ~10MB)
# bs=512 (sector size)
# skip=64 (start at sector 64, where idbloader begins)
# count=20480 (extract 20480 sectors = 10,485,760 bytes = 10MB)

dd if=/Users/johnsonhuynh/Downloads/rk3588-sd-ubuntu-noble-minimal-6.1-arm64-20251222.img \
   of=friendlyelec-bootloader.bin \
   bs=512 skip=64 count=20480

# Verify extraction
ls -la friendlyelec-bootloader.bin
# Expected: exactly 10,485,760 bytes (10MB)
```

### Step 3: Verify Extraction

```bash
# Check file size
stat -f%z friendlyelec-bootloader.bin
# Expected: 10485760

# Verify idbloader signature (first bytes should be "3b 8c...")
# The RK3588 idbloader has a specific magic header
xxd friendlyelec-bootloader.bin | head -5
```

**Expected signature pattern:**
- First bytes contain Rockchip idbloader identification
- Non-zero data throughout (not empty sectors)

### Step 4: (Optional) Compare with Our Build

```bash
# Compare file sizes
ls -la friendlyelec-bootloader.bin _out/u-boot-rockchip.bin

# Binary comparison (will show differences)
cmp friendlyelec-bootloader.bin _out/u-boot-rockchip.bin || echo "Files differ"

# Size comparison
echo "FriendlyELEC: $(stat -f%z friendlyelec-bootloader.bin) bytes"
echo "Our build:    $(stat -f%z _out/u-boot-rockchip.bin 2>/dev/null || echo 'not found') bytes"
```

## Flash Procedure

### CRITICAL: Identify Correct Device

```bash
# List all disks
diskutil list

# Look for your SD card (usually /dev/disk2 or /dev/disk3)
# Identify by size (e.g., 32GB) and partition layout
# NEVER flash to disk0 or disk1 (internal drives)
```

**Example output for SD card:**
```
/dev/disk2 (external, physical):
   #:                       TYPE NAME                    SIZE       IDENTIFIER
   0:     FDisk_partition_scheme                        *31.9 GB    disk2
   1:             Windows_NTFS   boot                    268.4 MB   disk2s1
   2:                    Linux                           31.6 GB    disk2s2
```

### Flash Bootloader to SD Card

```bash
# Unmount SD card (do NOT eject)
diskutil unmountDisk /dev/diskN

# Flash ONLY bootloader (preserves existing partitions)
# Using rdisk for faster raw writes (macOS specific)
# seek=64 writes at sector 64 (matching extraction offset)
sudo dd if=friendlyelec-bootloader.bin of=/dev/rdiskN bs=512 seek=64 status=progress

# Ensure data is written
sync

# Eject safely
diskutil eject /dev/diskN
```

**Replace `/dev/rdiskN` with your actual device (e.g., `/dev/rdisk2`)**

## Boot Test Procedure

1. Insert SD card into NanoPi M6
2. Connect HDMI monitor
3. Connect Ethernet cable
4. Power on the board
5. Observe for 2 minutes:
   - Watch for LED activity (SYS LED blinking)
   - Watch for HDMI output
   - Check for network activity (ping, DHCP lease)

## Expected Results

| If Result | Conclusion | Next Action |
|-----------|------------|-------------|
| **BOOTS** (LED activity, network up) | Issue is in **our build process** | Compare build configuration, toolchain, flags |
| **FAILS** (same as Attempts #1-4) | Issue is in **SD card boot path** | Test eMMC boot via MaskROM mode |

## Cleanup

After testing, optionally remove decompressed image to save space:

```bash
# Remove decompressed image (keeps original .gz)
rm /Users/johnsonhuynh/Downloads/rk3588-sd-ubuntu-noble-minimal-6.1-arm64-20251222.img

# Keep or remove extracted bootloader as needed
rm friendlyelec-bootloader.bin
```

## Safety Notes

1. **Triple-check device identifier** before `dd` command
2. **Never flash to disk0/disk1** (internal macOS drives)
3. **This only replaces bootloader sectors** - existing partitions are preserved
4. **Keep FriendlyELEC SD card as backup** - verified working baseline from Phase 1

## Related Documentation

- [MaskROM Recovery](MASKROM-RECOVERY.md) - For eMMC testing if SD boot fails
- [Boot Test Checklist](BOOT-TEST-CHECKLIST.md) - Document Attempt #5 results
- [Flash Workflow](FLASH-WORKFLOW.md) - General SD card flashing

---

*Procedure created: 2026-02-03*
*Source: FriendlyELEC official Ubuntu 24.04 Noble image*
