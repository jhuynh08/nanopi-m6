# Phase 3: Device Tree & Kernel - Research

**Researched:** 2026-02-03
**Domain:** Linux kernel device tree compilation and driver enablement for RK3588S SBC
**Confidence:** MEDIUM-HIGH

## Summary

This research investigates what is needed to boot a Linux kernel with essential hardware support on the NanoPi M6 (RK3588S). The phase requires device tree compilation, kernel configuration verification, and validation of three critical drivers: Ethernet (GMAC), USB host, and NVMe (for root filesystem since this unit has no eMMC).

The standard approach leverages the existing Talos SBC RK3588 project structure, which uses the Collabora Linux kernel fork (6.9.x) with RK3588 patches. The device tree for NanoPi M6 must be sourced from Armbian's linux-rockchip repository, which has board-specific support contributed by @efectn. The kernel config already includes NVMe, STMMAC Ethernet, and USB host support built-in, requiring only verification that PCIe lanes are properly configured in the device tree.

**Primary recommendation:** Add the NanoPi M6 device tree from Armbian's U-Boot patches to the existing kernel build, verify PCIe2x1 lane configuration for NVMe, and validate driver bring-up using HDMI console output.

## Standard Stack

The established tools and sources for this domain:

### Core
| Component | Version | Source | Purpose |
|-----------|---------|--------|---------|
| Linux Kernel | 6.9.x | Collabora gitlab (ref: 23bb9c65...) | RK3588 mainline support with patches |
| Device Tree | rk3588s-nanopi-m6.dts | Armbian build repo | Board-specific hardware description |
| DTC | In toolchain | Kernel build | Device tree compilation |
| Kernel Config | config-arm64 | artifacts/kernel/mainline/build/ | RK3588 kernel options |

### Supporting
| Component | Version | Purpose | When to Use |
|-----------|---------|---------|-------------|
| rk3588s.dtsi | Mainline | SoC base device tree include | Always included via #include |
| rk3588s-nanopi-r6.dtsi | Upstream | Similar board reference | Reference for M.2/PCIe config |
| FriendlyELEC kernel-rockchip | nanopi6-v6.1.y | Vendor kernel reference | Cross-check pinctrl if issues |

### Alternatives Considered
| Instead of | Could Use | Tradeoff |
|------------|-----------|----------|
| Armbian DTS | FriendlyELEC vendor kernel DTS | Vendor may have more features but less mainline compatibility |
| Collabora kernel | Rockchip BSP kernel (5.10/6.1) | BSP has more hardware support but drifts from mainline |

**No installation needed** - kernel sources fetched during Docker build via pkg.yaml.

## Architecture Patterns

### Recommended Project Structure
```
artifacts/
├── kernel/
│   └── mainline/
│       ├── build/
│       │   ├── config-arm64       # Kernel config (already exists)
│       │   └── pkg.yaml
│       ├── kernel/
│       │   └── pkg.yaml           # Installs dtb to /rootfs/dtb
│       └── prepare/
│           └── pkg.yaml           # Downloads kernel source
└── dts/                           # NEW: Board-specific DTS files
    └── rk3588s-nanopi-m6.dts      # NanoPi M6 device tree

installers/
└── rk3588/
    └── src/
        └── main.go                # Update for NanoPi M6 DTB path
```

### Pattern 1: Device Tree Inclusion
**What:** NanoPi M6 DTS includes the base rk3588s.dtsi and adds board-specific nodes
**When to use:** Always - this is how Rockchip device trees work
**Example:**
```dts
// Source: Armbian build repo rk3588s-nanopi-m6.dts
/dts-v1/;

#include "dt-bindings/gpio/gpio.h"
#include "dt-bindings/leds/common.h"
#include "dt-bindings/pinctrl/rockchip.h"
#include "dt-bindings/usb/pd.h"
#include "rk3588s.dtsi"

/ {
    model = "FriendlyElec NanoPi M6";
    compatible = "friendlyelec,nanopi-m6", "rockchip,rk3588s";

    aliases {
        ethernet0 = &gmac1;
        mmc0 = &sdhci;
        mmc1 = &sdmmc;
    };
    // ... board-specific nodes
};
```

### Pattern 2: PCIe Configuration for NVMe
**What:** Enable PCIe2x1 lane for M.2 NVMe slot
**When to use:** This device uses M.2 for root filesystem (no eMMC)
**Example:**
```dts
// Source: rk3588s-nanopi-r6.dtsi reference for PCIe config
&pcie2x1l1 {
    reset-gpios = <&gpio1 RK_PA7 GPIO_ACTIVE_HIGH>;
    vpcie3v3-supply = <&vcc_3v3_pcie20>;
    status = "okay";
};
```

### Pattern 3: Kernel DTB Installation Path
**What:** DTBs are installed to /rootfs/dtb/rockchip/ during kernel build
**When to use:** Always - Talos installer expects DTBs at specific path
**Example from artifacts/kernel/mainline/kernel/pkg.yaml:**
```yaml
cd ./arch/arm64/boot/dts
for vendor in $(find . -not -path . -type d); do
  dest="/rootfs/dtb/$vendor"
  mkdir -v $dest
  find ./$vendor/* -type f -name "*.dtb" -exec cp {} $dest \;
done
```

