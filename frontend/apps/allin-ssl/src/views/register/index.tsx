// External Libraries
import { NInput, NButton, NForm, NFormItem, NIcon, NImage } from 'naive-ui'
import { UserOutlined, LockOutlined, CodeOutlined, MailOutlined } from '@vicons/antd'

// Absolute Internal Imports
import { useTheme, useThemeCssVar } from '@baota/naive-ui/theme'
import { $t } from '@locales/index'

// Relative Internal Imports
import { useController } from './useController'

// Side-effect Imports
import styles from './index.module.css'

export default defineComponent({
	name: 'RegisterView',
	setup() {
		const { loading, error, formData, confirmPassword, handleSubmit, handleKeyup, handleGetCode, codeImg, formRef, rules } =
			useController()
		const { isDark } = useTheme()
		const cssVar = useThemeCssVar(['textColor2', 'actionColor', 'errorColor', 'primaryColor', 'primaryColorSuppl'])

		return () => (
			<div style={cssVar.value}>
				<div
					class={styles.container}
					style={`background-image:${isDark.value ? 'url(/static/images/login-bg-dark.svg)' : 'url(/static/images/login-bg.svg)'};`}
				>
					<div class={styles.loginBox}>
						<div class={styles.leftSection}>
							<h2 class={styles.leftTitle}>
								<img src="/static/images/logo.png" alt="logo" class={styles.logo} />
								<span>ALLinSSL</span>
							</h2>
							<div class={styles.leftImageWrapper}>
								<img src="/static/images/login-display.svg" alt="注册" class={styles.leftImage} />
							</div>
						</div>
						<div class={styles.rightSection}>
							<div class={styles.formContainer}>
								<h1 class={styles.title}>创建账号</h1>
								<NForm ref={formRef} model={formData.value} rules={rules} onSubmit={handleSubmit} class={styles.formWrapper}>
									<div class={styles.formContent}>
										<div class={styles.formInputs}>
											<NFormItem show-label={false} path="username">
												<NInput
													v-model:value={formData.value.username}
													onKeyup={handleKeyup}
													disabled={loading.value}
													placeholder="请输入用户名（3-32个字符）"
													clearable
													size="large"
												>
													{{
														prefix: () => <NIcon component={UserOutlined} class={styles.icon} />,
													}}
												</NInput>
											</NFormItem>
											<NFormItem show-label={false} path="email">
												<NInput
													v-model:value={formData.value.email}
													onKeyup={handleKeyup}
													disabled={loading.value}
													placeholder="请输入邮箱（选填）"
													clearable
													size="large"
												>
													{{
														prefix: () => <NIcon component={MailOutlined} class={styles.icon} />,
													}}
												</NInput>
											</NFormItem>
											<NFormItem show-label={false} path="password">
												<NInput
													onKeyup={handleKeyup}
													disabled={loading.value}
													v-model:value={formData.value.password}
													type="password"
													placeholder="请输入密码（至少6位）"
													clearable
													size="large"
													showPasswordOn="click"
												>
													{{
														prefix: () => <NIcon component={LockOutlined} class={styles.icon} />,
													}}
												</NInput>
											</NFormItem>
											<NFormItem show-label={false} path="confirmPassword">
												<NInput
													onKeyup={handleKeyup}
													disabled={loading.value}
													v-model:value={confirmPassword.value}
													type="password"
													placeholder="请再次输入密码"
													clearable
													size="large"
													showPasswordOn="click"
												>
													{{
														prefix: () => <NIcon component={LockOutlined} class={styles.icon} />,
													}}
												</NInput>
											</NFormItem>
											<NFormItem show-label={false} path="code">
												<NInput
													onKeyup={handleKeyup}
													disabled={loading.value}
													v-model:value={formData.value.code}
													type="text"
													placeholder="请输入验证码"
													clearable
													size="large"
												>
													{{
														prefix: () => <NIcon component={CodeOutlined} class={styles.icon} />,
														suffix: () => (
															<span
																onClick={handleGetCode}
																title="点击刷新验证码"
																class={`w-[10rem] h-[4rem] mr-[-1.5rem] flex items-center justify-center relative z-[999] cursor-pointer bg-slate-400 rounded-r-[6px] ${styles.codeImageContainer}`}
															>
																<NImage
																	src={codeImg.value}
																	preview-disabled
																	class="max-w-full max-h-full object-contain"
																/>
															</span>
														),
													}}
												</NInput>
											</NFormItem>
										</div>

										<div class={styles.formActions}>
											{error.value && <div class={styles.error}>{error.value}</div>}
											<NButton
												type="primary"
												size="large"
												block
												loading={loading.value}
												attrType="submit"
												onClick={handleSubmit}
												class={styles.submitBtn}
											>
												{loading.value ? '注册中...' : '立即注册'}
											</NButton>
											<div class={styles.loginLink}>
												已有账号？
												<a href="/login" class={styles.link}>
													去登录
												</a>
											</div>
										</div>
									</div>
								</NForm>
							</div>
						</div>
					</div>
				</div>
			</div>
		)
	},
})
