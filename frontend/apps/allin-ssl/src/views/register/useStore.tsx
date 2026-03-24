/**
 * @file 注册模块状态管理
 */
// External Libraries
import { useMessage } from '@baota/naive-ui/hooks'
import { useError } from '@baota/hooks/error'

// Type Imports
import type { RegisterParams } from '@/types/register'

// Absolute Internal Imports
import { registerApi } from '@api/register'
import { getLoginCode } from '@api/public'

/** 消息提示 */
const { success } = useMessage()
const { handleError } = useError()

export const useRegisterStore = defineStore('register-store', () => {
	const loading = ref(false)
	const codeImg = ref('')
	const error = shallowRef<string | null>(null)

	const formData = ref<RegisterParams>({
		username: '',
		password: '',
		email: '',
		code: '',
	})

	const confirmPassword = ref('')

	// 初始化注册请求
	const { fetch, data } = registerApi()

	/**
	 * 获取验证码
	 */
	const handleGetCode = async (): Promise<void> => {
		try {
			const { data: codeData } = await getLoginCode()
			codeImg.value = codeData.data
		} catch (err) {
			handleError(err)
		}
	}

	/**
	 * 注册处理
	 */
	const handleRegister = async (params: RegisterParams): Promise<void> => {
		try {
			loading.value = true
			error.value = null
			await fetch(params)
			const { status, message } = data.value
			if (status) {
				success('注册成功，正在跳转...')
				setTimeout(() => (location.href = '/'), 1000)
			} else {
				throw new Error(message)
			}
		} catch (err: unknown) {
			error.value = (err as Error).message
			handleGetCode()
		} finally {
			loading.value = false
		}
	}

	/**
	 * 重置表单
	 */
	const resetForm = (): void => {
		formData.value = { username: '', password: '', email: '', code: '' }
		confirmPassword.value = ''
		error.value = null
	}

	return {
		loading,
		codeImg,
		error,
		formData,
		confirmPassword,
		handleRegister,
		handleGetCode,
		resetForm,
	}
})

export const useStore = () => {
	const store = useRegisterStore()
	return { ...store, ...storeToRefs(store) }
}
