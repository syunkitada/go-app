import React, {Component} from 'react';
import { connect } from 'react-redux';

import { withStyles } from '@material-ui/core/styles';
import Toolbar from '@material-ui/core/Toolbar';
import { lighten } from '@material-ui/core/styles/colorManipulator';

import { fade } from '@material-ui/core/styles/colorManipulator';

import Input from '@material-ui/core/Input';
import InputAdornment from '@material-ui/core/InputAdornment';
import FormControl from '@material-ui/core/FormControl';
import SearchIcon from '@material-ui/icons/Search';

import Grid from '@material-ui/core/Grid';


import TablePagination from './TablePagination'

class TableToolbar extends Component {
  render() {
    const { classes, onChangeSearchInput, count, rowsPerPage, page, onChangePage, onChangeRowsPerPage } = this.props

    return (
      <Toolbar>
        <Grid container justify="space-between" spacing={24}>
          <Grid item>
            <div>
              <FormControl className={classes.margin}>
                <Input
                  id="input-with-icon-adornment"
                  placeholder="Search"
                  onChange={onChangeSearchInput}
                  startAdornment={
                    <InputAdornment position="start">
                      <SearchIcon />
                    </InputAdornment>
                  }
                />
              </FormControl>
            </div>
          </Grid>

          <Grid item>
            <TablePagination
              count={count}
              rowsPerPage={rowsPerPage}
              page={page}
              onChangePage={onChangePage}
              onChangeRowsPerPage={onChangeRowsPerPage}
            />
          </Grid>

        </Grid>
      </Toolbar>
    )
  }
}

function mapStateToProps(state, ownProps) {
  return {}
}

function mapDispatchToProps(dispatch, ownProps) {
  return {}
}

const style = theme => ({
  root: {
    width: '100%',
  },
  table: {
    minWidth: 1020,
  },
  tableWrapper: {
    overflowX: 'auto',
  },
  margin: {
    margin: theme.spacing.unit * 2,
  },
  highlight:
    theme.palette.type === 'light'
      ? {
          color: theme.palette.secondary.main,
          backgroundColor: lighten(theme.palette.secondary.light, 0.85),
        }
      : {
          color: theme.palette.text.primary,
          backgroundColor: theme.palette.secondary.dark,
        },
  spacer: {
    flex: '1 1 100%',
  },
  actions: {
    color: theme.palette.text.secondary,
  },
  title: {
    flex: '0 0 auto',
  },
  search: {
    position: 'relative',
    borderRadius: theme.shape.borderRadius,
    backgroundColor: fade(theme.palette.common.white, 0.15),
    '&:hover': {
      backgroundColor: fade(theme.palette.common.white, 0.25),
    },
    marginRight: theme.spacing.unit * 2,
    marginLeft: 0,
    width: '100%',
    [theme.breakpoints.up('sm')]: {
      marginLeft: theme.spacing.unit * 3,
      width: 'auto',
    },
  },
  searchIcon: {
    width: theme.spacing.unit * 9,
    height: '100%',
    position: 'absolute',
    pointerEvents: 'none',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
  },
});

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(withStyles(style, {withTheme: true})(TableToolbar));