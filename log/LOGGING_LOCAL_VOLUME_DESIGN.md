# new-api 本地卷日志持久化方案设计

## 1. 目标与边界

### 1.1 目标

1. 将 `new-api` 日志从临时卷 `emptyDir` 改为节点本地持久目录（hostPath）。
2. 日志可跨 Pod 重建保留（同一节点上）。
3. 可控日志保留周期、磁盘上限、清理策略，避免打满节点磁盘。

### 1.2 边界与事实

1. 当前 Pod：`new-api/new-api-master-7d5bb66bdb-bwnmv`。
2. 当前 `/app/logs` 挂载是 `emptyDir`，Pod 重建会丢。
3. 本方案是“本地卷落盘”，不等价于跨节点高可用归档。

---

## 2. 方案总览

采用三层控制：

1. **存储层**：`hostPath` 挂载到节点目录（例如 `/var/lib/new-api/logs/master`）。
2. **写入层**：应用写文件到 `/app/logs`（或由 stdout 侧车落地）。
3. **治理层**：日志轮转 + 保留天数 + 磁盘水位告警。

建议日志保留：

1. 默认保留 **14 天**。
2. 重要环境（生产）可保留 **30 天**。
3. 若节点盘较小，先用 **7 天**，再按实际增量调大。

---

## 3. 存储与目录规划

### 3.1 节点目录建议

统一目录：

```text
/var/lib/new-api/logs/
  master/
    app.log
    app.log.1.gz
    app.log.2.gz
```

### 3.2 挂载方式

在 Deployment 使用 `hostPath`：

- `type: DirectoryOrCreate`
- 容器内挂载点：`/app/logs`

示例（核心片段）：

```yaml
spec:
  template:
    spec:
      containers:
      - name: new-api
        volumeMounts:
        - name: logs
          mountPath: /app/logs
      volumes:
      - name: logs
        hostPath:
          path: /var/lib/new-api/logs/master
          type: DirectoryOrCreate
```

---

## 4. 日志保留与容量设计

### 4.1 容量测算公式

```text
所需容量(GB) = 日均日志量(GB) x 保留天数 x 安全系数(1.3~1.5)
```

示例：

1. 日均 0.8 GB，保留 14 天，系数 1.4
2. 容量约 `0.8 x 14 x 1.4 = 15.68 GB`
3. 建议至少预留 20 GB

### 4.2 建议保留策略

1. 生产默认：14 天。
2. 排障高峰期：临时升到 30 天。
3. 当节点可用磁盘 < 20% 时，自动缩短保留到 7 天（应急策略）。

### 4.3 轮转参数建议

以单文件轮转为例：

1. `maxsize`: 200MB（达到即切割）。
2. `daily`: 每天至少轮转一次。
3. `rotate`: 70（约可覆盖 14 天，按频率可微调）。
4. `compress`: 开启 gzip 压缩。
5. `copytruncate`：应用不支持 reopen 时启用。

---

## 5. 实施方式（推荐）

推荐在 Pod 内增加一个 `logrotate` sidecar，专门管理 `/app/logs/*.log`。

优点：

1. 与应用解耦，不改应用代码。
2. 策略清晰，可版本化管理。
3. 可随 Deployment 一起发布回滚。

### 5.1 sidecar 设计

1. 主容器写 `/app/logs/app.log`。
2. sidecar 周期执行 `logrotate`（例如每 5 分钟）。
3. 同挂载 `hostPath`，对同一目录生效。

### 5.2 logrotate 示例配置

```conf
/app/logs/*.log {
  daily
  rotate 14
  maxsize 200M
  missingok
  notifempty
  compress
  delaycompress
  copytruncate
  dateext
  dateformat -%Y%m%d-%s
}
```

> 若要保留 30 天，将 `rotate 14` 改为 `rotate 30`，并同步评估磁盘容量。

---

## 6. 可直接参考的 Deployment Patch

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: new-api-master
  namespace: new-api
spec:
  template:
    spec:
      containers:
      - name: new-api
        volumeMounts:
        - name: logs
          mountPath: /app/logs
      - name: logrotate
        image: blacklabelops/logrotate:latest
        imagePullPolicy: IfNotPresent
        env:
        - name: LOGS_DIRECTORIES
          value: /app/logs
        - name: LOGROTATE_INTERVAL
          value: hourly
        - name: LOGROTATE_COPIES
          value: "14"
        - name: LOGROTATE_SIZE
          value: 200M
        - name: LOGROTATE_COMPRESSION
          value: compress
        volumeMounts:
        - name: logs
          mountPath: /app/logs
      volumes:
      - name: logs
        hostPath:
          path: /var/lib/new-api/logs/master
          type: DirectoryOrCreate
```

说明：

1. 该示例用于表达思路，具体 sidecar 镜像参数请与你们镜像文档对齐。
2. 如果你不希望引入 sidecar，可在主机上用 `logrotate.d` 管理同目录。

---

## 7. 风险与约束

1. **节点绑定风险**：Pod 漂移到其他节点后，新节点日志目录是另一份。
2. **主机故障风险**：节点坏盘会导致该节点日志丢失。
3. **扩容风险**：多副本时每个节点各自产生日志，检索分散。

规避建议：

1. 关键日志仍建议异步采集到集中平台（Loki/ELK）做长期归档。
2. 至少对错误日志做异地备份（每天归档）。

---

## 8. 监控与告警建议

至少加 3 类监控：

1. 目录容量：`/var/lib/new-api/logs` 使用率 > 70% 告警，> 85% 严重告警。
2. 轮转健康：24 小时内无新 `.gz` 文件告警（可能轮转失效）。
3. 写入健康：`app.log` 5 分钟无增长告警（排除业务低峰时段）。

---

## 9. 验收清单

1. 重启 Pod 后，旧日志文件仍在同节点目录。
2. 轮转后出现压缩包，且最新日志持续写入。
3. 节点磁盘压力下不会无限增长（有删旧策略）。
4. 按保留策略（14/30 天）抽查历史日志可读。

---

## 10. 推荐落地参数（初始值）

1. 保留天数：14 天。
2. 单文件上限：200MB。
3. 总容量预算：20GB（按日均 0.8GB 估算）。
4. 轮转频率：daily + size 双触发。
5. 告警阈值：70%/85%。

---

## 11. 后续演进（建议）

1. 第一步：先完成 hostPath + 轮转，解决“临时卷丢日志”。
2. 第二步：接入集中日志系统，解决“跨节点检索与长期归档”。
3. 第三步：按合规要求增加归档存储（对象存储/冷存储）。
