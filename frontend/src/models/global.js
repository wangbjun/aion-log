import {changePlayerType, getTask, queryLogList, queryPlayer, queryRankList} from "@/services/api";
import moment from "moment";

const GlobalModel = {
  namespace: 'global',
  state: {
    logList: [],
    rankList: [],
    playerList: [],
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
      list.forEach(v => {
        v.rate = v.count / v.all_count
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
      yield put({
        type: 'saveDefault',
        payload: {
          playerList: list,
        },
      });
    },
    * changePlayerType({payload}, {call, put, select}) {
      return yield call(changePlayerType, payload);
    },
    * getTask({payload}, {call, put, select}) {
      const result = yield call(getTask, payload);
      yield put({
        type: 'saveDefault',
        payload: {
          taskStatus: result.data,
        },
      });
    },
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
