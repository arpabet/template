/*
 * Copyright (c) 2023 Zander Schwid & Co. LLC.
 * SPDX-License-Identifier: BUSL-1.1
 */

package service

import "github.com/pkg/errors"

var (

	ErrNotImplemented = errors.New("not implemented")

	ErrIntegrityDB = errors.New("db integrity")

	ErrUserAlreadyExist = errors.New("user already exist")
	ErrUserNotFound = errors.New("user not found")
	ErrUserInvalidPassword = errors.New("wrong password")

	ErrInvalidRecoverCode = errors.New("invalid recover code")

	ErrPageNotFound = errors.New("page not found")
)


