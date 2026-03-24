export interface RegisterParams {
	username: string
	password: string
	email: string
	code: string
}

export interface RegisterResponse {
	message: string
	code: number
	status: boolean
	data: null
}

export interface CheckUsernameResponse {
	code: number
	status: boolean
	data: {
		exists: boolean
	}
	message: string
}
