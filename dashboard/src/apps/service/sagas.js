import { delay } from 'redux-saga'
import { put, call, takeEvery, all } from 'redux-saga/effects'
import actions from '../../actions'
import modules from '../../modules'

function* post(action) {
  const {payload, error} = yield call(modules.service.post, action.payload)

  if (error) {
    yield put(actions.service.servicePostFailure(action, error, null))
  } else if (payload.error && payload.error != "") {
    yield put(actions.service.servicePostFailure(action, null, payload.error))
  } else {
    yield put(actions.service.servicePostSuccess(action, payload))
  }
}

function* watchGetIndex() {
  yield takeEvery(actions.service.serviceGetIndex, post)
}

export default {
  watchGetIndex,
}
