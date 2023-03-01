/*
 * SPDX-License-Identifier: AGPL-3.0-only
 * Copyright (c) 2023, v2rayA Organization <team@v2raya.org>
 */

package db

import (
	"fmt"
	"time"
)

var (
	BadLinkFormatError = fmt.Errorf("not a valid link")
)

type Subscription struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	UpdatedAt time.Time `gorm:"not null"`
	Link      string    `gorm:"not null"`
	Status    string    `gorm:"not null"` // Latency, error info, etc.
	Info      string    `gorm:"not null"` // Maybe include some info from provider

	Tag *string `gorm:"unique"`

	Node []Node
}
