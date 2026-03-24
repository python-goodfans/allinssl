// Type imports
import type { useAxiosReturn } from '@baota/hooks/axios'
import type {
	CreateOrderParams,
	CreateOrderResponse,
	Order,
	OrderListParams,
	PaymentConfig,
	Plan,
} from '@/types/payment'
import type { AxiosResponseData } from '@/types/public'

// Relative internal imports
import { useApi } from '@api/index'

/**
 * @description 创建支付订单
 */
export const createOrder = (params?: CreateOrderParams): useAxiosReturn<CreateOrderResponse, CreateOrderParams> =>
	useApi<CreateOrderResponse, CreateOrderParams>('/v1/payment/create_order', params)

/**
 * @description 获取订单状态
 */
export const getOrderStatus = (params?: { order_no: string }): useAxiosReturn<Order, { order_no: string }> =>
	useApi<Order, { order_no: string }>('/v1/payment/get_order_status', params)

/**
 * @description 获取订单列表
 */
export const getOrderList = (params?: OrderListParams): useAxiosReturn<Order[], OrderListParams> =>
	useApi<Order[], OrderListParams>('/v1/payment/get_order_list', params)

/**
 * @description 获取套餐列表
 */
export const getPlanList = (): useAxiosReturn<Plan[], Record<string, unknown>> =>
	useApi<Plan[], Record<string, unknown>>('/v1/payment/get_plan_list')

/**
 * @description 保存套餐
 */
export const savePlan = (params?: Partial<Plan>): useAxiosReturn<AxiosResponseData, Partial<Plan>> =>
	useApi<AxiosResponseData, Partial<Plan>>('/v1/payment/save_plan', params)

/**
 * @description 删除套餐
 */
export const delPlan = (params?: { id: string }): useAxiosReturn<AxiosResponseData, { id: string }> =>
	useApi<AxiosResponseData, { id: string }>('/v1/payment/del_plan', params)

/**
 * @description 获取支付配置
 */
export const getPaymentConfig = (): useAxiosReturn<PaymentConfig, Record<string, unknown>> =>
	useApi<PaymentConfig, Record<string, unknown>>('/v1/payment/get_payment_config')

/**
 * @description 保存支付配置
 */
export const savePaymentConfig = (
	params?: Partial<PaymentConfig>
): useAxiosReturn<AxiosResponseData, Partial<PaymentConfig>> =>
	useApi<AxiosResponseData, Partial<PaymentConfig>>('/v1/payment/save_payment_config', params)
