package istio

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/mizhexiaoxiao/k8s-api-service/controllers/k8s"
	networkingV1beta1 "istio.io/api/networking/v1beta1"
	"istio.io/client-go/pkg/apis/networking/v1beta1"
	versionedClient "istio.io/client-go/pkg/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type VSRoute struct {
	AppName            string                                `json:"appName,omitempty"`
	Version            string                                `json:"version,omitempty"`
	Category           string                                `json:"category,omitempty"`
	DeployNamespace    string                                `json:"deployNamespace,omitempty"`
	CanaryWeight       int32                                 `json:"weight,omitempty"`
	CanaryWeightSwitch bool                                  `json:"switch,omitempty"`
	HttpMatch          []*networkingV1beta1.HTTPMatchRequest `json:"match,omitempty"`
}

func replaceVersion(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, ".", "-"), "_", "-")
}

// When delete canary route rule, check all HTTPRouteDestination hava no canary version
func checkVsCanarySubsetExists(virtualService *v1beta1.VirtualService, namespace, appName, version string) error {
	subset := replaceVersion(version)
	canaryHost := fmt.Sprintf("%s-canary.%s.svc.cluster.local", appName, namespace)
	for _, m := range virtualService.Spec.Http {
		for _, n := range m.Route {
			if n.Destination.Host == canaryHost && n.Destination.Subset == subset {
				return fmt.Errorf("not all subset = %s is deleted, please check", subset)
			}
		}
	}
	return nil
}

// Get Uri if HTTPRoute.Match exists
func getVsMatchUri(httpRoute *networkingV1beta1.HTTPRoute) *networkingV1beta1.StringMatch {
	routeCopy := httpRoute.DeepCopy()
	if routeCopy == nil || len(routeCopy.Match) == 0 {
		return nil
	}
	for _, v := range routeCopy.Match {
		if v.Uri != nil {
			return v.Uri
		}
	}
	return nil
}

// Get HTTPRouteDestination index and weight of related route version
func getVsRouteDstIndexAndWeight(hDst []*networkingV1beta1.HTTPRouteDestination, version string) (int, int32) {
	subset := replaceVersion(version)
	for i, v := range hDst {
		if v.Destination.Subset == subset {
			return i, v.Weight
		}
	}
	return -1, -1
}

// Get HTTPRoute index by route name
func getRouteIndex(virtualService *v1beta1.VirtualService, name string) (int, error) {
	var (
		exists     bool
		routeIndex int
	)
	httpRoutes := virtualService.Spec.Http
	for idx, route := range httpRoutes {
		if route.Name == name {
			exists = true
			routeIndex = idx
			break
		}
	}
	if exists {
		return routeIndex, nil
	} else {
		return -1, fmt.Errorf("VirtualService HTTPRoute %q not found", name)
	}
}

// Insert a HTTPRoute into the specified index location
func insertRoute(s []*networkingV1beta1.HTTPRoute, index int, elem *networkingV1beta1.HTTPRoute) []*networkingV1beta1.HTTPRoute {
	ss := append([]*networkingV1beta1.HTTPRoute{}, s[index:]...)
	s = append(s[:index], elem)
	s = append(s, ss...)
	return s
}

func NewIstioClient(cluster string) (*versionedClient.Clientset, error) {
	k8sClient, err := k8s.GetClient(cluster)
	if err != nil {
		return nil, err
	}
	return versionedClient.NewForConfigOrDie(k8sClient.RestConfig), nil
}

type VSHttpRouteInterface interface {
	GetVS(ctx context.Context, name string) (*v1beta1.VirtualService, error)
	UpdateVS(ctx context.Context, virtualService *v1beta1.VirtualService) (*v1beta1.VirtualService, error)
	List(ctx context.Context, vsName, appName string) (*v1beta1.VirtualService, error)
	Get(ctx context.Context, vsName, routeName string) (*v1beta1.VirtualService, error)
	Create(ctx context.Context, vsName string, vr *VSRoute) (*v1beta1.VirtualService, error)
	Update(ctx context.Context, vsName, routeName string, vr *VSRoute) (*v1beta1.VirtualService, error)
	Delete(ctx context.Context, vsName, routeName string, vr *VSRoute) (*v1beta1.VirtualService, error)
}

type VSHttpRouteOperation struct {
	cs *versionedClient.Clientset
	ns string
}

func NewVSHttpRouteOperation(cs *versionedClient.Clientset, namespace string) VSHttpRouteInterface {
	return VSHttpRouteOperation{
		cs: cs,
		ns: namespace,
	}
}

