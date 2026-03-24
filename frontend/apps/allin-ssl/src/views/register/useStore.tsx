/**
 * @file 注册模块状态管理
 */
// External Libraries
import { useMessage } from '@baota/naive-ui/hooks'
import { useError } from '@baota/hooks/error'

// API Imports
import { getLoginCode } from '@api/public'
import { registerApi } from '@api/register'
import { ShallowRef } from 'vue'

const { success } = useMessage()
const { handleError } = useError()

interface RegisterFormData {
	username: string
	email: string
	password: string
	confirmPassword: string
	code: string
}

interface RegisterStoreExposes {
	loading: Ref<boolean>
	codeImg: Ref<string>
	error: ShallowRef<Error | string | null>
	registerData: Ref<RegisterFormData>
	handleRegister: (params: { username: string; password: string; email: string; code: string }) => Promise<void>
	handleGetCode: () => Promise<void>
}

export const useRegisterStore = defineStore('register-store', (): RegisterStoreExposes => {
	const codeImg = ref('')

	const registerData = ref<RegisterFormData>({
		username: '',
		email: '',
		password: '',
		confirmPassword: '',
		code: '',
	})

	const { fetch, error, data, message, loading } = registerApi()

	const handleRegister = async (params: {
		username: string
		password: string
		email: string
		code: string
	}): Promise<void> => {
		try {
			error.value = null
			message.value = true
			await fetch(params)
			const { status } = data.value
			if (status) {
				success('注册成功，正在跳转中...')
				setTimeout(() => (location.href = '/'), 1000)
			} else {
				throw new Error(data.value.message)
			}
		} catch (err: unknown) {
			error.value = (err as Error).message
		}
	}

	const handleGetCode = async (): Promise<void> => {
		try {
			const { data: codeData } = await getLoginCode()
			codeImg.value = codeData.data
		} catch (err) {
			handleError(err)
		}
	}

	return {
		loading,
		codeImg,
		error,
		registerData,
		handleRegister,
		handleGetCode,
	}
})

export const useStore = () => {
	const store = useRegisterStore()
	return { ...store, ...storeToRefs(store) }
}
