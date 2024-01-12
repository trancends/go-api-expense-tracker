package delivery

import (
	"database/sql"
	"expense-tracker/config"
	"expense-tracker/delivery/controller"
	"expense-tracker/repository"
	"expense-tracker/usecase"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
)

type Server struct {
	expenseUC usecase.ExpenseUsecase
	engine    *gin.Engine
	host      string
}

func (s *Server) initRoute() {
	rg := s.engine.Group("/api/v1")
	controller.NewExpenseHandler(s.expenseUC, rg).Route()
}

func (s *Server) Run() {
	s.initRoute()
	if err := s.engine.Run(s.host); err != nil {
		panic(fmt.Errorf("failed to start server %v", err))
	}
}

func NewServer() *Server {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Welcome to the Todo APP")
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database)

	db, err := sql.Open(cfg.Driver, psqlInfo)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to Database")
	if err != nil {
		log.Fatal(fmt.Errorf("config error: %v", err))
	}

	expenseRepository := repository.NewExpenseRepository(db)
	expenseUseCase := usecase.NewExpenseUsecase(expenseRepository)

	engine := gin.Default()
	host := fmt.Sprintf(":%s", cfg.ApiPort)

	return &Server{
		expenseUC: expenseUseCase,
		engine:    engine,
		host:      host,
	}
}
