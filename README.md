# KubePing

It is a Kubernetes solution designed to monitor the availability of external endpoints from each node of the Kubernetes cluster over TCP, HTTP, and ICMP. It exports Prometheus metrics and has a user-friendly web interface that helps you save time instead of making telnet/curl on each node.

<p align="center">
    <img src="kubeping-low.drawio.svg">
</p>

## Use case
Imagine a situation where your pods were evicted because of a node failure. The pods were then relocated to new nodes that had recently been added to the cluster. Unfortunately, this caused errors and resulted in service unavailability due to a lack of access to essential external endpoints. It was later discovered that the security department had failed to apply the appropriate access rules to the new cluster nodes.

This is just one scenario that highlights how KubePing can help you identify potential issues before they escalate into major problems.

## How It Works
The solution runs a lightweight DaemonSet in Kubernetes, ensuring that each node has a running instance. These instances probe external endpoints over:

__TCP__ – Checking port availability (e.g., database:5432, api:443)\
__HTTP__ – Ensuring services respond with the expected status codes\
__ICMP (Ping)__ – Verifying network reachability

The results are aggregated and exposed as Prometheus metrics:
```
probe_result{address="api.example.com:8080", instance="worker-node-1", job="kubeping", module="tcp", target="target1"}=1
probe_result{address="api.example.com:8080", instance="worker-node-2", job="kubeping", module="tcp", target="target1"}=0
probe_result{address="api.example.com:8080", instance="worker-node-3", job="kubeping", module="tcp", target="target1"}=1

probe_result{address="https://example.com", instance="worker-node-1", job="kubeping", module="http", target="target2"}=0
probe_result{address="https://example.com", instance="worker-node-2", job="kubeping", module="http", target="target2"}=1
probe_result{address="https://example.com", instance="worker-node-3", job="kubeping", module="http", target="target2"}=1

probe_result{address="192.168.0.1", instance="worker-node-1", job="kubeping", module="icmp", target="target3"}=1
probe_result{address="192.168.0.1", instance="worker-node-2", job="kubeping", module="icmp", target="target3"}=0
probe_result{address="192.168.0.1", instance="worker-node-3", job="kubeping", module="icmp", target="target3"}=1
```

And instead of SSH-ing into nodes, you can simply visit the web UI, where you can test connectivity to any endpoint:
<p align="center">
    <img src="kubeping-web.gif" width="70%" height="70%">
</p>

Here is how KubePing can be integrated into your workflow:
<p align="center">
    <img src="kubeping-high.drawio.svg">
</p>

## Installation
### Helm
Clone repository
```
git clone https://github.com/teymurgahramanov/kubeping.git && cd kubeping
```
Install Helm chart
```
helm upgrade --install kubeping ./helm
```
Test Web UI
```
kubectl port-forward svc/kubeping-web 8000:8000
```
To configure the exporter with static targets, refer to [values.yaml](./helm/values.yaml). Here is an example:
```yaml
exporter:
  config:
    exporter:
      defaultProbeInterval: 31
      defaultProbeTimeout: 13
    targets:
      target1:
        address: api.example.com:8080
        module: tcp
        timeout: 15
      target2:
        address: https://example.com
        module: http
        interval: 60
      target3:
        address: 192.168.0.1
        module: icmp
```

### Prometheus
Example job configuration:
```yaml
- job_name: kubeping
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
