import { useStore } from './useStore'
import type { Plan } from '@/types/payment'

export const useController = () => {
	const store = useStore()
	const {
		plans,
		selectedPlan,
		currentOrder,
		payUrl,
		paymentMethod,
		orderStatus,
		loading,
		error,
		fetchPlans,
		selectPlan,
		handleCreateOrder,
		stopPolling,
	} = store

	const showQrCode = ref(false)
	const showAlipayRedirect = ref(false)

	const handleSelectPlan = (plan: Plan): void => {
		selectPlan(plan)
	}

	const handlePay = async (): Promise<void> => {
		await handleCreateOrder()
		if (paymentMethod.value === 'wechat' && payUrl.value) {
			showQrCode.value = true
		} else if (paymentMethod.value === 'alipay' && payUrl.value) {
			showAlipayRedirect.value = true
			window.open(payUrl.value, '_blank')
		}
	}

	const handleCloseQrCode = (): void => {
		showQrCode.value = false
		stopPolling()
	}

	onMounted(() => {
		fetchPlans()
	})

	onUnmounted(() => {
		stopPolling()
	})

	return {
		...store,
		showQrCode,
		showAlipayRedirect,
		handleSelectPlan,
		handlePay,
		handleCloseQrCode,
	}
}
