import React, {Component} from 'react';
import { connect } from 'react-redux';

import Dashboard from '../../components/Dashboard'
import Index from './components/Index'


class Service extends Component {
  render() {
    const {match, auth} = this.props

    if (!auth.user) {
      return null
    }

    return (
      <Dashboard match={match}>
        <h2>{match.params.service}</h2>
        <Index match={match} />
      </Dashboard>
    );
  }
}

function mapStateToProps(state, ownProps) {
  const auth = state.auth
  const match = ownProps.match

  return {
    match: match,
    auth: auth,
  }
}

function mapDispatchToProps(dispatch, ownProps) {
  return {}
}

export default connect(
  mapStateToProps,
  mapDispatchToProps,
)(Service)
