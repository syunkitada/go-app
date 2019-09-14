import {call, put, takeEvery} from 'redux-saga/effects';

import actions from '../../actions';
import modules from '../../modules';

function* loginWithToken(action) {
  const {payload, error} = yield call(modules.auth.loginWithToken);

  if (error) {
    yield put(actions.auth.authLoginFailure({error: error.message}));
  } else if (payload.Error && payload.Error !== '') {
    yield put(actions.auth.authLoginFailure({error: payload.Error}));
  } else {
    yield put(
      actions.auth.authLoginSuccess({
        authority: payload.Data.Login.Authority,
        username: payload.Data.Login.Authority.Name,
      }),
    );
  }
}

function* login(action) {
  const {payload, error} = yield call(modules.auth.login, action.payload);

  if (error) {
    yield put(actions.auth.authLoginFailure(error.message));
  } else if (payload.Error && payload.Error !== '') {
    yield put(actions.auth.authLoginFailure(payload.error));
  } else {
    yield put(
      actions.auth.authLoginSuccess({
        authority: payload.Data.Login.Authority,
        username: payload.Data.Login.Authority.Name,
      }),
    );
  }
}

function* logout(action) {
  const {error} = yield call(modules.auth.logout);

  if (error) {
    yield put(actions.auth.authLogoutFailure(error));
  } else {
    yield put(actions.auth.authLogoutSuccess());
  }
}

function* watchLogin() {
  yield takeEvery(actions.auth.authLogin, login);
}

function* watchLoginWithToken() {
  yield takeEvery(actions.auth.authLoginWithToken, loginWithToken);
}

function* watchLogout() {
  yield takeEvery(actions.auth.authLogout, logout);
}

export default {
  watchLogin,
  watchLoginWithToken,
  watchLogout,
};
