package framework

import (
	"errors"
	"strings"
)

/**
	一旦引入了动态路由匹配的规则，之前使用的哈希规则就无法使用了。因为有通配符，在匹配 Request-URI 的时候，请求 URI 的某个字符或者某些字符是动态变化的，
无法使用 URI 做为 key 来匹配。那么，我们就需要其他的算法来支持路由匹配。
	如果你对算法比较熟悉，会联想到这个问题本质是一个字符串匹配，而字符串匹配，比较通用的高效方法就是字典树，也叫 trie 树。
*/

// Tree 字典树，也叫做 tris 树，也叫做前缀树
// trie 树不同于二叉树，它是多叉的树形结构，根节点一般是空字符串，而叶子节点保存的通常是字符串，一个节点的所有子孙节点都有相同的字符串前缀。
type Tree struct {
	root *node
}

func NewTree() *Tree {
	return &Tree{root: newNode()}
}

// FindHandler 匹配路由
func (tree *Tree) FindHandler(uri string) ControllerHandler {
	matchNode := tree.root.matchNode(uri)
	if matchNode == nil {
		return nil
	}
	return matchNode.handler
}

// AddRouter 添加路由，需要确认路由是否冲突。
// 我们先检查要增加的路由规则是否在树中已经有可以匹配的节点了。如果有的话，代表当前待增加的路由和已有路由存在冲突（ 我们用到了刚刚定义的 matchNode ）
// eg：/user/name ，/user/:id 这两个路由就是冲突的
// AddRouter 可以实现递归版本
func (tree *Tree) AddRouter(uri string, handler ControllerHandler) error {
	n := tree.root
	// 避免路由冲突
	if n.matchNode(uri) != nil {
		return errors.New("route exist：" + uri)
	}
	segments := strings.Split(uri, "/")
	for index, segment := range segments {
		// 如果不是通配符就大写
		if !isWildSegment(segment) {
			segment = strings.ToUpper(segment)
		}
		isLast := index == len(segments)-1
		var objNode *node
		children := n.filterChildNodes(segment)
		if len(children) > 0 {
			// 注意这里的逻辑：当chile.segment或者segment是通配符的时候，objNode不被赋值
			for _, child := range children {
				if child.segment == segment {
					objNode = child
					break
				}
			}
		}
		// 找不到节点，则新建
		if objNode == nil {
			objNode = newNode()
			objNode.segment = segment
			if isLast {
				objNode.isLast = isLast
				objNode.handler = handler
			}
			// 将新建的节点添加到当前的节点孩子中
			n.children = append(children, objNode)
		}
		// 将 n 重新设置
		n = objNode
	}
	return nil
}

// AddRouterRecurrence 递归版本
func (tree *Tree) AddRouterRecurrence(uri string, handler ControllerHandler) error {
	return nil
}

type node struct {
	// isLast 用于区别这个树中的节点是否有实际的路由含义
	isLast   bool              // 代表这个节点是否可以成为最终的路由规则。该节点是否能成为一个独立的uri, 是否自身就是一个终极节点
	segment  string            // uri中的字符串，代表这个节点表示的路由中某个段的字符串
	handler  ControllerHandler // 代表这个节点中包含的控制器，用于最终加载调用
	children []*node           // 代表这个节点下的子节点
}

func newNode() *node {
	return &node{
		isLast:   false,
		segment:  "",
		children: []*node{},
	}
}

// isWildSegment 判断一个segment是否是通用的segment，以:开头的
func isWildSegment(segment string) bool {
	return strings.HasPrefix(segment, ":")
}

// filterChildNodes 匹配符合和子节点
func (n *node) filterChildNodes(segment string) []*node {
	if len(n.children) <= 0 {
		return nil
	}
	// 如果 segment 是通配符，那么当前节点的所有子节点否满足要求
	// 理论上来讲如果是匹配路由的话segment是一定不可能是通配符的，为什么这里还要加上这个判断呢？
	// 因为 matchNode 方法在添加路由的时候也会使用到
	if isWildSegment(segment) {
		return n.children
	}
	nodes := make([]*node, 0, len(n.children))
	for _, child := range n.children {
		// 如果当前节点的segment是通配符，则符合要求
		if isWildSegment(child.segment) {
			nodes = append(nodes, child)
		} else if child.segment == segment { // 如果当前节点的segment和传递的segment等价，则符合要求
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// matchNode 判断路由是否已经在节点中，存在则返回节点
func (n *node) matchNode(uri string) *node {
	// 使用分隔符将uri切割成两部分
	segments := strings.SplitN(uri, "/", 2)

	// 先获取第一部分
	segment := segments[0]
	if !isWildSegment(segment) {
		segment = strings.ToUpper(segment)
	}
	// 匹配符合的节点
	children := n.filterChildNodes(segment)
	// 如果没有符合的子节点就说明匹配不到路由
	if len(children) <= 0 {
		return nil
	}
	// 如果只有一个segment，就不需要再往后走了
	if len(segments) == 1 {
		for _, child := range children {
			if !child.isLast {
				continue
			}
			return child
		}
		return nil
	}
	// 如果有2个segment，则递归从每一个子节点开始继续查找
	for _, child := range children {
		return child.matchNode(segments[1])
	}
	return nil
}
