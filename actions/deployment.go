package actions

import (
	"errors"
	"strconv"

	"github.com/CenturyLinkLabs/panamaxcli/config"
)

func ListDeployments(remote config.Remote) (Output, error) {
	c := DefaultAgentClientFactory.New(remote)
	deps, err := c.ListDeployments()
	if err != nil {
		return PlainOutput{}, err
	}

	if len(deps) == 0 {
		return PlainOutput{"No Deployments"}, nil
	}

	o := ListOutput{Labels: []string{"ID", "Name"}}
	for _, d := range deps {
		o.AddRow(map[string]string{
			"ID":   strconv.Itoa(d.ID),
			"Name": d.Name,
		})
	}

	return &o, nil
}

func DescribeDeployment(remote config.Remote, id string) (Output, error) {
	c := DefaultAgentClientFactory.New(remote)
	if id == "" {
		return PlainOutput{}, errors.New("Empty ID")
	}

	desc, err := c.DescribeDeployment(id)
	if err != nil {
		return PlainOutput{}, err
	}

	o := DetailOutput{
		Details: map[string]string{
			"Name":         desc.Name,
			"ID":           strconv.Itoa(desc.ID),
			"Redeployable": strconv.FormatBool(desc.Redeployable),
		},
	}

	return &o, nil
	return PlainOutput{}, nil
}
