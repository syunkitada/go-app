import * as React from 'react';
// import PropTypes from 'prop-types';

import classNames from 'classnames';
import CheckCircleIcon from '@material-ui/icons/CheckCircle';
import ErrorIcon from '@material-ui/icons/Error';
import InfoIcon from '@material-ui/icons/Info';
import CloseIcon from '@material-ui/icons/Close';
import green from '@material-ui/core/colors/green';
import amber from '@material-ui/core/colors/amber';
import IconButton from '@material-ui/core/IconButton';
import Snackbar from '@material-ui/core/Snackbar';
import SnackbarContent from '@material-ui/core/SnackbarContent';
import WarningIcon from '@material-ui/icons/Warning';
import {withStyles} from '@material-ui/core/styles';

const variantIcon = {
  success: CheckCircleIcon,
  warning: WarningIcon,
  error: ErrorIcon,
  info: InfoIcon,
};

const styles = theme => ({
  success: {
    backgroundColor: green[600],
  },
  error: {
    backgroundColor: theme.palette.error.dark,
  },
  info: {
    backgroundColor: theme.palette.primary.dark,
  },
  warning: {
    backgroundColor: amber[700],
  },
  icon: {
    fontSize: 20,
  },
  iconVariant: {
    opacity: 0.9,
    marginRight: theme.spacing.unit,
  },
  message: {
    display: 'flex',
    alignItems: 'center',
  },
});

interface IMsgSnackbar {
  classes;
  open;
  onClose;
  variant;
  vertical;
  horizontal;
  msg;
}

class MsgSnackbar extends React.Component<IMsgSnackbar> {
  render() {
    const {
      classes,
      className,
      open,
      onClose,
      variant,
      vertical,
      horizontal,
      msg,
    } = this.props;

    const Icon = variantIcon[variant];

    return (
      <Snackbar
        anchorOrigin={{vertical, horizontal}}
        open={open}
        onClose={onClose}
        ContentProps={{
          'aria-describedby': 'message-id',
        }}>
        <SnackbarContent
          className={classNames(classes[variant], className)}
          aria-describedby="client-snackbar"
          message={
            <span id="client-snackbar" className={classes.message}>
              <Icon className={classNames(classes.icon, classes.iconVariant)} />
              {msg}
            </span>
          }
          action={[
            <IconButton
              key="close"
              aria-label="Close"
              color="inherit"
              className={classes.close}
              onClick={onClose}>
              <CloseIcon className={classes.icon} />
            </IconButton>,
          ]}
        />
      </Snackbar>
    );
  }
}

// MsgSnackbar.propTypes = {
//   classes: PropTypes.object.isRequired,
//   className: PropTypes.string,
//   open: PropTypes.bool,
//   onClose: PropTypes.func,
//   variant: PropTypes.string,
//   vartical: PropTypes.string,
//   horizontal: PropTypes.string,
//   msg: PropTypes.string,
// };

export default withStyles(styles)(MsgSnackbar);
