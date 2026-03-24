// External Libraries
import {
	NButton,
	NCard,
	NDataTable,
	NModal,
	NSpin,
	NTabs,
	NTabPane,
	NTag,
	NForm,
	NFormItem,
	NInput,
	NInputNumber,
	NSwitch,
	NSpace,
} from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'

// Absolute Internal Imports
import { useThemeCssVar } from '@baota/naive-ui/theme'

// Type imports
import type { Order, Plan, PaymentConfig } from '@/types/payment'

// Relative Internal Imports
import { useController } from './useController'

// Side-effect Imports
import styles from './index.module.css'

export default defineComponent({
	name: 'PaymentManageView',
	setup() {
		const {
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
			handleSavePlan,
			handleDelPlan,
			handleSaveConfig,
			handleTabChange,
			handleAddPlan,
			handleEditPlan,
		} = useController()

		const cssVar = useThemeCssVar(['textColor2', 'actionColor', 'primaryColor'])

		// 订单列表列定义
		const orderColumns: DataTableColumns<Order> = [
			{ title: '订单号', key: 'order_no', ellipsis: true },
			{ title: '金额', key: 'amount', render: (row) => `¥${row.amount.toFixed(2)}` },
			{
				title: '支付方式',
				key: 'payment_method',
				render: (row) => (
					<NTag type={row.payment_method === 'wechat' ? 'success' : 'info'}>
						{row.payment_method === 'wechat' ? '微信' : '支付宝'}
					</NTag>
				),
			},
			{
				title: '状态',
				key: 'status',
				render: (row) => {
					const typeMap: Record<string, 'success' | 'warning' | 'error' | 'default'> = {
						paid: 'success',
						pending: 'warning',
						failed: 'error',
						refunded: 'default',
					}
					const labelMap: Record<string, string> = {
						paid: '已支付',
						pending: '待支付',
						failed: '失败',
						refunded: '已退款',
					}
					return <NTag type={typeMap[row.status] || 'default'}>{labelMap[row.status] || row.status}</NTag>
				},
			},
			{ title: '创建时间', key: 'create_time' },
			{ title: '支付时间', key: 'pay_time' },
		]

		// 套餐列表列定义
		const planColumns: DataTableColumns<Plan> = [
			{ title: '套餐名称', key: 'name' },
			{ title: '价格', key: 'price', render: (row) => `¥${row.price.toFixed(2)}` },
			{ title: '有效天数', key: 'duration', render: (row) => `${row.duration} 天` },
			{ title: '描述', key: 'description', ellipsis: true },
			{
				title: '状态',
				key: 'status',
				render: (row) => <NTag type={row.status === 1 ? 'success' : 'default'}>{row.status === 1 ? '启用' : '禁用'}</NTag>,
			},
			{
				title: '操作',
				key: 'actions',
				render: (row) => (
					<NSpace>
						<NButton size="small" onClick={() => handleEditPlan(row as unknown as Record<string, unknown>)}>
							编辑
						</NButton>
						<NButton size="small" type="error" onClick={() => handleDelPlan(row.id)}>
							删除
						</NButton>
					</NSpace>
				),
			},
		]

		return () => (
			<div style={cssVar.value} class={styles.container}>
				<div class={styles.header}>
					<h1 class={styles.title}>支付管理</h1>
				</div>

				<NCard>
					<NTabs
						value={activeTab.value}
						onUpdateValue={(v) => handleTabChange(v as 'orders' | 'plans' | 'config')}
						type="line"
					>
						{/* 订单管理 */}
						<NTabPane name="orders" tab="订单管理">
							<NSpin show={loading.value}>
								<NDataTable
									columns={orderColumns}
									data={orders.value}
									pagination={{
										page: orderPage.value,
										pageSize: orderPageSize.value,
										itemCount: orderTotal.value,
										onChange: (p) => {
											orderPage.value = p
											fetchOrders()
										},
									}}
									striped
								/>
							</NSpin>
						</NTabPane>

						{/* 套餐管理 */}
						<NTabPane name="plans" tab="套餐管理">
							<div class={styles.planToolbar}>
								<NButton type="primary" onClick={handleAddPlan}>
									新增套餐
								</NButton>
							</div>
							<NSpin show={loading.value}>
								<NDataTable columns={planColumns} data={plans.value} striped />
							</NSpin>
						</NTabPane>

						{/* 支付配置 */}
						<NTabPane name="config" tab="支付配置">
							<NSpin show={loading.value}>
								<NForm labelPlacement="left" labelWidth="180px" class={styles.configForm}>
									<h3 class={styles.configSection}>微信支付配置</h3>
									<NFormItem label="AppID">
										<NInput v-model:value={(paymentConfig.value as PaymentConfig).wechat_app_id} placeholder="wx***" />
									</NFormItem>
									<NFormItem label="商户号(MchID)">
										<NInput v-model:value={(paymentConfig.value as PaymentConfig).wechat_mch_id} placeholder="商户号" />
									</NFormItem>
									<NFormItem label="API Key">
										<NInput
											v-model:value={(paymentConfig.value as PaymentConfig).wechat_api_key}
											type="password"
											placeholder="API密钥"
											showPasswordOn="click"
										/>
									</NFormItem>
									<NFormItem label="回调通知URL">
										<NInput
											v-model:value={(paymentConfig.value as PaymentConfig).wechat_notify_url}
											placeholder="https://yourdomain.com/v1/pay_notify/wechat"
										/>
									</NFormItem>

									<h3 class={styles.configSection}>支付宝配置</h3>
									<NFormItem label="AppID">
										<NInput v-model:value={(paymentConfig.value as PaymentConfig).alipay_app_id} placeholder="202***" />
									</NFormItem>
									<NFormItem label="应用私钥">
										<NInput
											v-model:value={(paymentConfig.value as PaymentConfig).alipay_private_key}
											type="textarea"
											placeholder="RSA私钥"
											rows={4}
										/>
									</NFormItem>
									<NFormItem label="支付宝公钥">
										<NInput
											v-model:value={(paymentConfig.value as PaymentConfig).alipay_public_key}
											type="textarea"
											placeholder="支付宝公钥"
											rows={4}
										/>
									</NFormItem>
									<NFormItem label="回调通知URL">
										<NInput
											v-model:value={(paymentConfig.value as PaymentConfig).alipay_notify_url}
											placeholder="https://yourdomain.com/v1/pay_notify/alipay"
										/>
									</NFormItem>
									<NFormItem label="沙箱模式">
										<NSwitch v-model:value={(paymentConfig.value as PaymentConfig).alipay_sandbox} />
									</NFormItem>

									<NFormItem>
										<NButton type="primary" loading={loading.value} onClick={handleSaveConfig}>
											保存配置
										</NButton>
									</NFormItem>
								</NForm>
							</NSpin>
						</NTabPane>
					</NTabs>
				</NCard>

				{/* 套餐编辑弹窗 */}
				<NModal
					show={showPlanModal.value}
					onUpdateShow={(v) => (showPlanModal.value = v)}
					title={editingPlan.value.id ? '编辑套餐' : '新增套餐'}
					preset="card"
					style="width: 600px"
				>
					<NForm labelPlacement="left" labelWidth="120px">
						<NFormItem label="套餐名称">
							<NInput v-model:value={editingPlan.value.name} placeholder="请输入套餐名称" />
						</NFormItem>
						<NFormItem label="套餐描述">
							<NInput v-model:value={editingPlan.value.description} type="textarea" placeholder="套餐描述" rows={3} />
						</NFormItem>
						<NFormItem label="价格（元）">
							<NInputNumber v-model:value={editingPlan.value.price} min={0} precision={2} placeholder="0.00" />
						</NFormItem>
						<NFormItem label="有效天数">
							<NInputNumber v-model:value={editingPlan.value.duration} min={1} placeholder="30" />
						</NFormItem>
						<NFormItem label="功能描述">
							<NInput v-model:value={editingPlan.value.features} type="textarea" placeholder="JSON格式或文本" rows={3} />
						</NFormItem>
						<NFormItem label="排序">
							<NInputNumber v-model:value={editingPlan.value.sort_order} min={0} placeholder="0" />
						</NFormItem>
					</NForm>
					{{
						footer: () => (
							<NSpace justify="end">
								<NButton onClick={() => (showPlanModal.value = false)}>取消</NButton>
								<NButton type="primary" loading={loading.value} onClick={handleSavePlan}>
									保存
								</NButton>
							</NSpace>
						),
					}}
				</NModal>
			</div>
		)
	},
})
