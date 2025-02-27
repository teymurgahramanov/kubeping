# KubePing

It is a ready-to-use solution for Kubernetes designed to monitor the availability of external endpoints from each node of the Kubernetes cluster over TCP, HTTP, and ICMP and export Prometheus metrics.

<p align="center">
    <img src="kubeping-low.drawio.svg">
</p>

It also has a user-friendly web interface that helps you to save time instead of making telnet/curl on each node.
<p align="center">
    <img src="kubeping-web.gif">
</p>

## Use case
Imagine a situation where your pods were evicted because of a node failure. The pods were then relocated to new nodes that had recently been added to the cluster. Unfortunately, this caused errors and resulted in service unavailability due to a lack of access to essential external endpoints. It was later discovered that the security department had failed to apply the appropriate access rules to the new cluster nodes.

This is just one scenario that highlights how KubePing can help you identify potential issues before they escalate into major problems.

Here is the how KubePing can be integrated in your workflow:
<p align="center">
    <img src="kubeping-high.drawio.svg">
</p>

## Installation
### Helm
To install the Helm chart, refer to the [values.yaml](./helm/values.yaml) file.

### Prometheus
Configure the Prometheus job:
```yaml
- job_name: kubeping-exporter
  kubernetes_sd_configs:
    - role: endpoints
  relabel_configs:
    - source_labels: [__meta_kubernetes_endpoints_name]
      regex: kubeping-exporter
      action: keep
    - source_labels: [__meta_kubernetes_endpoint_node_name]
      action: replace
      target_label: instance
```
