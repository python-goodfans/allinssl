export interface Plan {
	id: string
	name: string
	description: string
	price: number
	duration: number
	features: string
	status: number
	sort_order: number
	create_time: string
}

export interface Order {
	id: string
	order_no: string
	user_id: string
	plan_id: string
	amount: number
	payment_method: 'wechat' | 'alipay'
	status: 'pending' | 'paid' | 'failed' | 'refunded'
	trade_no: string
	create_time: string
	update_time: string
	pay_time: string
}

export interface PaymentConfig {
	wechat_app_id: string
	wechat_mch_id: string
	wechat_api_key: string
	wechat_notify_url: string
	alipay_app_id: string
	alipay_private_key: string
	alipay_public_key: string
	alipay_notify_url: string
	alipay_sandbox: boolean
}

export interface CreateOrderParams {
	plan_id: string
	payment_method: 'wechat' | 'alipay'
}

export interface OrderListParams {
	p: number
	limit: number
	user_id?: string
	status?: string
	payment_method?: string
}

export interface CreateOrderResponse {
	order_no: string
	pay_url: string
	payment_method: 'wechat' | 'alipay'
	amount: number
}
