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
	// 加载客户端证书和密钥（双向认证时需要）
	cert, err := tls.LoadX509KeyPair("/home/brianclark/learn/rabbitMq/safe/ssl/client_certificate.pem", "/home/brianclark/learn/rabbitMq/safe/ssl/client_key.pem")
	if err != nil {
		log.Fatalf("加载客户端证书失败: %v", err)
	}

	// 加载CA证书（验证服务器证书）
	caCert, err := os.ReadFile("/home/brianclark/learn/rabbitMq/safe/ssl/ca_certificate.pem")
	if err != nil {
		log.Fatalf("读取CA证书失败: %v", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// 配置TLS
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert}, // 客户端证书（双向认证时需要）
		RootCAs:            caCertPool,              // 信任的CA
		InsecureSkipVerify: false,                   // 必须验证服务器证书
	}

	// 连接到 RabbitMQ 服务器
	conn, err := amqp.DialTLS("amqps://admin:123456@localhost:5671/", tlsConfig)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// 创建通道
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	queueName := "xxx"

	// 声明队列
	q, err := ch.QueueDeclare(
		queueName, // 队列名称
		false,     // 持久性
		false,     // 自动删除
		false,     // 排他性
		false,     // 非阻塞
		nil,       // 其他参数
	)
	if err != nil {
		panic(err)
	}

	// 发送消息
	body := fmt.Sprintf("Hello RabbitMQ! %d", 222)
	err = ch.Publish(
		"",     // 交换机
		q.Name, // 路由键
		false,  // 强制性
		false,  // 立即发送
		amqp.Publishing{ // 消息内容
			Body: []byte(body),
		})
	if err != nil {
		log.Fatalf("发送失败: %v", err)
		return
	}
	fmt.Printf("成功发送消息: %s\n", body)
}
