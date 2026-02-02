# Boot Test Checklist

Template for tracking boot attempts during NanoPi M6 bootloader bring-up.

## Overview

Without UART console access, boot verification relies on observable indicators:
- **LED activity** - Board LEDs indicate SPL/U-Boot/kernel stages
- **HDMI output** - Only available after Linux kernel boots (U-Boot has no HDMI driver)
- **Network connectivity** - Confirms kernel booted and networking initialized

**Critical understanding:** Mainline U-Boot for RK3588 does NOT support HDMI output. There is no VOP2 video driver. The U-Boot stage is "blind" - success can only be verified indirectly through LED activity patterns or eventual kernel boot.

## Iteration Strategy

### Tier 1: Quick Iterations (3 attempts max)

Test basic configurations with current known-good blob versions:
1. Primary defconfig with current blobs
2. Primary defconfig with updated blobs (if available)
3. Alternative defconfig as backup

**Exit criteria:** Move to Tier 2 if no LED activity observed in any attempt.

### Tier 2: Configuration Investigation (2 attempts max)

Deep analysis when Tier 1 fails:
4. Extract and analyze reference defconfig (Armbian M6)
5. Apply targeted patches based on analysis

**Exit criteria:** Move to Tier 3 if still no progress indicators.

### Tier 3: Decision Point

After 5 failed attempts with no LED/network activity:
- Document all configurations tested
- Consider whether UART acquisition is necessary
- May need to park phase for hardware debugging

## Success Indicators

| Indicator | What It Means | Stage Reached |
|-----------|---------------|---------------|
| Power LED on | Board has power | Power |
| SYS LED single blink | BL31/SPL initializing | TPL/SPL |
| SYS LED rapid blinking | U-Boot running | U-Boot |
| SYS LED steady/periodic | Kernel booted | Linux |
| HDMI shows output | Display driver loaded | Linux (6.15+) |
| Network ping responds | Networking initialized | Linux |
| Talosctl responds | Talos running | Talos |

**Failure indicators:**
- No LED activity after 10s = DDR/SPL failure
- LED blinks but stops = U-Boot crash or kernel fail
- HDMI blank but LED steady = Display config issue

---

## Boot Attempt Log

### Attempt #1

**Date:** ____-__-__
**Time:** __:__

#### Configuration

| Setting | Value |
|---------|-------|
| Defconfig | |
| DDR blob | |
| BL31 version | |
| U-Boot source | |
| SD card | |

#### Observation Timeline

| Time Window | Observation | Notes |
|-------------|-------------|-------|
| 0-10s | LED activity? [ ] Yes [ ] No | |
| 10-30s | LED pattern? | |
| 30-60s | HDMI output? [ ] Yes [ ] No | |
| 60-120s | Network ping? [ ] Yes [ ] No | |
| 120s+ | Talosctl? [ ] Yes [ ] No | |

#### Result

- [ ] SUCCESS: Kernel booted, indicators positive
- [ ] PARTIAL: Some activity but incomplete boot
- [ ] FAILURE: No activity observed

#### Observations

```
(Describe what was observed, any patterns, error indicators)
```

#### Next Steps

```
(What to try next based on this result)
```

---

### Attempt #2

**Date:** ____-__-__
**Time:** __:__

#### Configuration

| Setting | Value |
|---------|-------|
| Defconfig | |
| DDR blob | |
| BL31 version | |
| U-Boot source | |
| SD card | |

#### Observation Timeline

| Time Window | Observation | Notes |
|-------------|-------------|-------|
| 0-10s | LED activity? [ ] Yes [ ] No | |
| 10-30s | LED pattern? | |
| 30-60s | HDMI output? [ ] Yes [ ] No | |
| 60-120s | Network ping? [ ] Yes [ ] No | |
| 120s+ | Talosctl? [ ] Yes [ ] No | |

#### Result

- [ ] SUCCESS: Kernel booted, indicators positive
- [ ] PARTIAL: Some activity but incomplete boot
- [ ] FAILURE: No activity observed

#### Observations

```
(Describe what was observed, any patterns, error indicators)
```

#### Next Steps

```
(What to try next based on this result)
```

---

### Attempt #3

**Date:** ____-__-__
**Time:** __:__

#### Configuration

| Setting | Value |
|---------|-------|
| Defconfig | |
| DDR blob | |
| BL31 version | |
| U-Boot source | |
| SD card | |

#### Observation Timeline

| Time Window | Observation | Notes |
|-------------|-------------|-------|
| 0-10s | LED activity? [ ] Yes [ ] No | |
| 10-30s | LED pattern? | |
| 30-60s | HDMI output? [ ] Yes [ ] No | |
| 60-120s | Network ping? [ ] Yes [ ] No | |
| 120s+ | Talosctl? [ ] Yes [ ] No | |

#### Result

- [ ] SUCCESS: Kernel booted, indicators positive
- [ ] PARTIAL: Some activity but incomplete boot
- [ ] FAILURE: No activity observed

#### Observations

```
(Describe what was observed, any patterns, error indicators)
```

#### Next Steps

```
(What to try next based on this result)
```

---

## Configuration Quick Reference

### Blob Versions (Current Project)

| Component | File | Version |
|-----------|------|---------|
| DDR | rk3588_ddr_lp4_2112MHz_lp5_2400MHz_v1.16.bin | v1.16 |
| BL31 | rk3588_bl31_v1.45.elf | v1.45 |

### Available Defconfigs

| Defconfig | Source | Notes |
|-----------|--------|-------|
| nanopi-r6c-rk3588s_defconfig | Mainline U-Boot | Same RK3588S SoC |
| nanopi-r6s-rk3588s_defconfig | Mainline U-Boot | Same RK3588S SoC |
| nanopi-m6-rk3588s_defconfig | Armbian patches | M6-specific, needs extraction |
| rock5a-rk3588s_defconfig | Mainline U-Boot | Different board layout |

### Pre-Test Checklist

Before each boot attempt:
- [ ] SD card freshly flashed (not reused from failed attempt without reflash)
- [ ] Power supply is 5V 4A capable
- [ ] HDMI connected before power on
- [ ] Network cable connected (for ping test)
- [ ] Armbian recovery SD card on hand
- [ ] Timer/stopwatch ready

## Recovery Procedures

If boot fails:

1. **No activity at all:** Board may be in MaskROM mode or eMMC is interfering
   - See: [MaskROM Recovery](MASKROM-RECOVERY.md)

2. **Some LED activity then stops:** U-Boot or kernel crash
   - Try different defconfig
   - Try different blob versions

3. **LED activity but no HDMI:** Display configuration issue
   - Verify with network ping
   - Kernel may be running without display

4. **Need to reset eMMC:** Use MaskROM to erase eMMC boot sectors
   - See: [MaskROM Recovery](MASKROM-RECOVERY.md)

## Summary Table

Use this table to track all attempts at a glance:

| # | Date | Defconfig | DDR | BL31 | LED | HDMI | Net | Result |
|---|------|-----------|-----|------|-----|------|-----|--------|
| 1 | | | | | | | | |
| 2 | | | | | | | | |
| 3 | | | | | | | | |
| 4 | | | | | | | | |
| 5 | | | | | | | | |

---

*Related: [MaskROM Recovery](MASKROM-RECOVERY.md) | [Flash Workflow](FLASH-WORKFLOW.md)*
