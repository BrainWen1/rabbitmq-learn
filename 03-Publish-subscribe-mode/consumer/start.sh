#! /bin/bash

# 定义启动的订阅者数量
SUBSCRIBER_NUM=5

echo "===== 启动 ${SUBSCRIBER_NUM} 个发布订阅模式消费者 ====="

# 循环启动订阅者，后台运行，日志独立输出
for ((i=1; i<=${SUBSCRIBER_NUM}; i++))
do
    echo "启动第 $i 个订阅者，日志输出至 subscriber_${i}.log"
    nohup go run consumer.go > subscriber_${i}.log 2>&1 &
    sleep 0.2
done

echo "===== 全部 ${SUBSCRIBER_NUM} 个订阅者启动完成 ====="
echo "提示：所有订阅者都会收到完整的广播消息"