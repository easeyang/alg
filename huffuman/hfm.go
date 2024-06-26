package huffuman

import (
	"fmt"
	"sort"
)

type Node struct {
	Left   *Node // 左节点
	Right  *Node // 右节点
	Weight int   // 权重
	Value  byte  // 值
}

type HuffmanTree struct {
	Root  *Node // 根节点
	Bytes []byte
}

// 遍历哈夫曼树
func Traverse(node *Node) {
	if node == nil {
		return
	}
	fmt.Printf("Value: %c, Weight: %d\n", node.Value, node.Weight)
	Traverse(node.Left)
	Traverse(node.Right)
}

// 编码
func Encode(bytes []byte) []byte {
	nodes := make([]*Node, 0, len(bytes))
	// 扫码字节统计权重
	weights := make(map[byte]int)
	for _, b := range bytes {
		weights[b]++
	}
	// 构建哈夫曼树节点
	for b, w := range weights {
		nodes = append(nodes, &Node{Weight: w, Value: b})
	}
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Weight > nodes[j].Weight
	})
	// 构建哈夫曼树
	for len(nodes) > 1 {
		node := &Node{}
		node.Left = nodes[len(nodes)-1]
		node.Right = nodes[len(nodes)-2]
		node.Weight = node.Left.Weight + node.Right.Weight
		nodes = nodes[:len(nodes)-2]
		nodes = append(nodes, node)
		sort.Slice(nodes, func(i, j int) bool {
			return nodes[i].Weight > nodes[j].Weight
		})
	}
	rootNode := nodes[0]
	// 构建哈夫曼编码
	codes := make(map[byte]string)
	buildCode(codes, rootNode, "")
	fmt.Printf("codes: %v\n", codes)
	// 编码
	var result []byte
	for _, b := range bytes {
		result = append(result, []byte(codes[b])...)
	}
	return nil
}

func buildCode(codes map[byte]string, node *Node, code string) {
	if node == nil {
		return
	}
	if node.Left == nil && node.Right == nil {
		codes[node.Value] = code
	}
	buildCode(codes, node.Left, code+"0")
	buildCode(codes, node.Right, code+"1")
}

type HaffumanFile {
	version 
}
