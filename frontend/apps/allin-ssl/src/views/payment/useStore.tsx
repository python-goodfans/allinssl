/**
 * @file 支付状态管理
 */
import { useMessage } from '@baota/naive-ui/hooks'
import { useError } from '@baota/hooks/error'
import type { Order, Plan } from '@/types/payment'
import { createOrder, getOrderStatus, getPlanList } from '@api/payment'
import { ShallowRef } from 'vue'

const { success, error: msgError } = useMessage()
const { handleError } = useError()

interface PaymentStoreExposes {
	loading: Ref<boolean>
	plans: Ref<Plan[]>
	selectedPlan: Ref<Plan | null>
	currentOrder: Ref<Order | null>
	payUrl: Ref<string>
	paymentMethod: Ref<'wechat' | 'alipay'>
	orderStatus: Ref<string>
	error: ShallowRef<string | null>
	fetchPlans: () => Promise<void>
	selectPlan: (plan: Plan) => void
	handleCreateOrder: () => Promise<void>
	pollOrderStatus: () => void
	stopPolling: () => void
}

export const usePaymentStore = defineStore('payment-store', (): PaymentStoreExposes => {
	const plans = ref<Plan[]>([])
	const selectedPlan = ref<Plan | null>(null)
	const currentOrder = ref<Order | null>(null)
	const payUrl = ref('')
	const paymentMethod = ref<'wechat' | 'alipay'>('wechat')
	const orderStatus = ref('')
	const loading = ref(false)
	const error = shallowRef<string | null>(null)
	let pollingTimer: ReturnType<typeof setInterval> | null = null
	let pollingTimeout: ReturnType<typeof setTimeout> | null = null

	const fetchPlans = async (): Promise<void> => {
		try {
			const { fetch, data } = getPlanList()
			await fetch({})
			if (data.value?.status) {
				plans.value = (data.value.data as unknown as Plan[]) || []
			}
		} catch (err) {
			handleError(err)
		}
	}

	const selectPlan = (plan: Plan): void => {
		selectedPlan.value = plan
	}

	const handleCreateOrder = async (): Promise<void> => {
		if (!selectedPlan.value) {
			error.value = '请选择套餐'
			return
		}
		loading.value = true
		error.value = null
		try {
			const { fetch, data } = createOrder()
			await fetch({ plan_id: selectedPlan.value.id, payment_method: paymentMethod.value })
			if (data.value?.status) {
				const orderData = data.value.data as unknown as { order_no: string; pay_url: string; amount: number }
				currentOrder.value = {
					order_no: orderData.order_no,
					payment_method: paymentMethod.value,
					amount: orderData.amount,
					status: 'pending',
				} as Order
				payUrl.value = orderData.pay_url
				orderStatus.value = 'pending'
				if (paymentMethod.value === 'wechat') {
					pollOrderStatus()
				}
			} else {
				throw new Error(data.value?.message || '创建订单失败')
			}
		} catch (err: unknown) {
			error.value = (err as Error).message
		} finally {
			loading.value = false
		}
	}

	const pollOrderStatus = (): void => {
		if (!currentOrder.value) return

		pollingTimer = setInterval(async () => {
			try {
				const { fetch, data } = getOrderStatus()
				await fetch({ order_no: currentOrder.value!.order_no })
				if (data.value?.status) {
					const order = data.value.data as unknown as Order
					orderStatus.value = order.status
					if (order.status === 'paid') {
						stopPolling()
						success('支付成功！')
					} else if (order.status === 'failed') {
						stopPolling()
						msgError('支付失败')
					}
				}
			} catch (err) {
				handleError(err)
			}
		}, 3000)

		// Auto stop polling after 5 minutes
		pollingTimeout = setTimeout(() => {
			stopPolling()
		}, 5 * 60 * 1000)
	}

	const stopPolling = (): void => {
		if (pollingTimer) {
			clearInterval(pollingTimer)
			pollingTimer = null
		}
		if (pollingTimeout) {
			clearTimeout(pollingTimeout)
			pollingTimeout = null
		}
	}

	return {
		loading,
		plans,
		selectedPlan,
		currentOrder,
		payUrl,
		paymentMethod,
		orderStatus,
		error,
		fetchPlans,
		selectPlan,
		handleCreateOrder,
		pollOrderStatus,
		stopPolling,
	}
})

export const useStore = () => {
	const store = usePaymentStore()
	return { ...store, ...storeToRefs(store) }
}
