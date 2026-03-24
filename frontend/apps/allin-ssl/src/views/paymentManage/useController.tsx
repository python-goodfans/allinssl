/**
 * @file 支付管理控制器
 */
import { useStore } from './useStore'

export const useController = () => {
	const store = useStore()
	const {
		activeTab,
		fetchOrders,
		fetchPlans,
		fetchPaymentConfig,
		showPlanModal,
		editingPlan,
	} = store

	/**
	 * 切换标签页
	 */
	const handleTabChange = (tab: 'orders' | 'plans' | 'config'): void => {
		activeTab.value = tab
		if (tab === 'orders') fetchOrders()
		else if (tab === 'plans') fetchPlans()
		else if (tab === 'config') fetchPaymentConfig()
	}

	/**
	 * 打开新增套餐弹窗
	 */
	const handleAddPlan = (): void => {
		editingPlan.value = { name: '', description: '', price: 0, duration: 30, features: '', status: 1, sort_order: 0 }
		showPlanModal.value = true
	}

	/**
	 * 打开编辑套餐弹窗
	 */
	const handleEditPlan = (plan: Record<string, unknown>): void => {
		editingPlan.value = { ...plan } as typeof editingPlan.value
		showPlanModal.value = true
	}

	onMounted(() => {
		fetchOrders()
	})

	return {
		...store,
		handleTabChange,
		handleAddPlan,
		handleEditPlan,
	}
}
