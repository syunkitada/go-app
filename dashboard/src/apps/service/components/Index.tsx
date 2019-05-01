import * as React from 'react';
import {connect} from 'react-redux';

import actions from '../../../actions';
import IndexForm from './forms/IndexForm';
import IndexTable from './tables/IndexTable';
import RoutePanels from './panels/RoutePanels';
import Tabs from './tabs/Tabs';
import RequestErrSnackbar from './snackbars/RequestErrSnackbar';
import GetMsgSnackbar from './snackbars/GetMsgSnackbar';
import SubmitMsgSnackbar from './snackbars/SubmitMsgSnackbar';

import logger from '../../../lib/logger';

function renderIndex(routes, data, index) {
  if (!index) {
    return <div>Not Found</div>;
  }
  logger.info(['Index', 'renderIndex:', index.Kind, index.Name]);
  switch (index.Kind) {
    case 'Msg':
      return <div>{index.Name}</div>;
    case 'RoutePanels':
      return (
        <RoutePanels
          render={renderIndex}
          routes={routes}
          data={data}
          index={index}
        />
      );
    case 'RouteTabs':
      return (
        <Tabs render={renderIndex} routes={routes} data={data} index={index} />
      );
    case 'Table':
      return <IndexTable routes={routes} index={index} data={data} />;
    case 'Form':
      return <IndexForm routes={routes} index={index} data={data} />;
    default:
      return <div>Unsupported Kind: {index.Kind}</div>;
  }
}

interface IIndex {
  match;
  service;
  serviceName;
  projectName;
  getIndex;
}

class Index extends React.Component<IIndex> {
  state = {
    openAlertSnackbar: true,
    traceMsgMap: {},
  };

  componentWillMount() {
    logger.info(['Index', 'componentWillMount()']);
    const {match, getIndex} = this.props;
    getIndex(match.params);
  }

  handleCloseAlertSnackbar = () => {
    this.setState({openAlertSnackbar: false});
  };

  render() {
    const {match, service, serviceName, projectName, getIndex} = this.props;
    logger.info(['Index', 'render', projectName, serviceName]);

    if (
      service.serviceName !== serviceName ||
      service.projectName !== projectName
    ) {
      getIndex(match.params);
      return null;
    }

    let state: any = null;
    if (projectName) {
      state = service.projectServiceMap[projectName][serviceName];
    } else {
      state = service.serviceMap[serviceName];
    }

    if (state.isFetching) {
      return <div>Fetching...</div>;
    }

    console.log('DEBUG HOGElwlwlw');

    const routes = [this.props];
    let html = renderIndex(routes, state.Data, state.Index);

    return (
      <div>
        {html}
        <RequestErrSnackbar />
        <GetMsgSnackbar />
        <SubmitMsgSnackbar />
      </div>
    );
  }
}

function mapStateToProps(state, ownProps) {
  const match = ownProps.match;
  const auth = state.auth;
  const service = state.service;

  return {
    match: match,
    auth: auth,
    service: service,
    serviceName: match.params.service,
    projectName: match.params.project,
  };
}

function mapDispatchToProps(dispatch, ownProps) {
  return {
    getIndex: params => {
      dispatch(actions.service.serviceGetIndex(params));
    },
  };
}

export default connect(
  mapStateToProps,
  mapDispatchToProps,
)(Index);