import actionCreatorFactory from 'typescript-fsa';

const actionCreator = actionCreatorFactory();

export const serviceGetIndex = actionCreator<{params: any}>(
  'SERVICE_GET_INDEX',
);
export const serviceStartBackgroundSync = actionCreator(
  'SERVICE_START_BACKGROUND_SYNC',
);
export const serviceStopBackgroundSync = actionCreator(
  'SERVICE_STOP_BACKGROUND_SYNC',
);
export const serviceGetQueries = actionCreator<{
  queries: any;
  isSync: any;
  params: any;
}>('SERVICE_GET_QUERIES');
export const serviceSubmitQueries = actionCreator<{
  queryKind: any;
  action: any;
  fieldMap: any;
  targets: any;
  params: any;
}>('SERVICE_SUBMIT_QUERIES');
export const serviceCloseErr = actionCreator('SERVICE_CLOSE_ERR');
export const serviceCloseGetQueriesTctx = actionCreator(
  'SERVICE_CLOSE_GET_QUERIES_TCTX',
);
export const serviceCloseSubmitQueriesTctx = actionCreator(
  'SERVICE_CLOSE_SUBMIT_QUERIES_TCTX',
);
export const servicePostSuccess = actionCreator<{action: any; data: any}>(
  'SERVICE_POST_SUCCESS',
);
export const servicePostFailure = actionCreator<{action: any; error: any}>(
  'SERVICE_POST_FAILURE',
);