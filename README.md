# Opni = AIOps for Kubernetes + Observability Tools

Opni currently features log anomaly detection for Kubernetes.

##### What does Opni give me?
* AI generated insights on your cluster's log messages
  * **Control Plane & etcd** insights
    * Pretrained models maintained by Rancher Labs
    * Only for RKE1, RKE2, k3s clusters
  * **Workload & application** insights
    * Automatically learns what steady-state is in your workloads & applications
    * For any Kubernetes cluster  
* Every log message sent to Opni will be marked as:
  * **Normal**
  * **Suspicious** - Operators may want to investigate
  * **Anomalous** - Operators definitely should investigate  
* Open Distro for Elasticsearch + Kibana 
  * Opni dashboard to consume log insights & explore logs 
  * Ability to setup & send alerts (slack/email/etc) based on Opni log insights

![alt text](https://opni-public.s3.us-east-2.amazonaws.com/opni-inside-cluster-diagram.png)

----

### Deprecation Notice
The Opnidemo API has been removed in the v0.2.1 release.

Opnictl is currently deprecated.
## Getting started with Opni

### Full Install Opni in your Kubernetes cluster:

#### Manifests install (recommended):
Prerequisites:
  * *1 Nvidia GPU* required if you want the AI to learn from your workloads
    * You will need the nvidia [k8s-device-plugin](https://github.com/NVIDIA/k8s-device-plugin) deployed to the cluster.  The simplest way is to use the nvidia operator.  [This blog post](https://rancher.com/blog/2020/get-up-and-running-with-nvidia-gpus) contains instructions on how to deploy it.  If you are using a host that already has nvidia drivers (e.g an EKS cluster) deploy the Daemonset with the following command:
      ```bash
      kubectl create -f https://raw.githubusercontent.com/NVIDIA/k8s-device-plugin/v0.9.0/nvidia-device-plugin.yml
      ```
    * The Opni manager has alpha support for automatically setting up GPU configuration.  For more details please visit the [GPU Configuration](https://opni.io/setup/gpu/) page.
  * Cert manager installed.  This can be installed with the following command:
    ```bash
    kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.5.3/cert-manager.yaml
    ```

Installation:
  1) Run the deploy helper script:
     ```bash
     curl -sfL https://raw.githubusercontent.com/rancher/opni/main/deploy/deploy.sh | sh -
     ```
     OR
     Deploy the manifests in [deploy/manifests](https://github.com/rancher/opni/tree/main/deploy/manifests) in order
  1) Deploy a matching logAdapter from [deploy/examples/logAdapters](https://github.com/rancher/opni/tree/main/deploy/examples/logAdapters)


*If you want to deploy the GPU service edit the opnicluster resource and set the [deploy option](https://github.com/rancher/opni/blob/main/deploy/manifests/20_cluster.yaml#L31) to true*
#### Opnictl install (deprecated)
* Download the `opnictl` binary from the [latest release](https://github.com/rancher/opni/releases/tag/v0.1.3)
* Install Opni using `opnictl`
  ```
  opnictl install
  opnictl create demo
  ```
  * Will use your current kubeconfig context
  * Cluster Hardware requirements: 4 vCPUs, 16GB RAM
    * *1 Nvidia GPU* required if you want the AI to learn from your workloads (recommended)

Consume insights from the Opni Dashboard in Kibana. You will need to expose the Kibana service or port forward to do this.

### Demo Opni in a sandbox environment
* What you need: **an Ubuntu VM with 4 vCPUs & 16 GB RAM**
* 1-command installer that creates an RKE2 cluster with Opni installed and simulates [a failure](https://github.com/rancher/opni-docs/blob/22ed683e2b9e810b04561967d65682654350d787/quickstart_files/install_opni.sh#L72)
  ```
  curl -sfL https://raw.githubusercontent.com/rancher/opni-docs/main/quickstart_files/install_opni.sh | sh -
  ```

  * Experiment with injecting your own control plane failures to see how Opni responds
    * Can refer to [these failures](https://github.com/rancher/opni-docs/blob/main/examples/fault-injection.md) and this [anomaly injection script](https://github.com/rancher/opni-docs/blob/main/quickstart_files/errors_injection.sh) as starting points

The default username and password is admin/admin You must be in the Global Tenant mode if you are not already. Navigate to `Dashboard` then `Opni Logs Dashboard`.
 
----

**Watch a demo of Opni:**

[![](https://opni-public.s3.us-east-2.amazonaws.com/opni_youtube_gh.png)](https://youtu.be/DQVBwMaO_o0)
____
#### What's next?

 * v0.1.1 (Released) allows you to view Opni's log anomaly insights **only** on a demo environment created on a VM
 * v0.1.2 (Released) allows you install Opni into your existing Kubernetes cluster and consume log insights from it
 * v0.1.3 (August 2021) - only 1 GPU required, changes to the Opni operator, log anomaly optimizations
 * v0.2.0 (Fall 2021) will introduce a custom UI, AI applied to metrics, kubernetes events, audit logs, and more! 


![alt text](https://opni-public.s3.us-east-2.amazonaws.com/Opni-user-scenarios.png)

----


## License

Copyright (c) 2014-2020 [Rancher Labs, Inc.](http://rancher.com)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

