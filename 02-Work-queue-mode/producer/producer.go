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
		ServerName:         "rabbitmq-server", // 和证书里的DNS一致
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

	// 声明持久化队列（和消费者保持一致）
	queueName := "xxx"
	q, err := ch.QueueDeclare(
		queueName, // 队列名
		true,      // 持久化：RabbitMQ重启后队列不消失
		false,     // 自动删除
		false,     // 排他
		false,     // 不等待
		nil,
	)
	if err != nil {
		log.Fatalf("声明队列失败: %v", err)
	}

	// 发送10条任务消息
	for i := range 10 {
		body := fmt.Sprintf("work queue mode %d", i+1)
		err = ch.Publish(
			"",     // 交换机（默认）
			q.Name, // 路由键=队列名
			false,  // 不强制
			false,  // 不立即
			amqp.Publishing{
				DeliveryMode: amqp.Persistent, // 消息持久化：RabbitMQ重启后消息不丢失
				ContentType:  "text/plain",
				Body:         []byte(body),
			})
		if err != nil {
			log.Fatalf("发送消息失败: %v", err)
		}
		fmt.Printf("✅ 发送任务: %s\n", body)
	}
}
