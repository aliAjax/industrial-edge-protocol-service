# Operations Runbook

## SLO

API 可用性目标 99.9%，采集成功率目标 99.5%，规则评估 p99 小于 50ms。WAL 达到 80% 容量时触发告警，达到 95% 时拒绝新写入并保留现有 segment。

## 故障处理

网关断联时保留有界 WAL；恢复后按 sequence 顺序上传，服务端按 digest 和 sequence 去重。配置灰度失败时将 rollout 标记 aborted，旧版本继续运行。危险命令默认拒绝，审批和设备回执均写入审计链。

## 灾备与演练

定期复制 telemetry segment 和迁移版本，恢复后先运行 WAL replay 与 digest 校验，再开放采集。演练包括磁盘耗尽、时钟漂移、半包 Modbus、MQTT 重连、网关重启和命令超时。
