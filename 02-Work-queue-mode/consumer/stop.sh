#!/bin/bash

echo "===== 正在停止所有消费者进程 ====="

# 筛选出 go run consumer.go 相关进程并终止
pkill -f "go run consumer.go"

# 二次兜底：如果是编译后的二进制程序，也可以一起杀掉（备用）
pkill -f "./consumer"

echo "===== 所有消费者进程已停止 ====="

# 清理历史日志文件
rm -f consumer_*.log