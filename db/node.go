/*
 * SPDX-License-Identifier: AGPL-3.0-only
 * Copyright (c) 2023, daeuniverse Organization <team@v2raya.org>
 */

package db

import "time"

type Node struct {
	ID       uint   `gorm:"primaryKey;autoIncrement"`
	Link     string `gorm:"not null"`
	Name     string `gorm:"not null"`
	Address  string `gorm:"not null"`
	Protocol string `gorm:"not null"`

	Tag       *string   `gorm:"unique"`
	UpdatedAt time.Time `gorm:"null"`

	// Foreign keys.
	// Nil SubscriptionID indicates nodes belonging to no subscription.
	SubscriptionID *uint
	Subscription   *Subscription
}
