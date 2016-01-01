package generator

func (sd StructDescription) GenerateProxy(property Property, srcDir string) bool {

	pkList := sd.GetPKFieldList()
	if len(pkList) == 1 {
		return sd.generateProxyPK1(pkList[0], property, srcDir)
	} else if len(pkList) == 2 {
		// return sd.generateProxyPK2(pkList[0], pkList[1], property, srcDir)
	}
	return true
}
