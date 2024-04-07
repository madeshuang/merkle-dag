package merkledag

import (
	"encoding/json"
	"strings"
)

// Hash2File 从 Merkle DAG 存储中检索文件内容
func Hash2File(store KVStore, hash []byte, path string, hp HashPool) []byte {
	pathSegments := strings.Split(path, "/")
	if len(pathSegments) == 0 { // 如果路径为空
		return nil
	}

	objBytes, _ := store.Get(hash)

	var obj Object
	json.Unmarshal(objBytes, &obj)

	return recursiveSearch(store, obj, pathSegments, hp)
}

func recursiveSearch(store KVStore, obj Object, pathSegments []string, hp HashPool) []byte {
	if len(pathSegments) == 0 { // 如果路径为空
		return nil
	}

	for _, value := range obj.Links {
		switch pathSegments[0] {
		case "blob":
			blobValue, _ := store.Get(CalHash(value, hp))
			return store.Get(blobValue)
		case "link":
			// 将所有内容入栈
			return recursiveSearch(store, getObject(store, value.Hash), pathSegments, hp)
		case "tree":
			// 入栈 tree 后，递归处理剩余部分
			return recursiveSearch(store, getObject(store, value.Hash), pathSegments[1:], hp)
		}
	}

	return nil
}

//[辅助函数] 

func CalHash(data []byte, h hash.Hash) []byte {
	h.Reset()
	hash := h.Sum(data)
	h.Reset()
	return hash
}

func getObject(store KVStore, hash []byte) Object {
	objBytes, _ := store.Get(hash)

	var obj Object
	json.Unmarshal(objBytes, &obj)
	return obj
}
