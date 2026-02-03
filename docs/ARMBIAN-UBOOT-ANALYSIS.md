# Armbian NanoPi M6 U-Boot Analysis

**Purpose:** Extract and document Armbian's NanoPi M6 U-Boot configuration to enable successful boot on NanoPi M6 hardware.

**Context:** Boot test with rock5a-rk3588s_defconfig failed (no LED, no HDMI, no network). Armbian boots successfully on same hardware, proving correct configuration exists.

## Source Location

- **Repository:** https://github.com/armbian/build
- **Defconfig:** `patch/u-boot/v2025.10/defconfig/nanopi-m6-rk3588s_defconfig`
- **SPI Defconfig:** `patch/u-boot/v2025.10/defconfig/nanopi-m6-spi-rk3588s_defconfig`
- **Device Tree:** `patch/u-boot/v2025.10/dt_upstream_rockchip/rk3588s-nanopi-m6.dts`
- **Board Config:** `config/boards/nanopi-m6.conf`

**U-Boot Version:** v2025.10 (mainline, not vendor fork)

## Key CONFIG Differences from rock5a

| CONFIG Option | rock5a Value | M6 Value | Impact |
|---------------|--------------|----------|--------|
| `CONFIG_DEFAULT_DEVICE_TREE` | `"rockchip/rk3588s-rock-5a"` | `"rockchip/rk3588s-nanopi-m6"` | **Critical:** Board-specific device tree |
| `CONFIG_TARGET_EVB_RK3588` | (not set) | `y` | Uses generic RK3588 EVB target |
| `CONFIG_TARGET_ROCK5A_RK3588` | `y` | (not set) | Rock5A-specific target disabled |
| `CONFIG_DEFAULT_FDT_FILE` | `"rockchip/rk3588s-rock-5a.dtb"` | `"rockchip/rk3588s-nanopi-m6.dtb"` | **Critical:** Boot loads wrong DTB |
| `CONFIG_SF_DEFAULT_SPEED` | (not set) | `24000000` | SPI flash speed for SPI boot |
| `CONFIG_SF_DEFAULT_MODE` | (not set) | `0x2000` | SPI flash mode |
| `CONFIG_SF_DEFAULT_BUS` | (not set) | `5` | SPI bus selection |
| `CONFIG_SPL_SPI` | (not set) | `y` | SPL SPI support |
| `CONFIG_SPL_SPI_LOAD` | (not set) | `y` | SPL can load from SPI |
| `CONFIG_SYS_SPI_U_BOOT_OFFS` | (not set) | `0x60000` | U-Boot offset in SPI flash |
| `CONFIG_PCI` | (not set) | `y` | PCI support enabled |
| `CONFIG_AHCI` | (not set) | `y` | AHCI SATA support |
| `CONFIG_AHCI_PCI` | (not set) | `y` | AHCI over PCI |
| `CONFIG_DWC_AHCI` | (not set) | `y` | DesignWare AHCI driver |
| `CONFIG_CMD_PCI` | (not set) | `y` | PCI commands in U-Boot |
| `CONFIG_LED` | `y` | (not set) | LED subsystem (uses GPIO directly) |
| `CONFIG_LED_GPIO` | `y` | (not set) | GPIO LED driver |
| `CONFIG_NVME_PCI` | (not set) | `y` | NVMe support over PCI |
| `CONFIG_PCIE_DW_ROCKCHIP` | (not set) | `y` | Rockchip PCIe driver |
| `CONFIG_SCSI` | (not set) | `y` | SCSI subsystem |
| `CONFIG_DM_ETH` | (not set) | `y` | Ethernet device model |
| `CONFIG_DM_MDIO` | (not set) | `y` | MDIO device model |
| `CONFIG_PHY_GIGE` | (not set) | `y` | Gigabit PHY support |
| `CONFIG_ETH_DESIGNWARE` | (not set) | `y` | DesignWare ethernet |
| `CONFIG_RGMII` | (not set) | `y` | RGMII interface |
| `CONFIG_MII` | (not set) | `y` | MII interface |
| `CONFIG_GMAC_ROCKCHIP` | (not set) | `y` | Rockchip GMAC driver |
| `CONFIG_RTL8169` | (not set) | `y` | Realtek 8169 driver |
| `CONFIG_PHYLIB` | (not set) | `y` | PHY library |

