package verifier

import (
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"log"
	"sort"
	"time"
)

type node struct {
	root       bool
	childNodes []*node
	span       tracetest.SpanStub
}

func WaitAndAssertTraces(traceVerifiers ...func([]tracetest.SpanStubs)) {
	numTraces := len(traceVerifiers)
	traces := waitForTraces(numTraces)
	for _, v := range traceVerifiers {
		v(traces)
	}
}

func waitForTraces(numberOfTraces int) []tracetest.SpanStubs {
	// 最多等20s
	finish := false
	var traces []tracetest.SpanStubs
	var i int
	for !finish {
		select {
		case <-time.After(20 * time.Second):
			log.Printf("Timeout waiting for traces!")
			finish = true
		default:
			traces = groupAndSortTrace()
			if len(traces) >= numberOfTraces {
				finish = true
			}
			i++
		}
		if i == 10 {
			break
		}
	}
	return traces
}

func groupAndSortTrace() []tracetest.SpanStubs {
	spans := GetTestSpans()
	defer ResetTestSpans()
	traceMap := make(map[string][]tracetest.SpanStub)
	for _, span := range *spans {
		if span.SpanContext.HasTraceID() && span.SpanContext.TraceID().IsValid() {
			traceId := span.SpanContext.TraceID().String()
			spans, ok := traceMap[traceId]
			if !ok {
				spans = make([]tracetest.SpanStub, 0)
			}
			spans = append(spans, span)
			traceMap[traceId] = spans
		}
	}
	return sortTrace(traceMap)
}

func sortTrace(traceMap map[string][]tracetest.SpanStub) []tracetest.SpanStubs {
	traces := make([][]tracetest.SpanStub, 0)
	for _, trace := range traceMap {
		traces = append(traces, trace)
	}
	// 按开始时间从小到大排
	sort.Slice(traces, func(i, j int) bool {
		return traces[i][0].StartTime.UnixNano() < traces[j][0].StartTime.UnixNano()
	})
	for i, _ := range traces {
		traces[i] = sortSingleTrace(traces[i])
	}
	stubs := make([]tracetest.SpanStubs, 0)
	for i, _ := range traces {
		stubs = append(stubs, traces[i])
	}
	return stubs
}

func sortSingleTrace(stubs []tracetest.SpanStub) []tracetest.SpanStub {
	// 同一条trace的按span的父子关系排
	lookup := make(map[string]*node)
	for _, stub := range stubs {
		lookup[stub.SpanContext.SpanID().String()] = &node{
			root:       true,
			childNodes: make([]*node, 0),
			span:       stub,
		}
	}
	for _, stub := range stubs {
		n, ok := lookup[stub.SpanContext.SpanID().String()]
		if !ok {
			panic("no span id in stub " + stub.Name)
		}
		// 发现了父节点，就添加到父节点的子节点列表里面去
		if n.span.Parent.SpanID().IsValid() {
			parentSpanId := n.span.Parent.SpanID().String()
			parentNode, ok := lookup[parentSpanId]
			if ok {
				parentNode.childNodes = append(parentNode.childNodes, n)
				n.root = false
			}
		}
	}
	// 寻找根节点
	rootNodes := make([]*node, 0)
	for _, stub := range stubs {
		n, ok := lookup[stub.SpanContext.SpanID().String()]
		if !ok {
			panic("no span id in stub " + stub.Name)
		}
		sort.Slice(n.childNodes, func(i, j int) bool {
			return n.childNodes[i].span.StartTime.Unix() < n.childNodes[j].span.StartTime.Unix()
		})
		if n.root {
			rootNodes = append(rootNodes, n)
		}
	}
	sort.Slice(rootNodes, func(i, j int) bool {
		return rootNodes[i].span.StartTime.Unix() < rootNodes[j].span.StartTime.Unix()
	})
	// 层序遍历，获取排序后的span
	t := make([]tracetest.SpanStub, 0)
	for _, rootNode := range rootNodes {
		traversePreOrder(rootNode, &t)
	}
	return t
}

func traversePreOrder(n *node, acc *[]tracetest.SpanStub) {
	*acc = append(*acc, n.span)
	for _, child := range n.childNodes {
		traversePreOrder(child, acc)
	}
}
