// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	_ "embed"
	"errors"
	"fmt"
	"github.com/siderolabs/go-copy/copy"
	"github.com/siderolabs/talos/pkg/machinery/overlay"
	"github.com/siderolabs/talos/pkg/machinery/overlay/adapter"
	"golang.org/x/sys/unix"
	"os"
	"path/filepath"
)

const (
	// Standard offset for u-boot-rockchip.bin (sector 64)
	ubootOffset int64 = 512 * 64
	// Offset for vendor U-Boot uboot.img (sector 16384 = 8MB)
	ubootImgOffset int64 = 512 * 16384
)

func main() {
	adapter.Execute[rk3588ExtraOpts](&RK3588Installer{})
}

type RK3588Installer struct{}

type rk3588ExtraOpts struct {
	Board   string `json:"board"`
	Chipset string `json:"chipset"`
}

func ChipsetName(o rk3588ExtraOpts) string {
	if o.Chipset != "" {
		return o.Chipset
	}
	switch o.Board {
	case "rock-5a":
		return "rk3588s"
	case "rock-5b":
		return "rk3588"
	case "nanopi-m6":
		return "rk3588s"
	}
	return ""
}

func (i *RK3588Installer) GetOptions(extra rk3588ExtraOpts) (overlay.Options, error) {
	if extra.Board == "" {
		return overlay.Options{}, errors.New("board variant required")
	}

	kernelArgs := []string{
		"sysctl.kernel.kexec_load_disabled=1",
		"talos.dashboard.disabled=1",
		"slab_nomerge",
		"earlycon=uart8250,mmio32,0xfeb50000",
		"console=ttyFIQ0,1500000n8",
		"consoleblank=0",
		"console=ttyS2,1500000n8",
		"console=tty1",
		"loglevel=7",
		"cgroup_enable=cpuset",
		"swapaccount=1",
		"irqchip.gicv3_pseudo_nmi=0",
		"coherent_pool=2M",
	}

	return overlay.Options{
		Name:       extra.Board,
		KernelArgs: kernelArgs,
		PartitionOptions: overlay.PartitionOptions{
			Offset: 2048 * 10,
		},
	}, nil
}

func (i *RK3588Installer) Install(options overlay.InstallOptions[rk3588ExtraOpts]) error {
	if options.ExtraOptions.Board == "" {
		return errors.New("board variant required")
	}
	if options.ExtraOptions.Chipset == "" {
		return errors.New("chipset variant required")
	}

	var f *os.File
	f, err := os.OpenFile(options.InstallDisk, os.O_RDWR|unix.O_CLOEXEC, 0o666)
	if err != nil {
		return fmt.Errorf("opening install disk: %w", err)
	}
	defer f.Close() //nolint:errcheck

	// NanoPi M6 uses vendor U-Boot format with two files at different offsets
	if options.ExtraOptions.Board == "nanopi-m6" {
		// Write idbloader.img at sector 64 (same offset as standard u-boot)
		idbloader, err := os.ReadFile(filepath.Join(options.ArtifactsPath,
			fmt.Sprintf("arm64/u-boot/%s/idbloader.img", options.ExtraOptions.Board)))
		if err != nil {
			return fmt.Errorf("reading idbloader: %w", err)
		}
		if _, err = f.WriteAt(idbloader, ubootOffset); err != nil {
			return fmt.Errorf("writing idbloader: %w", err)
		}

		// Write uboot.img at sector 16384 (8MB offset)
		ubootImg, err := os.ReadFile(filepath.Join(options.ArtifactsPath,
			fmt.Sprintf("arm64/u-boot/%s/uboot.img", options.ExtraOptions.Board)))
		if err != nil {
			return fmt.Errorf("reading uboot.img: %w", err)
		}
		if _, err = f.WriteAt(ubootImg, ubootImgOffset); err != nil {
			return fmt.Errorf("writing uboot.img: %w", err)
		}
	} else {
		// Standard boards use single u-boot-rockchip.bin
		uboot, err := os.ReadFile(filepath.Join(options.ArtifactsPath,
			fmt.Sprintf("arm64/u-boot/%s/u-boot-rockchip.bin", options.ExtraOptions.Board)))
		if err != nil {
			return fmt.Errorf("reading u-boot: %w", err)
		}
		if _, err = f.WriteAt(uboot, ubootOffset); err != nil {
			return fmt.Errorf("writing u-boot: %w", err)
		}
	}

	// NB: In the case that the block device is a loopback device, we sync here
	// to ensure that the file is written before the loopback device is
	// unmounted.
	err = f.Sync()
	if err != nil {
		return err
	}

	dtb := filepath.Join("rockchip", fmt.Sprintf("%s-%s.dtb", ChipsetName(options.ExtraOptions), options.ExtraOptions.Board))
	src := filepath.Join(options.ArtifactsPath, "arm64/dtb", dtb)
	dst := filepath.Join(options.MountPrefix, "/boot/EFI/dtb", dtb)

	err = os.MkdirAll(filepath.Dir(dst), 0o600)
	if err != nil {
		return err
	}

	return copy.File(src, dst)
}
