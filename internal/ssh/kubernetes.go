// Copyright © 2019 Alvaro Saurin
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

package ssh

import (
	"fmt"
	"log"
	"net/url"
	"os/exec"
	"strings"
)

// isValidUrl tests a string to determine if it is a url or not.
func isValidUrl(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	} else {
		return true
	}
}

// DoLocalKubectl runs a local kubectl command
func DoLocalKubectl(kubeconfig string, args ...string) ApplyFunc {
	kubectl, err := exec.LookPath("kubectl")
	if err != nil {
		log.Fatal("kubectl not available")
	}

	f := append([]string{fmt.Sprintf("--kubeconfig=%s", kubeconfig)}, args...)
	return DoLocalExec(kubectl, f...)
}

// DoRemoteKubectl runs a local kubectl command in a remote machine
func DoRemoteKubectl(kubeconfig string, args ...string) ApplyFunc {
	cmd := append([]string{"kubectl", fmt.Sprintf("--kubeconfig=%s", kubeconfig)}, args...)
	return DoExec(strings.Join(cmd, " "))
}

// DoLocalKubectlApply applies some manifests with a local kubectl
func DoLocalKubectlApply(kubeconfig string, manifests []string) ApplyFunc {
	loaders := []Applyer{}
	for _, manifest := range manifests {
		if len(manifest) == 0 {
			continue
		}
		loaders = append(loaders, DoLocalKubectl(kubeconfig, "apply", "-f", manifest))
	}
	return DoComposed(loaders...)
}

// DoRemoteKubectlApply applies some manifests with a remote kubectl
func DoRemoteKubectlApply(kubeconfig string, manifests []string) ApplyFunc {
	actions := []Applyer{}

	// upload the local kubeconfig to some temporary remote file
	remoteKubeconfig, err := randomPath("tmpfile", "tmp")
	if err != nil {
		panic(err)
	}
	actions = append(actions, DoUploadFileToFile(kubeconfig, remoteKubeconfig))

	for _, manifest := range manifests {
		if len(manifest) == 0 {
			continue
		}

		if isValidUrl(manifest) {
			// it is an URL: just run the `kubectl apply`
			actions = append(actions, DoRemoteKubectl(remoteKubeconfig, "apply", "-f", manifest))
		} else {
			// it is a file: upload the file to a temporary, remote file and then `kubectl apply -f` it
			remoteManifest, err := randomPath("tmpfile", "tmp")
			if err != nil {
				panic(err)
			}

			actions = append(actions,
				DoUploadFileToFile(manifest, remoteManifest),
				DoRemoteKubectl(remoteKubeconfig, "apply", "-f", remoteManifest))
		}
	}
	return DoComposed(actions...)
}
