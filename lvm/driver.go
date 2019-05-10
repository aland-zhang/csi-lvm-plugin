/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package lvm

import (
	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/zdnscloud/gok8s/cache"
	"github.com/zdnscloud/gok8s/client"

	"github.com/kubernetes-csi/drivers/pkg/csi-common"
	"github.com/zdnscloud/cement/log"
)

type lvm struct {
	driver *csicommon.CSIDriver
	client client.Client

	ids *identityServer
	ns  *nodeServer
	cs  *controllerServer

	cap   []*csi.VolumeCapability_AccessMode
	cscap []*csi.ControllerServiceCapability
}

var (
	lvmDriver     *lvm
	vendorVersion = "0.3.0"
)

func GetLVMDriver(client client.Client) *lvm {
	return &lvm{client: client}
}

func (lvm *lvm) Run(driverName, nodeID, endpoint string, vgName string, cache cache.Cache) {
	lvm.driver = csicommon.NewCSIDriver(driverName, vendorVersion, nodeID)
	if lvm.driver == nil {
		log.Fatalf("Failed to initialize CSI Driver.")
	}

	lvm.driver.AddControllerServiceCapabilities([]csi.ControllerServiceCapability_RPC_Type{csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME})
	lvm.driver.AddVolumeCapabilityAccessModes([]csi.VolumeCapability_AccessMode_Mode{csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER})

	lvm.ids = NewIdentityServer(lvm.driver)
	lvm.ns = NewNodeServer(lvm.driver, lvm.client, nodeID, vgName)
	lvm.cs = NewControllerServer(lvm.driver, lvm.client, vgName, cache)

	server := csicommon.NewNonBlockingGRPCServer()
	server.Start(endpoint, lvm.ids, lvm.cs, lvm.ns)
	server.Wait()
}