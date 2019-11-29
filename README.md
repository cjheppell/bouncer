# bouncer

A Kubernetes Operator to support dynamically exposing NodePort services across a range of cloud providers

## What problem does this solve?

Inspired by [this StackOverflow question](https://stackoverflow.com/q/58341852/5235388), there doesn't seem to be a good way to dynamically expose non-HTTP services in Kubernetes for public consumption.

NodePort's are a way to do this, but they require public access to the nodes in the cluster for users outside of the cluster to consume them.

## How does it solve it?

Bouncer will effectively wraps a "NodePort" service, with knowledge of the cloud provider in use, such that the NodePort can be exposed for public consumption without the user needing to configure things such as firewalls, or load balancers, etc.

We imagine this might be solved via the following possible mechanisms:
- Exposing NodePorts in the cluster firewall so they're publicly accessible
- Automatically provisioning a cloud provider load balancer and mapping them to the NodePorts for provisioned pods
- Creating a TCP/UDP proxy for any of the pods within the cluster, to enable them to be publicly accessible