import React, {Component} from 'react'
import { connect } from 'react-redux'
import { Route, Redirect } from 'react-router-dom';

import logger from '../lib/logger';

class AuthRoute extends Component {
  render() {
    const { component: Component, auth, ...rest } = this.props
    logger.info('AuthRoute', 'render()')
    return (
      <Route {...rest}
        render={props =>
          auth.user ? (
            <Component {...props} />
          ) : (
            <Redirect
              to={{
                pathname: '/login',
                state: { from: props.location }
              }}
            />
          )
        }
      />
    );
  }
}

function mapStateToProps(state, ownProps) {
  const auth = state.auth

  return {
    auth: auth,
  }
}

export default connect(
  mapStateToProps,
)(AuthRoute)