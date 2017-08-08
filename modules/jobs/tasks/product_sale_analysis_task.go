package tasks

import (
	"matrix/modules/erp/service/product_sale_analysis_service"
)

type ProductSaleAnalysisTask struct {
}

func (t ProductSaleAnalysisTask) Name() string {
	return "商品销售统计分析"
}

func (t ProductSaleAnalysisTask) Run() {
	product_sale_analysis_service.ProductSaleAnalyze()
}
