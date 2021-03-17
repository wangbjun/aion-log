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

export async function queryStat(params) {
  return request('/api/v1/stats', {
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

export async function changePlayerType(params) {
  return request('/api/v1/changeType', {
    params,
    method: "GET",
    headers: {'Authorization': localStorage.getItem("token")}
  });
}

export async function getTask(params) {
  return request('/api/v1/getTask', {
    params,
    method: "GET"
  });
}
