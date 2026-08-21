# Bug Reproduction

采集任务把已过期的 context 传给后续读取和重试，取消/超时信号污染下一批。运行 collection.json 中的 4 条 acquisition 测试可观察该生命周期错误。
