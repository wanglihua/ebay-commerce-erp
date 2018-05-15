{% site_header "订单退款申请" . %}
<style type="text/css">
    .form-group.has-error span.refund_product_id:after {
        color: #D16E6C;
        display: inline;
        content: "必填";
    }
</style>
<div class="main-container" id="main-container">
    <div id="order_detail" style="width: 1024px;margin-left: auto;margin-right: auto;">
        <div class="widget-box">
            <div class="widget-header">
                <h4 class="smaller">订单信息</h4>
            </div>
            <div class="widget-body">
                <div class="widget-main">
                    <form id="refundForm">
                    <div class="row">
                        <div class="col-xs-3">
                            <div class="form-group">
                                <label for="order_code" class="col-xs-4 control-label">订单编号：</label>
                                <div class="col-xs-8">
                                    <input id="order_code" name="order_code" type="text" class="input-sm form-control" disabled="disabled" v-model="detail.order.code"/>
                                </div>
                            </div>
                        </div>
                        <div class="col-xs-3">
                            <div class="form-group">
                                <label for="created_time" class="col-xs-4 control-label">建单时间：</label>
                                <div class="col-xs-8">
                                    <input id="created_time" name="created_time" type="text" class="input-sm form-control" disabled="disabled" v-model="detail.order.created_time"/>
                                </div>
                            </div>
                        </div>
                        <div class="col-xs-3">
                            <div class="form-group">
                                <label for="status_name" class="col-xs-4 control-label">状态：</label>
                                <div class="col-xs-8">
                                    <input id="status_name" name="status_name" type="text" class="input-sm form-control" disabled="disabled" v-model="detail.status.name"/>
                                </div>
                            </div>
                        </div>
                        <div class="col-xs-3">
                            <div class="form-group">
                                <label for="order_type_name" class="col-xs-4 control-label">订单类型：</label>
                                <div class="col-xs-8">
                                    <input id="order_type_name" name="order_type" type="text" class="input-sm form-control" disabled="disabled" v-model="detail.order_type.name"/>
                                </div>
                            </div>
                        </div>
                    </div>
                    <div class="row">
                        <div class="col-xs-3">
                            <div class="form-group">
                                <label for="order_number" class="col-xs-4 control-label">订单代码：</label>
                                <div class="col-xs-8">
                                    <input id="order_number" name="order_number" type="text" class="input-sm form-control" disabled="disabled" v-model="detail.order.full_tran_id"/>
                                </div>
                            </div>
                        </div>
                        <div class="col-xs-3">
                            <div class="form-group">
                                <label for="order_paid_time" class="col-xs-4 control-label">支付时间：</label>
                                <div class="col-xs-8">
                                    <input id="order_paid_time" name="order_paid_time" type="text" class="input-sm form-control" disabled="disabled" v-model="detail.order.paid_time"/>
                                </div>
                            </div>
                        </div>
                        <div class="col-xs-3">
                            <div class="form-group">
                                <label for="order_hold_type" class="col-xs-4 control-label">hold类型：</label>
                                <div class="col-xs-8">
                                    <input id="order_hold_type" name="order_hold_type" type="text" class="input-sm form-control" disabled="disabled" v-model="detail.order_hold_type.name"/>
                                </div>
                            </div>
                        </div>
                        <div class="col-xs-3">
                            <div class="form-group">
                                <label for="waybill_number" class="col-xs-4 control-label">运单号：</label>
                                <div class="col-xs-8">
                                    <input id="waybill_number" name="waybill_number" type="text" class="input-sm form-control" disabled="disabled" v-model="detail.order.waybill_number"/>
                                </div>
                            </div>
                        </div>
                    </div>
                    <div class="row">
                        <div class="col-xs-12">
                            <div class="form-group">
                                <label for="order_desc" class="col-xs-1 control-label">描述：</label>
                                <div class="col-xs-11">
                                    <textarea id="order_desc" name="order_desc" rows="2" class="input-sm form-control" disabled="disabled" v-model="detail.order.desc"></textarea>
                                </div>
                            </div>
                        </div>
                    </div>
                    <hr style="margin-top: 0px;margin-bottom: 10px;"/>
                    <div class="row">
                        <div class="col-xs-12">
                            <table id="order_price_report" class="table table-bordered">
                                <thead>
                                <tr>
                                    <th class="dtRight">结账信息</th>
                                    <th class="dtRight">结算汇率</th>
                                    <th class="dtRight">收款小计</th>
                                    <th class="dtRight">运费收入</th>
                                    <th class="dtRight">GST税</th>
                                    <th class="dtRight">平台交易费</th>
                                    <th class="dtRight">付款交易费</th>
                                    <th class="dtRight">拣货费</th>
                                    <th class="dtRight">预估物流费</th>
                                    <th class="dtRight">实际物流费</th>
                                </tr>
                                </thead>
                                <tbody>
                                <tr>
                                    <td class="dtRight">收款货币</td>
                                    <td class="dtRight">
                                        AUD {{ format_amount(1.00) }}
                                    </td>
                                    <td class="dtRight">
                                        AUD {{ format_amount(detail.order.subtotal_amount) }}
                                    </td>
                                    <td class="dtRight">
                                        AUD {{ format_amount(detail.order.express_fee) }}
                                    </td>
                                    <td class="dtRight">
                                        AUD {{ format_amount(detail.order_gst) }}
                                    </td>
                                    <td class="dtRight">
                                        AUD {{ format_amount(detail.order.ptfy) }}
                                    </td>
                                    <td class="dtRight">
                                        AUD {{ format_amount(detail.order.jyfy) }}
                                    </td>
                                    <td class="dtRight">
                                        AUD {{ format_amount(detail.order_logistics.pick_fee) }}
                                    </td>
                                    <td class="dtRight">
                                        AUD {{ format_amount(detail.order_logistics.total_fee) }}
                                    </td>
                                    <td class="dtRight">
                                        AUD {{ format_amount(detail.order_logistics.actual_fee) }}
                                    </td>
                                </tr>
                                <tr>
                                    <td class="dtRight">结算人民币</td>
                                    <td class="dtRight">
                                        RMB {{ format_amount(detail.order.exchange_rate) }}
                                    </td>
                                    <td class="dtRight">
                                        RMB {{ format_amount(detail.order.subtotal_amount_rmb) }}
                                    </td>
                                    <td class="dtRight">
                                        RMB {{ format_amount(detail.order.express_fee * detail.order.exchange_rate)}}
                                    </td>
                                    <td class="dtRight">
                                        RMB {{ format_amount(detail.order_gst * detail.order.exchange_rate) }}
                                    </td>
                                    <td class="dtRight">
                                        RMB {{ format_amount(detail.order.ptfy  * detail.order.exchange_rate) }}
                                    </td>
                                    <td class="dtRight">
                                        RMB {{ format_amount(detail.order.jyfy  * detail.order.exchange_rate) }}
                                    </td>
                                    <td class="dtRight">
                                        RMB {{ format_amount(detail.order_logistics.pick_fee*detail.order.exchange_rate) }}
                                    </td>
                                    <td class="dtRight">
                                        RMB {{ format_amount(detail.order_logistics.total_fee*detail.order.exchange_rate) }}
                                    </td>
                                    <td class="dtRight">
                                        RMB {{ format_amount(detail.order_logistics.actual_fee*detail.order.exchange_rate) }}
                                    </td>
                                </tr>
                                </tbody>
                                <tfoot>
                                <tr>
                                    <td></td>
                                    <td></td>
                                    <td></td>
                                    <td></td>
                                    <td colspan="2">商品澳仓成本
                                        RMB {{ format_amount(detail.order_logistics.cost_amount * detail.order.exchange_rate) }}
                                    </td>
                                    <td colspan="2">物流成本 RMB {{ format_amount(detail.logistics_fee) }}
                                    </td>
                                    <td colspan="2">毛利 RMB {{ format_amount(detail.order_profit) }}</td>
                                </tr>
                                </tfoot>
                            </table>
                        </div>
                    </div>
                    <hr style="margin-top: 0px;margin-bottom: 10px;"/>
                    <div class="row">
                        <div class="col-xs-12">
                            <table id="order_item_list" class="table table-striped table-bordered">
                                <thead>
                                <tr>
                                    <th class="center" style="width: 50px;">选择</th>
                                    <th class="center" style="width: 60px;">图片</th>
                                    <th style="width: 120px;">商品名称</th>
                                    <th class="center" style="width: 60px;">数量</th>
                                    <th class="center" style="width: 100px;">售价</th>
                                    <th class="center" style="width: 100px;">澳仓成本价</th>
                                    <th class="center" style="width: 80px;">原始售价</th>
                                    <th class="center" style="width: 80px;">毛利</th>
                                    <th>留言</th>
                                </tr>
                                </thead>
                                <tbody>
                                <template v-for="order_item in detail.order_item_list">
                                    <tr>
                                        <td class="center">
                                            <div class="form-group">
                                            <label class="pos-rel">
                                                <input type="radio" name="refund_product_id" class="ace" :value="order_item.p.id"/>
                                                <span class="lbl refund_product_id"></span>
                                            </label>
                                            </div>
                                        </td>
                                        <td class="center">
                                            <a :href="order_item.psi.view_item_url" target="_blank">
                                                <img :src="order_item.p.image_url" alt="" style="width: 60px;height: 60px;"/>
                                            </a>
                                        </td>
                                        <td>
                                            <div :title="order_item.p.cn_name" style="width: 150px; white-space:normal; word-break: break-all; word-wrap:break-word;">
                                                {{order_item.p.cn_name}}
                                            </div>
                                            <div>
                                                <a :href="'/product/detail/' + order_item.p.id" target="product_detail">
                                                    {{order_item.p.sku}}
                                                </a>
                                            </div>
                                            <div>
                                                <a :href="order_item.psi.view_item_url" target="_blank">
                                                    {{order_item.oi.item_id}}
                                                </a>
                                            </div>
                                        </td>
                                        <td class="center">
                                            {{order_item.oi.quantity}}
                                        </td>
                                        <td class="center">
                                            <div>RMB
                                                {{ format_amount(order_item.oi.transaction_price_rmb)}} × {{order_item.oi.quantity}}
                                            </div>
                                            <div>=RMB
                                                {{format_amount(order_item.oi.transaction_price_rmb * order_item.oi.quantity)}}
                                            </div>
                                        </td>
                                        <td class="center">
                                            RMB {{ format_amount(order_item.oi.cost_price * detail.order.exchange_rate) }}
                                        </td>
                                        <td class="center">
                                            {{order_item.oi.transaction_price_currency}}
                                            {{ format_amount(order_item.oi.transaction_price)}}
                                        </td>
                                        <td class="center">
                                            RMB {{ format_amount(order_item.oi.gross_profit * detail.order.exchange_rate) }}
                                        </td>
                                        <td></td>
                                    </tr>
                                </template>
                                </tbody>
                            </table>
                        </div>
                    </div>
                    <hr style="margin-top: 0px;margin-bottom: 10px;"/>
                    <div class="row">
                        <div class="col-xs-3">
                            <input type="hidden" id="refund_order_id" name="refund_order_id" value="{% .order_id %}" />
                            <div class="form-group">
                                <label for="refund_type_id" class="col-xs-4 control-label">退款类型：</label>
                                <div class="col-xs-8">
                                    <select id="refund_type_id" name="refund_type_id" class="input-sm form-control" >
                                    <option value="">--退款类型--</option>
                                    {% range $index, $refund_type := .order_refund_type_list %}
                                    <option value="{% $refund_type.Id %}">{% $refund_type.Name %}</option>
                                    {% end %}
                                    </select>
                                </div>
                            </div>
                        </div>
                        <div class="col-xs-3">
                            <div class="form-group">
                                <label for="refund_reason_id" class="col-xs-4 control-label">退款原因：</label>
                                <div class="col-xs-8">
                                   <select id="refund_reason_id" name="refund_reason_id" class="input-sm form-control" >
                                    <option value="">--退款原因--</option>
                                    {% range $index, $refund_reason := .order_refund_reason_list %}
                                    <option value="{% $refund_reason.Id %}">{% $refund_reason.Name %}</option>
                                    {% end %}
                                    </select>
                                </div>
                            </div>
                        </div>
                        <div class="col-xs-3">
                            <div class="form-group">
                                <label for="refund_method_id" class="col-xs-4 control-label">退款渠道：</label>
                                <div class="col-xs-8">
                                   <select id="refund_method_id" name="refund_method_id" class="input-sm form-control" >
                                    <option value="">--退款渠道--</option>
                                    {% range $index, $refund_method := .order_refund_method_list %}
                                    <option value="{% $refund_method.Id %}">{% $refund_method.Name %}</option>
                                    {% end %}
                                    </select>
                                </div>
                            </div>
                        </div>
                        <div class="col-xs-3">
                            <div class="form-group">
                                <label for="refund_amount" class="col-xs-4 control-label">退款金额：</label>
                                <div class="col-xs-8">
                                    <input id="refund_amount" name="refund_amount" type="text" class="input-sm form-control" placeholder="AUD" />
                                </div>
                            </div>
                        </div>
                    </div>
                    <div class="row">
                        <div class="col-xs-12">
                            <div class="form-group">
                                <label for="refund_description" class="col-xs-1 control-label">退款描述：</label>
                                <div class="col-xs-11">
                                    <textarea id="refund_description" name="refund_description" rows="2" class="input-sm form-control"></textarea>
                                </div>
                            </div>
                        </div>
                    </div>
                    <hr style="margin-top: 5px;margin-bottom: 5px;"/>
                    <div class="row">
                        <div class="col-xs-12" style="text-align: right;">
                            <button class="btn btn-sm btn-purple" onclick="orderRefundApply()">申请退款</button>
                        </div>
                    </div>
                    </form>
                </div>
            </div>
        </div>
    </div>
