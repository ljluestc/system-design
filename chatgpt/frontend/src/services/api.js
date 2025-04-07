import axios from 'axios';

const api = axios.create({ baseURL: 'http://localhost:3001/api' });

export async function login(data) {
  const res = await api.post('/users/login', data);
  return res.data;
}

export async function register(data) {
  await api.post('/users/register', data);
}

export async function authenticate(token) {
  try {
    await api.get('/users/auth', { headers: { Authorization: `Bearer ${token}` } });
    return true;
  } catch {
    return false;
  }
}