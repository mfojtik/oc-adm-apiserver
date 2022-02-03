package requests

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	v1 "github.com/openshift/api/apiserver/v1"

	tm "github.com/buger/goterm"
	"github.com/spf13/cobra"

	apiserverclient "github.com/openshift/client-go/apiserver/clientset/versioned/typed/apiserver/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type RequestsOptions struct {
	genericclioptions.IOStreams
	configFlags *genericclioptions.ConfigFlags

	apiRequestCountClient apiserverclient.APIRequestCountInterface
}

func NewCmdTopRequests(streams genericclioptions.IOStreams) *cobra.Command {
	o := &RequestsOptions{
		IOStreams:   streams,
		configFlags: genericclioptions.NewConfigFlags(true),
	}

	cmd := &cobra.Command{
		Use:          "requests [flags]",
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

func (o *RequestsOptions) Complete(command *cobra.Command, args []string) error {
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

func (o *RequestsOptions) Validate() error {
	if o.apiRequestCountClient == nil {
		return errors.New("unable to initialize API request counts client")
	}
	return nil
}

func (o *RequestsOptions) collect(ctx context.Context, limit int) ([]v1.APIRequestCount, error) {
	current, err := o.apiRequestCountClient.List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// limit the amount of API request counts by sorting the list using current hour request counts.
	items := current.Items
	sort.Slice(items, func(i, j int) bool {
		return items[i].Status.CurrentHour.RequestCount >= items[j].Status.CurrentHour.RequestCount
	})
	if limit > 0 && len(items) > limit {
		items = items[0:limit]
	}

	// limit the amount of node records in API request count, by sorting the node records based on request counts.
	for i := range items {
		byNode := items[i].Status.CurrentHour.ByNode
		sort.Slice(byNode, func(i, j int) bool {
			return byNode[i].RequestCount >= byNode[j].RequestCount
		})
		if limit > 0 && len(byNode) > limit {
			byNode = byNode[0:limit]
		}
		for j := range byNode {
			byUser := byNode[j].ByUser
			sort.Slice(byUser, func(i, j int) bool {
				return byUser[i].RequestCount >= byUser[j].RequestCount
			})
			if limit > 0 && len(byUser) > limit {
				byUser = byUser[0:limit]
			}
			for g := range byUser {
				byVerb := byUser[g].ByVerb
				sort.Slice(byVerb, func(i, j int) bool {
					return byVerb[i].RequestCount >= byVerb[j].RequestCount
				})
				byUser[g].ByVerb = byVerb
			}

			byNode[j].ByUser = byUser
		}
		items[i].Status.CurrentHour.ByNode = byNode
	}

	return items, nil
}

func (o *RequestsOptions) Run() error {
	tm.Clear()
	for {
		tm.MoveCursor(1, 1)
		counts, err := o.collect(context.TODO(), 20)
		if err != nil {
			tm.Flush()
			return err
		}

		apiRequestPerResourceTable := tm.NewTable(0, 10, 5, ' ', 0)
		fmt.Fprintf(apiRequestPerResourceTable, boldWhite("Resource Name\tRequests Last Hour\tRequests Total\tTop Node")+"\n")
		for _, c := range counts {
			nodes := []string{}
			users := []string{}
			for nodeCount, n := range c.Status.CurrentHour.ByNode {
				for userCount, u := range n.ByUser {
					verbs := []string{}
					for verbsCount, v := range u.ByVerb {
						verbs = append(verbs, fmt.Sprintf("%sx%d", strings.ToUpper(v.Verb), v.RequestCount))
						if verbsCount >= 0 {
							break
						}
					}
					users = append(users, fmt.Sprintf("%s/%s [%d reqests](%s)", u.UserName, u.UserAgent, u.RequestCount, strings.Join(verbs, ", ")))
					if userCount >= 0 {
						break
					}
				}
				nodes = append(nodes, fmt.Sprintf("%s [%d requests] (%s)", n.NodeName, n.RequestCount, strings.Join(users, ", ")))
				// TODO: make the node report count configurable
				if nodeCount >= 0 {
					break
				}
			}
			fmt.Fprintf(apiRequestPerResourceTable, "%s\t%d\t%d\t%s\n", c.GetName(), c.Status.CurrentHour.RequestCount, c.Status.RequestCount, strings.Join(nodes, ","))
		}
		tm.Println(apiRequestPerResourceTable)

		tm.Flush()
		// TODO: make update interval configurable
		time.Sleep(2 * time.Second)
	}
}

var (
	boldWhite = color("\033[1;37m%s\033[0m")
)

func color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}
