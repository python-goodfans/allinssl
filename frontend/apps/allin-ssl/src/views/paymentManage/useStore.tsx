/**
 * @file 支付管理状态管理
 */
import { useMessage } from '@baota/naive-ui/hooks'
import { useError } from '@baota/hooks/error'
import type { Order, PaymentConfig, Plan } from '@/types/payment'
import { delPlan, getOrderList, getPaymentConfig, getPlanList, savePaymentConfig, savePlan } from '@api/payment'

const { success } = useMessage()
const { handleError } = useError()

interface PaymentManageStoreExposes {
	loading: Ref<boolean>
	orders: Ref<Order[]>
	orderTotal: Ref<number>
	orderPage: Ref<number>
	plans: Ref<Plan[]>
	paymentConfig: Ref<PaymentConfig>
	activeTab: Ref<string>
	fetchOrders: (page?: number) => Promise<void>
	fetchPlans: () => Promise<void>
	fetchPaymentConfig: () => Promise<void>
	handleSavePlan: (plan: Partial<Plan>) => Promise<void>
	handleDelPlan: (id: string) => Promise<void>
	handleSavePaymentConfig: () => Promise<void>
}

const defaultConfig: PaymentConfig = {
	wechat_app_id: '',
	wechat_mch_id: '',
	wechat_api_key: '',
	wechat_notify_url: '',
	alipay_app_id: '',
	alipay_private_key: '',
	alipay_public_key: '',
	alipay_notify_url: '',
	alipay_sandbox: false,
}

export const usePaymentManageStore = defineStore('payment-manage-store', (): PaymentManageStoreExposes => {
	const loading = ref(false)
	const orders = ref<Order[]>([])
	const orderTotal = ref(0)
	const orderPage = ref(1)
	const plans = ref<Plan[]>([])
	const paymentConfig = ref<PaymentConfig>({ ...defaultConfig })
	const activeTab = ref('orders')

	const fetchOrders = async (page = 1): Promise<void> => {
		loading.value = true
		try {
			const { fetch, data } = getOrderList()
			await fetch({ p: page, limit: 20 })
			if (data.value?.status) {
				orders.value = (data.value.data as unknown as Order[]) || []
				orderTotal.value = data.value.count || 0
				orderPage.value = page
			}
		} catch (err) {
			handleError(err)
		} finally {
			loading.value = false
		}
	}

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

	const fetchPaymentConfig = async (): Promise<void> => {
		try {
			const { fetch, data } = getPaymentConfig()
			await fetch({})
			if (data.value?.status) {
				paymentConfig.value = { ...defaultConfig, ...(data.value.data as unknown as PaymentConfig) }
			}
		} catch (err) {
			handleError(err)
		}
	}

	const handleSavePlan = async (plan: Partial<Plan>): Promise<void> => {
		loading.value = true
		try {
			const { fetch, data } = savePlan()
			await fetch(plan)
			if (data.value?.status) {
				success('保存成功')
				await fetchPlans()
			} else {
				throw new Error(data.value?.message || '保存失败')
			}
		} catch (err) {
			handleError(err)
		} finally {
			loading.value = false
		}
	}

	const handleDelPlan = async (id: string): Promise<void> => {
		loading.value = true
		try {
			const { fetch, data } = delPlan()
			await fetch({ id })
			if (data.value?.status) {
				success('删除成功')
				await fetchPlans()
			} else {
				throw new Error(data.value?.message || '删除失败')
			}
		} catch (err) {
			handleError(err)
		} finally {
			loading.value = false
		}
	}

	const handleSavePaymentConfig = async (): Promise<void> => {
		loading.value = true
		try {
			const { fetch, data } = savePaymentConfig()
			await fetch(paymentConfig.value)
			if (data.value?.status) {
				success('配置保存成功')
			} else {
				throw new Error(data.value?.message || '保存失败')
			}
		} catch (err) {
			handleError(err)
		} finally {
			loading.value = false
		}
	}

	return {
		loading,
		orders,
		orderTotal,
		orderPage,
		plans,
		paymentConfig,
		activeTab,
		fetchOrders,
		fetchPlans,
		fetchPaymentConfig,
		handleSavePlan,
		handleDelPlan,
		handleSavePaymentConfig,
	}
})

export const useStore = () => {
	const store = usePaymentManageStore()
	return { ...store, ...storeToRefs(store) }
}