func (o VSHttpRouteOperation) GetVS(ctx context.Context, name string) (*v1beta1.VirtualService, error) {
	return o.cs.NetworkingV1beta1().VirtualServices(o.ns).Get(ctx, name, metav1.GetOptions{})
}

func (o VSHttpRouteOperation) UpdateVS(ctx context.Context, vs *v1beta1.VirtualService) (*v1beta1.VirtualService, error) {
	return o.cs.NetworkingV1beta1().VirtualServices(o.ns).Update(ctx, vs, metav1.UpdateOptions{})
}

func (o VSHttpRouteOperation) List(ctx context.Context, vsName, appName string) (*v1beta1.VirtualService, error) {
	vs, err := o.GetVS(ctx, vsName)
	if err != nil {
		return nil, fmt.Errorf("VSHttpRouteOperation of GetVS failed, err: %s", err)
	}

	var routes []*networkingV1beta1.HTTPRoute
	httpRoutes := vs.Spec.Http
	stableRouteName := fmt.Sprintf("%s-stable", appName)
	canaryRoutePrefix := fmt.Sprintf("%s-canary-v", appName)
	for _, route := range httpRoutes {
		if strings.HasPrefix(route.Name, canaryRoutePrefix) || route.Name == stableRouteName {
			routes = append(routes, route)
		}
	}
	vs.Spec.Http = routes
	return vs, nil
}

func (o VSHttpRouteOperation) Get(ctx context.Context, vsName, routeName string) (*v1beta1.VirtualService, error) {
	vs, err := o.GetVS(ctx, vsName)
	if err != nil {
		return nil, fmt.Errorf("VSHttpRouteOperation of GetVS failed, err: %s", err)
	}
	var routes []*networkingV1beta1.HTTPRoute
	httpRoutes := vs.Spec.Http
	for _, v := range httpRoutes {
		if v.Name == routeName {
			routes = append(routes, v)
		}
	}
	vs.Spec.Http = routes
	return vs, nil
}

func (o VSHttpRouteOperation) Create(ctx context.Context, vsName string, vr *VSRoute) (*v1beta1.VirtualService, error) {
	var (
		canaryExists     bool
		defaultHeader    string
		firstCanaryIndex int
	)
	vs, err := o.GetVS(ctx, vsName)
	if err != nil {
		return nil, fmt.Errorf("VSHttpRouteOperation of GetVS failed, err: %s", err)
	}

	routeName := fmt.Sprintf("%s-%s", vr.AppName, replaceVersion(vr.Version))
	_, err = getRouteIndex(vs, routeName)
	if err == nil {
		return nil, fmt.Errorf("VirtualService HTTPRoute %q already exists", routeName)
	}

	stableRouteIndex, err := getRouteIndex(vs, fmt.Sprintf("%s-stable", vr.AppName))
	if err != nil {
		return nil, fmt.Errorf("getRouteIndex failed, err: %s", err)
	}

	switch vr.Category {
	case "backend":
		defaultHeader = "x-weike-forward"
	case "frontend":
		defaultHeader = "x-weike-fe-forward"
	default:
		defaultHeader = "x-weike-forward"
	}

	defaultHttpMatch := &networkingV1beta1.HTTPMatchRequest{
		Headers: map[string]*networkingV1beta1.StringMatch{
			defaultHeader: {
				MatchType: &networkingV1beta1.StringMatch_Exact{
					Exact: vr.Version,
				},
			},
		},
	}
	defaultHttpRouteDestination := []*networkingV1beta1.HTTPRouteDestination{
		{
			Destination: &networkingV1beta1.Destination{
				Host:   fmt.Sprintf("%s-canary.%s.svc.cluster.local", vr.AppName, vr.DeployNamespace),
				Subset: replaceVersion(vr.Version),
			},
			Weight: 100,
		},
		{
			Destination: &networkingV1beta1.Destination{
				Host:   fmt.Sprintf("%s.%s.svc.cluster.local", vr.AppName, vr.DeployNamespace),
				Subset: "stable",
			},
			Weight: 0,
		},
	}
	httpRoutes := vs.Spec.Http
	stableRoute := httpRoutes[stableRouteIndex]
	stableUri := getVsMatchUri(stableRoute)
	if stableUri != nil {
		defaultHttpMatch.Uri = stableUri
	}
	canaryHttpRoute := &networkingV1beta1.HTTPRoute{
		Name: routeName,
		Match: []*networkingV1beta1.HTTPMatchRequest{
			defaultHttpMatch,
		},
		Route: defaultHttpRouteDestination,
	}

	for i, v := range httpRoutes {
		if strings.HasPrefix(v.Name, fmt.Sprintf("%s-canary-v", vr.AppName)) {
			canaryExists = true
			firstCanaryIndex = i
			break
		}
	}
	if canaryExists {
		httpRoutes = insertRoute(httpRoutes, firstCanaryIndex, canaryHttpRoute)
	} else {
		httpRoutes = insertRoute(httpRoutes, stableRouteIndex, canaryHttpRoute)
	}
	vs.Spec.Http = httpRoutes
	result, err := o.UpdateVS(ctx, vs)
	if err != nil {
		return nil, fmt.Errorf("VSHttpRouteOperation of Update virtualService failed, err: %s", err)
	}
	return result, nil
}

