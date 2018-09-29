package model_api

import (
	"encoding/hex"
	"errors"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/scrypt"

	"github.com/syunkitada/goapp/pkg/authproxy/model"
)

func CreateUser(name string, password string) error {
	db, err := gorm.Open("mysql", Conf.AuthproxyDatabase.Connection)
	defer db.Close()
	if err != nil {
		return err
	}

	var user model.User

	if err := db.Debug().Where("name = ?", name).First(&user).Error; err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return err
		}

		hashedPassword, hashedErr := GenerateHashFromPassword(name, password)
		if hashedErr != nil {
			return hashedErr
		}

		user = model.User{
			Name:     name,
			Password: hashedPassword,
		}
		db.Debug().Create(&user)

		return nil
	}

	return nil
}

func CreateRole(name string, projectName string) error {
	db, err := gorm.Open("mysql", Conf.AuthproxyDatabase.Connection)
	defer db.Close()
	if err != nil {
		return err
	}

	var role model.Role
	var project model.Project

	if err := db.Debug().First(&project, "name = ?", projectName).Error; err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return err
		}
	}

	if err := db.Debug().Where("name = ?", name).First(&role).Error; err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return err
		}

		role = model.Role{
			Name:      name,
			ProjectID: project.ID,
		}
		db.Debug().Create(&role)

		return nil
	}

	return nil
}

func AssignRole(userName string, roleName string) error {
	db, err := gorm.Open("mysql", Conf.AuthproxyDatabase.Connection)
	defer db.Close()
	if err != nil {
		return err
	}

	var user model.User
	var role model.Role

	db.Debug().Where("name = ?", roleName).First(&role)

	db.Debug().Preload("Roles").First(&user, "name = ?", userName)
	db.Debug().Model(&user).Association("Roles").Append(&role)
	return nil
}

func CreateProject(name string, projectRoleName string) error {
	db, err := gorm.Open("mysql", Conf.AuthproxyDatabase.Connection)
	defer db.Close()
	if err != nil {
		return err
	}

	var project model.Project
	var projectRole model.ProjectRole

	if err := db.Debug().First(&projectRole, "name = ?", projectRoleName).Error; err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return err
		}
	}

	if err := db.Debug().Where("name = ?", name).First(&project).Error; err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return err
		}

		project = model.Project{
			Name:          name,
			ProjectRoleID: projectRole.ID,
		}
		db.Debug().Create(&project)

		return nil
	}

	return nil
}

func CreateProjectRole(name string) error {
	db, err := gorm.Open("mysql", Conf.AuthproxyDatabase.Connection)
	defer db.Close()
	if err != nil {
		return err
	}

	var projectRole model.ProjectRole

	if err := db.Debug().Where("name = ?", name).First(&projectRole).Error; err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return err
		}

		projectRole = model.ProjectRole{
			Name: name,
		}
		db.Debug().Create(&projectRole)

		return nil
	}

	return nil

}

func AssignProjectRole(projectName string, projectRoleName string) error {
	db, err := gorm.Open("mysql", Conf.AuthproxyDatabase.Connection)
	defer db.Close()
	if err != nil {
		return err
	}

	var project model.Project
	var projectRole model.ProjectRole

	db.Debug().Where("name = ?", projectRoleName).First(&projectRole)

	db.Debug().Preload("ProjectRoles").First(&project, "name = ?", projectName)
	db.Debug().Model(&project).Association("ProjectRoles").Append(&projectRole)
	return nil
}

func CreateService(name string, scope string) error {
	db, err := gorm.Open("mysql", Conf.AuthproxyDatabase.Connection)
	defer db.Close()
	if err != nil {
		return err
	}

	var service model.Service

	if err := db.Debug().Where("name = ?", name).First(&service).Error; err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return err
		}

		service = model.Service{
			Name:  name,
			Scope: scope,
		}
		db.Debug().Create(&service)

		return nil
	}

	return nil
}

func AssignService(projectRoleName string, serviceName string) error {
	db, err := gorm.Open("mysql", Conf.AuthproxyDatabase.Connection)
	defer db.Close()
	if err != nil {
		return err
	}

	var projectRole model.ProjectRole
	var service model.Service

	db.Debug().Where("name = ?", serviceName).First(&service)

	db.Debug().Preload("Services").First(&projectRole, "name = ?", projectRoleName)
	db.Debug().Model(&projectRole).Association("Services").Append(&service)

	return nil
}

func GetAuthUser(authRequest *model.AuthRequest) (*model.User, error) {
	db, err := gorm.Open("mysql", Conf.AuthproxyDatabase.Connection)
	defer db.Close()
	if err != nil {
		return nil, err
	}

	var users []model.User
	if err := db.Debug().Where("name = ?", authRequest.Username).Find(&users).Error; err != nil {
		return nil, err
	}

	if len(users) != 1 {
		return nil, errors.New("Invalid User")
	}

	hashedPassword, hashedErr := GenerateHashFromPassword(authRequest.Username, authRequest.Password)
	if hashedErr != nil {
		return nil, hashedErr
	}

	user := users[0]
	if user.Password != hashedPassword {
		return nil, errors.New("Invalid Password")
	}

	return &user, nil
}

func GenerateHashFromPassword(username string, password string) (string, error) {
	converted, err := scrypt.Key([]byte(password), []byte(Conf.Admin.Secret+username), 16384, 8, 1, 32)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(converted[:]), nil
}

func GetUserAuthority(username string) (*model.UserAuthority, error) {
	db, err := gorm.Open("mysql", Conf.AuthproxyDatabase.Connection)
	defer db.Close()
	if err != nil {
		return nil, err
	}

	var users []model.CustomUser
	if err := db.Debug().Raw(sqlSelectUser+"WHERE u.name = ?", username).Scan(&users).Error; err != nil {
		return nil, err
	}

	serviceMap := map[string]bool{}
	projectServiceMap := map[string]model.ProjectService{}
	for _, user := range users {
		switch user.ServiceScope {
		case "user":
			serviceMap[user.ServiceName] = true
		case "project":
			glog.Info(user)
			if projectService, ok := projectServiceMap[user.ProjectName]; ok {
				projectService.ServiceMap[user.ServiceName] = true
			} else {
				projectService := model.ProjectService{
					RoleName:        user.RoleName,
					ProjectName:     user.ProjectName,
					ProjectRoleName: user.ProjectRoleName,
					ServiceMap:      map[string]bool{},
				}
				projectService.ServiceMap[user.ServiceName] = true
				projectServiceMap[user.ProjectName] = projectService
			}
		}
	}

	userAuthority := model.UserAuthority{
		ServiceMap:        serviceMap,
		ProjectServiceMap: projectServiceMap,
	}

	return &userAuthority, nil
}
