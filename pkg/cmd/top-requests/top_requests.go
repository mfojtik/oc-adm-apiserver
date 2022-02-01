package top_requests

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	tm "github.com/buger/goterm"
	apiserverclient "github.com/openshift/client-go/apiserver/clientset/versioned/typed/apiserver/v1"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type TopRequestsOptions struct {
	genericclioptions.IOStreams
	configFlags *genericclioptions.ConfigFlags

	apiRequestCountClient apiserverclient.APIRequestCountInterface
}

func NewCmdTopRequests(streams genericclioptions.IOStreams) *cobra.Command {
	o := &TopRequestsOptions{
		IOStreams:   streams,
		configFlags: genericclioptions.NewConfigFlags(true),
	}

	cmd := &cobra.Command{
		Use:          "watch [flags]",
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			if err := o.Complete(c, args); err != nil {
				return err
			}
			if err := o.Validate(); err != nil {
				return err
			}
			if err := o.Run(); err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}

func (o *TopRequestsOptions) Complete(command *cobra.Command, args []string) error {
	config, err := o.configFlags.ToRESTConfig()
	if err != nil {
		return err
	}
	client, err := apiserverclient.NewForConfig(config)
	if err != nil {
		return err
	}
	o.apiRequestCountClient = client.APIRequestCounts()

	return nil
}

func (o *TopRequestsOptions) Validate() error {
	if o.apiRequestCountClient == nil {
		return errors.New("unable to initialize API request counts client")
	}
	return nil
}

type apiCounter struct {
	name       string
	countHour  int64
	countTotal int64
}

func (o *TopRequestsOptions) collect(ctx context.Context) ([]apiCounter, error) {
	result, err := o.apiRequestCountClient.List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	counter := []apiCounter{}
	for _, counts := range result.Items {
		counter = append(counter, apiCounter{
			name:       counts.GetName(),
			countHour:  counts.Status.CurrentHour.RequestCount,
			countTotal: counts.Status.RequestCount,
		})
	}
	sort.Slice(counter, func(i, j int) bool {
		return counter[i].countHour >= counter[j].countHour
	})
	return counter, nil
}

func (o *TopRequestsOptions) Run() error {
	tm.Clear()
	for {
		totals := tm.NewTable(0, 10, 5, ' ', 0)
		tm.MoveCursor(1, 1)
		fmt.Fprintf(totals, "Name\tRequests Last Hour\tRequests Total\n")

		result, err := o.collect(context.TODO())
		//fmt.Printf("result: %+v\n", result)
		if err != nil {
			tm.Flush()
			return err
		}
		for i, r := range result {
			fmt.Fprintf(totals, "%s\t%d\t%d\n", r.name, r.countHour, r.countTotal)
			if i > 20 {
				break
			}
		}
		tm.Println(totals)
		tm.Flush()
		time.Sleep(2 * time.Second)
	}
}