func (o VSHttpRouteOperation) Update(ctx context.Context, vsName, routeName string, vr *VSRoute) (*v1beta1.VirtualService, error) {
	vs, err := o.GetVS(ctx, vsName)
	if err != nil {
		return nil, fmt.Errorf("VSHttpRouteOperation of GetVS failed, err: %s", err)
	}

	routeIndex, err := getRouteIndex(vs, routeName)
	if err != nil {
		return nil, fmt.Errorf("getRouteIndex failed, err: %s", err)
	}
	stableRouteIndex, err := getRouteIndex(vs, fmt.Sprintf("%s-stable", vr.AppName))
	if err != nil {
		return nil, fmt.Errorf("getRouteIndex failed, err: %s", err)
	}

	httpRoutes := vs.Spec.Http
	route := httpRoutes[routeIndex]
	stableRoute := httpRoutes[stableRouteIndex]
	dstIndex, dstWeight := getVsRouteDstIndexAndWeight(route.Route, vr.Version)
	// if canary weight changed, then update weight and use stable match replace canary match
	switch vr.CanaryWeightSwitch {
	case true:
		if vr.CanaryWeight != dstWeight && len(route.Route) == 2 {
			route.Route[dstIndex].Weight = vr.CanaryWeight
			route.Route[1-dstIndex].Weight = 100 - vr.CanaryWeight
			route.Match = stableRoute.Match
		} else {
			log.Printf("canary weight not changed, do nothing")
		}
	default:
		if stableUri := getVsMatchUri(stableRoute); stableUri != nil {
			for _, match := range vr.HttpMatch {
				if match.Uri == nil {
					match.Uri = stableUri
				}
			}
		}
		route.Match = vr.HttpMatch

		// reset canary weight back to default 100, stable weight back to default 0
		if dstWeight != 100 && len(route.Route) == 2 {
			route.Route[dstIndex].Weight = 100
			route.Route[1-dstIndex].Weight = 0
		}
	}

	httpRoutes[routeIndex] = route
	vs.Spec.Http = httpRoutes
	result, err := o.UpdateVS(ctx, vs)
	if err != nil {
		return nil, fmt.Errorf("VSHttpRouteOperation of Update virtualService failed, err: %s", err)
	}
	return result, nil
}

func (o VSHttpRouteOperation) Delete(ctx context.Context, vsName, routeName string, vr *VSRoute) (*v1beta1.VirtualService, error) {
	vs, err := o.GetVS(ctx, vsName)
	if err != nil {
		return nil, fmt.Errorf("VSHttpRouteOperation of GetVS failed, err: %s", err)
	}

	if strings.HasSuffix(routeName, "stable") {
		return nil, fmt.Errorf("stable vs rule %q cannot be deleted", routeName)
	}

	routeIndex, err := getRouteIndex(vs, routeName)
	if err != nil {
		return nil, fmt.Errorf("getRouteIndex failed, err: %s", err)
	}

	httpRoutes := vs.Spec.Http
	httpRoutes = append(httpRoutes[:routeIndex], httpRoutes[routeIndex+1:]...)
	if err = checkVsCanarySubsetExists(vs, o.ns, vr.AppName, vr.Version); err != nil {
		return nil, fmt.Errorf("checkVsCanarySubsetExists failed, err: %s", err)
	}
	vs.Spec.Http = httpRoutes
	result, err := o.UpdateVS(ctx, vs)
	if err != nil {
		return nil, fmt.Errorf("VSHttpRouteOperation of Update virtualService failed, err: %s", err)
	}
	return result, nil
}
