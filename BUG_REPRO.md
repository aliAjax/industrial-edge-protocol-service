# Bug Reproduction

回放失败路径未及时关闭 segment，且上传错误被 defer 覆盖，事务确认点在回滚后仍可能前移。运行 collection.json 中的 4 条 replay 相关测试可复现。
