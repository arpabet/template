/*
 * Copyright (c) 2025 Karagatan LLC.
 * SPDX-License-Identifier: BUSL-1.1
 */

package service

import "golang.org/x/xerrors"

var (

	ErrNotImplemented = xerrors.New("not implemented")

	ErrIntegrityDB = xerrors.New("db integrity")

	ErrUserAlreadyExist = xerrors.New("user already exist")
	ErrUserNotFound = xerrors.New("user not found")
	ErrUserInvalidPassword = xerrors.New("wrong password")

	ErrInvalidRecoverCode = xerrors.New("invalid recover code")

	ErrPageNotFound = xerrors.New("page not found")
)


