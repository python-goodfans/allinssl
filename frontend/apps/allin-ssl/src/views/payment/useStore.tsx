/**
 * @file 支付状态管理
 */
// External Libraries
import { useMessage } from '@baota/naive-ui/hooks'
import { useError } from '@baota/hooks/error'

// Type Imports
import type { Plan, Order, CreateOrderResponse } from '@/types/payment'

// API Imports
import { getPlanList, createOrder, getOrderStatus } from '@api/payment'

const { success, error: msgError } = useMessage()
const { handleError } = useError()

export const usePaymentStore = defineStore('payment-store', () => {
	const loading = ref(false)
	const plans = ref<Plan[]>([])
	const selectedPlan = ref<Plan | null>(null)
	const currentOrder = ref<CreateOrderResponse | null>(null)
	const orderStatus = ref<string>('')
	const showQrModal = ref(false)
	let pollTimer: ReturnType<typeof setInterval> | null = null

	/**
	 * 获取套餐列表
	 */
	const fetchPlans = async (): Promise<void> => {
		try {
			loading.value = true
			const { fetch, data } = getPlanList()
			await fetch()
			if (data.value?.status) {
				plans.value = data.value.data || []
			}
		} catch (err) {
			handleError(err)
		} finally {
			loading.value = false
		}
	}

	/**
	 * 选择套餐
	 */
	const selectPlan = (plan: Plan): void => {
		selectedPlan.value = plan
	}

	/**
	 * 创建支付订单
	 */
	const handleCreateOrder = async (paymentMethod: 'wechat' | 'alipay'): Promise<CreateOrderResponse | null> => {
		if (!selectedPlan.value) {
			msgError('请先选择套餐')
			return null
		}
		try {
			loading.value = true
			const { fetch, data } = createOrder({ plan_id: selectedPlan.value.id, payment_method: paymentMethod })
			await fetch()
			if (data.value?.status) {
				currentOrder.value = data.value.data as unknown as CreateOrderResponse
				orderStatus.value = 'pending'
				return currentOrder.value
			} else {
				throw new Error(data.value?.message || '创建订单失败')
			}
		} catch (err) {
			handleError(err)
			return null
		} finally {
			loading.value = false
		}
	}

	/**
	 * 开始轮询订单状态（每3秒一次，最多5分钟）
	 */
	const startPollOrderStatus = (orderNo: string): void => {
		stopPollOrderStatus()
		let elapsed = 0
		const maxTime = 5 * 60 * 1000 // 5分钟
		const interval = 3000

		pollTimer = setInterval(async () => {
			elapsed += interval
			try {
				const { fetch, data } = getOrderStatus({ order_no: orderNo })
				await fetch()
				const order = data.value?.data as unknown as Order
				if (order) {
					orderStatus.value = order.status
					if (order.status === 'paid') {
						success('支付成功！')
						stopPollOrderStatus()
						showQrModal.value = false
						setTimeout(() => (location.href = '/'), 2000)
					} else if (order.status === 'failed') {
						msgError('支付失败')
						stopPollOrderStatus()
					}
				}
			} catch (err) {
				// 忽略轮询错误
			}

			if (elapsed >= maxTime) {
				stopPollOrderStatus()
			}
		}, interval)
	}

	/**
	 * 停止轮询
	 */
	const stopPollOrderStatus = (): void => {
		if (pollTimer !== null) {
			clearInterval(pollTimer)
			pollTimer = null
		}
	}

	return {
		loading,
		plans,
		selectedPlan,
		currentOrder,
		orderStatus,
		showQrModal,
		fetchPlans,
		selectPlan,
		handleCreateOrder,
		startPollOrderStatus,
		stopPollOrderStatus,
	}
})

export const useStore = () => {
	const store = usePaymentStore()
	return { ...store, ...storeToRefs(store) }
}
