package logistics_fee_type

const (
	FEE_TYPE1 = 1
	FEE_TYPE2 = 2
)

func Get_fee_type_list() []LogisticsFeeTypeInfo {
	var fee_type_list = make([]LogisticsFeeTypeInfo, 0)

	var fee_type1 LogisticsFeeTypeInfo
	fee_type1.Id = FEE_TYPE1
	fee_type1.Name = "费率类型1"
	fee_type_list = append(fee_type_list, fee_type1)

	var fee_type2 LogisticsFeeTypeInfo
	fee_type2.Id = FEE_TYPE2
	fee_type2.Name = "费率类型2"
	fee_type_list = append(fee_type_list, fee_type2)

	return fee_type_list
}
