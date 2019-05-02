import fetch from 'cross-fetch';

interface IResponse {
  payload: any;
  error: any;
}

function post({
  serviceName,
  actionName,
  projectName,
  queries,
}): Promise<IResponse> {
  const body = JSON.stringify({
    Action: {
      Name: actionName,
      ProjectName: projectName,
      Queries: queries,
      ServiceName: serviceName,
    },
  });

  return fetch(process.env.REACT_APP_AUTHPROXY_URL + '/' + serviceName, {
    body,
    credentials: 'include',
    method: 'POST',
    mode: 'cors',
  })
    .then(resp => {
      if (!resp.ok) {
        return resp.json().then(payload => {
          const result: IResponse = {
            error: {
              err: payload.Err,
              errCode: resp.status,
            },
            payload: null,
          };
          return result;
        });
      }

      return resp.json().then(payload => {
        const result: IResponse = {
          error: null,
          payload,
        };
        return result;
      });
    })
    .catch(error => {
      const result: IResponse = {
        error: {
          err: error,
        },
        payload: null,
      };
      return result;
    });
}

export default {
  post,
};
