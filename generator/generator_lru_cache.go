package generator

func (sd StructDescription) GenerateLRUCache(property Property, srcDir string) bool {

	pkList := sd.GetPKFieldList()
	if len(pkList) == 1 {
		return sd.generateLRUCache1PK(pkList[0], property, srcDir)
	} else if len(pkList) == 2 {
		return sd.generateLRUCache2PK(pkList[0], pkList[1], property, srcDir)
	}
	return true
}

func (sd StructDescription) getLRUCacheType(pkFiled FieldDescriptoin) (s string) {
	switch pkFiled.FieldGoType {
	case "int":
		return "LRUCacheInt"
	case "int32":
		return "LRUCacheInt32"
	case "int64":
		return "LRUCacheInt64"
	case "uint32":
		return "LRUCacheUint32"
	case "uint64":
		return "LRUCacheUint64"
	case "string":
		return "LRUCacheString"
	case "key.KeyUint64":
		return "ShardLRUCacheKeyUint64"
	// case "key.KeyInt":
	case "key.KeyInt32":
		return "LRUCacheKeyInt32"
	case "key.String":
		return "LRUCacheKeyString"
		// case "key.KeyString":

	}
	return
}