### Critical Differences

1. **Device Tree Selection:**
   - rock5a loads `rk3588s-rock-5a.dtb`
   - M6 needs `rk3588s-nanopi-m6.dtb`
   - This is the **root cause** of boot failure - wrong GPIO/pinmux/regulator configuration

2. **Target Selection:**
   - rock5a uses `CONFIG_TARGET_ROCK5A_RK3588`
   - M6 uses generic `CONFIG_TARGET_EVB_RK3588`
   - This affects board-specific initialization sequences

3. **Peripheral Support:**
   - M6 defconfig enables PCIe, NVMe, SATA, networking
   - These may or may not affect early boot, but are hardware-specific

## DDR/BL31 Blob Requirements

### From rockchip64_common.inc (Armbian defaults for RK3588):

```bash
DDR_BLOB="rk35/rk3588_ddr_lp4_2112MHz_lp5_2400MHz_v1.18.bin"
BL31_BLOB="rk35/rk3588_bl31_v1.48.elf"
```

### Currently in our pkg.yaml:

```bash
ROCKCHIP_TPL: rk35/rk3588_ddr_lp4_2112MHz_lp5_2400MHz_v1.16.bin
BL31: rk35/rk3588_bl31_v1.45.elf
```

### Recommendation:

Update to latest blob versions:
- **DDR:** v1.16 -> v1.18 (optional, v1.16 should work)
- **BL31:** v1.45 -> v1.48 (optional, v1.45 should work)

**Note:** The DDR blob version is less critical than the defconfig - the blob performs memory training, but the defconfig determines board initialization and device tree selection.

## Armbian Build Process

From `nanopi-m6.conf`:

```bash
BOOTCONFIG="nanopi-m6-rk3588s_defconfig"
BOOT_SOC="rk3588"
BOOT_SCENARIO="spl-blobs"
BOOTSOURCE="https://github.com/u-boot/u-boot.git"
BOOTBRANCH="tag:v2025.10"
BOOTPATCHDIR="v2025.10"
```

Armbian:
1. Clones mainline U-Boot v2025.10
2. Applies patches from `patch/u-boot/v2025.10/`
3. Copies `nanopi-m6-rk3588s_defconfig` to configs/
4. Copies `rk3588s-nanopi-m6.dts` to arch/arm/dts/
5. Builds with `make nanopi-m6-rk3588s_defconfig`

## Device Tree Analysis

The M6 device tree (`rk3588s-nanopi-m6.dts`) defines:

### LEDs:
```dts
leds {
    sys_led: led-0 {
        gpios = <&gpio1 RK_PA4 GPIO_ACTIVE_HIGH>;  // GPIO1_A4
    };
    user_led: led-1 {
        gpios = <&gpio1 RK_PA6 GPIO_ACTIVE_HIGH>;  // GPIO1_A6
    };
};
```

### Ethernet:
```dts
&gmac1 {
    phy-handle = <&rgmii_phy1>;
    phy-mode = "rgmii-rxid";
    tx_delay = <0x42>;
};
```

### PMIC (RK806 over SPI2):
```dts
&spi2 {
    pmic@0 {
        compatible = "rockchip,rk806";
        // All voltage regulators defined here
    };
};
```

This device tree is **completely different** from rock5a, explaining why boot fails.

## Recommended Changes for Plan 02-05

### Option A: Use Armbian Defconfig Directly (RECOMMENDED)

1. Download `nanopi-m6-rk3588s_defconfig` from Armbian
2. Download `rk3588s-nanopi-m6.dts` from Armbian
3. Place in appropriate locations
4. Update pkg.yaml to use `nanopi-m6-rk3588s_defconfig`

