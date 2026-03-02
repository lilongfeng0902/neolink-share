# Kubernetes Prometheus 监控系统安装文档

本文档介绍如何在 Kubernetes 集群中安装 kube-prometheus 监控系统。

## 目录

- [环境要求](#环境要求)
- [组件介绍](#组件介绍)
- [安装步骤](#安装步骤)
- [注意事项](#注意事项)
- [访问监控面板](#访问监控面板)
- [故障排查](#故障排查)
- [卸载说明](#卸载说明)

---

## 环境要求

- Kubernetes 版本要求：1.30+
- kubectl 已配置并可访问集群
- 足够的集群资源（建议：CPU 4核+，内存 8GB+）
- kubeconfig 文件路径：`C:\Users\li.longfeng\.kube\config`

---

## 组件介绍

### 核心组件

| 组件 | 作用 | 资源类型 |
|------|------|----------|
| **Prometheus Operator** | 管理 Prometheus、Alertmanager 等自定义资源的控制器，简化监控系统的部署和运维 | Deployment |
| **Prometheus** | 核心监控系统，负责采集和存储时序数据，并提供查询接口 | StatefulSet |
| **Alertmanager** | 处理 Prometheus 发送的告警，支持告警分组、去重、路由和通知 | StatefulSet |
| **Grafana** | 可视化监控面板，提供丰富的图表和仪表板 | Deployment |
| **Node Exporter** | 收集节点级别的硬件和操作系统指标（DaemonSet 模式部署） | DaemonSet |
| **Kube State Metrics** | 收集 Kubernetes 对象的状态指标（Pod、Deployment、Node 等） | Deployment |
| **cAdvisor** | 收集容器级别的资源使用指标（CPU、内存、文件系统、网络） | DaemonSet |
| **Blackbox Exporter** | 探针式监控，支持 HTTP、HTTPS、DNS、TCP、ICMP 等协议探测 | Deployment |
| **Prometheus Adapter** | 将 Prometheus 指标转换为 Kubernetes API，支持 HPA 自动扩缩容 | Deployment |

### 自定义资源定义（CRD）

| CRD | 作用 |
|-----|------|
| **Prometheus** | 定义 Prometheus 实例配置 |
| **Alertmanager** | 定义 Alertmanager 实例配置 |
| **ServiceMonitor** | 声明式服务发现，定义服务监控目标 |
| **PodMonitor** | 声明式 Pod 监控，定义 Pod 监控目标 |
| **PrometheusRule** | 定义 Prometheus 和 Alertmanager 的告警规则 |
| **Probe** | 定义 Blackbox Exporter 的探测目标 |
| **AlertmanagerConfig** | 定义 Alertmanager 的路由和通知配置 |

---

## 安装步骤

### 前置检查

确认 kubectl 可用：

```bash
kubectl --kubeconfig="C:\Users\li.longfeng\.kube\config" cluster-info
kubectl --kubeconfig="C:\Users\li.longfeng\.kube\config" version
```

### 步骤 1：安装 CRD 和 Namespace

创建自定义资源定义和 monitoring 命名空间：

```bash
# 使用 server-side apply 避免 CRD annotation 过大问题
kubectl --kubeconfig="C:\Users\li.longfeng\.kube\config" apply --server-side=true -f manifests/setup/
```

**预期输出：**
```
customresourcedefinition.apiextensions.k8s.io/alertmanagerconfigs.monitoring.coreos.com serverside-applied
customresourcedefinition.apiextensions.k8s.io/alertmanagers.monitoring.coreos.com serverside-applied
customresourcedefinition.apiextensions.k8s.io/podmonitors.monitoring.coreos.com serverside-applied
...
namespace/monitoring serverside-applied
```

**重要说明：**
- 必须使用 `--server-side=true` 参数，因为部分 CRD 文件较大（>500KB），会超过 Kubernetes annotation 262144 字节的限制

### 步骤 2：安装监控组件

安装所有监控组件：

```bash
kubectl --kubeconfig="C:\Users\li.longfeng\.kube\config" apply -f manifests/
```

**预期输出：**
```
alertmanager.monitoring.coreos.com/main created
deployment.apps/grafana created
prometheus.monitoring.coreos.com/k8s created
...
```

### 步骤 3：验证安装状态

等待所有 Pod 启动完成：

```bash
kubectl --kubeconfig="C:\Users\li.longfeng\.kube\config" get pods -n monitoring
```

**正常状态示例：**
```
NAME                                   READY   STATUS    RESTARTS   AGE
alertmanager-main-0                    2/2     Running   0          1m
alertmanager-main-1                    2/2     Running   0          1m
alertmanager-main-2                    2/2     Running   0          1m
blackbox-exporter-78f88d798c-xxx       3/3     Running   0          1m
grafana-7f6476f698-xxx                 1/1     Running   0          1m
kube-state-metrics-5f8f8d6cd7-xxx      3/3     Running   0          1m
node-exporter-xxx                      2/2     Running   0          1m
prometheus-adapter-5c89794f6c-xxx      1/1     Running   0          1m
prometheus-k8s-0                       2/2     Running   0          1m
prometheus-k8s-1                       2/2     Running   0          1m
prometheus-operator-5955df45f7-xxx     2/2     Running   0          1m
```

查看服务状态：

```bash
kubectl --kubeconfig="C:\Users\li.longfeng\.kube\config" get svc -n monitoring
```

---

## 注意事项

### 1. CRD 安装必须使用 Server-Side Apply

**问题：** 直接使用 `kubectl apply` 创建 CRD 会报错：
```
CustomResourceDefinition is invalid: metadata.annotations: Too long: must have at most 262144 bytes
```

**原因：** 部分 CRD 文件超过 500KB，kubectl apply 会将完整配置存储在 annotation 中，超过限制。

**解决方案：** 必须使用 `--server-side=true` 参数：
```bash
kubectl --kubeconfig="C:\Users\li.longfeng\.kube\config" apply --server-side=true -f manifests/setup/
```

### 2. 镜像仓库配置

当前配置使用私有镜像仓库：`220.181.114.186:80/system_containers/`

**完整镜像清单：**

| 组件 | 镜像:标签 | 作用 |
|------|-----------|------|
| **Alertmanager** | alertmanager:v0.30.1 | 告警处理 |
| **Prometheus** | prometheus:v3.9.1 | 指标采集和存储 |
| **Prometheus Operator** | prometheus-operator:v0.88.0 | 管理 CRD 控制器 |
| **Grafana** | grafana:12.3.1 | 可视化面板 |
| **Node Exporter** | node-exporter:v1.10.2 | 节点指标采集 |
| **Kube State Metrics** | kube-state-metrics:v2.18.0 | K8s 对象状态指标 |
| **cAdvisor** | cadvisor:v0.51.0 | 容器资源指标 |
| **Blackbox Exporter** | blackbox-exporter:v0.28.0 | 探针式监控 |
| **Prometheus Adapter** | prometheus-adapter:v0.12.0 | 指标 API 适配器 |
| **kube-rbac-proxy** | kube-rbac-proxy:v0.20.2 | RBAC 代理（边车容器）|
| **prometheus-config-reloader** | prometheus-config-reloader:v0.88.0 | 配置热重载 |
| **configmap-reload** | configmap-reload:v0.15.0 | ConfigMap 重载 |

如果遇到镜像拉取失败（ImagePullBackOff），需要：
- 确保镜像仓库可访问
- 配置正确的 imagePullSecret 认证信息
- 如需使用公共镜像仓库，需修改所有 deployment/daemonset/statefulset 中的镜像地址

**验证镜像使用情况：**
```bash
kubectl --kubeconfig="C:\Users\li.longfeng\.kube\config" get pods -n monitoring -o jsonpath='{range .items[*]}{.metadata.name}{"\t"}{range .spec.containers[*]}{.image}{"\n"}{end}{end}' | sort -u
```

### 3. 资源要求

完整安装后资源占用（参考值）：
- CPU：约 2-4 核
- 内存：约 4-8 GB
- 存储：Prometheus 默认使用持久化存储，确保 PV 可用

### 4. 网络策略

部署文件中包含 NetworkPolicy，限制了 Pod 间的网络访问。如遇到访问问题，可以：
- 检查 CNI 网络插件是否支持 NetworkPolicy
- 临时删除 NetworkPolicy 进行测试

### 5. 监控数据持久化

- Prometheus 数据默认存储在 EmptyDir（重启会丢失）
- 生产环境建议配置持久化存储（PVC）

---

## 访问监控面板

### 方法 1：Port Forward（推荐用于测试）

**访问 Grafana：**
```bash
kubectl --kubeconfig="C:\Users\li.longfeng\.kube\config" port-forward -n monitoring svc/grafana 3000:3000
```
浏览器访问：http://localhost:3000
- 默认用户名：`admin`
- 默认密码：`admin`（首次登录会提示修改）

**访问 Prometheus：**
```bash
kubectl --kubeconfig="C:\Users\li.longfeng\.kube\config" port-forward -n monitoring svc/prometheus-k8s 9090:9090
```
浏览器访问：http://localhost:9090

**访问 Alertmanager：**
```bash
kubectl --kubeconfig="C:\Users\li.longfeng\.kube\config" port-forward -n monitoring svc/alertmanager-main 9093:9093
```
浏览器访问：http://localhost:9093

### 方法 2：Ingress（生产环境推荐）

配置 Ingress 规则暴露服务，示例：

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: grafana-ingress
  namespace: monitoring
spec:
  rules:
  - host: grafana.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: grafana
            port:
              number: 3000
```

### 方法 3：LoadBalancer/NodePort

修改 Service 类型：

```bash
kubectl --kubeconfig="C:\Users\li.longfeng\.kube\config" patch svc grafana -n monitoring -p '{"spec":{"type":"LoadBalancer"}}'
```

---

## 故障排查

### Pod 无法启动（ImagePullBackOff）

**现象：**
```bash
kubectl get pods -n monitoring
# STATUS: ImagePullBackOff
```

**排查步骤：**

1. 查看 Pod 详细信息：
```bash
kubectl --kubeconfig="C:\Users\li.longfeng\.kube\config" describe pod <pod-name> -n monitoring
```

2. 检查镜像地址是否正确：
```bash
kubectl --kubeconfig="C:\Users\li.longfeng\.kube\config" get pod <pod-name> -n monitoring -o jsonpath='{.spec.containers[*].image}'
```

3. 确认镜像仓库认证：
   - 检查是否需要创建 imagePullSecret
   - 确认镜像仓库可访问性

### Pod 启动但状态异常

**检查 Pod 日志：**
```bash
kubectl --kubeconfig="C:\Users\li.longfeng\.kube\config" logs <pod-name> -n monitoring
```

**检查 Pod 事件：**
```bash
kubectl --kubeconfig="C:\Users\li.longfeng\.kube\config" describe pod <pod-name> -n monitoring
```

### Prometheus 无法采集数据

**检查 ServiceMonitor：**
```bash
kubectl --kubeconfig="C:\Users\li.longfeng\.kube\config" get servicemonitor -n monitoring
```

**检查 Prometheus targets：**
```bash
kubectl --kubeconfig="C:\Users\li.longfeng\.kube\config" port-forward -n monitoring svc/prometheus-k8s 9090:9090
# 访问 http://localhost:9090/targets
```

### 常见错误及解决方案

| 错误 | 原因 | 解决方案 |
|------|------|----------|
| `Too long: must have at most 262144 bytes` | CRD annotation 过大 | 使用 `--server-side=true` |
| `ImagePullBackOff` | 镜像拉取失败 | 检查镜像仓库配置和认证 |
| `CrashLoopBackOff` | 容器启动失败 | 查看容器日志排查配置错误 |
| `Pending` | 资源不足或 PVC 不可用 | 检查节点资源和存储配置 |

---

## 卸载说明

### 完全卸载

**步骤 1：删除所有监控组件**
```bash
kubectl --kubeconfig="C:\Users\li.longfeng\.kube\config" delete -f manifests/
```

**步骤 2：删除 CRD 和 Namespace**
```bash
kubectl --kubeconfig="C:\Users\li.longfeng\.kube\config" delete -f manifests/setup/
```

**验证卸载：**
```bash
# 检查 CRD
kubectl --kubeconfig="C:\Users\li.longfeng\.kube\config" get crd | grep monitoring.coreos.com

# 检查 namespace
kubectl --kubeconfig="C:\Users\li.longfeng\.kube\config" get ns monitoring
```

### 部分卸载

如只删除特定组件，可单独删除对应的 YAML 文件：
```bash
kubectl --kubeconfig="C:\Users\li.longfeng\.kube\config" delete -f manifests/<component>-<file>.yaml
```

---

## 附录

### 配置文件位置

- CRD 定义：`manifests/setup/`
- 监控组件：`manifests/`
- 主配置文件：`manifests/prometheus-prometheus.yaml`
- Grafana 配置：`manifests/grafana-*.yaml`
- 告警规则：`manifests/*-prometheusRule.yaml`

### 默认端口

| 服务 | 端口 | 用途 |
|------|------|------|
| Grafana | 3000 | Web UI |
| Prometheus | 9090 | Web UI |
| Alertmanager | 9093 | Web UI |
| Node Exporter | 9100 | 指标采集 |
| cAdvisor | 8080 | 指标采集 |
| Blackbox Exporter | 9115 | 指标采集 |
| Kube State Metrics | 8443/9443 | 指标采集 |

### 参考资源

- [Prometheus 官方文档](https://prometheus.io/docs/)
- [Prometheus Operator 文档](https://prometheus-operator.dev/)
- [Grafana 文档](https://grafana.com/docs/)
- [kube-prometheus 项目](https://github.com/prometheus-operator/kube-prometheus)

---

**文档版本：** v1.0
**更新日期：** 2026-03-02
**维护者：** Claude Code
