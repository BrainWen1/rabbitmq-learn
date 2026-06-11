package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/streadway/amqp"
)

func main() {
	// 加载客户端证书和密钥
	cert, err := tls.LoadX509KeyPair("/home/brianclark/learn/rabbitMq/safe/ssl/client_certificate.pem", "/home/brianclark/learn/rabbitMq/safe/ssl/client_key.pem")
	if err != nil {
		log.Fatalf("加载客户端证书失败: %v", err)
	}

	// 加载CA证书
	caCert, err := os.ReadFile("/home/brianclark/learn/rabbitMq/safe/ssl/ca_certificate.pem")
	if err != nil {
		log.Fatalf("加载CA证书失败: %v", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// TLS配置
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
		ServerName:         "rabbitmq-server",
		InsecureSkipVerify: false,
	}

	// 连接加密RabbitMQ
	conn, err := amqp.DialTLS("amqps://admin:123456@localhost:5671/", tlsConfig)
	if err != nil {
		log.Fatalf("连接RabbitMQ失败: %v", err)
	}
	defer conn.Close()

	// 创建通道
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("创建通道失败: %v", err)
	}
	defer ch.Close()

	// 声明持久化队列（和生产者必须完全一致）
	queueName := "xxx"
	q, err := ch.QueueDeclare(
		queueName,
		true,  // 持久化（和生产者一致）
		false, // 自动删除
		false, // 排他
		false, // 不等待
		nil,
	)
	if err != nil {
		log.Fatalf("声明队列失败: %v", err)
	}

	// 关键配置：开启公平分发（限制每个消费者同时处理1条消息）
	err = ch.Qos(
		1,     // 预取数：每个消费者最多同时处理1条消息
		0,     // 预取大小（字节）
		false, // 全局设置
	)
	if err != nil {
		log.Fatalf("设置QoS失败: %v", err)
	}

	// 消费消息（关闭自动确认，手动ack）
	msgs, err := ch.Consume(
		q.Name,
		"",    // 消费者标签
		false, // 关闭自动确认（autoAck=false，必须手动ack）
		false, // 排他
		false, // 不等待
		false, // 无本地
		nil,
	)
	if err != nil {
		log.Fatalf("注册消费者失败: %v", err)
	}

	// 处理消息
	fmt.Println("🟢 消费者启动成功，等待任务...")
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			fmt.Printf("📥 收到任务: %s\n", d.Body)

			// 模拟耗时任务（比如导入单词）
			time.Sleep(2 * time.Second)

			// 手动确认消息（必须！处理完再ack，避免消息丢失）
			err := d.Ack(false)
			if err != nil {
				log.Printf("❌ 消息确认失败: %v", err)
			}
			fmt.Printf("✅ 任务完成: %s\n", d.Body)
		}
	}()

	<-forever // 阻塞，保持消费者运行
}
