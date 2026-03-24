import { useStore } from './useStore'
import type { Plan } from '@/types/payment'

export const useController = () => {
	const store = useStore()
	const { fetchOrders, fetchPlans, fetchPaymentConfig, handleSavePlan, handleDelPlan } = store

	const showPlanModal = ref(false)
	const editingPlan = ref<Partial<Plan>>({
		name: '',
		description: '',
		price: 0,
		duration: 30,
		features: '',
		status: 1,
		sort_order: 0,
	})

	const handleOpenAddPlan = (): void => {
		editingPlan.value = {
			name: '',
			description: '',
			price: 0,
			duration: 30,
			features: '',
			status: 1,
			sort_order: 0,
		}
		showPlanModal.value = true
	}

	const handleOpenEditPlan = (plan: Plan): void => {
		editingPlan.value = { ...plan }
		showPlanModal.value = true
	}

	const handleSavePlanSubmit = async (): Promise<void> => {
		await handleSavePlan(editingPlan.value)
		showPlanModal.value = false
	}

	const handleConfirmDelPlan = async (id: string): Promise<void> => {
		await handleDelPlan(id)
	}

	onMounted(() => {
		fetchOrders()
		fetchPlans()
		fetchPaymentConfig()
	})

	return {
		...store,
		showPlanModal,
		editingPlan,
		handleOpenAddPlan,
		handleOpenEditPlan,
		handleSavePlanSubmit,
		handleConfirmDelPlan,
	}
}
