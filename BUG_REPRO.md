# Bug Reproduction

命令拒绝后确认会把状态重新放回 inflight，重试状态机和终态判断不一致，队列长期显示执行中。运行 collection.json 中的 4 条 checks/q8 测试可复现；修复后全部通过。
