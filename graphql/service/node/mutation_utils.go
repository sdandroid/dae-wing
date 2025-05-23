/*
 * SPDX-License-Identifier: AGPL-3.0-only
 * Copyright (c) 2023, daeuniverse Organization <team@v2raya.org>
 */

package node

import (
	"context"
	"errors"
	"fmt"

	"github.com/daeuniverse/dae-wing/common"
	"github.com/daeuniverse/dae-wing/db"
	"github.com/daeuniverse/dae-wing/graphql/internal"
	"github.com/graph-gophers/graphql-go"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var DuplicatedError = fmt.Errorf("node already exists")

type ImportResult struct {
	Link  string
	Error *string
	Node  *Resolver
}

func importNode(d *gorm.DB, subscriptionId *uint, arg *internal.ImportArgument) (m *db.Node, err error) {
	if err = arg.ValidateTag(); err != nil {
		return nil, err
	}
	m, err = db.NewNodeModel(arg.Link, arg.Tag, subscriptionId)
	if err != nil {
		return nil, err
	}
	var existsNode db.Node
	d.Model(&db.Node{}).
		Where("name = ?", m.Name).
		Where("subscription_id = ?", subscriptionId).First(&existsNode)

	if existsNode.ID > 0 {
		if existsNode.Link != arg.Link {
			newModel, err := db.NewNodeModel(arg.Link, nil, nil)
			if err != nil {
				return nil, err
			}
			q := d.Model(&db.Node{}).Where("name =? and subscription_id = ? ", m.Name, subscriptionId).Updates(newModel)
			if err = q.Error; err != nil {
				return nil, err
			}
		}

	} else {
		if err = d.Create(m).Error; err != nil {
			return nil, err
		}
	}
	return m, nil
}

// Import nodes. If abortError is false, err will always be nil.
func Import(d *gorm.DB, abortError bool, subscriptionId *uint, argument []*internal.ImportArgument) (rs []*ImportResult, err error) {
	for _, arg := range argument {
		var m *db.Node
		if m, err = importNode(d, subscriptionId, arg); err != nil {
			if abortError && !errors.Is(err, DuplicatedError) {
				return nil, err
			}
			info := err.Error()
			rs = append(rs, &ImportResult{
				Link:  arg.Link,
				Error: &info,
				Node:  nil,
			})
			continue
		}
		rs = append(rs, &ImportResult{
			Link:  arg.Link,
			Error: nil,
			Node: &Resolver{
				Node: m,
			},
		})
	}
	return rs, nil
}

func Update(d *gorm.DB, _id graphql.ID, link string) (r *Resolver, err error) {
	id, err := common.DecodeCursor(_id)
	if err != nil {
		return nil, err
	}
	newModel, err := db.NewNodeModel(link, nil, nil)
	if err != nil {
		return nil, err
	}
	q := d.Model(&db.Node{ID: id}).Updates(newModel)
	if err = q.Error; err != nil {
		return nil, err
	}
	if q.RowsAffected == 0 {
		return nil, fmt.Errorf("no such node")
	}
	return &Resolver{
		Node: newModel,
	}, nil
}

func AutoUpdateVersionByIds(d *gorm.DB, ids []uint) (err error) {
	var sys db.System
	if err = d.Model(&db.System{}).
		FirstOrCreate(&sys).Error; err != nil {
		return err
	}
	if !sys.Running {
		return nil
	}

	if err = d.Exec(`update groups
                set version = groups.version + 1
                from groups g
                    inner join group_nodes
                    on g.system_id = ? and g.id = group_nodes.group_id and group_nodes.node_id in ?
				where g.id = groups.id`, sys.ID, ids).Error; err != nil {
		return err
	}

	return nil
}

func Remove(ctx context.Context, _ids []graphql.ID) (n int32, err error) {
	ids, err := common.DecodeCursorBatch(_ids)
	if err != nil {
		return 0, err
	}
	tx := db.BeginTx(ctx)
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	// Update modified if any nodes are referenced by running config.
	if err = AutoUpdateVersionByIds(tx, ids); err != nil {
		return 0, err
	}

	// Remove.
	q := tx.Where("id in ?", ids).
		Select(clause.Associations).
		Delete(&db.Node{})
	if q.Error != nil {
		return 0, q.Error
	}

	return int32(q.RowsAffected), nil
}

func Tag(ctx context.Context, _id graphql.ID, tag string) (n int32, err error) {
	if err = common.ValidateTag(tag); err != nil {
		return 0, err
	}
	id, err := common.DecodeCursor(_id)
	if err != nil {
		return 0, err
	}
	q := db.DB(ctx).Model(&db.Node{}).
		Where("id = ?", id).
		Update("tag", tag)
	if q.Error != nil {
		return 0, q.Error
	}
	return int32(q.RowsAffected), nil
}
