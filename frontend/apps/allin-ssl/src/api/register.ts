// Type imports
import type { useAxiosReturn } from '@baota/hooks/axios'
import type {
	CheckUsernameParams,
	CheckUsernameResponse,
	RegisterParams,
} from '@/types/register'
import type { AxiosResponseData } from '@/types/public'

// Relative internal imports
import { useApi } from '@api/index'

/**
 * @description 用户注册
 * @param {RegisterParams} [params] 注册参数
 * @returns {useAxiosReturn<AxiosResponseData, RegisterParams>} 注册操作的组合式 API 调用封装
 */
export const registerApi = (params?: RegisterParams): useAxiosReturn<AxiosResponseData, RegisterParams> =>
	useApi<AxiosResponseData, RegisterParams>('/v1/register/submit', params)

/**
 * @description 检查用户名是否可用
 * @param {CheckUsernameParams} [params] 检查参数
 * @returns {useAxiosReturn<CheckUsernameResponse, CheckUsernameParams>} 检查结果的组合式 API 调用封装
 */
export const checkUsernameApi = (
	params?: CheckUsernameParams
): useAxiosReturn<CheckUsernameResponse, CheckUsernameParams> =>
	useApi<CheckUsernameResponse, CheckUsernameParams>('/v1/register/check_username', params)
