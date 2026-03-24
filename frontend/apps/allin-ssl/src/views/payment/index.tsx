// External Libraries
import { NButton, NCard, NModal, NSpin, NTag, NDivider } from 'naive-ui'

// Absolute Internal Imports
import { useThemeCssVar } from '@baota/naive-ui/theme'

// Relative Internal Imports
import { useController } from './useController'

// Side-effect Imports
import styles from './index.module.css'

export default defineComponent({
	name: 'PaymentView',
	setup() {
		const {
			loading,
			plans,
			selectedPlan,
			currentOrder,
			orderStatus,
			showQrModal,
			qrCodeDataUrl,
			selectPlan,
			handleWechatPay,
			handleAlipay,
			handleCloseQrModal,
		} = useController()

		const cssVar = useThemeCssVar(['textColor2', 'actionColor', 'primaryColor', 'primaryColorSuppl'])

		return () => (
			<div style={cssVar.value} class={styles.container}>
				<div class={styles.header}>
					<h1 class={styles.title}>选择套餐</h1>
					<p class={styles.subtitle}>选择适合您的套餐，享受 ALLinSSL 全功能服务</p>
				</div>

				<NSpin show={loading.value}>
					<div class={styles.planGrid}>
						{plans.value.map((plan) => (
							<NCard
								key={plan.id}
								class={[styles.planCard, selectedPlan.value?.id === plan.id ? styles.selectedCard : '']}
								onClick={() => selectPlan(plan)}
							>
								<div class={styles.planHeader}>
									<h2 class={styles.planName}>{plan.name}</h2>
									<div class={styles.planPrice}>
										<span class={styles.priceSymbol}>¥</span>
										<span class={styles.priceValue}>{plan.price.toFixed(2)}</span>
									</div>
									<p class={styles.planDuration}>{plan.duration} 天有效期</p>
								</div>
								<NDivider />
								<p class={styles.planDesc}>{plan.description}</p>
								{selectedPlan.value?.id === plan.id && (
									<NTag type="success" class={styles.selectedTag}>
										已选择
									</NTag>
								)}
							</NCard>
						))}
					</div>
				</NSpin>

				{selectedPlan.value && (
					<div class={styles.paySection}>
						<h2 class={styles.payTitle}>
							选择支付方式 - {selectedPlan.value.name}（¥{selectedPlan.value.price.toFixed(2)}）
						</h2>
						<div class={styles.payButtons}>
							<NButton
								type="success"
								size="large"
								loading={loading.value}
								onClick={handleWechatPay}
								class={styles.payBtn}
							>
								微信支付
							</NButton>
							<NButton
								type="info"
								size="large"
								loading={loading.value}
								onClick={handleAlipay}
								class={styles.payBtn}
							>
								支付宝支付
							</NButton>
						</div>
					</div>
				)}

				{/* 微信支付二维码弹窗 */}
				<NModal
					show={showQrModal.value}
					onUpdateShow={(v) => { if (!v) handleCloseQrModal() }}
					title="微信扫码支付"
					preset="card"
					style="width: 400px"
				>
					<div class={styles.qrContainer}>
						{orderStatus.value === 'pending' && (
							<>
								<p class={styles.qrTip}>请使用微信扫描以下二维码完成支付</p>
								<div class={styles.qrCode}>
									<p class={styles.qrCodeUrl}>{qrCodeDataUrl.value}</p>
									<p class={styles.qrCodeHint}>（请复制链接在微信中打开，或使用二维码扫描工具）</p>
								</div>
								<div class={styles.orderInfo}>
									<p>订单号: {currentOrder.value?.order_no}</p>
									<p>金额: ¥{currentOrder.value?.amount?.toFixed(2)}</p>
								</div>
								<NSpin size="small" />
								<p class={styles.polling}>正在等待支付结果...</p>
							</>
						)}
						{orderStatus.value === 'paid' && (
							<div class={styles.success}>✅ 支付成功！正在跳转...</div>
						)}
						{orderStatus.value === 'failed' && (
							<div class={styles.failed}>❌ 支付失败，请重试</div>
						)}
					</div>
					{{
						footer: () => (
							<NButton onClick={handleCloseQrModal}>关闭</NButton>
						),
					}}
				</NModal>
			</div>
		)
	},
})
