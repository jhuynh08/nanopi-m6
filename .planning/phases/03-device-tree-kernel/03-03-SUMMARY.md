---
phase: 03-device-tree-kernel
plan: 03
subsystem: bootloader
tags: [u-boot, dtb, rockchip, rk3588s, friendlyelec, bldr, overlay]

# Dependency graph
requires:
  - phase: 03-01
    provides: "Pre-extracted vendor DTB in artifacts/dtb"
  - phase: 03-02
    provides: "RK3588 installer updated for vendor U-Boot format"
provides:
  - "Vendor U-Boot binaries committed at artifacts/u-boot/nanopi-m6/vendor/"
  - "Simplified pkg.yaml copying vendor binaries (no build)"
  - "Successful overlay build with NanoPi M6 artifacts"
  - "DTB context path fix for Bldr"
affects:
  - phase-03-04
  - phase-04-talos-integration

# Tech tracking
tech-stack:
  added:
    - "FriendlyELEC idbloader.img (vendor)"
    - "FriendlyELEC uboot.img (vendor)"
  patterns:
    - "Vendor binary copy pattern in Bldr pkg.yaml"
    - "/pkg context path pattern for Bldr"

key-files:
  created:
    - "artifacts/u-boot/nanopi-m6/vendor/idbloader.img"
    - "artifacts/u-boot/nanopi-m6/vendor/uboot.img"
  modified:
    - "artifacts/u-boot/nanopi-m6/pkg.yaml"
    - "artifacts/dtb/pkg.yaml"

key-decisions:
  - "Move vendor binaries inside nanopi-m6/ directory for Bldr context inclusion"
  - "Use /pkg path prefix for context-based file access in pkg.yaml"
  - "Overlay build validates integration; raw image requires CI/CD"

patterns-established:
  - "Vendor binary context pattern: Place vendor files in subdirectory of pkg.yaml location"
  - "Bldr context path: /pkg/ for context directory, /src/ for source stages"

# Metrics
duration: 26min
completed: 2026-02-03
---

# Phase 03 Plan 03: Full Image Build Summary

**Talos SBC overlay builds successfully with FriendlyELEC vendor U-Boot binaries and DTB, validating NanoPi M6 integration in Bldr build system**

## Performance

- **Duration:** 26 min
- **Started:** 2026-02-03T06:39:07Z
- **Completed:** 2026-02-03T07:05:00Z
- **Tasks:** 2/3 completed (Task 3 deferred)
- **Files modified:** 4

## Accomplishments

- Committed FriendlyELEC vendor U-Boot binaries to artifacts/u-boot/nanopi-m6/vendor/
- Simplified pkg.yaml to copy vendor binaries (removed mainline U-Boot build steps)
- Fixed Bldr context paths for both DTB and U-Boot pkg.yaml files
- Full overlay build completes successfully with NanoPi M6 artifacts
- All build artifacts verified: DTB, idbloader.img, uboot.img, installer binary

## Task Commits

Each task was committed atomically:

1. **Task 1a: Commit vendor U-Boot binaries** - `6769b8d` (feat)
2. **Task 1b: Update pkg.yaml for vendor binaries** - `f24a8d6` (feat)
3. **Task 1c: Fix Bldr context paths** - `e5ef392` (fix)

Task 2 verified build success (no commit - verification only)
Task 3 deferred (requires CI/CD for raw image generation)

**Plan metadata:** Committed with this summary

## Files Created/Modified

- `artifacts/u-boot/nanopi-m6/vendor/idbloader.img` - FriendlyELEC bootloader stage 1
- `artifacts/u-boot/nanopi-m6/vendor/uboot.img` - FriendlyELEC U-Boot proper
- `artifacts/u-boot/nanopi-m6/pkg.yaml` - Simplified to copy vendor binaries
- `artifacts/dtb/pkg.yaml` - Fixed context path (/pkg/rockchip/ instead of /pkg/artifacts/dtb/rockchip/)

## Decisions Made

