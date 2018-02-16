// shippy-consignment-service/repository.go

package main

import (
	pb "github.com/seiji-thirdbridge/shippy-consignment-service/proto/consignment"
	mgo "gopkg.in/mgo.v2"
)

const (
	dbName                = "shippy"
	consignmentCollection = "consignments"
)

type repository interface {
	Create(*pb.Consignment) error
	GetAll() ([]*pb.Consignment, error)
	Close()
}

// ConsignmentRepository - implementation of Repository interface
type ConsignmentRepository struct {
	session *mgo.Session
}

// Create - create a new consignment
func (repo *ConsignmentRepository) Create(consignment *pb.Consignment) error {
	return repo.collection().Insert(consignment)
}

// GetAll - returns all consignments
func (repo *ConsignmentRepository) GetAll() ([]*pb.Consignment, error) {
	var consignments []*pb.Consignment

	err := repo.collection().Find(nil).All(&consignments)
	return consignments, err
}

// Close the database session
func (repo *ConsignmentRepository) Close() {
	repo.session.Close()
}

func (repo *ConsignmentRepository) collection() *mgo.Collection {
	return repo.session.DB(dbName).C(consignmentCollection)
}
