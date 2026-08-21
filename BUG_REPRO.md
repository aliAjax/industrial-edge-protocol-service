# Bug Reproduction

并发查询遥测时，返回切片与存储或 HTTP scratch 缓冲共享底层数组，写入和读取会产生数据竞争并污染快照。运行 collection.json 中的 4 条 race 测试可在埋错基线复现；修复后全部通过。