**Steps:**
```bash
# Download defconfig
curl -o configs/nanopi-m6-rk3588s_defconfig \
  "https://raw.githubusercontent.com/armbian/build/main/patch/u-boot/v2025.10/defconfig/nanopi-m6-rk3588s_defconfig"

# Download device tree
curl -o arch/arm/dts/rk3588s-nanopi-m6.dts \
  "https://raw.githubusercontent.com/armbian/build/main/patch/u-boot/v2025.10/dt_upstream_rockchip/rk3588s-nanopi-m6.dts"

# Build with correct defconfig
make nanopi-m6-rk3588s_defconfig
make -j$(nproc)
```

**pkg.yaml changes:**
```yaml
prepare:
  - |
    cd /src
    # Download NanoPi M6 defconfig from Armbian
    curl -o configs/nanopi-m6-rk3588s_defconfig \
      "https://raw.githubusercontent.com/armbian/build/main/patch/u-boot/v2025.10/defconfig/nanopi-m6-rk3588s_defconfig"
    # Download NanoPi M6 device tree from Armbian
    curl -o arch/arm/dts/rk3588s-nanopi-m6.dts \
      "https://raw.githubusercontent.com/armbian/build/main/patch/u-boot/v2025.10/dt_upstream_rockchip/rk3588s-nanopi-m6.dts"
    make nanopi-m6-rk3588s_defconfig
```

### Option B: Patch rock5a Defconfig

Apply sed commands to transform rock5a defconfig to M6:

```bash
make rock5a-rk3588s_defconfig

# Change device tree references
sed -i 's/rk3588s-rock-5a/rk3588s-nanopi-m6/g' .config

# Change target
sed -i 's/CONFIG_TARGET_ROCK5A_RK3588=y/# CONFIG_TARGET_ROCK5A_RK3588 is not set/' .config
echo "CONFIG_TARGET_EVB_RK3588=y" >> .config

# Add M6-specific configs
echo "CONFIG_SF_DEFAULT_SPEED=24000000" >> .config
echo "CONFIG_SF_DEFAULT_MODE=0x2000" >> .config
# ... many more options needed

make olddefconfig  # Resolve conflicts
```

**Problems with Option B:**
- Requires maintaining complex sed/patch logic
- Easy to miss required options
- Device tree still missing
- More fragile than using Armbian's tested config

### Option C: Store Configs Locally

1. Store `nanopi-m6-rk3588s_defconfig` in `artifacts/u-boot/nanopi-m6/`
2. Store `rk3588s-nanopi-m6.dts` in `artifacts/u-boot/nanopi-m6/`
3. Copy files during build

**Advantage:** No network dependency during build
**Disadvantage:** Configs may drift from Armbian updates

## Recommended Option: A (with fallback to C)

**Primary: Option A**
- Download configs directly from Armbian during build
- Always gets latest tested configuration
- Minimal maintenance

**Fallback: Option C**
- If network access is problematic during builds
- Store local copies with version tracking
- Update periodically from Armbian

### Implementation for Plan 02-05

Update `artifacts/u-boot/nanopi-m6/pkg.yaml`:

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

        # Download NanoPi M6 defconfig from Armbian (tested configuration)
        curl -fsSL -o configs/nanopi-m6-rk3588s_defconfig \
          "https://raw.githubusercontent.com/armbian/build/main/patch/u-boot/v2025.10/defconfig/nanopi-m6-rk3588s_defconfig"

        # Download NanoPi M6 device tree from Armbian
        curl -fsSL -o arch/arm/dts/rk3588s-nanopi-m6.dts \
          "https://raw.githubusercontent.com/armbian/build/main/patch/u-boot/v2025.10/dt_upstream_rockchip/rk3588s-nanopi-m6.dts"

        # Use the correct M6 defconfig
        make nanopi-m6-rk3588s_defconfig
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

## Summary

The boot failure was caused by using `rock5a-rk3588s_defconfig` which:
1. Loads wrong device tree (`rk3588s-rock-5a.dtb` instead of `rk3588s-nanopi-m6.dtb`)
2. Has different board initialization sequence
3. Misses M6-specific peripheral configurations

**Solution:** Use Armbian's `nanopi-m6-rk3588s_defconfig` and `rk3588s-nanopi-m6.dts` which are proven to work on NanoPi M6 hardware.

---
*Analysis completed: 2026-02-03*
*Source: Armbian build repository @ main branch*
