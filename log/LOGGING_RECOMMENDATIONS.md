# Kubernetes 日志现状与落地方案（`new-api/new-api-master-7d5bb66bdb-bwnmv`）

## 1. 结论摘要

当前并非“完全没有落盘”，而是：

1. **有写到容器标准输出日志（stdout/stderr）**，所以 `kubectl logs` 能看到日志。
2. **没有做可持续保留的持久化日志落盘方案**，Pod 重建或节点日志轮转后历史日志会丢。
3. `/app/logs` 虽有挂载，但挂载类型是 `emptyDir`，属于临时卷，不具备持久化能力。

---

## 2. 已核查证据（2026-03-03）

1. 上下文确认  
`kubectl config current-context` 返回：`kubernetes-admin@kubernetes`

2. Pod 状态正常  
`kubectl -n new-api get pod new-api-master-7d5bb66bdb-bwnmv -o wide` 显示 `Running`

3. Pod 描述中卷配置  
`kubectl -n new-api describe pod new-api-master-7d5bb66bdb-bwnmv` 显示：
- `/app/logs` 来自卷 `logs`
- `logs` 的类型为 `EmptyDir`

4. Deployment 配置再次确认  
`kubectl -n new-api get deploy new-api-master -o yaml` 显示：
- `volumeMounts`: `mountPath: /app/logs`
- `volumes`: `name: logs`, `emptyDir: {}`

5. 容器内目录实际情况  
`kubectl -n new-api exec ... -- ls -lah /app/logs` 目录为空

6. 标准输出日志持续存在  
`kubectl -n new-api logs ... --tail=120` 有持续业务日志与系统日志输出

---

## 3. 问题本质

1. **应用日志当前主通道是 stdout**，符合 K8s 常见实践。
2. **缺少集中日志系统**，导致日志可观测性依赖节点本地容器日志文件。
3. **`emptyDir` 不适合日志持久化**，Pod 删除/重建后数据即失。
4. 即使短期在节点有日志文件，也会受 kubelet/container runtime 日志轮转策略影响。

---

## 4. 方案建议（按推荐顺序）

## 方案 A（推荐）：stdout + 集中日志采集

适用：你希望日志可检索、可告警、可保留 30/90/180 天，且不受 Pod 生命周期影响。

1. 保持应用继续写 stdout/stderr（无需强制改代码写文件）。
2. 集群部署日志采集 DaemonSet：`fluent-bit` / `vector` / `promtail`。
3. 后端存储选型：`Loki`、`ELK`、`OpenSearch`、`ClickHouse`。
4. 配置按 namespace/pod/container 打标签。
5. 配置保留策略和告警规则。

优点：
1. 与 K8s 原生模式一致。
2. 不依赖单 Pod 或单节点。
3. 后续扩容多副本无需改日志写入逻辑。

风险与注意：
1. 需规划日志存储成本和保留周期。
2. 需做敏感字段脱敏（token、密钥、手机号等）。

---

## 方案 B：文件日志持久化（PVC）

适用：你有合规或审计要求，明确要求“日志必须写入文件并持久保存”。

1. 将 `/app/logs` 从 `emptyDir` 改为 `PersistentVolumeClaim`。
2. 创建专用 PVC（如 `new-api-logs`）。
3. 应用确保真正写文件到 `/app/logs`。
4. 增加日志轮转策略（应用内 rotate 或 sidecar/logrotate）。
5. 多副本场景考虑 RWX 存储或按 Pod 分卷。

优点：
1. 满足“文件落盘”要求。
2. Pod 重建不立即丢日志（取决于 PVC 策略）。

风险与注意：
1. 存储增长快，需要配额与清理策略。
2. 多副本共享写入可能带来竞争和性能问题。
3. 仅有 PVC 不等于“可检索”，通常仍建议配合集中采集。

---

## 方案 C：临时兜底（不推荐长期使用）

适用：短期先稳住，不立即改架构。

1. 保持 stdout 方案。
2. 调整节点日志轮转参数（增大单文件大小/保留数量）。
3. 增加定期导出归档任务。

局限：
1. 仍依赖节点生命周期。
2. 故障排查和跨时间检索能力弱。

---

## 5. 实施优先级建议

1. **第一阶段（本周）**：上线方案 A 的最小可用版本  
目标：先保证日志不因 Pod 重建丢失，并可按 Pod/时间检索。

2. **第二阶段（下周）**：完善保留/告警/脱敏  
目标：形成运维可用的日志治理基线。

3. **第三阶段（按合规需求）**：如必须文件审计，再叠加方案 B  
目标：满足“文件持久化”审计要求。

---

## 6. 验收标准（建议）

1. 重启 `new-api-master` Pod 后，历史日志仍可在日志平台查询。
2. 可按 `namespace=new-api`、`pod=new-api-master-*` 检索。
3. 保留策略生效（如 30 天）。
4. 关键错误日志可触发告警（例如 5xx 激增、panic、DB 连接失败）。
5. 日志中敏感字段已脱敏或不输出。

---

## 7. 你这个环境的直接判断

基于当前配置（`/app/logs -> emptyDir` + `kubectl logs` 有输出）：

1. **“应用没产生日志”这个判断不成立**。
2. **“日志没有持久化落盘/归档”这个判断成立**。
3. 推荐优先走 **方案 A（stdout + 集中采集）**，这是 Kubernetes 下更稳的长期方案。

---

## 8. 可选下一步

如需，我可以继续提供：

1. `new-api` 命名空间可用的 `fluent-bit + Loki` 最小化 YAML。
2. 把 `Deployment` 从 `emptyDir` 改 `PVC` 的 patch 示例（可直接 `kubectl apply`）。
