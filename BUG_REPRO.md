# Bug Reproduction

MQTT 解码、鉴权、QoS 和路由层把错误链格式化成普通字符串，调用方无法区分坏包、权限拒绝和临时网络错误。运行 collection.json 中的 4 条 checks/m9 测试可复现。
