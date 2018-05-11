{% site_header "订单详情" . %}
<div class="main-container" id="main-container">
    <div id="order_detail" style="width: 1024px;margin-left: auto;margin-right: auto;">
        <div class="widget-box">
            <div class="widget-header">
                <h4 class="smaller">订单信息</h4>
            </div>
            <div class="widget-body">
                <div class="widget-main">
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
                    {% if ne .order_feedback_list_len 0 %}
                    <div class="row">
                        <div class="col-xs-12">
                            <div class="form-group">
                                <label for="order_desc" class="col-xs-1 control-label">订单评价：</label>
                                <div class="col-xs-11">
                                    {% range $index, $feedback := .order_feedback_list %}
                                    <div>
                                        <span style="padding-left:25px;line-height: 24px;" class="feedback_detail_{% $feedback.OrderFeedback.CommentTypeText %}">
{% $feedback.OrderFeedback.CommentTypeText %}，SKU：{% $feedback.OrderItem.ItemSku %}, 评价时间：{% $feedback.OrderFeedback.CommentTimeText %}，评论：{% $feedback.OrderFeedback.CommentText %}
                                        </span>
                                    </div>
                                    {% end %}
                                </div>
                            </div>
                        </div>
                    </div>
                    {% end %}
                    <div class="row">
                        <div class="col-xs-12">
                            <div class="form-group">
                                <label for="order_desc" class="col-xs-1 control-label">描述：</label>
                                <div class="col-xs-11">
                                    <textarea id="order_desc" name="order_desc" rows="2" class="input-sm form-control" v-model="detail.order.desc"></textarea>
                                </div>
                            </div>
                        </div>
                    </div>

                    <hr style="margin-top: 5px;margin-bottom: 5px;"/>

                    <div class="row">
                        <div class="col-xs-12" style="text-align:right; margin-top: 5px;">
                            <button class="btn btn-sm" style="border-width: 2px; padding: 3px 5px;" onclick="save_order_desc()">保存描述</button>
                            {% if .can_reship %}
                            <button class="btn btn-sm btn-primary" style="border-width: 2px; padding: 3px 5px;" onclick="reship_order()">重发订单</button>
                            {% end %}
                        </div>
                    </div>

                    <div class="row">
                        <div class="col-xs-12" style="text-align: right;">
                            {% if or .is_hold .is_dph %}
                            {% if .is_express %}
                            <label for="unlock_express">解锁快递：</label>
                            <select id="unlock_express" class="input-sm" style="width: 80px;">
                                <option :value="0">否</option>
                                <option :value="1">是</option>
                            </select>
                            &nbsp;
                            {% end %}
                            <label for="prefer_stock">指定匹配仓库：</label>
                            <select id="prefer_stock" class="input-sm" style="width: 120px;" onchange="load_logistics_model_list()">
                                <option :value="0">当前不指定</option>
                                <template v-for="stock in detail.au_stock_list">
                                    <option :value="stock.id">{{stock.name}}</option>
                                </template>
                            </select>
                            &nbsp;<label for="prefer_mode">指定物流方式：</label>
                            <select id="prefer_mode" class="input-sm" style="width: 120px;">
                                <option :value="0">当前不指定</option>
                                <template v-for="mode in logistics_mode_list">
                                    <option :value="mode.id">{{mode.name}}</option>
                                </template>
                            </select>
                            <button class="btn btn-sm btn-purple" onclick="rematch_order_logistics()">
                                重新匹配物流
                            </button>
                            {% end %}
                            {% if .can_edit_waybill_number %}
                            <button class="btn btn-sm btn-purple" @click="show_order_waybill_number_edit">修改运单号</button>
                            {% end %}
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <div class="widget-box">
            <div class="widget-header">
                <h4 class="smaller">客户信息</h4>
            </div>
            <div class="widget-body">
                <div class="widget-main">
                    <div class="row">
                        <div class="col-xs-3">
                            <div class="form-group">
                                <label for="shop_name" class="col-xs-4 control-label">客户来源：</label>
                                <div class="col-xs-8">
                                    <input id="shop_name" name="detail.shop.name" type="text" class="input-sm form-control" disabled="disabled" v-model="detail.shop.name"/>
                                </div>
                            </div>
                        </div>
                        <div class="col-xs-3">
                            <div class="form-group">
                                <label for="buyer_user_id" class="col-xs-4 control-label">买家帐号：</label>
                                <div class="col-xs-8">
                                    <input id="buyer_user_id" name="detail.order.buyer_user_id" type="text" class="input-sm form-control" disabled="disabled" v-model="detail.order.buyer_user_id"/>
                                </div>
                            </div>
                        </div>
                        <div class="col-xs-3">
                            <div class="form-group">
                                <label for="buyer_user_name" class="col-xs-4 control-label">客户姓名：</label>
                                <div class="col-xs-8">
                                    <input id="buyer_user_name" name="detail.order.buyer_user_name" type="text" class="input-sm form-control" disabled="disabled" v-model="detail.shipping_address.name"/>
                                </div>
                            </div>
                        </div>
                        <div class="col-xs-3">
                            <div class="form-group">
                                <label for="country_name" class="col-xs-4 control-label">国家：</label>
                                <div class="col-xs-8">
                                    <input id="country_name" name="detail.shipping_address.country_name" type="text" class="input-sm form-control" disabled="disabled" v-model="detail.shipping_address.country_name"/>
                                </div>
                            </div>
                        </div>
                        <div class="col-xs-3">
                            <div class="form-group">
                                <label for="phone" class="col-xs-4 control-label">联系电话：</label>
                                <div class="col-xs-8">
                                    <input id="phone" name="detail.shipping_address.country_name" type="text" class="input-sm form-control" disabled="disabled" v-model="detail.shipping_address.phone"/>
                                </div>
                            </div>
                        </div>
                        <div class="col-xs-9" style="height: 32px;"></div>

                        <div class="col-xs-6">
                            <div class="form-group">
                                <label for="order_buyer_email" class="col-xs-2 control-label">电子邮件：</label>
                                <input id="order_buyer_email" class="col-xs-10" type="text" disabled="disabled" v-model="detail.order.buyer_email"/>
                            </div>
                        </div>
                        <div class="col-xs-6" style="height: 32px;"></div>
                        <div class="col-xs-6">
                            <div class="form-group">
                                <label for="ebay_address" class="col-xs-2 control-label">平台地址：</label>
                                <div class="col-xs-10">
                                    <div id="ebay_address" class="address-box">
                                        <div>
                                            <label for="ebay_address_name">Name:</label>
                                            <input id="ebay_address_name" type="text" v-model="detail.ebay_shipping_address.name"/>
                                        </div>
                                        <div>
                                            <label>Address:</label>
                                            <span>{{detail.ebay_shipping_address.street1}},{{detail.ebay_shipping_address.street2}}</span>
                                        </div>
                                        <div>
                                            <label for="ebay_address_city">City:</label>
                                            <input id="ebay_address_city" type="text" v-model="detail.ebay_shipping_address.city_name"/>
                                        </div>
                                        <div>
                                            <label for="ebay_address_state">State:</label>
                                            <input id="ebay_address_state" type="text" v-model="detail.ebay_shipping_address.province"/>
                                        </div>
                                        <div>
                                            <label for="ebay_address_country">Country:</label>
                                            <input id="ebay_address_country" type="text" v-model="detail.ebay_shipping_address.country"/>
                                        </div>
                                        <div>
                                            <label for="ebay_address_country_name">Country Name:</label>
                                            <input id="ebay_address_country_name" type="text" v-model="detail.ebay_shipping_address.country_name"/>
                                        </div>
                                        <div>
                                            <label for="ebay_address_postal_code">Post Code:</label>
                                            <input id="ebay_address_postal_code" type="text" v-model="detail.ebay_shipping_address.postal_code"/>
                                        </div>
                                        <div>
                                            <label for="ebay_address_phone">Phone:</label>
                                            <input id="ebay_address_phone" type="text" v-model="detail.ebay_shipping_address.phone"/>
                                        </div>
                                    </div>

                                </div>
                            </div>
                        </div>
                        <div class="col-xs-6">
                            <div class="form-group">
                                <label for="paypal_address" class="col-xs-2 control-label">PayPal<br/>地址：</label>
                                <div class="col-xs-10">
                                    <div id="paypal_address" class="address-box">
                                        <div>
                                            <label for="paypal_address_name">Name:</label>
                                            <input id="paypal_address_name" type="text" v-model="detail.paypal_shipping_address.name"/>
                                        </div>
                                        <div>
                                            <label>Address:</label>
                                            <span>{{detail.paypal_shipping_address.street1}},{{detail.paypal_shipping_address.street2}}</span>
                                        </div>
                                        <div>
                                            <label for="paypal_address_city">City:</label>
                                            <input id="paypal_address_city" type="text" v-model="detail.paypal_shipping_address.city_name"/>
                                        </div>
                                        <div>
                                            <label for="paypal_address_state">State:</label>
                                            <input id="paypal_address_state" type="text" v-model="detail.paypal_shipping_address.province"/>
                                        </div>
                                        <div>
                                            <label for="paypal_address_country">Country:</label>
                                            <input id="paypal_address_country" type="text" v-model="detail.paypal_shipping_address.country"/>
                                        </div>
                                        <div>
                                            <label for="paypal_address_country_name">Country Name:</label>
                                            <input id="paypal_address_country_name" type="text" v-model="detail.paypal_shipping_address.country_name"/>
                                        </div>
                                        <div>
                                            <label for="paypal_address_postal_code">Post Code:</label>
                                            <input id="paypal_address_postal_code" type="text" v-model="detail.paypal_shipping_address.postal_code"/>
                                        </div>
                                        <div>
                                            <label for="paypal_address_phone">Phone:</label>
                                            <input id="paypal_address_phone" type="text" v-model="detail.paypal_shipping_address.phone"/>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                    {% if .can_modify %}
                    <hr style="margin-top: 5px;margin-bottom: 5px;"/>
                    <div class="row">
                        <div class="col-xs-12" style="text-align: right;">
                            <button class="btn btn-sm btn-purple" @click="show_address_edit">修改</button>
                        </div>
                    </div>
                    {% end %}
                </div>
            </div>
        </div>

        <div class="widget-box">
            <div class="widget-header">
                <h4 class="smaller">订单报表</h4>
            </div>
            <div class="widget-body">
                <div class="widget-main">
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
                    {% if .can_refund %}
                    <hr style="margin-top: 5px;margin-bottom: 5px;"/>
                    <div class="row">
                        <div class="col-xs-12" style="text-align: right;">
                            <button class="btn btn-sm btn-purple" onclick="showRefund()">退款</button>
                        </div>
                    </div>
                    {% end %}
                </div>
            </div>
        </div>

        <div class="widget-box" v-if="detail.refund_list.length != 0">
            <div class="widget-header">
                <h4 class="smaller">订单退款</h4>
            </div>
            <div class="widget-body">
                <div class="widget-main">
                    <div class="row">
                        <div class="col-xs-12">
                            <table id="order_refund_list" class="table table-bordered">
                                <thead>
                                <tr>
                                    <th class="dtCenter">Transaction Id</th>
                                    <th class="dtRight">退款金额</th>
                                    <th>备注</th>
                                    <th class="dtCenter">退款时间</th>
                                </tr>
                                </thead>
                                <tbody>
                                <tr v-for="refund in detail.refund_list">
                                    <td class="dtCenter">{{ refund.transaction_id }}</td>
                                    <td class="dtRight">{{ refund.amount }}</td>
                                    <td>{{ refund.description }}</td>
                                    <td class="dtCenter">{{ refund.create_time }}</td>
                                </tr>
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <div class="widget-box">
            <div class="widget-header">
                <h4 class="smaller">订单详情</h4>
            </div>
            <div class="widget-body">
                <div class="widget-main">
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
                                            <label class="pos-rel">
                                                <input type="checkbox" class="ace" :value="order_item.oi.id"/>
                                                <span class="lbl"></span>
                                            </label>
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
                    {% if .can_modify %}
                    <hr style="margin-top: 5px;margin-bottom: 5px;"/>
                    <div class="row">
                        <div class="col-xs-12" style="text-align: right;">
                            <button class="btn btn-sm btn-success" @click="add_order_item">添加</button>
                            <button class="btn btn-sm btn-purple" @click="show_order_item">修改</button>
                            <button class="btn btn-sm btn-danger" @click="delete_order_item">删除</button>
                            <button class="btn btn-sm btn-pink" onclick="showOrderSplit()">拆分</button>
                        </div>
                    </div>
                    {% end %}
                </div>
            </div>
        </div>

        {% if .is_pickup %}
        <div class="widget-box">
            <div class="widget-header">
                <h4 class="smaller">PICK UP 信息</h4>
            </div>
            <div class="widget-body">
                <div class="widget-main">
                    <form id="pickupForm">
                        <div class="row">
                            <div class="col-xs-3">
                                <div class="form-group">
                                    <label for="creator" class="col-xs-5 control-label">建单人员：</label>
                                    <div class="col-xs-7">
                                        <input id="creator" name="creator" type="text" class="input-sm form-control" disabled="disabled" v-model="pickup.creator"/>
                                    </div>
                                </div>
                            </div>
                            <div class="col-xs-3">
                                <div class="form-group">
                                    <label for="payee" class="col-xs-5 control-label">收款人员：</label>
                                    <div class="col-xs-7">
                                        <input id="payee" name="payee" type="text" class="input-sm form-control" disabled="disabled" v-model="pickup.payee"/>
                                    </div>
                                </div>
                            </div>
                            <div class="col-xs-3">
                                <div class="form-group">
                                    <label for="is_paid" class="col-xs-5 control-label">已收款：</label>
                                    <div class="col-xs-7">
                                        <select id="is_paid" name="is_paid" class="input-sm form-control" v-model.number="pickup.is_paid" disabled="disabled">
                                            <option :value="0">否</option>
                                            <option :value="1">是</option>
                                        </select>
                                    </div>
                                </div>
                            </div>
                            <div class="col-xs-3">
                                <div class="form-group">
                                    <label for="pay_method_id" class="col-xs-5 control-label">收款方式：</label>
                                    <div class="col-xs-7">
                                        <select id="pay_method_id" name="pay_method_id" class="input-sm form-control" v-model.number="pickup.pay_method_id" disabled="disabled">
                                            <template v-for="pay_method in pickup_pay_method_list">
                                                <option :value="pay_method.id">{{pay_method.name}}</option>
                                            </template>
                                        </select>
                                    </div>
                                </div>
                            </div>
                        </div>

                        {% if eq .is_pickup_paid 0 %}
                        {% if .can_modify %}
                        <div class="row">
                            <div class="col-xs-3">
                                <div class="form-group">
                                    <label for="pickup_actual_is_paid" class="col-xs-5 control-label">现场已收款：</label>
                                    <div class="col-xs-7">
                                        <select id="pickup_actual_is_paid" name="pickup_actual_is_paid" class="input-sm form-control" v-model.number="pickup.actual_is_paid">
                                            <option :value="0">否</option>
                                            <option :value="1">是</option>
                                        </select>
                                    </div>
                                </div>
                            </div>
                            <div class="col-xs-3">
                                <div class="form-group">
                                    <label for="pickup_actual_pay_method_id" class="col-xs-5 control-label">现场收款方式：</label>
                                    <div class="col-xs-7">
                                        <select id="pickup_actual_pay_method_id" name="pickup_actual_pay_method_id" class="input-sm form-control" v-model.number="pickup.actual_pay_method_id">
                                            <template v-for="pay_method in pickup_actual_pay_method_list">
                                                <option :value="pay_method.id">{{pay_method.name}}</option>
                                            </template>
                                        </select>
                                    </div>
                                </div>
                            </div>
                            <div class="col-xs-3">
                                <div class="form-group">
                                    <label for="pickup_actual_amount" class="col-xs-5 control-label">实际收款澳币：</label>
                                    <div class="col-xs-7">
                                        <input id="pickup_actual_amount" name="pickup_actual_amount" class="input-sm form-control" v-model.number="pickup.actual_amount"/>
                                    </div>
                                </div>
                            </div>
                            <div class="col-xs-3">
                                <div class="form-group">
                                    <label for="pickup_description" class="col-xs-5 control-label">描述：</label>
                                    <div class="col-xs-7">
                                        <input id="pickup_description" name="pickup_description" class="input-sm form-control" v-model="pickup.description"/>
                                    </div>
                                </div>
                            </div>
                            <!--description-->
                        </div>
                        {% else %}
                        <div class="row">
                            <div class="col-xs-3">
                                <div class="form-group">
                                    <label for="pickup_actual_is_paid" class="col-xs-5 control-label">现场已收款：</label>
                                    <div class="col-xs-7">
                                        <select id="pickup_actual_is_paid" name="pickup_actual_is_paid" class="input-sm form-control" v-model.number="pickup.actual_is_paid" disabled="disabled">
                                            <option :value="0">否</option>
                                            <option :value="1">是</option>
                                        </select>
                                    </div>
                                </div>
                            </div>
                            <div class="col-xs-3">
                                <div class="form-group">
                                    <label for="pickup_actual_pay_method_id" class="col-xs-5 control-label">现场收款方式：</label>
                                    <div class="col-xs-7">
                                        <select id="pickup_actual_pay_method_id" name="pickup_actual_pay_method_id" class="input-sm form-control" v-model.number="pickup.actual_pay_method_id" disabled="disabled">
                                            <template v-for="pay_method in pickup_actual_pay_method_list">
                                                <option :value="pay_method.id">{{pay_method.name}}</option>
                                            </template>
                                        </select>
                                    </div>
                                </div>
                            </div>
                            <div class="col-xs-3">
                                <div class="form-group">
                                    <label for="pickup_actual_amount" class="col-xs-5 control-label">实际收款澳币：</label>
                                    <div class="col-xs-7">
                                        <input id="pickup_actual_amount" name="pickup_actual_amount" class="input-sm form-control" v-model.number="pickup.actual_amount" disabled="disabled"/>
                                    </div>
                                </div>
                            </div>
                            <div class="col-xs-3">
                                <div class="form-group">
                                    <label for="pickup_description" class="col-xs-5 control-label">描述：</label>
                                    <div class="col-xs-7">
                                        <input id="pickup_description" name="pickup_description" class="input-sm form-control" v-model="pickup.description" disabled="disabled"/>
                                    </div>
                                </div>
                            </div>
                        </div>
                        {% end %}
                    </form>
                    {% if .can_modify %}
                    <hr style="margin-top: 5px;margin-bottom: 5px;"/>
                    <div class="row">
                        <div class="col-xs-12" style="text-align: right;">
                            <button class="btn btn-sm btn-purple" @click="savePickupActual">保存</button>
                        </div>
                    </div>
                    {% end %}
                    {% end %}
                </div>
            </div>
        </div>
        {% end %}

        <div class="widget-box">
            <div class="widget-header">
                <h4 class="smaller">订单商品库存</h4>
            </div>
            <div class="widget-body">
                <div class="widget-main">
                    <div class="row">
                        <div class="col-xs-12">
                            订单已匹配物流仓库：&nbsp;<span style="color:blue;">{{detail.logistics_stock.name}}</span>&nbsp;&nbsp;&nbsp;
                            订单已匹配物流方式：&nbsp;<span style="color:blue;">{{detail.logistics_mode.name}}</span>
                        </div>
                    </div>
                    <div class="row">
                        <div class="col-xs-12">
                            <table id="product_inventory_list" class="table table-striped table-bordered">
                                <thead>
                                <tr>
                                    <th class="center" style="width: 100px;">SKU</th>
                                    <th class="center">商品名称</th>
                                    <th class="center" style="width: 120px;">仓库</th>
                                    <th class="center" style="width: 120px;">库位</th>
                                    <th class="center" style="width: 80px;">库存</th>
                                    <th class="center" style="width: 80px;">待发货</th>
                                    <th class="center" style="width: 100px;">成本</th>
                                    <th class="center" style="width: 80px;">到仓时间</th>
                                </tr>
                                </thead>
                                <tbody>
                                <template v-for="inventory in detail.inventory_list">
                                    <tr v-if="inventory.st.id == detail.logistics_stock.id" style="background-color: #dde6f0;">
                                        <td class="center" style="vertical-align: middle">
                                            {{inventory.p.sku}}
                                        </td>
                                        <td class="wrap-word" style="vertical-align: middle">
                                            <div>{{inventory.p.cn_name}}</div>
                                        </td>
                                        <td class="center" style="vertical-align: middle">
                                            {{inventory.st.name}}
                                        </td>
                                        <td class="center" style="vertical-align: middle">
                                            {{inventory.sl.name}}
                                        </td>
                                        <td class="center" style="vertical-align: middle">
                                            {{inventory.nvt.qty}}
                                        </td>
                                        <td class="center" style="vertical-align: middle">
                                            {{inventory.nvt.pending_qty}}
                                        </td>
                                        <td class="center" style="vertical-align: middle">
                                            {{ format_amount(inventory.nvt.cost_price)}}
                                        </td>
                                        <td class="center" style="vertical-align: middle">
                                            {{ inventory.stock_time }}
                                        </td>
                                    </tr>
                                    <tr v-else>
                                        <td class="center" style="vertical-align: middle">
                                            {{inventory.p.sku}}
                                        </td>
                                        <td class="wrap-word" style="vertical-align: middle">
                                            <div>{{inventory.p.cn_name}}</div>
                                        </td>
                                        <td class="center" style="vertical-align: middle">
                                            {{inventory.st.name}}
                                        </td>
                                        <td class="center" style="vertical-align: middle">
                                            {{inventory.sl.name}}
                                        </td>
                                        <td class="center" style="vertical-align: middle">
                                            {{inventory.nvt.qty}}
                                        </td>
                                        <td class="center" style="vertical-align: middle">
                                            {{inventory.nvt.pending_qty}}
                                        </td>
                                        <td class="center" style="vertical-align: middle">
                                            {{ format_amount(inventory.nvt.cost_price)}}
                                        </td>
                                        <td class="center" style="vertical-align: middle">
                                            <div v-if="inventory.transfer_id != 0">
                                                <a :href="'/container/transfer/detail/' + inventory.transfer_id" target="_blank">{{ inventory.stock_time }}</a>
                                            </div>
                                        </td>
                                    </tr>
                                </template>
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <div class="widget-box">
            <div class="widget-header">
                <h4 class="smaller">操作日志</h4>
            </div>
            <div class="widget-body">
                <div class="widget-main">
                    <table id="orderHistoryList" class="table table-striped table-bordered table-hover"></table>
                </div>
            </div>
        </div>

        <div class="modal fade" id="orderWaybillNumberDialog">
            <div class="modal-dialog" style="width: 480px;">
                <div class="modal-content">
                    <div class="modal-header">
                        <button type="button" class="close" data-dismiss="modal" aria-hidden="True">&times;</button>
                        <h4 class="modal-title">订单运单号</h4>
                    </div>
                    <div class="modal-body clearfix">
                        <div class="row">
                            <div class="col-xs-12">
                                <form class="form-horizontal" role="form" id="orderWaybillNumberForm">
                                    <input id="waybill_number_order_id" name="detail.order.id" type="hidden" v-model.number="detail.order.id"/>

                                    <div class="row">
                                        <div class="col-xs-12">
                                            <div class="form-group">
                                                <label for="order_waybill_number" class="col-xs-3 control-label">运单号</label>
                                                <div class="col-xs-9">
                                                    <input id="order_waybill_number" name="detail.order.waybill_number" type="text" class="input-sm form-control" v-model="detail.order.waybill_number"/>
                                                </div>
                                            </div>
                                        </div>
                                    </div>

                                </form>
                            </div>
                        </div>
                    </div>
                    <div class="modal-footer">
                        <button type="button" class="btn btn-success btn-sm" @click="saveOrderWaybillNumber()">
                            <i class="ace-icon fa fa-floppy-o"></i>保存
                        </button>
                        <button type="button" class="btn btn-default btn-sm" data-dismiss="modal">
                            <i class="ace-icon fa fa-times"></i>关闭
                        </button>
                    </div>
                </div>
            </div>
        </div>

        <div class="modal fade" id="orderShippingAddressDetailDialog">
            <div class="modal-dialog" style="width: 680px;">
                <div class="modal-content">
                    <div class="modal-header">
                        <button type="button" class="close" data-dismiss="modal" aria-hidden="True">&times;</button>
                        <h4 class="modal-title">订单送货地址详细</h4>
                    </div>
                    <div class="modal-body clearfix">
                        <div class="row">
                            <div class="col-xs-12">
                                <form class="form-horizontal" role="form" id="orderShippingAddressDetailForm">
                                    <input id="order_shipping_address_id" name="detail.shipping_address.id" type="hidden" v-model.number="detail.shipping_address.id"/>
                                    <input id="order_id" name="detail.shipping_address.order_id" type="hidden" v-model.number="detail.shipping_address.order_id"/>

                                    <div class="row">
                                        <div class="col-xs-6">
                                            <div class="form-group">
                                                <label for="name" class="col-xs-4 control-label">Name</label>
                                                <div class="col-xs-8">
                                                    <input id="name" name="detail.shipping_address.name" type="text" class="input-sm form-control" v-model="detail.shipping_address.name"/>
                                                </div>
                                            </div>
                                        </div>
                                    </div>

                                    <div class="row">
                                        <div class="col-xs-6">
                                            <div class="form-group">
                                                <label for="street1" class="col-xs-4 control-label">Street1</label>
                                                <div class="col-xs-8">
                                                    <input id="street1" name="detail.shipping_address.street1" type="text" class="input-sm form-control" v-model="detail.shipping_address.street1"/>
                                                </div>
                                            </div>
                                        </div>

                                        <div class="col-xs-6">
                                            <div class="form-group">
                                                <label for="street2" class="col-xs-4 control-label">Street2</label>
                                                <div class="col-xs-8">
                                                    <input id="street2" name="detail.shipping_address.street2" type="text" class="input-sm form-control" v-model="detail.shipping_address.street2"/>
                                                </div>
                                            </div>
                                        </div>
                                    </div>

                                    <div class="row">
                                        <div class="col-xs-6">
                                            <div class="form-group">
                                                <label for="city_name" class="col-xs-4 control-label">City</label>
                                                <div class="col-xs-8">
                                                    <input id="city_name" name="detail.shipping_address.city_name" type="text" class="input-sm form-control" v-model="detail.shipping_address.city_name"/>
                                                </div>
                                            </div>
                                        </div>

                                        <div class="col-xs-6">
                                            <div class="form-group">
                                                <label for="province" class="col-xs-4 control-label">State</label>
                                                <div class="col-xs-8">
                                                    <input id="province" name="detail.shipping_address.province" type="text" class="input-sm form-control" v-model="detail.shipping_address.province"/>
                                                </div>
                                            </div>
                                        </div>
                                    </div>

                                    <div class="row">
                                        <div class="col-xs-6">
                                            <div class="form-group">
                                                <label for="country" class="col-xs-4 control-label">Country</label>
                                                <div class="col-xs-8">
                                                    <input id="country" name="detail.shipping_address.country" type="text" class="input-sm form-control" v-model="detail.shipping_address.country"/>
                                                </div>
                                            </div>
                                        </div>

                                        <div class="col-xs-6">
                                            <div class="form-group">
                                                <label for="country_name" class="col-xs-4 control-label">Country
                                                    Name</label>
                                                <div class="col-xs-8">
                                                    <input id="shipping_address_country_name" name="detail.shipping_address.country_name" type="text" class="input-sm form-control" v-model="detail.shipping_address.country_name"/>
                                                </div>
                                            </div>
                                        </div>
                                    </div>

                                    <div class="row">
                                        <div class="col-xs-6">
                                            <div class="form-group">
                                                <label for="postal_code" class="col-xs-4 control-label">Post
                                                    Code</label>
                                                <div class="col-xs-8">
                                                    <input id="postal_code" name="detail.shipping_address.postal_code" type="text" class="input-sm form-control" v-model="detail.shipping_address.postal_code"/>
                                                </div>
                                            </div>
                                        </div>

                                        <div class="col-xs-6">
                                            <div class="form-group">
                                                <label for="phone" class="col-xs-4 control-label">Phone</label>
                                                <div class="col-xs-8">
                                                    <input id="shipping_address_phone" name="detail.shipping_address.phone" type="text" class="input-sm form-control" v-model="detail.shipping_address.phone"/>
                                                </div>
                                            </div>
                                        </div>
                                    </div>

                                </form>
                            </div>
                        </div>
                    </div>
                    <div class="modal-footer">
                        <button type="button" class="btn btn-success btn-sm" @click="saveOrderShippingAddress()">
                            <i class="ace-icon fa fa-floppy-o"></i>保存
                        </button>
                        <button type="button" class="btn btn-default btn-sm" data-dismiss="modal">
                            <i class="ace-icon fa fa-times"></i>关闭
                        </button>
                    </div>
                </div>
            </div>
        </div>

        <div class="modal fade" id="orderItemDetailDialog">
            <div class="modal-dialog" style="width: 680px;">
                <div class="modal-content">
                    <div class="modal-header">
                        <button type="button" class="close" data-dismiss="modal" aria-hidden="True">&times;</button>
                        <h4 class="modal-title">订单明细详细</h4>
                    </div>
                    <div class="modal-body clearfix">
                        <div class="row">
                            <div class="col-xs-12">
                                <form class="form-horizontal" role="form" id="orderItemDetailForm">
                                    <input id="order_item_id" name="order_item_detail.item_id" type="hidden" v-model.number="order_item_detail.item_id"/>
                                    <input id="order_item_order_id" name="order_item_detail.order_id" type="hidden" v-model.number="order_item_detail.order_id"/>

                                    <div class="row">
                                        <div class="col-xs-6">
                                            <div class="form-group">
                                                <label for="order_item_sku" class="col-xs-4 control-label">SKU</label>
                                                <div class="col-xs-8">
                                                    <input id="order_item_sku" name="order_item_detail.sku" type="text" class="input-sm form-control" v-model="order_item_detail.sku"/>
                                                </div>
                                            </div>
                                        </div>

                                        <div class="col-xs-6">
                                            <div class="form-group">
                                                <label for="order_item_qty" class="col-xs-4 control-label">数量</label>
                                                <div class="col-xs-8">
                                                    <input id="order_item_qty" name="order_item_detail.qty" type="text" class="input-sm form-control" v-model.number="order_item_detail.qty"/>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                    <div class="row" v-if="order_item_detail.is_order_sgcj">
                                        <div class="col-xs-6">
                                            <div class="form-group">
                                                <label for="price" class="col-xs-4 control-label">价格（AUD）</label>
                                                <div class="col-xs-8">
                                                    <input id="price" name="order_item_detail.price" type="text" class="input-sm form-control" v-model.number="order_item_detail.price"/>
                                                </div>
                                            </div>
                                        </div>
                                        <div class="col-xs-6"></div>
                                    </div>
                                </form>
                            </div>
                        </div>
                    </div>
                    <div class="modal-footer">
                        <button type="button" class="btn btn-success btn-sm" @click="save_order_item()">
                            <i class="ace-icon fa fa-floppy-o"></i>保存
                        </button>
                        <button type="button" class="btn btn-default btn-sm" data-dismiss="modal">
                            <i class="ace-icon fa fa-times"></i>关闭
                        </button>
                    </div>
                </div>
            </div>
        </div>

        <div class="modal fade" id="orderRefundDialog">
            <div class="modal-dialog" style="width: 680px;">
                <div class="modal-content">
                    <div class="modal-header">
                        <button type="button" class="close" data-dismiss="modal" aria-hidden="True">&times;</button>
                        <h4 class="modal-title">订单退款</h4>
                    </div>
                    <div class="modal-body clearfix">
                        <div class="row">
                            <div class="col-xs-12">
                                <form class="form-horizontal" role="form" id="refundForm">
                                    <input id="refund_order_id" name="refund.order_id" type="hidden" value="{% .order_id %}"/>
                                    <div class="row">
                                        <div class="col-xs-12">
                                            <div class="form-group">
                                                <label for="refund_amount" class="col-xs-4 control-label">退款金额</label>
                                                <div class="col-xs-8">
                                                    <input id="refund_amount" name="refund.amount" type="text" class="input-sm form-control" />
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                    <div class="row">
                                        <div class="col-xs-12">
                                            <div class="form-group">
                                                <label for="refund_memo" class="col-xs-4 control-label">备注</label>
                                                <div class="col-xs-8">
                                                    <textarea id="refund_memo" name="refund.memo" rows="3" class="input-sm form-control"></textarea>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                </form>
                            </div>
                        </div>
                    </div>
                    <div class="modal-footer">
                        <button type="button" class="btn btn-success btn-sm" @click="save_order_item()">
                            <i class="ace-icon fa fa-floppy-o"></i>保存
                        </button>
                        <button type="button" class="btn btn-default btn-sm" data-dismiss="modal">
                            <i class="ace-icon fa fa-times"></i>关闭
                        </button>
                    </div>
                </div>
            </div>
        </div>
    </div>
