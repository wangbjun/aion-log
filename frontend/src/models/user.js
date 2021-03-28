import {query as queryUsers, queryCurrent} from '@/services/user';
import {setAuthority} from "@/utils/authority";

const UserModel = {
  namespace: 'user',
  state: {
    currentUser: {},
  },
  effects: {
    * fetch(_, {call, put}) {
      const response = yield call(queryUsers);
      yield put({
        type: 'save',
        payload: response,
      });
    },

    * fetchCurrent(_, {call, put}) {
      const response = yield call(queryCurrent);
      yield put({
        type: 'saveCurrentUser',
        payload: response,
      });
    },
  },
  reducers: {
    saveCurrentUser(state, action) {
      setAuthority(action.payload.data.currentAuthority);
      return {...state, currentUser: action.payload.data || {}};
    },
  },
};
export default UserModel;
