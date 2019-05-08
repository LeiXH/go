package server

import (
	"fmt"
	"pkg/services/printer"

	"pkg/database"
	"pkg/logger"
)

type Server struct {
	config Option
}


func (srv *Server) Prepare() (err error) {
	if err = database.InitMySQLConn(&srv.config.MysqlConfig); err != nil {
		return
	}
	if !srv.config.PrinterConfig.Enabled {
		logger.Info("printer disabled by user")
		return
	}
	if err = printer.InitPrinter(srv.config.PrinterConfig); err != nil {
		return
	}
	return
}

func New(configFile string) error {
	ss := &Server{}
	if err := ss.loadConfig(configFile); err != nil {
		return fmt.Errorf("loading server config get err %s", err)
	}

	ss.Prepare()

	//defer ss.GracefulShutdown()

	return nil
}


func (srv *Server) GracefulShutdown() {
	database.CloseConn()
	printer.ClosePrinter()
}