</div>
<div id="orderSplitDialogWrapper"></div>
{%template "erp/ddlc/order/order_detail.var" .%}
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
                "shipping_address": {},
                "ebay_shipping_address": {},
                "paypal_shipping_address": {},
                "order_logistics": {},
                "logistics_setting": {},
                "logistics_stock": {},
                "logistics_mode": {},
                "order_item_list": [],
                "inventory_list": [],
                "au_stock_list": [],
                "refund_list": []
            },
            "pickup": {},
            "pickup_pay_method_list": [],
            "pickup_actual_pay_method_list": [],
            "order_item_detail": {},
            "logistics_mode_list": []
        },
        mounted: function () {
            var this_vm = this;

            $.ajax(
                {
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
                        this_vm.detail = data;
                        this_vm.order_item_detail = {};
                        this_vm.order_item_detail.order_id = order_id;
                        this_vm.pickup = data.pickup;
                        this_vm.pickup_pay_method_list = data.pickup_pay_method_list;
                        this_vm.pickup_actual_pay_method_list = data.pickup_actual_pay_method_list;
                    }
                });
        },
        methods: {
            format_amount: function (amount) {
                return dtFormatFloat(amount);
            },
            show_address_edit: function () {
                $("#orderShippingAddressDetailDialog").modal(modalOptions);
            },
            saveOrderShippingAddress: function () {
                showMask("提交中...");
                var form = $("#orderShippingAddressDetailForm");
                var valid = validate(form, {
                    'detail.shipping_address.name': {
                        required: true
                    },
                    'detail.shipping_address.street1': {
                        required: true
                    },
                    'detail.shipping_address.city_name': {
                        required: true
                    },
                    'detail.shipping_address.province': {
                        required: true
                    },
                    'detail.shipping_address.country_name': {
                        required: true
                    },
                    'detail.shipping_address.postal_code': {
                        required: true
                    }
                });

                if (!valid) {
                    hideMask();
                } else {
                    $.ajax({
                            url: '{%url "ErpOrderShippingAddress.Save"%}',
                            type: "POST",
                            processData: false,
                            contentType: 'application/json',
                            data: JSON.stringify(vm.detail.shipping_address),
                            success: function (data) {
                                hideMask();
                                if (data.success === false) {
                                    $.msg.error(data.message);
                                    //return;
                                }
                                else {
                                    $.msg.show(data.message);
                                    $('#orderShippingAddressDetailDialog').modal('hide');
                                    load_order_detail();
                                    //return;
                                }
                            }
                        }
                    );
                }
            },
            show_order_waybill_number_edit: function () {
                $("#orderWaybillNumberDialog").modal(modalOptions);
            },
            saveOrderWaybillNumber: function () {
                showMask("提交中...");
                var form = $("#orderWaybillNumberForm");
                var valid = validate(form, {
                    'detail.order.id': {
                        required: true
                    },
                    'detail.order.waybill_number': {
                        required: true
                    }
                });

                if (!valid) {
                    hideMask();
                } else {
                    $.ajax({
                            url: '{%url "ErpOrderWaybillNumber.Save"%}',
                            type: "POST",
                            data: {
                                "order_id": vm.detail.order.id,
                                "waybill_number": vm.detail.order.waybill_number
                            },
                            success: function (data) {
                                hideMask();
                                if (data.success === false) {
                                    $.msg.error(data.message);
                                    //return;
                                }
                                else {
                                    $.msg.show(data.message);
                                    $('#orderWaybillNumberDialog').modal('hide');
                                    load_order_detail();
                                    //return;
                                }
                            }
                        }
                    );
                }
            },
            add_order_item: function () {
                $.ajax({
                    url: '{%url "ErpOrderItem.DetailData"%}',
                    type: "POST",
                    data: {
                        "item_id": 0,
                        "order_id": order_id
                    },
                    success: function (data) {
                        if (data.success === false) {
                            $.msg.error(data.message);
                        }
                        else {
                            vm.order_item_detail = data.order_item_detail;
                            $("#orderItemDetailDialog").modal(modalOptions);
                        }
                    }
                });
            },
            show_order_item: function () {
                var id_list = [];
                $('#order_item_list tbody tr>td:first-child input:checkbox').each(function () {
                    if (this.checked) {
                        id_list.push($(this).val());
                    }
                });

                if (id_list.length === 0 || id_list.length > 1) {
                    $.msg.alert("请选择一行数据！");
                    return;
                }

                $.ajax({
                    url: '{%url "ErpOrderItem.DetailData"%}',
                    type: "POST",
                    data: {
                        "item_id": id_list[0],
                        "order_id": order_id
                    },
                    success: function (data) {
                        if (data.success === false) {
                            $.msg.error(data.message);
                        }
                        else {
                            vm.order_item_detail = data.order_item_detail;
                            $("#orderItemDetailDialog").modal(modalOptions);
                        }
                    }
                });
            },
            save_order_item: function () {
                showMask("提交中...");
                var form = $("#orderItemDetailForm");
                var valid = validate(form, {
                    'order_item_detail.sku': {
                        required: true
                    },
                    'order_item_detail.qty': {
                        required: true, digits: true
                    }
                });

                if (!valid) {
                    hideMask();
                } else {
                    $.ajax({
                            url: '{%url "ErpOrderItem.Save"%}',
                            type: "POST",
                            processData: false,
                            contentType: 'application/json',
                            data: JSON.stringify(vm.order_item_detail),
                            success: function (data) {
                                hideMask();
                                if (data.success === false) {
                                    $.msg.error(data.message);
                                    //return;
                                }
                                else {
                                    $.msg.show(data.message);
                                    $('#orderItemDetailDialog').modal('hide');
                                    load_order_detail();
                                    //return;
                                }
                            }
                        }
                    );
                }
            },
            savePickupActual: function () {
                showMask("提交中...");
                var form = $("#pickupForm");
                var valid = validate(form, {
                    'pick.actual_is_paid': {
                        required: true, number: true
                    },
                    'pick.actual_pay_method_id': {
                        required: true, number: true
                    },
                    'pick.actual_amount': {
                        required: true, number: true
                    }
                });

                if (!valid) {
                    hideMask();
                } else {
                    $.ajax({
                            url: '{%url "ErpOrder.SavePickupActual"%}',
                            type: "POST",
                            processData: false,
                            contentType: 'application/json',
                            data: JSON.stringify(vm.pickup),
                            success: function (data) {
                                hideMask();
                                if (data.success === false) {
                                    $.msg.error(data.message);
                                    //return;
                                }
                                else {
                                    $.msg.show(data.message);
                                    load_order_detail();
                                    //return;
                                }
                            }
                        }
                    );
                }
            },
            delete_order_item: function () {
                var id_list = [];
                $('#order_item_list tbody tr>td:first-child input:checkbox').each(function () {
                    if (this.checked) {
                        id_list.push($(this).val());
                    }
                });

                if (id_list.length === 0 || id_list.length > 1) {
                    $.msg.alert("请选择一行数据！");
                    return;
                }

                $.msg.confirm("确认要删除吗？", function () {
                    $.ajax({
                        url: '{%url "ErpOrderItem.Delete"%}',
                        type: "POST",
                        data: {"item_id": id_list[0]},
                        success: function (data) {
                            if (data.success === false) {
                                $.msg.error(data.message);
                                return;
                            }
                            else {
                                $.msg.show(data.message);
                            }
                            load_order_detail();
                        }
                    });
                });
            }
        }
    });

    function reship_order() {
        $.ajax({
            "dataType": 'json',
            "type": "POST",
            "url": '{% url "ErpOrder.ReShip"%}',
            "data": {
                "order_id": order_id
            },
            "success": function (data) {
                if (data.success === false) {
                    $.msg.error(data.message);
                    return;
                }

                setTimeout("window.location.href = '/order/detail/" + data.new_order_id + "'", 1000);
                $.msg.show(data.message)
            }
        });
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

    function rematch_order_logistics() {
        var order_list = [];
        order_list.push(order_id);

        $.ajax({
            "dataType": 'json',
            "url": '{% url "ErpOrder.ReMatchHoldOrderLogistics"%}',
            "data": {
                "order_id_list": order_list,
                "prefer_stock_id": $("#prefer_stock").val(),
                "prefer_mode_id": $("#prefer_mode").val(),
                "unlock_express": $("#unlock_express").val()
            },
            "success": function (data) {
                if (data.success === false) {
                    $.msg.error(data.message);
                }
                else {
                    setTimeout("window.location.reload();", 1000);
                    $.msg.show(data.message)
                }
            }
        });
    }

    function save_order_desc() {
        $.ajax({
            "dataType": 'json',
            "url": '{% url "ErpOrder.SaveDesc" %}',
            "data": {
                "order_id": order_id,
                "order_desc": $("#order_desc").val()
            },
            "success": function (data) {
                if (data.success === false) {
                    $.msg.error(data.message);
                }
                else {
                    $.msg.show(data.message);
                }
            }
        });
    }

    function load_logistics_model_list() {
        $.ajax({
            "dataType": 'json',
            "url": '{% url "ErpLogisticsMode.GetListByStock" %}',
            "data": {
                "stock_id": $("#prefer_stock").val()
            },
            "success": function (data) {
                if (data.success === false) {
                    $.msg.error(data.message);
                }
                else {
                    vm.logistics_mode_list = data.mode_list;
                }
            }
        });
    }

    $(function () {
        initDataTable(
            '#orderHistoryList',
            {
                "sAjaxSource": '{%url "ErpOrderHistory.ListData"%}',
                "sDom": "<'row'<'col-xs-12't>>",
                "aaSorting": [[3, "asc"]],
                "aoColumns": [
                    {"mDataProp": "order_code", "title": "订单编码", "bSearchable": true, "bSortable": false, "width": 120},
                    {"mDataProp": "status_name", "title": "状态名称", "bSearchable": true, "bSortable": false, "width": 120},
                    {"mDataProp": "description", "title": "描述", "bSearchable": true, "bSortable": false},
                    {"mDataProp": "create_time", "title": "操作时间", "bSearchable": false, "bSortable": false, "width": 120}
                ],
                "fnServerData": function (sSource, aoData, fnCallback) {
                    showMask("数据加载中，请稍候");

                    aoData.push({'name': 'order_id', 'value': order_id});

                    $.ajax({
                        "dataType": 'json',
                        "url": sSource,
                        "data": aoData,
                        "success": function (data) {
                            hideMask();

                            if (data.success === false) {
                                $.msg.error(data.message);
                                return;
                            }

                            fnCallback(data);
                        }
                    });
                }
            }
        );
    });

    function showOrderSplit() {
        $.ajax({
            url: '{%url "ErpOrderSplit.Show" %}',
            type: "POST",
            data: {
                "order_id": vm.detail.order.id
            },
            success: function (result) {
                if (result.success === false) {
                    $.msg.error(result.message);
                    return;
                }
                $("#orderSplitDialogWrapper").html(result);
                $('#orderSplitDialog').modal(modalOptions);
            }
        });
    }

    function saveOrderSplit() {
        var form = $("#orderSplitForm");

        showMask("数据加载中，请稍候");
        $.ajax({
            "dataType": 'json',
            "url": '{% url "ErpOrderSplit.Save"%}',
            "data": form.serialize(),
            "success": function (data) {
                hideMask();
                if (data.success === false) {
                    $.msg.error(data.message);
                    return;
                }

                $('#orderSplitDialog').modal("hide");
                load_order_detail();
                $.msg.show(data.message);
            }
        });
    }

    function showRefund() {
        $("#refundDialog").modal(modalOptions);
    }

    function refund () {
        showMask("提交中...");
        var form = $("#refundForm");
        var valid = validate(form, {
            'refund.order_id': {
                required: true, number: true
            },
            'refund.amount': {
                required: true, number: true
            },
            'refund.memo': {
                required: false
            }
        });

        if (!valid) {
            hideMask();
        } else {
            $.ajax({
                    url: '{%url "ErpOrderRefund.Refund"%}',
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

    function loadRefundList() {
        $.ajax({
            "dataType": 'json',
            "url": '{% url "ErpOrderRefund.ListData" %}',
            "data": {
                "order_id": order_id
            },
            "success": function (data) {
                if (data.success === false) {
                    $.msg.error(data.message);
                }
                else {
                    vm.detail.refund_list = data.refund_list;
                }
            }
        });
    }

    $(function(){
        loadRefundList();
    });
</script>
{% site_footer %}

