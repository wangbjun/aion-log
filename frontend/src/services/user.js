import request from '@/utils/request';
export async function query() {
  return request('/api/v1/users');
}
export async function queryCurrent() {
  return request('/api/v1/user/current');
}
export async function Login(params) {
  return request('/api/v1/user/login', {
    method: 'POST',
    data: params,
  });
}
