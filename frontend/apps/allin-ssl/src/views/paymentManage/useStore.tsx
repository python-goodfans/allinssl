/**
 * @file 支付管理状态管理
 */
// External Libraries
import { useMessage } from '@baota/naive-ui/hooks'
import { useError } from '@baota/hooks/error'

// Type Imports
import type { Plan, Order, PaymentConfig } from '@/types/payment'

// API Imports
import { getPlanList, getOrderList, savePlan, delPlan, getPaymentConfig, savePaymentConfig } from '@api/payment'

const { success, error: msgError } = useMessage()
const { handleError } = useError()

export const usePaymentManageStore = defineStore('payment-manage-store', () => {
	const loading = ref(false)
	const activeTab = ref<'orders' | 'plans' | 'config'>('orders')

	// 订单相关
	const orders = ref<Order[]>([])
	const orderTotal = ref(0)
	const orderPage = ref(1)
	const orderPageSize = ref(20)

	// 套餐相关
	const plans = ref<Plan[]>([])
	const showPlanModal = ref(false)
	const editingPlan = ref<Partial<Plan>>({})

	// 支付配置
	const paymentConfig = ref<Partial<PaymentConfig>>({})

	/**
	 * 获取订单列表
	 */
	const fetchOrders = async (): Promise<void> => {
		try {
			loading.value = true
			const { fetch, data } = getOrderList({ p: orderPage.value, limit: orderPageSize.value })
			await fetch()
			if (data.value?.status) {
				orders.value = (data.value.data as unknown as Order[]) || []
				orderTotal.value = data.value.count || 0
			}
		} catch (err) {
			handleError(err)
		} finally {
			loading.value = false
		}
	}

	/**
	 * 获取套餐列表
	 */
	const fetchPlans = async (): Promise<void> => {
		try {
			loading.value = true
			const { fetch, data } = getPlanList()
			await fetch()
			if (data.value?.status) {
				plans.value = (data.value.data as unknown as Plan[]) || []
			}
		} catch (err) {
			handleError(err)
		} finally {
			loading.value = false
		}
	}

	/**
	 * 保存套餐
	 */
	const handleSavePlan = async (): Promise<void> => {
		try {
			loading.value = true
			const { fetch, data } = savePlan(editingPlan.value)
			await fetch()
			if (data.value?.status) {
				success('保存成功')
				showPlanModal.value = false
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

	/**
	 * 删除套餐
	 */
	const handleDelPlan = async (id: string): Promise<void> => {
		try {
			loading.value = true
			const { fetch, data } = delPlan({ id })
			await fetch()
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

	/**
	 * 获取支付配置
	 */
	const fetchPaymentConfig = async (): Promise<void> => {
		try {
			loading.value = true
			const { fetch, data } = getPaymentConfig()
			await fetch()
			if (data.value?.status) {
				paymentConfig.value = (data.value.data as unknown as PaymentConfig) || {}
			}
		} catch (err) {
			handleError(err)
		} finally {
			loading.value = false
		}
	}

	/**
	 * 保存支付配置
	 */
	const handleSaveConfig = async (): Promise<void> => {
		try {
			loading.value = true
			const { fetch, data } = savePaymentConfig(paymentConfig.value)
			await fetch()
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
		activeTab,
		orders,
		orderTotal,
		orderPage,
		orderPageSize,
		plans,
		showPlanModal,
		editingPlan,
		paymentConfig,
		fetchOrders,
		fetchPlans,
		handleSavePlan,
		handleDelPlan,
		fetchPaymentConfig,
		handleSaveConfig,
	}
})

export const useStore = () => {
	const store = usePaymentManageStore()
	return { ...store, ...storeToRefs(store) }
}
