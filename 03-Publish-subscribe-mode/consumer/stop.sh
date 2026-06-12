#!/bin/bash

echo "===== 停止所有发布订阅消费者 ====="

# 批量终止 consumer.go 对应的进程
pkill -f "go run consumer.go"

# 删除日志文件
rm -f subscriber_*.log

echo "===== 所有订阅者已停止，临时队列将自动删除 ====="