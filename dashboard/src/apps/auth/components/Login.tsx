import * as React from 'react';
import {connect} from 'react-redux';
import {Redirect} from 'react-router-dom';

import {Theme} from '@material-ui/core/styles/createMuiTheme';
import createStyles from '@material-ui/core/styles/createStyles';
import withStyles, {
  StyleRules,
  WithStyles,
} from '@material-ui/core/styles/withStyles';

import Avatar from '@material-ui/core/Avatar';
import Button from '@material-ui/core/Button';
import CssBaseline from '@material-ui/core/CssBaseline';
import FormControl from '@material-ui/core/FormControl';
import Input from '@material-ui/core/Input';
import InputLabel from '@material-ui/core/InputLabel';
import Paper from '@material-ui/core/Paper';
import Typography from '@material-ui/core/Typography';

import LockOutlinedIcon from '@material-ui/icons/LockOutlined';

import actions from '../../../actions';
import MsgSnackbar from '../../../components/snackbars/MsgSnackbar';

const styles = (theme: Theme): StyleRules =>
  createStyles({
    avatar: {
      backgroundColor: theme.palette.secondary.main,
      margin: theme.spacing.unit,
    },
    form: {
      marginTop: theme.spacing.unit,
      width: '100%', // Fix IE11 issue.
    },
    layout: {
      display: 'block', // Fix IE11 issue.
      marginLeft: theme.spacing.unit * 3,
      marginRight: theme.spacing.unit * 3,
      width: 'auto',
      [theme.breakpoints.up(400 + theme.spacing.unit * 3 * 2)]: {
        marginLeft: 'auto',
        marginRight: 'auto',
        width: 400,
      },
    },
    paper: {
      alignItems: 'center',
      display: 'flex',
      flexDirection: 'column',
      marginTop: theme.spacing.unit * 8,
      padding: `${theme.spacing.unit * 2}px ${theme.spacing.unit * 3}px ${theme
        .spacing.unit * 3}px`,
    },
    submit: {
      marginTop: theme.spacing.unit * 3,
    },
  });

interface ILogin extends WithStyles<typeof styles> {
  auth;
  history;
  onSubmit;
}

class Login extends React.Component<ILogin> {
  public render() {
    const {classes, auth, history, onSubmit} = this.props;
    const {from} = history.location.state || {from: {pathname: '/'}};

    if (auth.redirectToReferrer) {
      return <Redirect to={from} />;
    }

    if (auth.user) {
      return <Redirect to={{pathname: '/'}} />;
    }

    if (auth.isFetching) {
      return <div>During authentication</div>;
    }

    let msgHtml: any = null;
    if (auth.error != null && auth.error !== '') {
      const variant = 'error';
      const vertical = 'bottom';
      const horizontal = 'left';

      msgHtml = (
        <MsgSnackbar
          open={true}
          onClose={this.handleClose}
          vertical={vertical}
          horizontal={horizontal}
          variant={variant}
          msg={auth.error}
        />
      );
    }

    return (
      <React.Fragment>
        <CssBaseline />
        {msgHtml}
        <main className={classes.layout}>
          <Paper className={classes.paper}>
            <Avatar className={classes.avatar}>
              <LockOutlinedIcon />
            </Avatar>
            <Typography variant="headline">Sign in</Typography>
            <form className={classes.form} onSubmit={onSubmit}>
              <FormControl margin="normal" required={true} fullWidth={true}>
                <InputLabel htmlFor="name">Name</InputLabel>
                <Input id="name" name="name" autoFocus={true} />
              </FormControl>
              <FormControl margin="normal" required={true} fullWidth={true}>
                <InputLabel htmlFor="password">Password</InputLabel>
                <Input
                  name="password"
                  type="password"
                  id="password"
                  autoComplete="current-password"
                />
              </FormControl>
              <Button
                type="submit"
                fullWidth={true}
                variant="raised"
                color="primary"
                className={classes.submit}>
                Sign in
              </Button>
            </form>
          </Paper>
        </main>
      </React.Fragment>
    );
  }

  private handleClose = (event, reason) => {
    return;
  };
}

function mapStateToProps(state, ownProps) {
  const auth = state.auth;

  return {auth};
}

function mapDispatchToProps(dispatch, ownProps) {
  return {
    onSubmit: e => {
      e.preventDefault();
      const name = e.target.name.value.trim();
      const password = e.target.password.value.trim();
      dispatch(actions.auth.authLogin(name, password));
    },
  };
}

export default connect(
  mapStateToProps,
  mapDispatchToProps,
)(withStyles(styles)(Login));