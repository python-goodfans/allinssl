// External Libraries
import md5 from 'crypto-js/md5'
import type { FormInst } from 'naive-ui'

// Absolute Internal Imports
import { useError } from '@baota/hooks/error'

// Relative Internal Imports
import { useStore } from './useStore'

/**
 * @file 注册控制器
 * @description 处理注册页面的业务逻辑，包括表单验证、密码加密等功能
 */

// ==================== 工具函数 ====================
const encryptPassword = (password: string): string => {
	return md5(`${password}_bt_all_in_ssl`).toString()
}

// ==================== 控制器逻辑 ====================
export const useController = () => {
	const store = useStore()
	const { handleError } = useError()
	const { error, registerData, handleRegister: storeHandleRegister, handleGetCode } = store

	const formRef = ref<FormInst | null>(null)

	const handleRegisterBusiness = async (): Promise<void> => {
		try {
			const { password, confirmPassword } = registerData.value
			if (password !== confirmPassword) {
				error.value = '两次输入的密码不一致'
				return
			}
			const encryptedPassword = encryptPassword(password)
			await storeHandleRegister({
				...registerData.value,
				password: encryptedPassword,
			})
		} catch (err) {
			handleError(err)
		}
	}

	const handleSubmit = async (event: Event): Promise<void> => {
		event.preventDefault()
		if (!formRef.value) return
		try {
			await formRef.value.validate()
			await handleRegisterBusiness()
		} catch (validationErrors) {
			console.log('表单校验失败:', validationErrors)
		}
	}

	const handleKeyup = (event: KeyboardEvent): void => {
		if (event.key === 'Enter') {
			handleSubmit(event as unknown as Event)
		}
	}

	const scope = effectScope()
	scope.run(() => {
		watch(error, (newValue) => {
			if (newValue) {
				setTimeout(() => {
					error.value = null
				}, 5000)
			}
		})
		onScopeDispose(() => {
			scope.stop()
		})
	})

	onMounted(() => {
		handleGetCode()
	})

	return {
		...store,
		formRef,
		handleSubmit,
		handleKeyup,
	}
}
