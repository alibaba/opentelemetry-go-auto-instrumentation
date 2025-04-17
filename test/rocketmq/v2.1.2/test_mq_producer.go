// Copyright (c) 2024 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2"
	"log"
	"sync"
	"time"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/test/verifier"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func main() {
	// 初始化生产者
	initTopic()
	p := initRocketMQ()
	defer p.Shutdown()
	err := p.Start()
	if err != nil {
		panic(err)
	}

	log.Println("===== 测试同步发送成功 =====")
	testSyncSendSuccess(p)

	log.Println("\n===== 测试同步发送失败 =====")
	testSyncSendFailure(p)

	log.Println("\n===== 测试异步发送成功 =====")
	testAsyncSendSuccess(p)

	log.Println("\n===== 测试单向发送成功 =====")
	testOneWaySendSuccess(p)

}

// 测试同步发送成功场景
func testSyncSendSuccess(p rocketmq.Producer) {
	// 创建并发送消息
	msg := primitive.NewMessage(topicName, []byte("测试同步发送"))
	msg.WithTag("test_sync")

	result, err := p.SendSync(context.Background(), msg)
	if err != nil {
		log.Printf("同步发送失败: %v\n", err)
		return
	}

	log.Printf("同步发送成功: %s\n", result.String())

	// 验证Span
	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		span := stubs[0][0]
		VerifyRocketMQProduceAttributes(span, topicName, "test_sync", "", "publish", false)
	}, 1)
}

// 测试同步发送失败场景
func testSyncSendFailure(p rocketmq.Producer) {
	// 使用超短超时触发错误
	msg := primitive.NewMessage(topicName, []byte("测试同步发送失败"))
	msg.WithTag("test_sync_failure")

	// 创建1微秒超时的上下文，确保超时发生
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Microsecond)
	defer cancel()

	// 等待一小段时间确保超时已触发
	time.Sleep(10 * time.Millisecond)

	_, err := p.SendSync(ctx, msg)
	if err == nil {
		log.Println("错误: 期望发送失败但实际成功")
		return
	}

	log.Printf("预期的同步发送失败: %v\n", err)

	// 验证错误Span
	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		span := stubs[0][0]
		VerifyRocketMQProduceAttributes(span, topicName, "test_sync_failure", "", "publish", true)
	}, 1)
}

// 测试异步发送成功场景
func testAsyncSendSuccess(p rocketmq.Producer) {
	// 创建并异步发送消息
	msg := primitive.NewMessage(topicName, []byte("测试异步发送"))
	msg.WithTag("test_async")

	var wg sync.WaitGroup
	wg.Add(1)

	err := p.SendAsync(context.Background(), func(ctx context.Context, result *primitive.SendResult, err error) {
		defer wg.Done()
		if err != nil {
			log.Printf("异步发送回调返回错误: %v\n", err)
		} else {
			log.Printf("异步发送成功: %s\n", result.String())
		}
	}, msg)

	if err != nil {
		log.Printf("异步发送请求失败: %v\n", err)
		return
	}

	wg.Wait() // 等待回调完成

	// 验证Span
	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		span := stubs[0][0]
		VerifyRocketMQProduceAttributes(span, topicName, "test_async", "", "publish", false)
	}, 1)
}

// 测试异步发送失败场景
/*func testAsyncSendFailure(p rocketmq.Producer) {

	msg := primitive.NewMessage("test_async_failure", []byte("测试异步发送失败"))
	msg.WithTag("test_async_failure")

	var wg sync.WaitGroup
	wg.Add(1)

	err := p.SendAsync(context.Background(), func(ctx context.Context, result *primitive.SendResult, err error) {
		defer wg.Done()
		result = nil
		err = errors.New("模拟异步发送失败")
		if err == nil {
			log.Println("错误: 期望异步回调返回错误但实际成功")
		} else {
			log.Printf("预期的异步回调错误: %v\n", err)
		}

	}, msg)

	if err != nil {
		// 注意：有些实现可能在发起请求时就返回错误，这也是可接受的
		log.Printf("异步发送请求失败: %v\n", err)
	}

	wg.Wait() // 等待回调完成

	// 验证Span
	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		span := stubs[0][0]
		VerifyRocketMQAttributes(span, "test_async_failure", "test_async_failure", "", "publish", true)
	}, 1)
}*/

// 测试单向发送成功场景
func testOneWaySendSuccess(p rocketmq.Producer) {
	// 创建并单向发送消息
	msg := primitive.NewMessage(topicName, []byte("测试单向发送"))
	msg.WithTag("test_oneway")

	err := p.SendOneWay(context.Background(), msg)
	if err != nil {
		log.Printf("单向发送失败: %v\n", err)
		return
	}

	log.Println("单向发送完成")

	// 验证Span
	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		span := stubs[0][0]
		VerifyRocketMQProduceAttributes(span, topicName, "test_oneway", "", "publish", false)
	}, 3)
}

// 测试单向发送失败场景
/*func testOneWaySendFailure(p rocketmq.Producer) {

	msg := primitive.NewMessage(topicName, []byte("测试单向发送失败"))
	msg.WithTag("test_oneway_failure")
	// 创建1微秒超时的上下文，确保超时发生
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Microsecond)
	defer cancel()

	// 等待一小段时间确保超时已触发
	time.Sleep(10 * time.Millisecond)
	err := p.SendOneWay(ctx, msg)
	if err == nil {
		// 注意：由于单向发送不关心结果，某些实现可能不会立即返回错误
		log.Println("单向发送请求成功，但实际结果未知")
	} else {
		log.Printf("单向发送失败: %v\n", err)
	}

	// 验证Span
	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		span := stubs[0][0]
		VerifyRocketMQAttributes(span, topicName, "test_oneway_failure", "", "publish", true) // 单向发送可能不会在span中标记错误
	}, 3)
}*/
