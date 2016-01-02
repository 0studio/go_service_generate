package generator

func (sd StructDescription) GenerateMC(property Property, srcDir string) bool {

	pkList := sd.GetPKFieldList()
	if len(pkList) == 1 {
		return sd.generateMCPK1(pkList[0], property, srcDir)
	} else if len(pkList) == 2 {
		return sd.generateMCPK2(pkList[0], pkList[1], property, srcDir)
	}
	return true
}
