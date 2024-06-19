// Copyright 2019-2024 The Inspektor Gadget authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tests

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/sys/unix"

	gadgettesting "github.com/inspektor-gadget/inspektor-gadget/gadgets/testing"
	igtesting "github.com/inspektor-gadget/inspektor-gadget/pkg/testing"
	"github.com/inspektor-gadget/inspektor-gadget/pkg/testing/containers"
	igrunner "github.com/inspektor-gadget/inspektor-gadget/pkg/testing/ig"
	"github.com/inspektor-gadget/inspektor-gadget/pkg/testing/match"
	"github.com/inspektor-gadget/inspektor-gadget/pkg/testing/utils"
	eventtypes "github.com/inspektor-gadget/inspektor-gadget/pkg/types"
)

type traceMountEvent struct {
	eventtypes.CommonData

	Delta     uint64 `json:"delta"`
	Flags     uint64 `json:"flags"`
	Pid       int    `json:"pid"`
	Tid       int    `json:"tid"`
	MountNsID uint64 `json:"mount_ns_id"`
	Timestamp uint64 `json:"timestamp"`
	Ret       int    `json:"ret"`
	Comm      string `json:"comm"`
	Fs        string `json:"fs"`
	Src       string `json:"src"`
	Dest      string `json:"dest"`
	Data      string `json:"data"`
	OpStr     string `json:"op_str"`
}

func TestTraceMount(t *testing.T) {
	gadgettesting.RequireEnvironmentVariables(t)
	utils.InitTest(t)

	containerFactory, err := containers.NewContainerFactory(utils.Runtime)
	require.NoError(t, err, "new container factory")
	containerName := "test-trace-mount"
	containerImage := "docker.io/library/busybox:latest"

	var ns string
	containerOpts := []containers.ContainerOption{containers.WithContainerImage(containerImage)}

	switch utils.CurrentTestComponent {
	case utils.KubectlGadgetTestComponent:
		ns = utils.GenerateTestNamespaceName(t, "test-trace-mount")
		containerOpts = append(containerOpts, containers.WithContainerNamespace(ns))
	case utils.IgLocalTestComponent:
		containerOpts = append(containerOpts, containers.WithPrivileged())
	}

	testContainer := containerFactory.NewContainer(
		containerName,
		"while true; do mount /mnt /mnt; sleep 0.1; done",
		containerOpts...,
	)

	testContainer.Start(t)
	t.Cleanup(func() {
		testContainer.Stop(t)
	})

	var runnerOpts []igrunner.Option
	var testingOpts []igtesting.Option
	commonDataOpts := []utils.CommonDataOption{utils.WithContainerImageName(containerImage), utils.WithContainerID(testContainer.ID())}

	switch utils.CurrentTestComponent {
	case utils.IgLocalTestComponent:
		runnerOpts = append(runnerOpts, igrunner.WithFlags(fmt.Sprintf("-r=%s", utils.Runtime), "--timeout=5"))
	case utils.KubectlGadgetTestComponent:
		runnerOpts = append(runnerOpts, igrunner.WithFlags(fmt.Sprintf("-n=%s", ns), "--timeout=5"))
		testingOpts = append(testingOpts, igtesting.WithCbBeforeCleanup(utils.PrintLogsFn(ns)))
		commonDataOpts = append(commonDataOpts, utils.WithK8sNamespace(ns))
	}

	runnerOpts = append(runnerOpts, igrunner.WithValidateOutput(
		func(t *testing.T, output string) {
			expectedEntry := &traceMountEvent{
				CommonData: utils.BuildCommonData(containerName, commonDataOpts...),
				Comm:       "mount",
				OpStr:      "MOUNT",
				Src:        "/mnt",
				Dest:       "/mnt",
				Ret:        -int(unix.ENOENT),
				Data:       "",

				// Check only the existence of these fields
				Flags:     utils.NormalizedInt,
				Timestamp: utils.NormalizedInt,
				Delta:     utils.NormalizedInt,
				Pid:       utils.NormalizedInt,
				Tid:       utils.NormalizedInt,
				MountNsID: utils.NormalizedInt,
				Fs:        utils.NormalizedStr,
			}

			normalize := func(e *traceMountEvent) {
				utils.NormalizeCommonData(&e.CommonData)
				utils.NormalizeInt(&e.Timestamp)
				utils.NormalizeInt(&e.Delta)
				utils.NormalizeInt(&e.Flags)
				utils.NormalizeInt(&e.Pid)
				utils.NormalizeInt(&e.Tid)
				utils.NormalizeInt(&e.MountNsID)
				utils.NormalizeString(&e.Fs)
			}

			match.ExpectEntriesToMatch(t, output, normalize, expectedEntry)
		},
	))

	traceMountCmd := igrunner.New("trace_mount", runnerOpts...)

	igtesting.RunTestSteps([]igtesting.TestStep{traceMountCmd}, t, testingOpts...)
}
