package huffuman

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strings"
	"time"
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
func Encode(bytes []byte) ([]byte, error) {
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
	// 构建哈夫曼编码表
	codes := make(map[byte]string)
	buildCode(codes, rootNode, "")
	// fmt.Printf("codes: %v\n", codes)
	// 编码
	encodeStrB := strings.Builder{}
	for _, b := range bytes {
		encodeStrB.WriteString(codes[b])
	}
	encodeStr := encodeStrB.String()
	headerBytes := NewHaffumanFileHeder(codes, int64(len(encodeStr))).Bytes()
	// fmt.Printf("headerBytes: %v, len=%v \n", headerBytes, len(headerBytes))
	contentBytes := make([]byte, int64(math.Ceil(float64(len(encodeStr))/8.0)))
	for i := 0; i < len(encodeStr); i++ {
		if encodeStr[i] == '1' {
			contentBytes[i/8] |= 1 << (7 - i%8)
		}
	}
	// fmt.Printf("contentBytes: %08b, len=%v \n", contentBytes, len(contentBytes))

	// 合并两个字节数组并且返回
	return append(headerBytes, contentBytes...), nil
}

// 解码
func Decode(bts []byte) ([]byte, error) {
	// 解析文件头
	h := &HaffumanFileHeder{}
	h.Version = bts[0]
	h.CreateTime = int64(binary.LittleEndian.Uint64(bts[1:9]))
	h.ContentBitLen = int64(binary.LittleEndian.Uint64(bts[9:17]))
	h.MapLen = int16(binary.LittleEndian.Uint16(bts[17:19]))
	codesBytes := bts[19 : 19+h.MapLen]
	err := json.Unmarshal(codesBytes, &h.Codes)
	if err != nil {
		return nil, err
	}
	// fmt.Printf("#header#\n %v\n", h)
	// 解码
	encodeBytes := bts[19+h.MapLen:]
	encodeStrB := strings.Builder{}
	for _, b := range encodeBytes {
		encodeStrB.WriteString(fmt.Sprintf("%08b", b))
	}
	encodeStr := encodeStrB.String()
	// fmt.Printf("encodeStr: %s, len=%v \n", encodeStr, len(encodeStr))
	// 解码
	buf := &bytes.Buffer{}
	prefix := ""
	for i := 0; i < len(encodeStr); i++ {
		prefix += string(encodeStr[i])
		if v, ok := h.Codes[prefix]; ok {
			buf.WriteByte(v)
			prefix = ""
		}
	}

	return buf.Bytes(), nil
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

type HaffumanFileHeder struct {
	CreateTime    int64           // 创建时间
	ContentBitLen int64           // 内容位数(bit)长度
	MapLen        int16           // 哈夫曼编码表长度(保存map转成json后的字符串字节长度)
	Codes         map[string]byte // 哈夫曼编码表
	Version       byte            // 版本
}

func (h *HaffumanFileHeder) String() string {
	return fmt.Sprintf("Version: %d, CreateTime: %d, ContentBitLen: %d, MapLen: %d, Codes: %v",
		h.Version, h.CreateTime, h.ContentBitLen, h.MapLen, h.Codes)
}

func (h *HaffumanFileHeder) Bytes() []byte {
	codesBytes, err := json.Marshal(h.Codes)
	if err != nil {
		panic(err)
	}

	// 文件头字节数量
	// 1(版本) + 8(创建时间) + 8(内容位数长度) + 2(编码表长度) + 编码表内容长度

	bf := &bytes.Buffer{}
	bf.WriteByte(h.Version)
	binary.Write(bf, binary.LittleEndian, h.CreateTime)
	binary.Write(bf, binary.LittleEndian, h.ContentBitLen)
	binary.Write(bf, binary.LittleEndian, int16(len(codesBytes)))
	bf.Write(codesBytes)
	return bf.Bytes()
}

func NewHaffumanFileHeder(codeMap map[byte]string, bitLen int64) *HaffumanFileHeder {
	h := &HaffumanFileHeder{}
	h.Version = 1
	h.CreateTime = time.Now().Unix()
	h.Codes = swapMapKeyValue(codeMap)
	h.ContentBitLen = bitLen
	return h
}

func swapMapKeyValue(m map[byte]string) map[string]byte {
	n := make(map[string]byte)
	for k, v := range m {
		n[v] = k
	}
	return n
}

// 压缩
func Zip(srcName, disName string) {
	// 打开srcFileName文件
	fSrc, err := os.Open(srcName)
	if err != nil {
		fmt.Println("打开文件失败", err)
		return
	}
	defer fSrc.Close()
	content, err := io.ReadAll(fSrc)
	if err != nil {
		fmt.Println("读取文件内容失败", err)
		return
	}
	contentEncoded, err := Encode(content)
	if err != nil {
		fmt.Println("压缩文件失败", err)
		return
	}
	// 创建disFileName文件
	fDis, err := os.Create(disName)
	if err != nil {
		fmt.Println("创建文件失败", err)
		return
	}
	defer fDis.Close()
	_, err = fDis.Write(contentEncoded)
	if err != nil {
		fmt.Println("写入文件失败", err)
		return
	}
}

// 解压
func Unzip(srcName, disName string) {
	// 打开srcFileName文件
	fSrc, err := os.Open(srcName)
	if err != nil {
		fmt.Println("打开文件失败", err)
		return
	}
	defer fSrc.Close()
	content, err := io.ReadAll(fSrc)
	if err != nil {
		fmt.Println("读取文件内容失败", err)
		return
	}
	contentDecoded, err := Decode(content)
	if err != nil {
		fmt.Println("解压文件失败", err)
		return
	}
	// 创建disFileName文件
	fDis, err := os.Create(disName)
	if err != nil {
		fmt.Println("创建文件失败", err)
		return
	}
	defer fDis.Close()
	_, err = fDis.Write(contentDecoded)
	if err != nil {
		fmt.Println("写入文件失败", err)
		return
	}
}
