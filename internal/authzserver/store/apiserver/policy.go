package apiserver

import (
	"context"
	"encoding/json"

	"github.com/AlekSi/pointer"
	"github.com/avast/retry-go"
	"github.com/changaolee/skeleton/internal/authzserver/store"
	"github.com/changaolee/skeleton/pkg/errors"
	"github.com/changaolee/skeleton/pkg/log"
	pb "github.com/changaolee/skeleton/pkg/proto/apiserver/v1"
	"github.com/ory/ladon"
)

type policyStore struct {
	ds *datastore
}

var _ store.PolicyStore = (*policyStore)(nil)

func newPolicies(ds *datastore) *policyStore {
	return &policyStore{ds: ds}
}

// List 返回所有的授权策略.
func (p *policyStore) List() (map[string][]*ladon.DefaultPolicy, error) {
	pols := make(map[string][]*ladon.DefaultPolicy)

	log.Infow("Loading policies")

	req := &pb.ListPoliciesRequest{
		Offset: pointer.ToInt64(0),
		Limit:  pointer.ToInt64(-1),
	}

	var resp *pb.ListPoliciesResponse
	err := retry.Do(
		func() error {
			var listErr error
			resp, listErr = p.ds.cli.ListPolicies(context.Background(), req)
			if listErr != nil {
				return listErr
			}

			return nil
		}, retry.Attempts(3),
	)
	if err != nil {
		return nil, errors.Wrap(err, "list policies failed")
	}

	log.Infof("Policies found (%d total)[username:name]:", len(resp.Items))

	for _, v := range resp.Items {
		log.Infof(" - %s:%s", v.Username, v.Name)

		var policy ladon.DefaultPolicy

		if err := json.Unmarshal([]byte(v.PolicyShadow), &policy); err != nil {
			log.Warnf("Failed to load policy for %s, error: %s", v.Name, err.Error())

			continue
		}

		pols[v.Username] = append(pols[v.Username], &policy)
	}

	return pols, nil
}
