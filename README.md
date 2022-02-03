# oc adm apiserver

Is an OpenShift `oc` command plugin that help cluster administrators or developers monitor Kubernetes API server usage and performance.
Currently only `requests` subcommand is implemented.

## Installation

`$ go install github.com/mfojtik/oc-adm-apiserver`

For debugging/developing, you can run `make` command in Git repo and copy the binary to your $PATH.

### Sub-commands

* `requests` - periodically query the OpenShift `apirequestcount` resources and format the output to human-readable form:

Example:
```
Resource Name                                                         Requests Last Hour     Requests Total     Top Node
configmaps.v1                                                                742                    21187              10.0.218.204 [742 requests] (system:serviceaccount:openshift-monitoring:cluster-monitoring-operator/Go-http-client/2.0 [158 reqests](GETx87))
ingresses.v1.config.openshift.io                                             717                    717                10.0.218.204 [427 requests] (system:serviceaccount:openshift-console-operator:console-operator/Go-http-client/2.0 [385 reqests](GETx383))
secrets.v1                                                                   300                    8744               10.0.132.235 [300 requests] (system:serviceaccount:openshift-kube-controller-manager-operator:kube-controller-manager-operator/Go-http-client/2.0 [131 reqests](GETx130))
serviceaccounts.v1                                                           173                    5972               10.0.218.204 [173 requests] (system:serviceaccount:openshift-kube-scheduler-operator:openshift-kube-scheduler-operator/Go-http-client/2.0 [47 reqests](GETx46))
apiservices.v1.apiregistration.k8s.io                                        151                    1935               10.0.218.204 [148 requests] (system:serviceaccount:openshift-apiserver-operator:openshift-apiserver-operator/Go-http-client/2.0 [108 reqests](GETx108))
subjectaccessreviews.v1.authorization.k8s.io                                 150                    9212               10.0.218.204 [150 requests] (system:serviceaccount:openshift-apiserver:openshift-apiserver-sa/Go-http-client/2.0 [142 reqests](CREATEx142))
services.v1                                                                  141                    6006               10.0.132.235 [141 requests] (system:serviceaccount:openshift-monitoring:prometheus-k8s/Prometheus/2.32.1 [39 reqests](WATCHx39))
netnamespaces.v1.network.openshift.io                                        137                    137                10.0.100.138 [85 requests] (system:serviceaccount:openshift-sdn:sdn-controller/openshift-sdn-controller/v0.0.0 [67 reqests](CREATEx61))
prometheusrules.v1.monitoring.coreos.com                                     130                    1560               10.0.218.204 [130 requests] (system:serviceaccount:openshift-kube-apiserver-operator:kube-apiserver-operator/cluster-kube-apiserver-operator/v0.0.0 [84 reqests](GETx84))
rolebindings.v1.rbac.authorization.k8s.io                                    116                    5019               10.0.218.204 [116 requests] (system:serviceaccount:openshift-kube-apiserver-operator:kube-apiserver-operator/Go-http-client/2.0 [42 reqests](GETx42))
events.v1                                                                    111                    2668               10.0.132.235 [95 requests] (system:serviceaccount:openshift-machine-config-operator:default/machine-config-operator/v0.0.0 [64 reqests](CREATEx64))
networks.v1.config.openshift.io                                              98                     98                 10.0.143.249 [58 requests] (system:serviceaccount:openshift-dns-operator:dns-operator/dns-operator/v0.0.0 [20 reqests](GETx20))
deployments.v1.apps                                                          86                     2514               10.0.218.204 [86 requests] (system:serviceaccount:openshift-monitoring:cluster-monitoring-operator/Go-http-client/2.0 [44 reqests](GETx28))
apirequestcounts.v1.apiserver.openshift.io                                   83                     5133               10.0.218.204 [81 requests] (system:apiserver/kube-apiserver/v1.23.3+b63be7f [81 reqests](GETx58))
podnetworkconnectivitychecks.v1alpha1.controlplane.operator.openshift.io     75                     291                10.0.132.235 [61 requests] (system:serviceaccount:openshift-network-operator:default/Go-http-client/2.0 [54 reqests](GETx54))
roles.v1.rbac.authorization.k8s.io                                           73                     4207               10.0.218.204 [73 requests] (system:serviceaccount:openshift-monitoring:cluster-monitoring-operator/Go-http-client/2.0 [38 reqests](UPDATEx16))
endpoints.v1                                                                 64                     1709               10.0.132.235 [64 requests] (system:serviceaccount:openshift-monitoring:prometheus-k8s/Prometheus/2.32.1 [41 reqests](WATCHx41))
flowschemas.v1beta2.flowcontrol.apiserver.k8s.io                             52                     1240               10.0.132.235 [39 requests] (system:apiserver/kube-apiserver/v1.23.3+b63be7f [39 reqests](GETx36))
authentications.v1.operator.openshift.io                                     49                     695                10.0.218.204 [49 requests] (system:serviceaccount:openshift-authentication-operator:authentication-operator/Go-http-client/2.0 [49 reqests](GETx48))
tokenreviews.v1.authentication.k8s.io                                        41                     1801               10.0.218.204 [41 requests] (system:serviceaccount:openshift-cluster-csi-drivers:aws-ebs-csi-driver-controller-sa/kube-rbac-proxy/v0.0.0 [8 reqests](CREATEx8))

```