#!/bin/bash

# 定义启动的消费者数量
CONSUMER_NUM=10

echo "===== 开始启动 ${CONSUMER_NUM} 个 RabbitMQ 消费者 ====="

# 循环启动消费者
for ((i=1; i<=${CONSUMER_NUM}; i++))
do
    echo "启动第 $i 个消费者，日志输出至 consumer_${i}.log"
    # 后台运行 go 程序，标准输出/错误重定向到独立日志文件
    nohup go run consumer.go > consumer_${i}.log 2>&1 &
    # 短暂延时，避免瞬间并发抢占资源
    sleep 0.2
done

echo "===== 全部 ${CONSUMER_NUM} 个消费者启动完成 ====="
echo "可查看对应 log 文件观察消费情况"