1. **Vendor binary location:** Moved from `artifacts/u-boot/nanopi-m6-vendor/` to `artifacts/u-boot/nanopi-m6/vendor/`
   - Bldr context is the pkg.yaml directory, so vendor binaries must be inside that directory to be accessible
   - Simplifies context management and follows single-directory-per-stage pattern

2. **Bldr context path pattern:**
   - Context files are accessible at `/pkg/` (not `/src/`)
   - Fixed both DTB and U-Boot pkg.yaml to use correct paths

3. **Raw image generation deferred:**
   - `make local-talos-sbc-rk3588-mainline` produces overlay, not raw image
   - Raw image generation requires Talos imager + registry (CI/CD workflow)
   - Overlay build validates all artifacts are correctly integrated

## Deviations from Plan

### Context Path Fixes

**1. [Rule 3 - Blocking] Bldr context paths were incorrect**
- **Found during:** Task 2 (build attempt)
- **Issue:** Build failed with "No such file or directory" for vendor binaries and DTB
  - pkg.yaml used `/src/artifacts/...` but context path is `/pkg/`
  - Vendor binaries were outside the pkg.yaml context directory
- **Fix:**
  - Moved vendor binaries to `artifacts/u-boot/nanopi-m6/vendor/`
  - Updated pkg.yaml to use `/pkg/vendor/` path
  - Fixed DTB pkg.yaml to use `/pkg/rockchip/` path
- **Files modified:** artifacts/u-boot/nanopi-m6/pkg.yaml, artifacts/dtb/pkg.yaml
- **Verification:** Build completes successfully
- **Committed in:** e5ef392

### Task 3 Scope Deviation

**2. [Rule 4 - Architectural] Raw image generation requires CI/CD infrastructure**
- **Found during:** Task 2 verification
- **Issue:** Plan expected `make local-talos-sbc-rk3588-mainline` to produce `.raw.xz` image
- **Actual behavior:** Build produces overlay artifacts, not raw disk image
  - Raw image requires: pushing overlay to registry, running Talos imager
  - This is the standard CI/CD workflow, not local development
- **Resolution:** Task 3 (flash to SD card) deferred to CI/CD or future local imager setup
- **Impact:** Overlay build validates integration; raw image is a packaging step

---

**Total deviations:** 2 (1 auto-fixed blocking, 1 scope clarification)
**Impact on plan:** Context path fixes necessary for build success. Scope clarification documents correct workflow.

## Issues Encountered

- **GNU Make version incompatibility:** System make (3.81) doesn't support `export define` syntax in Makefile
  - Resolution: Use `/opt/homebrew/bin/gmake` (GNU Make 4.4.1)
  - May want to add this to development setup documentation

## Build Output Verification

```
_out/
  artifacts/arm64/
    dtb/rockchip/rk3588s-nanopi-m6.dtb    (262KB - vendor DTB)
    u-boot/nanopi-m6/
      idbloader.img                        (4.0MB - vendor bootloader)
      uboot.img                            (4.0MB - vendor U-Boot)
      u-boot-rockchip.bin                  (9.2MB - mainline, not used)
  installers/rk3588                        (3.5MB - overlay installer binary)
  profiles/board-rk3588.yaml               (profile for Talos imager)
```

## User Setup Required

None - no external service configuration required.

## Next Phase Readiness

**Ready for Phase 03-04 (Hardware Validation) with caveats:**

1. **Overlay is complete** - All NanoPi M6 artifacts integrated and build passes
2. **Raw image generation options:**
   - **Option A (CI/CD):** Push overlay to registry, trigger CI workflow for image generation
   - **Option B (Local imager):** Set up local Talos imager (requires Talos repo checkout + patches)
   - **Option C (Existing image):** Use FriendlyELEC Ubuntu image with vendor U-Boot (validated in Phase 2)

3. **Recommended approach for Phase 04:**
   - Use CI/CD to generate raw image (aligns with project workflow)
   - Alternative: Continue with FriendlyELEC Ubuntu for driver validation

**Blockers:** None
**Concerns:** Raw image generation workflow should be documented for local development

---
*Phase: 03-device-tree-kernel*
*Plan: 03 (Full Image Build)*
*Completed: 2026-02-03*
