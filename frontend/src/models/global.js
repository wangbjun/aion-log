import {queryLogList, queryPlayer, queryRankList} from "@/services/api";
import moment from "moment";

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
      yield put({
        type: 'saveDefault',
        payload: {
          logList: result.data,
        },
      });
    },
    * fetchRankList({payload}, {call, put, select}) {
      const result = yield call(queryRankList, payload);
      let list = result.data.list
      if (payload.name) {
        list = list.filter(v => {
          return v.player.indexOf(payload.name) !== -1
        })
      }
      if (payload.pro !== undefined) {
        list = list.filter(v => {
          return v.pro === parseInt(payload.pro)
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
      let list = result.data.list
      if (payload.st) {
        list = list.filter(v => {
          return moment(v.time).isAfter(moment(payload.st))
        })
      }
      if (payload.et) {
        list = list.filter(v => {
          return moment(v.time).isBefore(moment(payload.et))
        })
      }
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
