# Bug Reproduction

灰度配置在零 capabilities 或 typed-nil transport 场景写入未初始化 map，首次启动会 panic。运行 collection.json 中的 4 条 configrollout 测试可复现；修复后全部通过。