### Anti-Patterns to Avoid
- **Modifying rk3588s.dtsi directly:** Never edit the SoC include file; override in board DTS
- **Disabling nodes in wrong scope:** Use status = "disabled" in board DTS, not dtsi
- **PCIe without power supply:** Always define vpcie3v3-supply or device won't initialize
- **Missing reset-gpios:** NVMe may appear but fail without proper reset GPIO config

## Don't Hand-Roll

Problems that have existing solutions:

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| Device tree for NanoPi M6 | Write from scratch | Armbian rk3588s-nanopi-m6.dts | Board-specific pinctrl, regulators already defined |
| Kernel config for RK3588 | Start from defconfig | Existing config-arm64 | Already tuned for Talos with RK3588 support |
| DTB compilation | Manual dtc invocation | Kernel build (`make dtbs`) | Handles includes, dependencies properly |
| PHY driver for ethernet | Custom MDIO code | STMMAC + RTL8211F in-tree | Tested, maintained upstream drivers |

**Key insight:** The existing Talos SBC RK3588 project already has the kernel config with all required drivers. The only new work is adding the NanoPi M6 device tree from Armbian.

## Common Pitfalls

### Pitfall 1: PCIe BAR Allocation Failures
**What goes wrong:** dmesg shows "BAR 1 fails to assign" or "can't assign; no space"
**Why it happens:** PCIe memory regions not properly sized in device tree or UEFI/bootloader
**How to avoid:** Verify pcie@fe150000 (pcie2x1l0) and pcie@fe160000 (pcie2x1l1) have proper ranges defined
**Warning signs:** NVMe detected but I/O errors; device appears then disappears

### Pitfall 2: Ethernet TX Issues with RTL8211F
**What goes wrong:** Link up but no traffic transmits; ping fails; DHCP times out
**Why it happens:** Incorrect RX/TX delay configuration in device tree
**How to avoid:** Use known-working values from Armbian DTS:
```dts
&gmac1 {
    clock_in_out = "output";
    phy-mode = "rgmii-rxid";
    tx_delay = <0x42>;
    // rx_delay = <0x00>; // RX delay disabled for RK3588
    status = "okay";
};
```
**Warning signs:** Link LED on, but dmesg shows "stmmac_tx_timeout" or zero TX packets

### Pitfall 3: USB Port Not Enumerating Devices
**What goes wrong:** USB ports powered but devices not detected
**Why it happens:** USB PHY or USB controller not enabled; power regulator disabled
**How to avoid:** Verify these nodes are status = "okay":
- usb_host0_ehci, usb_host0_ohci, usb_host0_xhci
- u2phy0, u2phy0_host
- vcc5v0_host_20 regulator
**Warning signs:** No USB-related messages in dmesg; no /dev/bus/usb

### Pitfall 4: NVMe Not Detected Despite PCIe Enabled
**What goes wrong:** PCIe controller initializes but no NVMe device appears
**Why it happens:** M.2 slot uses different PCIe lane than expected; power not applied
**How to avoid:**
1. Verify which PCIe controller the M.2 slot connects to (pcie2x1l1 or pcie2x1l2)
2. Ensure vcc_3v3_pcie20 regulator is enabled
3. Check reset-gpios matches actual board routing
**Warning signs:** lspci shows no devices; dmesg shows PCIe link training failure

### Pitfall 5: Wrong Device Tree Selected at Boot
**What goes wrong:** Kernel boots but wrong hardware configured
**Why it happens:** DTB filename in U-Boot doesn't match compiled DTB; extlinux.conf wrong
**How to avoid:**
1. DTB must be at /boot/EFI/dtb/rockchip/rk3588s-nanopi-m6.dtb
2. Installer copies DTB to correct location (update main.go ChipsetName function)
3. Verify with `cat /proc/device-tree/model`
**Warning signs:** Model shows wrong board; peripherals missing

### Pitfall 6: DTS Compilation Warnings Indicating Real Problems
**What goes wrong:** dtc warns about "unevaluated properties" or "failed to match schema"
**Why it happens:** DTS uses properties not in schema; node structure incorrect
**How to avoid:**
1. Most warnings are benign (disabled nodes)
2. Watch for: "unable to find property", "node references non-existent"
3. Critical: Any warning about phandle resolution failure
**Warning signs:** Warnings during kernel build; boot hangs at device probe

## Code Examples

Verified patterns from official sources:

