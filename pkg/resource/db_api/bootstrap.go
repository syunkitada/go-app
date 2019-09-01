package db_api

import (
	"fmt"

	"github.com/jinzhu/gorm"

	"github.com/syunkitada/goapp/pkg/base/base_db_model"
	"github.com/syunkitada/goapp/pkg/lib/exec_utils"
	"github.com/syunkitada/goapp/pkg/lib/logger"
)

func (api *Api) Bootstrap(tctx *logger.TraceContext, isRecreate bool) (err error) {
	startTime := logger.StartTrace(tctx)
	defer func() { logger.EndTrace(tctx, startTime, err, 0) }()
	if err = exec_utils.CreateDatabase(tctx, api.baseConf, api.databaseConf.Connection, isRecreate); err != nil {
		return err
	}

	var db *gorm.DB
	db, err = api.Open(tctx)
	if err != nil {
		return err
	}
	defer api.Close(tctx, db)

	if err = db.AutoMigrate(&base_db_model.User{}).Error; err != nil {
		return err
	}
	if err = db.AutoMigrate(&base_db_model.Role{}).Error; err != nil {
		return err
	}
	if err = db.AutoMigrate(&base_db_model.Project{}).Error; err != nil {
		return err
	}
	if err = db.AutoMigrate(&base_db_model.ProjectRole{}).Error; err != nil {
		return err
	}
	if err = db.AutoMigrate(&base_db_model.Service{}).Error; err != nil {
		return err
	}

	for _, projectRole := range api.appConf.Auth.DefaultProjectRoles {
		if err = api.CreateProjectRole(tctx, db, projectRole.Name); err != nil {
			return err
		}
		fmt.Printf("Created ProjectRole: %s\n", projectRole.Name)
	}

	for _, project := range api.appConf.Auth.DefaultProjects {
		if err = api.CreateProject(tctx, db, project.Name, project.ProjectRole); err != nil {
			return err
		}
		fmt.Printf("Created Project: %s\n", project.Name)
	}

	for _, role := range api.appConf.Auth.DefaultRoles {
		if err = api.CreateRole(tctx, db, role.Name, role.Project); err != nil {
			return err
		}
		fmt.Printf("Created Role: %s\n", role.Name)
	}

	for _, user := range api.appConf.Auth.DefaultUsers {
		if err = api.CreateUser(tctx, db, user.Name, user.Password); err != nil {
			return err
		}
		for _, roleName := range user.Roles {
			if err = api.AssignRoleToUser(tctx, db, roleName, user.Name); err != nil {
				return err
			}
		}

		fmt.Printf("Created User: %s\n", user.Name)
	}

	for _, service := range api.appConf.Auth.DefaultServices {
		if err = api.CreateOrUpdateService(tctx, db, &service); err != nil {
			return err
		}
	}

	return nil
}