// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	_ "embed"
	"os"
	"path/filepath"

	"github.com/siderolabs/go-copy/copy"
	"github.com/siderolabs/talos/pkg/machinery/overlay"
	"github.com/siderolabs/talos/pkg/machinery/overlay/adapter"
)

const (
	dtb = "rockchip/rk3566-radxa-cm3-io.dtb"
)

func main() {
	adapter.Execute(&radxaCM3Io{})
}

type radxaCM3Io struct{}

type radxaCM3IoExtraOptions struct{}

func (i *radxaCM3Io) GetOptions(extra radxaCM3IoExtraOptions) (overlay.Options, error) {
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

func (i *radxaCM3Io) Install(options overlay.InstallOptions[radxaCM3IoExtraOptions]) error {
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
