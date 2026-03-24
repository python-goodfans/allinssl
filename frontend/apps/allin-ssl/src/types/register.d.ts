export interface RegisterParams {
	username: string
	password: string
	email: string
	code: string
}

export interface CheckUsernameParams {
	username: string
}

export interface CheckUsernameResponse {
	exists: boolean
}
