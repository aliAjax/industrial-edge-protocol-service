# Bug Reproduction

告警批处理的 worker 生命周期开关组合错误：worker 注册、放行和结果收集顺序不一致，可能提前返回、丢结果或等待不结束。运行 collection.json 中的 4 条 race 测试可复现。
