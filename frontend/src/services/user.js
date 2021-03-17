import request from '@/utils/request';

export async function query() {
  return request('/api/v1/users');
}

export async function queryCurrent() {
  return request('/api/v1/user/current', {headers: {'Authorization': localStorage.getItem("token")}});
}

export async function login(params) {
  return request('/api/v1/user/login', {
    method: 'POST',
    data: params,
  });
}
