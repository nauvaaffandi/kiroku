import axios from 'axios'

export const api = axios.create({
  baseURL: '',
  withCredentials: true,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json'
  }
})

export const initApiClient = async () => {
  // No remote backend in offline app.
  return ''
}

export const refreshBackendUrl = async () => {
  return ''
}
