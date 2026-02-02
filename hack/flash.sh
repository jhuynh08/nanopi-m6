#!/usr/bin/env bash
#
# Flash an image to SD card (macOS)
# Usage: ./hack/flash.sh <image-file> [disk-number] [--bootloader]
#
# Examples:
#   ./hack/flash.sh _out/metal-arm64.raw 2              # Full disk image at sector 0
#   ./hack/flash.sh u-boot-rockchip.bin 2 --bootloader  # Bootloader at sector 64
#

set -euo pipefail

IMAGE="${1:-}"
DISK_NUM="${2:-}"
BOOTLOADER_MODE="${3:-}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

error() { echo -e "${RED}ERROR: $*${NC}" >&2; exit 1; }
warn() { echo -e "${YELLOW}WARNING: $*${NC}" >&2; }
info() { echo -e "${GREEN}$*${NC}"; }

# Validate image file
[[ -z "$IMAGE" ]] && error "Usage: $0 <image-file> [disk-number]"
[[ ! -f "$IMAGE" ]] && error "Image file not found: $IMAGE"

# Show available disks
echo "Available disks:"
diskutil list external

# Get or confirm disk number
if [[ -z "$DISK_NUM" ]]; then
    echo ""
    read -p "Enter disk number (e.g., 2 for /dev/disk2): " DISK_NUM
fi

DISK="/dev/disk${DISK_NUM}"
RDISK="/dev/rdisk${DISK_NUM}"

# Verify disk exists
[[ ! -e "$DISK" ]] && error "Disk not found: $DISK"

# Get disk info for confirmation
DISK_SIZE=$(diskutil info "$DISK" | grep "Disk Size" | awk '{print $3, $4}')
DISK_NAME=$(diskutil info "$DISK" | grep "Device / Media Name" | cut -d: -f2 | xargs)

# Safety check: refuse if disk looks like internal
if diskutil info "$DISK" | grep -q "Internal"; then
    error "Refusing to write to internal disk: $DISK"
fi

# Confirm before proceeding
echo ""
echo "================================================"
echo "Target: $DISK ($DISK_SIZE)"
echo "Name: $DISK_NAME"
echo "Image: $IMAGE ($(du -h "$IMAGE" | cut -f1))"
if [[ "$BOOTLOADER_MODE" == "--bootloader" ]]; then
    echo "Mode: BOOTLOADER (writing at sector 64)"
else
    echo "Mode: FULL DISK (writing at sector 0)"
fi
echo "================================================"
echo ""

if [[ "$BOOTLOADER_MODE" == "--bootloader" ]]; then
    warn "This will OVERWRITE the bootloader area on $DISK"
else
    warn "This will ERASE ALL DATA on $DISK"
fi
read -p "Type 'yes' to continue: " CONFIRM

[[ "$CONFIRM" != "yes" ]] && error "Aborted by user"

# Unmount disk
info "Unmounting $DISK..."
diskutil unmountDisk "$DISK" || warn "Unmount failed, continuing anyway"

# Flash image
if [[ "$BOOTLOADER_MODE" == "--bootloader" ]]; then
    # Rockchip bootloader: write at sector 64 (32KB offset)
    info "Flashing bootloader at sector 64..."
    echo "Command: sudo dd if=$IMAGE of=$RDISK seek=64 bs=512 status=progress"
    sudo dd if="$IMAGE" of="$RDISK" seek=64 bs=512 status=progress
else
    info "Flashing full disk image (this may take several minutes)..."
    echo "Command: sudo dd if=$IMAGE of=$RDISK bs=1m status=progress"
    sudo dd if="$IMAGE" of="$RDISK" bs=1m status=progress
fi

# Sync and eject
info "Syncing..."
sync

info "Ejecting $DISK..."
diskutil eject "$DISK"

info "Done! SD card is ready."
echo ""
echo "Next steps:"
echo "1. Insert SD card into NanoPi M6"
echo "2. Connect HDMI and/or network"
echo "3. Power on and watch for:"
echo "   - LED activity (SYS LED should blink)"
echo "   - HDMI output (U-Boot/kernel messages)"
echo "   - Network (ping after ~60-90 seconds)"
