// External library dependencies
import type { useAxiosReturn } from '@baota/hooks/axios'
import type { AxiosResponseData } from '@/types/public'

// Type imports
import type {
	Plan,
	Order,
	PaymentConfig,
	CreateOrderParams,
	CreateOrderResponse,
	OrderListParams,
} from '@/types/payment'

// Relative internal imports
import { useApi } from '@api/index'

/**
 * @description 创建支付订单
 */
export const createOrder = (params?: CreateOrderParams): useAxiosReturn<AxiosResponseData<CreateOrderResponse>, CreateOrderParams> =>
	useApi<AxiosResponseData<CreateOrderResponse>, CreateOrderParams>('/v1/payment/create_order', params)

/**
 * @description 查询订单状态
 */
export const getOrderStatus = (params?: { order_no: string }): useAxiosReturn<AxiosResponseData<Order>, { order_no: string }> =>
	useApi<AxiosResponseData<Order>, { order_no: string }>('/v1/payment/get_order_status', params)

/**
 * @description 获取订单列表
 */
export const getOrderList = (params?: OrderListParams): useAxiosReturn<AxiosResponseData<Order[]>, OrderListParams> =>
	useApi<AxiosResponseData<Order[]>, OrderListParams>('/v1/payment/get_order_list', params)

/**
 * @description 获取套餐列表
 */
export const getPlanList = (): useAxiosReturn<AxiosResponseData<Plan[]>, Record<string, unknown>> =>
	useApi<AxiosResponseData<Plan[]>, Record<string, unknown>>('/v1/payment/get_plan_list')

/**
 * @description 保存套餐
 */
export const savePlan = (params?: Partial<Plan>): useAxiosReturn<AxiosResponseData<null>, Partial<Plan>> =>
	useApi<AxiosResponseData<null>, Partial<Plan>>('/v1/payment/save_plan', params)

/**
 * @description 删除套餐
 */
export const delPlan = (params?: { id: string }): useAxiosReturn<AxiosResponseData<null>, { id: string }> =>
	useApi<AxiosResponseData<null>, { id: string }>('/v1/payment/del_plan', params)

/**
 * @description 获取支付配置
 */
export const getPaymentConfig = (): useAxiosReturn<AxiosResponseData<PaymentConfig>, Record<string, unknown>> =>
	useApi<AxiosResponseData<PaymentConfig>, Record<string, unknown>>('/v1/payment/get_payment_config')

/**
 * @description 保存支付配置
 */
export const savePaymentConfig = (params?: Partial<PaymentConfig>): useAxiosReturn<AxiosResponseData<null>, Partial<PaymentConfig>> =>
	useApi<AxiosResponseData<null>, Partial<PaymentConfig>>('/v1/payment/save_payment_config', params)
