package model

import (
	"myorm/log"
	"myorm/session"
)

type Account struct {
	ID int `myorm:"PRIMARY KEY"`
	Password string
}

type ForBenchmarkTest interface {
	ForBenchmarkTest() error
}

func (a *Account) ForBenchmarkTest() error {
	// 浪费一些时间
	count := 0
	for i := 0; i < 1000; i++ {
		count++
	}
	for i := 0; i < 1000; i++ {
		count--
	}
	return nil
}

func (a *Account) AfterQuery(s *session.Session) error {
	log.Info("after query", a)
	a.Password = "******"
	return nil
}