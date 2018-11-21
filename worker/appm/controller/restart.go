// RAINBOND, Application Management Platform
// Copyright (C) 2014-2017 Goodrain Co., Ltd.

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version. For any non-GPL usage of Rainbond,
// one or multiple Commercial Licenses authorized by Goodrain Co., Ltd.
// must be obtained first.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package controller

import (
	"fmt"
	"sync"

	"github.com/Sirupsen/logrus"
	v1 "github.com/goodrain/rainbond/worker/appm/types/v1"
)

type restartController struct {
	stopChan     chan struct{}
	controllerID string
	appService   []v1.AppService
	manager      *Manager
}

func (s *restartController) Begin() {
	var wait sync.WaitGroup
	for _, service := range s.appService {
		go func(service v1.AppService) {
			wait.Add(1)
			defer wait.Done()
			service.Logger.Info("App runtime begin restart app service "+service.ServiceAlias, getLoggerOption("starting"))
			if err := s.restartOne(service); err != nil {
				logrus.Errorf("restart service %s failure %s", service.ServiceAlias, err.Error())
			} else {
				service.Logger.Info(fmt.Sprintf("restart service %s success", service.ServiceAlias), getLastLoggerOption())
			}
		}(service)
	}
	wait.Wait()
	s.manager.callback(s.controllerID, nil)
}
func (s *restartController) restartOne(app v1.AppService) error {
	stopController := stopController{}
	if err := stopController.stopOne(app); err != nil {
		app.Logger.Error("(Restart)Stop app failure %s,you could waiting stoped and manual start it", getCallbackLoggerOption())
		return err
	}
	startController := startController{}
	if err := startController.startOne(app); err != nil {
		app.Logger.Error("(Restart)Start app failure %s,you could waiting it start success.or manual stop it", getCallbackLoggerOption())
		return err
	}
	return nil
}
func (s *restartController) Stop() error {
	close(s.stopChan)
	return nil
}
