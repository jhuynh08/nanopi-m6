# Flash and Verification Workflow

This document describes how to flash images to SD card and verify boot on NanoPi M6 hardware.

## Prerequisites

- macOS with SD card reader (or USB adapter)
- microSD card (32GB+ recommended, Samsung EVO or SanDisk Extreme preferred)
- NanoPi M6 board with power supply
- HDMI monitor and/or network connection for verification

## Flashing Images

### Using the Flash Script (Recommended)

```bash
# Flash a built image
./hack/flash.sh _out/metal-arm64.raw

# Or specify disk number directly (skip interactive prompt)
./hack/flash.sh _out/metal-arm64.raw 2
```

The script will:
1. Show available external disks
2. Ask for confirmation before writing
3. Refuse to write to internal disks
4. Use fast raw device (rdisk) for better performance

### Manual Flashing

If you prefer manual steps:

```bash
# 1. List disks to find your SD card
diskutil list external

# 2. Unmount (not eject!) the SD card
diskutil unmountDisk /dev/disk2

# 3. Flash using raw device for speed
sudo dd if=_out/metal-arm64.raw of=/dev/rdisk2 bs=1m status=progress

# 4. Sync and eject
sync
diskutil eject /dev/disk2
```

## Boot Verification (Without UART)

Since we don't have UART initially, use these methods to verify boot:

### 1. LED Indicators

| Time | LED Behavior | Meaning |
|------|-------------|---------|
| 0-5s | SYS LED blinks once | Power on, BL31/SPL starting |
| 5-10s | SYS LED blinking | U-Boot running |
| 10-30s | Rapid LED activity | Kernel booting |
| 30s+ | Steady or periodic blink | System running |

**No LED activity after 10 seconds = bootloader failed**

### 2. HDMI Output

Connect HDMI before powering on. You should see:
- U-Boot logo/text (within 5 seconds)
- Kernel boot messages
- Talos maintenance mode or login prompt

### 3. Network Ping

After ~60-90 seconds:
1. Check router DHCP leases for new device
2. Ping the assigned IP address
3. Try SSH if Talos is running: `talosctl --nodes <ip> version`

### 4. Built-in LCD (After Kernel Boots)

If you have the FriendlyElec 2.1" LCD touchscreen attached:
- Display will show boot progress once kernel display driver loads
- Useful for visual confirmation even without HDMI

## Troubleshooting

### Board appears dead (no LEDs)
- Check power supply (must be 5V 4A minimum for M6)
- Try different power cable
- Verify SD card is fully inserted
- Reflash SD card with known-good image

### LEDs blink but no HDMI
- Kernel/DTB mismatch likely
- U-Boot may be working but kernel fails
- Try different HDMI cable/monitor
- Check if Armbian baseline works

### HDMI shows kernel panic
- Driver or configuration issue
- Note the panic message for debugging
- Try Armbian baseline to isolate issue

### Network not reachable
- Check physical cable connection
- Verify router DHCP is enabled
- May need static IP configuration

## Baseline Testing

Before testing custom images, verify hardware with Armbian:

1. Download Armbian for NanoPi M6 from https://www.armbian.com/nanopi-m6/
2. Flash using the script: `./hack/flash.sh Armbian_*.img`
3. Boot and verify HDMI + network work
4. This proves hardware is functional

If Armbian works but custom image doesn't, the issue is in our image.
If Armbian also fails, the issue is hardware or SD card.

## SD Card Recommendations

For best reliability:
- **Samsung EVO Plus** - Fast, reliable, good endurance
- **SanDisk Extreme** - High write speeds, good for frequent reflashing
- Avoid: No-name brands, cards older than 3 years

Minimum class: UHS-I U1 (10 MB/s write)
Recommended: UHS-I U3 or better

## Quick Reference

```bash
# Flash image
./hack/flash.sh <image-file>

# Check network after boot
ping <ip-from-dhcp-leases>

# Verify Talos (once running)
talosctl --nodes <ip> version
talosctl --nodes <ip> dmesg
```
