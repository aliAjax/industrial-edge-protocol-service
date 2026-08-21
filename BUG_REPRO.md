# Bug Reproduction

遥测重采样、聚合和过滤复用了输入切片底层数组，空边界路径还会产生越界 panic。运行 collection.json 中的 4 条 checks/s6 测试可复现；修复后全部通过。
