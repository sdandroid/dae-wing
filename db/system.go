/*
 * SPDX-License-Identifier: AGPL-3.0-only
 * Copyright (c) 2023, daeuniverse Organization <team@v2raya.org>
 */

package db

import "time"

type System struct {
	ID                     uint      `gorm:"primaryKey;autoIncrement"`
	Running                bool      `gorm:"not null;default:false"`
	RunningConfigVersion   uint      `gorm:"not null;default:0"`
	RunningDnsVersion      uint      `gorm:"not null;default:0"`
	RunningRoutingVersion  uint      `gorm:"not null;default:0"`
	RunningGroupVersionSum uint      `gorm:"not null;default:0"`
	RunningGroupIds        string    `gorm:"not null;default:''"`
	RunningAt              time.Time `gorm:"null"`

	// Foreign keys.
	RunningConfigID  *uint
	RunningConfig    *Config
	RunningDnsID     *uint
	RunningDns       *Dns
	RunningRoutingID *uint
	RunningRouting   *Routing
	RunningGroups    []Group
}
