import {
	NButton,
	NCard,
	NDataTable,
	NForm,
	NFormItem,
	NInput,
	NInputNumber,
	NModal,
	NSpace,
	NSwitch,
	NTabs,
	NTabPane,
	NTag,
} from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { useThemeCssVar } from '@baota/naive-ui/theme'
import { useController } from './useController'
import type { Order, Plan } from '@/types/payment'
import styles from './index.module.css'

export default defineComponent({
	name: 'PaymentManageView',
	setup() {
		const {
			loading,
			orders,
			orderTotal,
			orderPage,
			plans,
			paymentConfig,
			activeTab,
			showPlanModal,
			editingPlan,
			fetchOrders,
			handleOpenAddPlan,
			handleOpenEditPlan,
			handleSavePlanSubmit,
			handleConfirmDelPlan,
			handleSavePaymentConfig,
		} = useController()

		const cssVar = useThemeCssVar(['textColor2', 'actionColor', 'primaryColor'])

		const orderStatusMap: Record<string, { label: string; type: 'default' | 'success' | 'error' | 'warning' }> = {
			pending: { label: '待支付', type: 'warning' },
			paid: { label: '已支付', type: 'success' },
			failed: { label: '支付失败', type: 'error' },
			refunded: { label: '已退款', type: 'default' },
		}

		const orderColumns: DataTableColumns<Order> = [
			{ title: '订单号', key: 'order_no', width: 200 },
			{ title: '用户', key: 'user_id', width: 120 },
			{ title: '金额', key: 'amount', width: 100, render: (row) => `¥${row.amount}` },
			{
				title: '支付方式',
				key: 'payment_method',
				width: 120,
				render: (row) => (row.payment_method === 'wechat' ? '微信支付' : '支付宝'),
			},
			{
				title: '状态',
				key: 'status',
				width: 100,
				render: (row) => {
					const s = orderStatusMap[row.status] || { label: row.status, type: 'default' as const }
					return <NTag type={s.type}>{s.label}</NTag>
				},
			},
			{ title: '创建时间', key: 'create_time', width: 180 },
		]

		const planColumns: DataTableColumns<Plan> = [
			{ title: '套餐名', key: 'name' },
			{ title: '价格', key: 'price', render: (row) => `¥${row.price}` },
			{ title: '时长(天)', key: 'duration' },
			{ title: '排序', key: 'sort_order' },
			{
				title: '操作',
				key: 'actions',
				render: (row) => (
					<NSpace>
						<NButton size="small" onClick={() => handleOpenEditPlan(row)}>
							编辑
						</NButton>
						<NButton size="small" type="error" onClick={() => handleConfirmDelPlan(row.id)}>
							删除
						</NButton>
					</NSpace>
				),
			},
		]

		return () => (
			<div style={cssVar.value} class={styles.container}>
				<h1 class={styles.title}>支付管理</h1>
				<NTabs v-model:value={activeTab.value}>
					<NTabPane name="orders" tab="订单管理">
						<NDataTable
							columns={orderColumns}
							data={orders.value}
							loading={loading.value}
							pagination={{
								page: orderPage.value,
								pageSize: 20,
								itemCount: orderTotal.value,
								onChange: (page: number) => fetchOrders(page),
							}}
						/>
					</NTabPane>

					<NTabPane name="plans" tab="套餐管理">
						<div class="mb-4">
							<NButton type="primary" onClick={handleOpenAddPlan}>
								添加套餐
							</NButton>
						</div>
						<NDataTable columns={planColumns} data={plans.value} loading={loading.value} />
					</NTabPane>

					<NTabPane name="config" tab="支付配置">
						<NCard class={styles.configCard}>
							<NForm>
								<h3 class={styles.configSection}>微信支付配置</h3>
								<NFormItem label="AppID">
									<NInput v-model:value={paymentConfig.value.wechat_app_id} placeholder="请输入微信 AppID" />
								</NFormItem>
								<NFormItem label="商户号">
									<NInput v-model:value={paymentConfig.value.wechat_mch_id} placeholder="请输入商户号" />
								</NFormItem>
								<NFormItem label="API密钥">
									<NInput
										v-model:value={paymentConfig.value.wechat_api_key}
										type="password"
										placeholder="请输入 API 密钥"
										showPasswordOn="click"
									/>
								</NFormItem>
								<NFormItem label="回调地址">
									<NInput
										v-model:value={paymentConfig.value.wechat_notify_url}
										placeholder="如: https://yourdomain.com/v1/pay_notify/wechat"
									/>
								</NFormItem>

								<h3 class={styles.configSection}>支付宝配置</h3>
								<NFormItem label="AppID">
									<NInput v-model:value={paymentConfig.value.alipay_app_id} placeholder="请输入支付宝 AppID" />
								</NFormItem>
								<NFormItem label="私钥">
									<NInput
										v-model:value={paymentConfig.value.alipay_private_key}
										type="textarea"
										placeholder="请输入应用私钥 (PKCS1 格式)"
										rows={4}
									/>
								</NFormItem>
								<NFormItem label="公钥">
									<NInput
										v-model:value={paymentConfig.value.alipay_public_key}
										type="textarea"
										placeholder="请输入支付宝公钥"
										rows={4}
									/>
								</NFormItem>
								<NFormItem label="回调地址">
									<NInput
										v-model:value={paymentConfig.value.alipay_notify_url}
										placeholder="如: https://yourdomain.com/v1/pay_notify/alipay"
									/>
								</NFormItem>
								<NFormItem label="沙箱模式">
									<NSwitch v-model:value={paymentConfig.value.alipay_sandbox} />
								</NFormItem>

								<NButton type="primary" loading={loading.value} onClick={handleSavePaymentConfig}>
									保存配置
								</NButton>
							</NForm>
						</NCard>
					</NTabPane>
				</NTabs>

				{/* Plan Edit Modal */}
				<NModal
					show={showPlanModal.value}
					onUpdateShow={(val: boolean) => (showPlanModal.value = val)}
					preset="card"
					title={editingPlan.value.id ? '编辑套餐' : '添加套餐'}
					class={styles.planModal}
				>
					<NForm>
						<NFormItem label="套餐名称">
							<NInput v-model:value={editingPlan.value.name} placeholder="请输入套餐名称" />
						</NFormItem>
						<NFormItem label="描述">
							<NInput v-model:value={editingPlan.value.description} placeholder="请输入套餐描述" />
						</NFormItem>
						<NFormItem label="价格(元)">
							<NInputNumber v-model:value={editingPlan.value.price} min={0} precision={2} />
						</NFormItem>
						<NFormItem label="时长(天)">
							<NInputNumber v-model:value={editingPlan.value.duration} min={1} />
						</NFormItem>
						<NFormItem label="功能说明">
							<NInput v-model:value={editingPlan.value.features} type="textarea" placeholder="功能说明 (JSON格式)" />
						</NFormItem>
						<NFormItem label="排序">
							<NInputNumber v-model:value={editingPlan.value.sort_order} min={0} />
						</NFormItem>
					</NForm>
					<div class="flex justify-end gap-3 mt-4">
						<NButton onClick={() => (showPlanModal.value = false)}>取消</NButton>
						<NButton type="primary" loading={loading.value} onClick={handleSavePlanSubmit}>
							保存
						</NButton>
					</div>
				</NModal>
			</div>
		)
	},
})
