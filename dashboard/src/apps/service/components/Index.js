import React, {Component} from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';

import { withStyles } from '@material-ui/core/styles';

import actions from '../../../actions'
import IndexTable from './tables/IndexTable'
import RoutePanels from './panels/RoutePanels'
import RouteTabs from './tabs/RouteTabs'


const styles = theme => ({
  nested: {
    paddingLeft: theme.spacing.unit * 4,
  },
});

function renderIndex(routes, data, index) {
  if (!index) {
    return <div>Not Found</div>
  }
  console.log("DEBUG: Index.renderIndex: ", index.Kind, index.Name)
  switch(index.Kind) {
    case "Msg":
      return <div>{index.Name}</div>
    case "RoutePanels":
      return <RoutePanels render={renderIndex} routes={routes} data={data} index={index} />
    case "RouteTabs":
      return <RouteTabs render={renderIndex} routes={routes} data={data} index={index} />
    case "Table":
      return <IndexTable routes={routes} columns={index.Columns} data={data[index.DataKey]} />
    default:
      return <div>Unsupported Kind: {index.Kind}</div>
  }
}

class Index extends Component {
  componentWillMount() {
    console.log("Index.componentWillMount")
    const {match, getIndex} = this.props
    getIndex(match.params)
  }


  render() {
    const {match, service, serviceName, projectName, getIndex} = this.props
		console.log("Index.reder", projectName, serviceName)

    if (service.serviceName !== serviceName || service.projectName !== projectName) {
      getIndex(match.params)
      return null
    }

    let state = null
    if (projectName) {
      state = service.projectServiceMap[projectName][serviceName]
    } else {
      state = service.serviceMap[serviceName]
    }

    if (state.isFetching) {
      return <div>Fetching...</div>
    }

    const routes = [this.props]
    let html = renderIndex(routes, state.Data, state.Index)

    return (
      <div>
        { html }
      </div>
    );
  }
}

Index.propTypes = {
  classes: PropTypes.object.isRequired,
  match: PropTypes.object.isRequired,
  auth: PropTypes.object.isRequired,
  service: PropTypes.object.isRequired,
  serviceName: PropTypes.string.isRequired,
  projectName: PropTypes.string.isRequired,
}

function mapStateToProps(state, ownProps) {
  const match = ownProps.match
  const auth = state.auth
  const service = state.service

  return {
    match: match,
    auth: auth,
    service: service,
    serviceName: match.params.service,
    projectName: match.params.project,
  }
}

function mapDispatchToProps(dispatch, ownProps) {
  return {
    getIndex: (params) => {
      dispatch(actions.service.serviceGetIndex(params));
    }
  }
}

export default connect(
  mapStateToProps,
  mapDispatchToProps,
)(withStyles(styles)(Index))