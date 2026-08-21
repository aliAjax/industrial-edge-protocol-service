# Bug Reproduction

BACnet 多层错误使用字符串格式重新包装，底层哨兵和类型信息不在错误链中。运行 collection.json 中的 4 条 checks/b3 测试可复现错误身份丢失；修复后全部通过。
