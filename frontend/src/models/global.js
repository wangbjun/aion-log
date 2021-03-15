import {addTask, changePlayerType, getTask, queryLogList, queryPlayer, queryRankList, queryStat} from "@/services/api";

const GlobalModel = {
  namespace: 'global',
  state: {
    logList: [],
    rankList: [],
    stat: {},
    playerList: [],
    isRuning: true
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
      yield put({
        type: 'saveDefault',
        payload: {
          rankList: list,
        },
      });
    },
    * fetchStat({payload}, {call, put, select}) {
      const result = yield call(queryStat, payload);
      yield put({
        type: 'saveDefault',
        payload: {
          stat: result.data,
        },
      });
    },
    * fetchPlayerList({payload}, {call, put, select}) {
      const result = yield call(queryPlayer, payload);
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
      yield put({
        type: 'saveDefault',
        payload: {
          playerList: list,
        },
      });
    },
    * changePlayerType({payload}, {call, put, select}) {
      yield call(changePlayerType, payload);
    },
    * getTask({payload}, {call, put, select}) {
      const result = yield call(getTask, payload);
      yield put({
        type: 'saveDefault',
        payload: {
          isRuning: result.data,
        },
      });
    },
    * addTask({payload}, {call, put, select}) {
      return yield call(addTask, payload);
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
