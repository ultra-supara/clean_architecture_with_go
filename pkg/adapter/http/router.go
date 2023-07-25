package http

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/yagikota/clean_architecture_wtih_go/pkg/domain/service"
	"github.com/yagikota/clean_architecture_wtih_go/pkg/infra"
	"github.com/yagikota/clean_architecture_wtih_go/pkg/infra/mysql"
	"github.com/yagikota/clean_architecture_wtih_go/pkg/usecase"
)

const (
	apiVersion      = "/v1"
	healthCheckRoot = "/health_check"
	// student
	studentsAPIRoot = apiVersion + "/students"
	studentIDParam  = "student_id"
)

func InitRouter() *echo.Echo {
	e := echo.New()
	e.Use(
		middleware.Logger(),
		middleware.Recover(),
	)

	// health check
	healthCheckGroup := e.Group(healthCheckRoot)
	{
		relativePath := ""
		healthCheckGroup.GET(relativePath, healthCheck)
	}

	// student , DI ,各層の情報を集約してきてhandlingしやすくさせるための準備
	mySQLConn := infra.NewMySQLConnector()
	studentRepository := mysql.NewStudentRepository(mySQLConn.Conn)
	studentService := service.NewStudentService(studentRepository)
	studentUsecase := usecase.NewUserUsecase(studentService)

	studentGroup := e.Group(studentsAPIRoot)
	{
		handler := NewStudentHandler(studentUsecase)
		// v1/students
		relativePath := ""
		studentGroup.GET(relativePath, handler.FindAllStudents())

		// v1/students/:student_id
		relativePath = "/:" + studentIDParam
		studentGroup.GET(relativePath, handler.FindStudentByID())
	}

	return e
}
