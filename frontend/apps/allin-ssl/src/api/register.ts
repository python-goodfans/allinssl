// External library dependencies
import type { useAxiosReturn } from '@baota/hooks/axios'
import type { AxiosResponseData } from '@/types/public'

// Type imports
import type { RegisterParams, RegisterResponse, CheckUsernameResponse } from '@/types/register'

// Relative internal imports
import { useApi } from '@api/index'

/**
 * @description 用户注册
 */
export const registerApi = (params?: RegisterParams): useAxiosReturn<RegisterResponse, RegisterParams> =>
	useApi<RegisterResponse, RegisterParams>('/v1/register/submit', params)

/**
 * @description 检查用户名是否可用
 */
export const checkUsernameApi = (params?: { username: string }): useAxiosReturn<CheckUsernameResponse, { username: string }> =>
	useApi<CheckUsernameResponse, { username: string }>('/v1/register/check_username', params)
