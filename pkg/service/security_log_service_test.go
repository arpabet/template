/*
 * Copyright (c) 2023 Zander Schwid & Co. LLC.
 * SPDX-License-Identifier: BUSL-1.1
 */

package service_test

import (
	"github.com/stretchr/testify/require"
	"github.com/sprintframework/template/pkg/service"
	"sync"
	"testing"
	"time"
)

func TestDateFormatter(t *testing.T) {

	var el eventList
	err := el.addEvent()
	require.NoError(t, err)
	err = el.addEvent()
	require.NoError(t, err)

	cnt := 0
	el.eventMap.Range(func(key, value interface{}) bool {
		//fmt.Printf("key = %v\n", key)
		cnt++
		return true
	})

	require.Equal(t, 2, cnt)

}

type eventList struct {
	eventMap sync.Map
}

func (t *eventList) addEvent() error {

	current := time.Now()
tryAgain:
	utc := current.UTC()

	if has, err := t.hasEvent(utc); err != nil {
		return err
	} else if has {
		current = current.Add(time.Millisecond)
		//fmt.Printf("current %v\n", current)
		goto tryAgain
	}

	key := utc.Format(service.DDMMYYYYhhmmss)
	t.eventMap.Store(key, true)
	return nil
}

func (t *eventList) hasEvent(utc time.Time) (bool, error) {
	key := utc.Format(service.DDMMYYYYhhmmss)
	if _ , ok := t.eventMap.Load(key); ok {
		return true, nil
	}
	return false, nil
}