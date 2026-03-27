/**
 * @file 注册控制器
 */
// External Libraries
import md5 from 'crypto-js/md5'
import type { FormInst } from 'naive-ui'

// Relative Internal Imports
import { useStore } from './useStore'

const encryptPassword = (password: string): string => {
	return md5(`${password}_bt_all_in_ssl`).toString()
}

export const useController = () => {
	const store = useStore()
	const { error, formData, confirmPassword, handleRegister, handleGetCode, loading, codeImg } = store

	const formRef = ref<FormInst | null>(null)

	/**
	 * 表单校验规则
	 */
	const rules = {
		username: [
			{ required: true, message: '请输入用户名', trigger: ['input', 'blur'] },
			{ min: 3, max: 32, message: '用户名长度须在3-32个字符之间', trigger: ['input', 'blur'] },
		],
		email: [
			{ required: false, type: 'email', message: '请输入有效的邮箱地址', trigger: ['input', 'blur'] },
		],
		password: [
			{ required: true, message: '请输入密码', trigger: ['input', 'blur'] },
			{ min: 6, message: '密码长度至少6个字符', trigger: ['input', 'blur'] },
		],
		confirmPassword: [
			{
				required: true,
				validator: (_rule: unknown, _value: string) => {
					if (confirmPassword.value !== formData.value.password) {
						return new Error('两次输入的密码不一致')
					}
					return true
				},
				trigger: ['input', 'blur'],
			},
		],
		code: [{ required: true, message: '请输入验证码', trigger: ['input', 'blur'] }],
	}

	/**
	 * 处理表单提交
	 */
	const handleSubmit = async (event: Event): Promise<void> => {
		event.preventDefault()
		if (!formRef.value) return

		try {
			await formRef.value.validate()
			const encryptedPassword = encryptPassword(formData.value.password)
			await handleRegister({ ...formData.value, password: encryptedPassword })
		} catch (validationErrors) {
			console.log('表单校验失败:', validationErrors)
		}
	}

	const handleKeyup = (event: KeyboardEvent): void => {
		if (event.key === 'Enter') {
			handleSubmit(event as unknown as Event)
		}
	}

	// 监听错误信息，5秒后清除
	const scope = effectScope()
	scope.run(() => {
		watch(error, (newValue) => {
			if (newValue) {
				setTimeout(() => {
					error.value = null
				}, 5000)
			}
		})
		onScopeDispose(() => scope.stop())
	})

	onMounted(() => {
		handleGetCode()
	})

	return {
		...store,
		formRef,
		rules,
		handleSubmit,
		handleKeyup,
	}
}
