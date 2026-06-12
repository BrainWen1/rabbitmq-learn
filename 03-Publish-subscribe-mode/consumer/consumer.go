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

	// 声明fanout交换机（和生产者保持一致）
	exchangeName := "exchange_fanout"
	err = ch.ExchangeDeclare(
		exchangeName,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("声明交换机失败: %v", err)
	}

	// 声明临时队列（发布订阅的核心）
	// - 不指定队列名，RabbitMQ会自动生成随机名称
	// - 排他、自动删除：消费者退出后队列自动消失
	q, err := ch.QueueDeclare(
		"",    // 队列名（空字符串=自动生成）
		false, // 不持久化（临时队列）
		true,  // 自动删除
		true,  // 排他：只有当前连接能访问这个队列
		false, // 不等待
		nil,
	)
	if err != nil {
		log.Fatalf("声明队列失败: %v", err)
	}

	// 将临时队列绑定到fanout交换机
	err = ch.QueueBind(
		q.Name,       // 队列名
		"",           // fanout交换机无视路由键，留空即可
		exchangeName, // 交换机名称
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("绑定队列到交换机失败: %v", err)
	}

	// 创建消费者，订阅广播消息
	msgs, err := ch.Consume(
		q.Name,
		"",    // 消费者标签
		true,  // 自动确认（广播场景常用，简化逻辑）
		true,  // 排他
		false, // 不等待
		false, // 无本地
		nil,
	)
	if err != nil {
		log.Fatalf("注册消费者失败: %v", err)
	}

	// 处理广播消息
	fmt.Println("🟢 事件订阅者启动成功，等待广播消息...")
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			fmt.Printf("📥 收到广播事件: %s\n", d.Body)
			// 这里可以写不同的业务逻辑，比如：
			// - 发送欢迎邮件
			// - 记录注册日志
			// - 同步用户数据到统计模块
		}
	}()

	<-forever // 阻塞，保持消费者运行
}
