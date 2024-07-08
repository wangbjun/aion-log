import {queryLogList, queryPlayer, queryRankList} from "@/services/api";
import {notification} from "antd";

const GlobalModel = {
  namespace: 'global',
  state: {
    visible: true,
    logList: [],
    rankList: [],
    playerList: [],
    stockList: [],
    taskStatus: {}
  },
  effects: {
    * fetchLogList({payload}, {call, put, select}) {
      const result = yield call(queryLogList, payload);
      if (result.code !== 200) {
        notification.error(result.msg);
        return
      }
      yield put({
        type: 'saveDefault',
        payload: {
          logList: result.data,
        },
      });
    },
    * fetchRankList({payload}, {call, put, select}) {
      const result = yield call(queryRankList, payload);
      if (result.code !== 200) {
        notification.error({
          message: result.msg
        });
        return
      }
      let list = result.data.list
      if (payload.name) {
        list = list.filter(v => {
          return v.player.indexOf(payload.name) !== -1
        })
      }
      if (payload.class !== undefined) {
        list = list.filter(v => {
          return v.class === parseInt(payload.class)
        })
      }
      list.forEach(v => {
        v.rate = v.counts / v.all_counts
      })
      yield put({
        type: 'saveDefault',
        payload: {
          rankList: list,
        },
      });
    },
    * fetchPlayerList({payload}, {call, put, select}) {
      const result = yield call(queryPlayer, payload);
      if (result.code !== 200) {
        notification.error({
          message: result.msg
        });
        return
      }
      let list = result.data.list
      if (payload.name) {
        list = list.filter(v => {
          return v.name.indexOf(payload.name) !== -1
        })
      }
      if (payload.type) {
        list = list.filter(v => {
          return v.type === parseInt(payload.type)
        })
      }
      if (payload.class !== undefined) {
        list = list.filter(v => {
          return v.class === parseInt(payload.class)
        })
      }
      yield put({
        type: 'saveDefault',
        payload: {
          playerList: list,
        },
      });
    },
    * closeModal({payload}, {call, put, select}) {
      yield put({
        type: 'saveDefault',
        payload: {
          visible: false,
        },
      });
      sessionStorage.setItem("modalClose", "true")
    }
  },
  reducers: {
    saveDefault(state, {payload}) {
      return {
        ...state, ...payload
      };
    },
  },
};
export default GlobalModel;
