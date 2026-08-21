# Bug Reproduction

Modbus 帧、寄存器解析和传输分片直接暴露或修改收包缓冲，下一帧会改写上一帧数据，短帧还可能越界。运行 collection.json 中的 4 条 checks/m10 测试可复现；修复后全部通过。
