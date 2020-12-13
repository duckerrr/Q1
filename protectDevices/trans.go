package protectDevices

type Trans struct {
	BaseDevice

	IfExternalPro bool			//是否配置了外部保护"busdifferential"
	IfInternalPro bool			//是否配置了内部保护"gas"
	Type_trans int				//双绕组（2）或三绕组（3）
	Bus []string				//对应的母线
	Gas_density float64			//瓦斯密度

}