package xsbb_models

// 这个 SaleAmountReportInfo 只用来接收查询结果，不创建数据库表
type SaleAmountReportProductInfo struct {
	Id                   int64   `xorm:"'id'" json:"id"`
	ProductId            int64   `xorm:"'product_id'" json:"product_id"`
	ProductSku           string  `xorm:"'product_sku'" json:"product_sku"`
	ProductName          string  `xorm:"'product_name'" json:"product_name"`
	ProductWeight        float64 `xorm:"'product_weight'" json:"product_weight"`
	CostPrice            float64 `xorm:"'cost_price'" json:"cost_price"`
	SaleUserId           int64   `xorm:"'sale_user_id'" json:"sale_user_id"`
	SaleUserName         string  `xorm:"'sale_user_name'" json:"sale_user_name"`
	ServiceUserId        int64   `xorm:"'sale_user_id'" json:"service_user_id"`
	ServiceUserName      string  `xorm:"'service_user_name'" json:"service_user_name"`
	DevUserId            int64   `xorm:"'dev_user_id'" json:"dev_user_id"`
	DevUserName          string  `xorm:"'dev_user_name'" json:"dev_user_name"`
	PaidAmount           float64 `xorm:"'paid_amount'" json:"paid_amount"`
	ExpressFee           float64 `xorm:"'express_fee'" json:"express_fee"`
	Income               float64 `xorm:"'income'" json:"income"`
	CostAmount           float64 `xorm:"'cost_amount'" json:"cost_amount"`
	EstimateLogisticsFee float64 `xorm:"'estimate_logistics_fee'" json:"estimate_logistics_fee"`
	ActualLogisticsFee   float64 `xorm:"'actual_logistics_fee'" json:"actual_logistics_fee"`
	JyfyAmount           float64 `xorm:"'jyfy_amount'" json:"jyfy_amount"`
	PtfyAmount           float64 `xorm:"'ptfy_amount'" json:"ptfy_amount"`
	PickFee              float64 `xorm:"'pick_fee'" json:"pick_fee"`
	GstFee               float64 `xorm:"'gst_fee'" json:"gst_fee"`
	Outcome              float64 `xorm:"'outcome'" json:"outcome"`
	Profit               float64 `xorm:"'profit'" json:"profit"`
	ProfitRate           float64 `xorm:"'profit_rate'" json:"profit_rate"`
	OrderQty             int64   `xorm:"'order_qty'" json:"order_qty"`
	ProductQty           int64   `xorm:"'product_qty'" json:"product_qty"`
}

/*

SELECT
    r.product_id                                    AS product_id,
    (SELECT p.sku FROM erp_product p
    WHERE p.id = r.product_id)                      AS product_sku,
    (SELECT p.cn_name FROM erp_product p
    WHERE p.id = r.product_id)                      AS product_name,
    (SELECT p.weight FROM erp_product p
    WHERE p.id = r.product_id)                      AS product_weight,
    (SELECT avg(nvt.cost_price) AS cost_price
     FROM erp_inventory nvt
         INNER JOIN erp_stock stk ON stk.id = nvt.stock_id
     WHERE stk.stock_type_id = 300
           AND nvt.product_id = r.product_id)       AS cost_price,
    r.sale_user_id                                  AS sale_user_id,
    (SELECT concat(u.login_name, ' ', u.nick_name)
     FROM auth_user u
     WHERE u.id = r.sale_user_id)                   AS sale_user_name,
    r.service_user_id                               AS service_user_id,
    (SELECT concat(u.login_name, ' ', u.nick_name)
     FROM auth_user u
     WHERE u.id = r.service_user_id)                AS service_user_name,
    r.dev_user_id                                   AS dev_user_id,
    (SELECT concat(u.login_name, ' ', u.nick_name)
     FROM auth_user u
     WHERE u.id = r.dev_user_id)                    AS dev_user_name,
    sum(r.paid_amount)                              AS paid_amount,
    sum(r.express_fee)                              AS express_fee,
    sum(r.income)                                   AS income,
    sum(r.cost_amount)                              AS cost_amount,
    sum(r.estimate_logistics_fee)                   AS estimate_logistics_fee,
    sum(r.actual_logistics_fee)                     AS actual_logistics_fee,
    sum(r.jyfy_amount)                              AS jyfy_amount,
    sum(r.ptfy_amount)                              AS ptfy_amount,
    sum(r.pick_fee)                                 AS pick_fee,
    sum(r.gst_fee)                                  AS gst_fee,
    sum(r.outcome)                                  AS outcome,
    sum(r.income - r.outcome)                       AS profit,
    sum(r.income - r.outcome) / sum(r.income) * 100 AS profit_rate,
    count(DISTINCT r.order_id)                      AS order_qty,
    sum(r.product_qty)                              AS product_qty
FROM erp_sale_amount_report r
WHERE r.platform_id = 0
      AND r.shop_id = 0
      AND r.paid_time >= '2006-01-01'
      AND r.paid_time < '2017-12-30'
      AND r.product_id = 0
      AND r.sale_user_id = 0
      AND r.service_user_id = 0
      AND r.dev_user_id = 0
GROUP BY r.product_id
ORDER BY product_id;

 */
