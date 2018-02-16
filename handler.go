// shippy-consignment-service/handler.go

package main

import (
	"context"
	"log"

	"gopkg.in/mgo.v2"

	pb "github.com/seiji-thirdbridge/shippy-consignment-service/proto/consignment"
	vesselProto "github.com/seiji-thirdbridge/shippy-vessel-service/proto/vessel"
)

type service struct {
	session      *mgo.Session
	vesselClient vesselProto.VesselServiceClient
}

func (s *service) getRepo() repository {
	return &ConsignmentRepository{s.session.Clone()}
}

// CreateConsignment - Stores the provided consignment
func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {
	repo := s.getRepo()
	defer repo.Close()

	vesselResponse, err := s.vesselClient.FindAvailable(context.Background(), &vesselProto.Specification{
		MaxWeight: req.Weight,
		Capacity:  int32(len(req.Containers)),
	})
	if err != nil {
		return err
	}
	log.Printf("Found vessel: %s\n", vesselResponse.Vessel.Name)

	req.VesselId = vesselResponse.Vessel.Id

	// Save the consignment
	err = repo.Create(req)
	if err != nil {
		return err
	}

	res.Created = true
	res.Consignment = req
	return nil
}

// GetConsignments - Get a list of all consignments
func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	repo := s.getRepo()
	defer repo.Close()

	consignments, err := repo.GetAll()
	if err != nil {
		return err
	}

	res.Consignments = consignments
	return nil
}
