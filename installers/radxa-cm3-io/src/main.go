// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	_ "embed"
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/siderolabs/go-copy/copy"
	"github.com/siderolabs/talos/pkg/machinery/overlay"
	"github.com/siderolabs/talos/pkg/machinery/overlay/adapter"
	"golang.org/x/sys/unix"
)

const (
	dtb = "rockchip/rk3566-radxa-cm3-io.dtb"
)

func main() {
	adapter.Execute(context.Background(), &radxaCM3Io{})
}

type radxaCM3Io struct{}

type radxaCM3IoExtraOptions struct{}

func (i *radxaCM3Io) GetOptions(ctx context.Context, extra radxaCM3IoExtraOptions) (overlay.Options, error) {
	kernelArgs := []string{
		"cma=128MB",
		"console=tty0",
		"console=ttyS2,115200",
		"sysctl.kernel.kexec_load_disabled=1",
		"talos.dashboard.disabled=1",
	}
	return overlay.Options{
		Name: "radxa-cm3-io",
		KernelArgs: kernelArgs,
		PartitionOptions: overlay.PartitionOptions{
			Offset: 2048 * 10,
		},
	}, nil
}

func (i *radxaCM3Io) Install(ctx context.Context, options overlay.InstallOptions[radxaCM3IoExtraOptions]) error {
	uBootBin := filepath.Join(options.ArtifactsPath, "arm64/u-boot/radxa-cm3-io/u-boot-rockchip.bin")

	if err := uBootLoaderInstall(uBootBin, options.InstallDisk); err != nil {
		return err
	}

	src := filepath.Join(options.ArtifactsPath, "arm64/dtb", dtb)
	dst := filepath.Join(options.MountPrefix, "boot/EFI/dtb", dtb)

	if err := copyFileAndCreateDir(src, dst); err != nil {
		return err
	}

	return nil
}

func copyFileAndCreateDir(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o600); err != nil {
		return err
	}

	return copy.File(src, dst)
}

func uBootLoaderInstall(uBootBin, installDisk string) error {
	f, err := os.OpenFile(installDisk, os.O_RDWR|unix.O_CLOEXEC, 0o666)
	if err != nil {
		return fmt.Errorf("failed to open %s: %w", installDisk, err)
	}

	defer f.Close() //nolint:errcheck

	uboot, err := os.ReadFile(uBootBin)
	if err != nil {
		return err
	}

	if _, err = f.WriteAt(uboot, off); err != nil {
		return err
	}

	// NB: In the case that the block device is a loopback device, we sync here
	// to ensure that the file is written before the loopback device is
	// unmounted.
	return f.Sync()
}
