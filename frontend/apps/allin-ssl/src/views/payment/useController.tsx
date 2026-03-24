/**
 * @file 支付控制器
 */
import { useStore } from './useStore'

export const useController = () => {
	const store = useStore()
	const {
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
	} = store

	const qrCodeDataUrl = ref('')

	/**
	 * 处理微信支付
	 */
	const handleWechatPay = async (): Promise<void> => {
		const order = await handleCreateOrder('wechat')
		if (!order) return

		if (order.code_url) {
			// 生成二维码（使用简单的占位方式展示code_url文本，实际可集成qrcode库）
			qrCodeDataUrl.value = order.code_url
			showQrModal.value = true
			startPollOrderStatus(order.order_no)
		}
	}

	/**
	 * 处理支付宝支付
	 */
	const handleAlipay = async (): Promise<void> => {
		const order = await handleCreateOrder('alipay')
		if (!order) return

		if (order.pay_url) {
			startPollOrderStatus(order.order_no)
			window.open(order.pay_url, '_blank')
		}
	}

	/**
	 * 关闭二维码弹窗
	 */
	const handleCloseQrModal = (): void => {
		showQrModal.value = false
		stopPollOrderStatus()
	}

	onMounted(() => {
		fetchPlans()
	})

	onUnmounted(() => {
		stopPollOrderStatus()
	})

	return {
		...store,
		qrCodeDataUrl,
		handleWechatPay,
		handleAlipay,
		handleCloseQrModal,
	}
}