### Adding NanoPi M6 DTS to Kernel Build
```yaml
# In artifacts/kernel/mainline/prepare/pkg.yaml (after tar extraction)
# Download NanoPi M6 DTS from Armbian
prepare:
  - |
    # Download NanoPi M6 device tree from Armbian
    curl -fsSL -o arch/arm64/boot/dts/rockchip/rk3588s-nanopi-m6.dts \
      "https://raw.githubusercontent.com/armbian/build/main/patch/u-boot/v2025.10/dt_upstream_rockchip/rk3588s-nanopi-m6.dts"

    # Add to Makefile for compilation
    sed -i '/rk3588s-nanopi-r6s.dtb/a dtb-$(CONFIG_ARCH_ROCKCHIP) += rk3588s-nanopi-m6.dtb' \
      arch/arm64/boot/dts/rockchip/Makefile
```

### Installer DTB Path Update
```go
// Source: installers/rk3588/src/main.go
func ChipsetName(o rk3588ExtraOpts) string {
    if o.Chipset != "" {
        return o.Chipset
    }
    switch o.Board {
    case "rock-5a":
        return "rk3588s"
    case "rock-5b":
        return "rk3588"
    case "nanopi-m6":  // ADD THIS
        return "rk3588s"
    }
    return ""
}
```

### NVMe Detection Validation
```bash
# During boot validation, check dmesg for NVMe:
dmesg | grep -i nvme
# Expected:
# nvme nvme0: pci function 0001:21:00.0
# nvme nvme0: 1/0/0 default/read/poll queues

lsblk
# Expected:
# nvme0n1  259:0    0  1000G  0 disk
```

### Ethernet Link Validation
```bash
# Check ethernet interface appears
ip link show eth0
# Expected: eth0: <BROADCAST,MULTICAST> ... state DOWN

# Check PHY detected
dmesg | grep -i "stmmac\|r8169\|rtl"
# Expected: stmmac_dvr_probe: ... RTL8211F Gigabit Ethernet

# Bring up interface
ip link set eth0 up
dmesg | grep "Link is Up"
# Expected: r8169 eth0: Link is Up - 1Gbps/Full
```

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| Rockchip BSP kernel 5.10 | Collabora mainline 6.9+ | 2023-2024 | Better upstream support, fewer out-of-tree patches |
| Separate idbloader+u-boot.itb | Vendor U-Boot binaries | Phase 2 discovery | NanoPi M6 requires vendor format |
| Board-specific kernel configs | Single RK3588 config | Talos SBC pattern | Simpler maintenance |
| UART0 console | HDMI/tty1 console | User decision | Available for debug, HDMI primary |

**Deprecated/outdated:**
- Rockchip kernel 4.19: Obsolete, lacks modern driver support
- Custom DTB overlays for NVMe: Not needed; base DTS should enable PCIe correctly
- EDK2 UEFI: Alternative bootloader, but we're using vendor U-Boot per Phase 2

## Open Questions

Things that couldn't be fully resolved:

1. **Exact PCIe lane for NanoPi M6 M.2 slot**
   - What we know: NanoPi R6C uses pcie2x1l1 or pcie2x1l2 for M.2
   - What's unclear: NanoPi M6 may use different lane assignment
   - Recommendation: Try pcie2x1l1 first (matches R6-series); adjust if NVMe not detected

2. **NanoPi M6 DTS mainline status**
   - What we know: Armbian has working DTS; not yet in mainline Linux
   - What's unclear: When/if it will be upstreamed
   - Recommendation: Use Armbian DTS; monitor mainline for eventual inclusion

3. **Collabora kernel branch compatibility with Armbian DTS**
   - What we know: Both target similar kernel versions (6.x)
   - What's unclear: May be pinctrl or clock naming differences
   - Recommendation: Test DTS; fix includes if compilation fails

## Sources

### Primary (HIGH confidence)
- Armbian rk3588s-nanopi-m6.dts - Board device tree with working peripheral configs
- Existing project config-arm64 - Kernel config with NVMe/STMMAC/USB enabled
- Talos overlay interface (pkg/machinery/overlay) - Install method signature

### Secondary (MEDIUM confidence)
- [milas/talos-sbc-rk3588](https://github.com/milas/talos-sbc-rk3588) - Overlay pattern and kernel packaging
- [armbian/linux-rockchip releases](https://github.com/armbian/linux-rockchip/releases) - NanoPi M6 support PR #258
- [torvalds/linux rk3588s-nanopi-r6.dtsi](https://github.com/torvalds/linux/blob/master/arch/arm64/boot/dts/rockchip/rk3588s-nanopi-r6.dtsi) - Reference PCIe/GMAC config

### Tertiary (LOW confidence - needs validation)
- WebSearch results for RK3588 NVMe issues - PCIe BAR allocation failure reports
- WebSearch results for STMMAC ethernet troubleshooting - TX delay configuration

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH - Existing project already has kernel/DTB build pipeline
- Architecture: HIGH - Talos overlay pattern well documented in existing code
- Device tree source: MEDIUM - Armbian DTS not yet tested with Collabora kernel
- Pitfalls: MEDIUM - Based on community reports, not direct experience

**Research date:** 2026-02-03
**Valid until:** 30 days (kernel/DTB are stable; Armbian DTS unlikely to change significantly)

---

*Phase: 03-device-tree-kernel*
*Research completed: 2026-02-03*
