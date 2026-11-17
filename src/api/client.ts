import axios from 'axios'

export const api = axios.create({
  baseURL: '',
  withCredentials: true,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json'
  }
})
