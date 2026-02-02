# Architecture Patterns: Talos SBC Overlay Projects

**Domain:** Talos Linux SBC overlay development
**Researched:** 2026-02-02
**Overall Confidence:** HIGH (verified via official repositories)

## Executive Summary

Talos SBC overlay projects follow a standardized structure derived from the [siderolabs/sbc-template](https://github.com/siderolabs/sbc-template). The architecture is container-based, using the `bldr` build system to produce OCI images that integrate with the Talos imager. Key components include: artifacts (firmware, bootloader, device trees), installers (Go binaries), and profiles (YAML configuration). The build system orchestrates multi-stage container builds with dependency resolution between packages.

## Directory Structure Overview

```
talos-sbc-overlay/
├── .github/
│   └── workflows/
│       └── ci.yaml              # CI/CD pipeline
├── artifacts/                   # Build artifacts (firmware, u-boot, dtb)
│   ├── arm-trusted-firmware/    # ATF per SoC (rk3328, rk3399, rk3568, rk3588)
│   ├── rkbin/                   # Rockchip proprietary binaries
│   ├── <board>/
│   │   └── u-boot/              # Board-specific U-Boot
│   │       └── pkg.yaml
│   └── dtb/                     # Device tree blobs (if separate)
├── installers/                  # Installer binaries (Go)
│   ├── <board>/
│   │   ├── src/
│   │   │   ├── main.go
│   │   │   ├── go.mod
│   │   │   └── go.sum
│   │   └── pkg.yaml
│   └── pkg.yaml
├── internal/
│   └── base/
│       └── pkg.yaml             # Base toolchain dependency
├── profiles/                    # Board profile definitions
│   ├── <board>/
│   │   └── <board>.yaml
│   └── pkg.yaml
├── hack/                        # Development scripts
├── .conform.yaml                # Commit message linting
├── .kres.yaml                   # Build pipeline configuration
├── Makefile                     # Build orchestration
├── Pkgfile                      # Root package definition (bldr)
└── go.work                      # Go workspace (multi-module)
```

## Component Boundaries

### 1. Artifacts (`/artifacts/`)

**Responsibility:** Contains firmware, bootloader, and device tree binaries required for board boot.

**Communicates With:** Installers (artifacts are deployed by installer), Build System (artifacts are built via pkg.yaml)

| Subdirectory | Purpose | Build Output |
|--------------|---------|--------------|
| `arm-trusted-firmware/<soc>/` | ARM Trusted Firmware per SoC | `bl31.elf` |
| `rkbin/` | Rockchip DDR init and BL31 binaries | Proprietary blobs |
| `<board>/u-boot/` | Board-specific U-Boot build | `u-boot-rockchip.bin` |
| `dtb/` | Device tree blobs (when separate from kernel) | `*.dtb` |

**Key Pattern:** Artifacts are organized hierarchically:
- SoC-level (ATF, rkbin) - shared across boards with same SoC
- Board-level (U-Boot) - specific to each board

### 2. Installers (`/installers/`)

**Responsibility:** Go binaries that execute during Talos installation, copying artifacts to target disk and configuring kernel arguments.

**Communicates With:** Artifacts (reads and deploys), Profiles (referenced by name), Talos Imager (invoked by)

**Interface Contract:**
```go
// Installer must implement these methods (from siderolabs/talos)
type Installer interface {
    GetOptions(ExtraOptions) (Options, error)  // Returns kernel args, etc.
    Install(InstallOptions) error              // Deploys artifacts to disk
}
```

**CLI Contract:**
- Accept YAML via stdin, output YAML to stdout
- Implement `install` and `get-options` subcommands
- Exit non-zero on errors
- Must be statically linked (for scratch containers)

**Structure per board:**
```
installers/<board>/
├── src/
│   ├── main.go          # Installer implementation
│   ├── go.mod           # Go module definition
│   └── go.sum           # Dependency checksums
└── pkg.yaml             # Build configuration for installer binary
```

### 3. Profiles (`/profiles/`)

**Responsibility:** YAML definitions specifying disk image parameters for each board.

**Communicates With:** Talos Imager (consumed by), Installers (paired with)

**Profile Schema:**
```yaml
arch: arm64
platform: metal
secureboot: false
output:
  kind: image
  imageOptions:
    diskSize: 1306525696      # ~1.2GB
    diskFormat: raw
  outFormat: .xz              # Compression format
bootloader: grub              # or u-boot
```

**Key Fields:**
| Field | Purpose |
|-------|---------|
| `arch` | Target architecture (arm64) |
| `platform` | Installation target (metal, cloud, etc.) |
| `diskSize` | Total disk image size in bytes |
| `diskFormat` | raw, qcow2, etc. |
| `outFormat` | Compression (.xz, .gz, none) |
| `bootloader` | grub or u-boot |

### 4. Internal Base (`/internal/base/`)

**Responsibility:** Provides shared build toolchain dependencies used by other packages.

**Communicates With:** All other packages (dependency)

**Typical Content:**
```yaml
name: base
variant: scratch
shell: /bin/bash
dependencies:
  - image: "{{ .BUILD_ARG_TOOLS_PREFIX }}/tools:{{ .BUILD_ARG_TOOLS }}"
    to: /rootfs
finalize:
  - from: /rootfs
    to: /
```

### 5. Build System Files

**Pkgfile (root):**
- Magic comment directing to bldr frontend
- Version definitions for upstream sources (U-Boot, kernel, ATF)
- SHA checksums for reproducible builds
- Variable definitions used in pkg.yaml templates

**pkg.yaml (per component):**
- Package name and variant (alpine or scratch)
- Dependencies (internal stages or external images)
- Build steps (sources, prepare, build, install, test)
- Finalize instructions (what to include in output)

**.kres.yaml:**
- Defines build targets for the Makefile
- Specifies build arguments (PKGS_PREFIX, PKGS, KERNEL_VARIANT)
- Controls what `make rekres` generates

**Makefile:**
- Auto-generated by `make rekres`
- Provides targets: `docker-*`, `local-*`, `target-*`
- Orchestrates bldr via Docker buildx

## Data Flow Diagram

```
                                    INPUTS
                                      |
    +----------------+----------------+----------------+
    |                |                |                |
    v                v                v                v
+--------+     +---------+     +----------+     +--------+
| Pkgfile|     | rkbin/  |     | Kernel   |     | U-Boot |
| (vars) |     | (blobs) |     | (upstream)|    | (src)  |
+--------+     +---------+     +----------+     +--------+
    |                |                |                |
    +----------------+-------+--------+----------------+
                             |
                             v
                    +----------------+
                    |  bldr build    |
                    |  (Docker/      |
                    |   buildkit)    |
                    +----------------+
                             |
         +-------------------+-------------------+
         |                   |                   |
         v                   v                   v
  +-------------+     +-------------+     +-------------+
  | artifacts/  |     | installers/ |     | profiles/   |
  | (u-boot,    |     | (Go bins)   |     | (YAML)      |
  | atf, dtb)   |     |             |     |             |
  +-------------+     +-------------+     +-------------+
         |                   |                   |
         +-------------------+-------------------+
                             |
                             v
                    +----------------+
                    | Container      |
                    | Image          |
                    | (OCI)          |
                    +----------------+
                             |
                             v
                    +----------------+
                    | Talos Imager   |
                    | (imager tool)  |
                    +----------------+
                             |
                             v
                    +----------------+
                    | Bootable       |
                    | Disk Image     |
                    | (.raw.xz)      |
                    +----------------+
```

## Build Order and Dependencies

### Dependency Graph

```
                    base
                      |
         +-----+-----+-----+-----+
         |     |     |     |     |
         v     v     v     v     v
       ATF   rkbin  DTB  U-Boot  Installer
         |     |     |     |        |
         +--+--+     +--+--+        |
            |           |           |
            v           v           |
         firmware   bootloader      |
            |           |           |
            +-----+-----+           |
                  |                 |
                  v                 v
               <board>          profiles
                  |                 |
                  +--------+--------+
                           |
                           v
                      sbc-overlay
                      (final image)
```

### Build Order (from leaf to root)

1. **base** - Build toolchain setup (first, no deps)
2. **arm-trusted-firmware** - Per SoC ATF (depends on base)
3. **rkbin** - Rockchip blobs (depends on base)
4. **<board>/u-boot** - Board-specific U-Boot (depends on ATF, rkbin)
5. **<board> installer** - Go installer binary (depends on base)
6. **profiles** - Aggregate all board profiles (depends on individual profiles)
7. **sbc-overlay** - Final image (aggregates all artifacts, installers, profiles)

### Parallelization Opportunities

```
                    [base]
                      |
    +-----------------+------------------+
    |                 |                  |
    v                 v                  v
[ATF/rkbin]    [installer build]    [profiles]
    |                 |                  |
    v                 |                  |
[u-boot]              |                  |
    |                 |                  |
    +--------+--------+------------------+
             |
             v
      [final overlay]
```

- ATF, rkbin, and installers can build in parallel after base
- U-Boot depends on ATF and rkbin
- Profiles can build independently
- Final overlay waits for all

## Key Files to Modify When Adding a New Board

### Required Changes

| File/Directory | Action | Purpose |
|----------------|--------|---------|
| `artifacts/<board>/u-boot/pkg.yaml` | CREATE | U-Boot build config |
| `installers/<board>/src/main.go` | CREATE | Installer implementation |
| `installers/<board>/src/go.mod` | CREATE | Go module |
| `installers/<board>/pkg.yaml` | CREATE | Installer build config |
| `profiles/<board>/<board>.yaml` | CREATE | Board profile |
| `installers/pkg.yaml` | MODIFY | Add board to aggregation |
| `profiles/pkg.yaml` | MODIFY | Add board profile to aggregation |
| `Pkgfile` | MODIFY (maybe) | Add variables if new SoC/kernel |

### Typical Workflow

1. **Create U-Boot artifact:**
   ```
   artifacts/<board>/u-boot/pkg.yaml
   ```
   Define U-Boot defconfig, patches, build steps

2. **Create installer:**
   ```
   installers/<board>/src/main.go    # Copy from template, customize
   installers/<board>/pkg.yaml       # Build configuration
   ```
   Implement GetOptions() and Install() methods

3. **Create profile:**
   ```
   profiles/<board>/<board>.yaml
   ```
   Define disk image parameters

4. **Update aggregation:**
   - Add board to `installers/pkg.yaml` finalize section
   - Add board to `profiles/pkg.yaml` finalize section

5. **Regenerate build system:**
   ```bash
   make rekres
   ```

6. **Build and test:**
   ```bash
   make docker-sbc-<overlay>
   ```

## Variant: milas/talos-sbc-rk3588 Differences

The milas project differs from official siderolabs structure:

| Aspect | siderolabs/sbc-rockchip | milas/talos-sbc-rk3588 |
|--------|------------------------|------------------------|
| Kernel | Uses Talos upstream | Custom kernel (6.9+) from Collabora |
| SoC Support | Multiple (RK33xx, RK35xx, RK3588) | RK3588 only |
| Build Targets | Single overlay | Mainline + BSP variants |
| Complexity | More boards, more abstraction | Simpler, focused |

**Key milas additions:**
- `artifacts/kernel/` - Custom kernel build
- Dual kernel variants: `mainline` and `bsp`
- Collabora U-Boot/kernel forks instead of mainline

## Patterns to Follow

### Pattern 1: Board-Specific pkg.yaml

Each board component gets its own pkg.yaml with clear boundaries:
```yaml
name: nanopi-m6-u-boot
variant: scratch
shell: /bin/bash
dependencies:
  - stage: arm-trusted-firmware-rk3588
  - stage: rkbin
  - stage: base
steps:
  - sources:
      - url: https://github.com/...
        sha256: ...
    build: |
      make nanopi-m6_defconfig
      make -j$(nproc)
    install: |
      cp u-boot-rockchip.bin /rootfs/artifacts/
```

### Pattern 2: Installer Template

Copy and customize from sbc-template:
```go
func (i *BoardInstaller) GetOptions(extra sbc.ExtraOptions) (sbc.Options, error) {
    return sbc.Options{
        KernelArgs: []string{
            "console=tty0",
            "sysctl.kernel.kexec_load_disabled=1",
            "talos.dashboard.disabled=1",
        },
    }, nil
}

func (i *BoardInstaller) Install(options sbc.InstallOptions) error {
    return options.ArtifactHelper.CopyFile(
        "/artifacts/u-boot/nanopi-m6/u-boot-rockchip.bin",
        filepath.Join(options.BootDir, "u-boot-rockchip.bin"),
    )
}
```

### Pattern 3: Profile Definition

Minimal but complete:
```yaml
arch: arm64
platform: metal
secureboot: false
output:
  kind: image
  imageOptions:
    diskSize: 1306525696
    diskFormat: raw
  outFormat: .xz
bootloader: grub
```

## Anti-Patterns to Avoid

### Anti-Pattern 1: Mixing SoC-Level and Board-Level

**Wrong:** Putting board-specific code in SoC-level packages
**Right:** ATF and rkbin are SoC-level; U-Boot config is board-level

### Anti-Pattern 2: Dynamic Linking in Installers

**Wrong:** Using CGO or dynamic dependencies in installer Go code
**Right:** Build with `CGO_ENABLED=0` for static linking (required for scratch variant)

### Anti-Pattern 3: Hardcoded Paths

**Wrong:** Hardcoding `/boot` or `/rootfs` paths
**Right:** Use options.BootDir, options.ArtifactHelper from installer interface

## Integration Points

### With Talos Imager

```bash
# Build overlay image
docker run --rm -t -v $PWD/_out:/out \
  ghcr.io/siderolabs/imager:v1.9.0 \
  --overlay-name=nanopi-m6 \
  --overlay-image=ghcr.io/username/talos-sbc-nanopi-m6:v1.9.0 \
  --arch=arm64
```

### With Talos Omni

Omni can consume custom overlay images when:
1. Image is published to accessible registry
2. Machine configuration references the overlay
3. Overlay is signed (for secure boot scenarios)

## Sources

- [siderolabs/sbc-rockchip](https://github.com/siderolabs/sbc-rockchip) - Official Rockchip overlay (HIGH confidence)
- [siderolabs/sbc-template](https://github.com/siderolabs/sbc-template) - Template for new overlays (HIGH confidence)
- [siderolabs/overlays](https://github.com/siderolabs/overlays) - Overlay registry (HIGH confidence)
- [milas/talos-sbc-rk3588](https://github.com/milas/talos-sbc-rk3588) - Custom RK3588 overlay (HIGH confidence)
- [siderolabs/bldr](https://github.com/siderolabs/bldr) - Build system documentation (HIGH confidence)
- [Talos Overlays Documentation](https://docs.siderolabs.com/talos/v1.11/build-and-extend-talos/custom-images-and-development/overlays) - Official docs (HIGH confidence)
