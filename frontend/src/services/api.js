import request from '@/utils/request';

export async function queryLogList(params) {
  return request('/api/v1/logs', {
    params,
    method: "GET"
  });
}

export async function queryRankList(params) {
  return request('/api/v1/ranks', {
    params,
    method: "GET"
  });
}

export async function queryPlayer(params) {
  return request('/api/v1/players', {
    params,
    method: "GET"
  });
}
