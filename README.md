# Talos Linux NanoPi M6 Overlay

Talos Linux overlay and custom kernel for the [FriendlyELEC NanoPi M6](https://wiki.friendlyelec.com/wiki/index.php/NanoPi_M6) (RK3588S) single-board computer.

## Status

This project is under active development. See [.planning/](.planning/) for detailed project planning and progress.

**Goal:** Boot Talos Linux on NanoPi M6 and register with Talos Omni to join a home Kubernetes cluster.

## Why does this exist?

The NanoPi M6 uses the Rockchip RK3588S SoC, which requires a newer kernel than the 6.6 LTS available in upstream Talos Linux. This project adapts the [milas/talos-sbc-rk3588](https://github.com/milas/talos-sbc-rk3588) overlay (which supports Rock 5 series boards) to add NanoPi M6 support.

Key differences from Rock 5:
- RK3588S (no PCIe 3.0 x4 slot variant)
- eMMC and microSD storage
- Different GPIO and peripheral layout
- FriendlyELEC-specific device tree requirements

## Build

Prerequisites:
- Docker
- Make
- Go 1.22+

```bash
# Build all artifacts
make

# Build specific target
make u-boot-nanopi-m6
make kernel
make installer-rk3588
```

## Device Support

Currently targeting:
- **NanoPi M6** (RK3588S) - in development

Upstream support (from milas/talos-sbc-rk3588):
- Rock 5B (RK3588)
- Rock 5A (RK3588S)

## Resources

- [NanoPi M6 Wiki](https://wiki.friendlyelec.com/wiki/index.php/NanoPi_M6)
- [Collabora RK3588 upstreaming](https://gitlab.collabora.com/hardware-enablement/rockchip-3588)
- [siderolabs/talos](https://github.com/siderolabs/talos/)
- [milas/talos-sbc-rk3588](https://github.com/milas/talos-sbc-rk3588) - upstream project

## Disclaimer

This is NOT supported or endorsed by Rockchip, FriendlyELEC, Sidero Labs, or Collabora - please do not contact them with support requests for this project.