</div>
{%template "erp/ppgl/order_refund/order_refund_apply.var" .%}
<script type="text/javascript">
    var vm = new Vue({
        el: "#order_detail",
        data: {
            "detail": {
                "order": {},
                "status": {},
                "shop": {},
                "order_type": {},
                "order_hold_type": {},
                "order_logistics": {},
                "order_item_list": []
            }
        },
        mounted: function () {
            var this_vm = this;

            $.ajax(
                {
                    "dataType": 'json',
                    "url": '{% url "ErpOrderRefund.ApplyShowData"%}',
                    "data": {
                        "id": order_id
                    },
                    "success": function (data) {
                        if (data.success === false) {
                            $.msg.error(data.message);
                            return;
                        }
                        this_vm.detail = data;
                    }
                });
        },
        methods: {
            format_amount: function (amount) {
                return dtFormatFloat(amount);
            }
        }
    });

    function orderRefundApply() {
        showMask("提交中...");
        var form = $("#refundForm");
        var valid = validate(form, {
            'refund_order_id': {
                required: true
            },
            'refund_product_id': {
                required: true
            },
            'refund_type_id': {
                required: true
            },
            'refund_reason_id': {
                required: true
            },
            'refund_method_id': {
                required: true
            },
            'refund_amount': {
                required: true, number: true
            },
            'refund_description': {
                required: false, minlength: 0, maxlength: 400
            }
        });

        if (!valid) {
            hideMask();
        } else {
            $.ajax({
                    url: '{%url "ErpOrderRefund.Apply"%}',
                    type: "POST",
                    data: form.serialize(),
                    success: function (data) {
                        hideMask();
                        if (data.success === false) {
                            $.msg.error(data.message);
                            //return;
                        }
                        else {
                            $.msg.show(data.message);
                            //return;
                        }
                    }
                }
            );
        }
    }

    function load_order_detail() {
        $.ajax({
            "dataType": 'json',
            "url": '{% url "ErpOrder.DetailData"%}',
            "data": {
                "id": order_id
            },
            "success": function (data) {
                if (data.success === false) {
                    $.msg.error(data.message);
                    return;
                }

                vm.detail = data;
                vm.order_item_detail = {};
                vm.order_item_detail.order_id = order_id;
            }
        });
    }

</script>
{% site_footer %}
