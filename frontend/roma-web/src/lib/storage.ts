const ACCESS_KEY = 'roma.access'
const REFRESH_KEY = 'roma.refresh'

export function saveTokens(access: string, refresh: string) {
  localStorage.setItem(ACCESS_KEY, access)
  localStorage.setItem(REFRESH_KEY, refresh)
}
export function getAccess() { return localStorage.getItem(ACCESS_KEY) }
export function getRefresh() { return localStorage.getItem(REFRESH_KEY) }
export function clearTokens() {
  localStorage.removeItem(ACCESS_KEY)
  localStorage.removeItem(REFRESH_KEY)
}
