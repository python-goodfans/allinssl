import { NButton, NCard, NModal, NRadioGroup, NRadioButton, NSpace, NSpin, NTag } from 'naive-ui'
import { useThemeCssVar } from '@baota/naive-ui/theme'
import { useController } from './useController'
import styles from './index.module.css'

export default defineComponent({
	name: 'PaymentView',
	setup() {
		const {
			plans,
			selectedPlan,
			paymentMethod,
			orderStatus,
			payUrl,
			loading,
			error,
			showQrCode,
			handleSelectPlan,
			handlePay,
			handleCloseQrCode,
		} = useController()

		const cssVar = useThemeCssVar(['textColor2', 'actionColor', 'primaryColor'])

		return () => (
			<div style={cssVar.value} class={styles.container}>
				<div class={styles.header}>
					<h1 class={styles.title}>选择套餐</h1>
					<p class={styles.subtitle}>选择适合您的套餐，享受更多功能</p>
				</div>

				<div class={styles.planGrid}>
					{plans.value.map((plan) => (
						<NCard
							key={plan.id}
							class={[styles.planCard, selectedPlan.value?.id === plan.id ? styles.selected : '']}
							onClick={() => handleSelectPlan(plan)}
						>
							<div class={styles.planName}>{plan.name}</div>
							<div class={styles.planPrice}>
								<span class={styles.price}>¥{plan.price}</span>
								<span class={styles.duration}>/{plan.duration}天</span>
							</div>
							<div class={styles.planDesc}>{plan.description}</div>
						</NCard>
					))}
				</div>

				{selectedPlan.value && (
					<div class={styles.paySection}>
						<h2 class={styles.sectionTitle}>选择支付方式</h2>
						<NRadioGroup v-model:value={paymentMethod.value} class={styles.payMethodGroup}>
							<NSpace>
								<NRadioButton value="wechat">微信支付</NRadioButton>
								<NRadioButton value="alipay">支付宝</NRadioButton>
							</NSpace>
						</NRadioGroup>

						<div class={styles.orderSummary}>
							<span>套餐: {selectedPlan.value.name}</span>
							<span class={styles.totalAmount}>合计: ¥{selectedPlan.value.price}</span>
						</div>

						{error.value && <div class={styles.error}>{error.value}</div>}

						<NButton
							type="primary"
							size="large"
							loading={loading.value}
							onClick={handlePay}
							class={styles.payButton}
						>
							立即支付
						</NButton>
					</div>
				)}

				{/* WeChat QR Code Modal */}
				<NModal
					show={showQrCode.value}
					onUpdateShow={(val: boolean) => !val && handleCloseQrCode()}
					preset="card"
					title="微信扫码支付"
					class={styles.qrModal}
				>
					<div class={styles.qrContent}>
						{orderStatus.value === 'paid' ? (
							<div class={styles.paySuccess}>
								<NTag type="success" size="large">
									支付成功！
								</NTag>
							</div>
						) : (
							<div class={styles.qrWrapper}>
								{payUrl.value ? (
									<div>
										<p class={styles.qrTip}>请使用微信扫描二维码完成支付</p>
										<div class={styles.qrCode}>
											<img src={`/v1/payment/qrcode?url=${encodeURIComponent(payUrl.value)}`} alt="QR Code" />
										</div>
										<div class={styles.pollingStatus}>
											<NSpin size="small" />
											<span>等待支付结果...</span>
										</div>
									</div>
								) : (
									<NSpin />
								)}
							</div>
						)}
					</div>
				</NModal>
			</div>
		)
	},
})
