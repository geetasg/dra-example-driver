/*
 * Copyright 2023 The Kubernetes Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"fmt"
	"math/rand"
	"os"

	resourceapi "k8s.io/api/resource/v1beta1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/utils/ptr"

	"github.com/google/uuid"
)

func randomDriverVersion() string {
	// Example publishes 8 devices. driverVersion will match for atleast 2 devices since we only sample between 1 and 5
	major := 1 + rand.Intn(5) // Random number between 1 and 5
	return fmt.Sprintf("%d.0.0", major)
}

func enumerateAllPossibleDevices(numGPUs int) (AllocatableDevices, error) {
	seed := os.Getenv("NODE_NAME")
	uuids := generateUUIDs(seed, numGPUs)

	alldevices := make(AllocatableDevices)
	for i, uuid := range uuids {
		device := resourceapi.Device{
			Name: fmt.Sprintf("gpu-%d", i),
			Basic: &resourceapi.BasicDevice{
				Attributes: map[resourceapi.QualifiedName]resourceapi.DeviceAttribute{
					"index": {
						IntValue: ptr.To(int64(i)),
					},
					"uuid": {
						StringValue: ptr.To(uuid),
					},
					"model": {
						StringValue: ptr.To("LATEST-GPU-MODEL"),
					},
					"driverVersion": {
						VersionValue: ptr.To(randomDriverVersion()),
					},
				},
				Capacity: map[resourceapi.QualifiedName]resourceapi.DeviceCapacity{
					"memory": {
						Value: resource.MustParse("80Gi"),
					},
				},
			},
		}
		alldevices[device.Name] = device
	}
	return alldevices, nil
}

func generateUUIDs(seed string, count int) []string {
	rand := rand.New(rand.NewSource(hash(seed)))

	uuids := make([]string, count)
	for i := 0; i < count; i++ {
		charset := make([]byte, 16)
		rand.Read(charset)
		uuid, _ := uuid.FromBytes(charset)
		uuids[i] = "gpu-" + uuid.String()
	}

	return uuids
}

func hash(s string) int64 {
	h := int64(0)
	for _, c := range s {
		h = 31*h + int64(c)
	}
	return h
}
