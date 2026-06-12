package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"

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

	// 声明fanout类型的交换机（发布订阅的核心）
	exchangeName := "exchange_fanout"
	err = ch.ExchangeDeclare(
		exchangeName, // 交换机名称
		"fanout",     // 交换机类型：fanout（广播）
		true,         // 持久化
		false,        // 自动删除
		false,        // 内部
		false,        // 不等待
		nil,
	)
	if err != nil {
		log.Fatalf("声明交换机失败: %v", err)
	}

	// 模拟用户注册事件，发送广播消息
	event := "用户注册事件: 用户ID=1001,用户名=brianclark,注册时间=2026-06-12"
	err = ch.Publish(
		exchangeName, // 交换机名称
		"",           // fanout类型交换机无视路由键，留空即可
		false,        // 不强制
		false,        // 不立即
		amqp.Publishing{
			DeliveryMode: amqp.Persistent, // 消息持久化
			ContentType:  "text/plain",
			Body:         []byte(event),
		})
	if err != nil {
		log.Fatalf("发送消息失败: %v", err)
	}
	fmt.Println("✅ 广播事件发送成功")
}
