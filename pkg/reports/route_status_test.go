package reports_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/kgateway-dev/kgateway/v2/pkg/reports"
)

// A delegated child route is reported against its parent HTTPRoute with the namespace set. If the
// child's spec parentRef omits the namespace (so it defaults to the child's namespace), the status
// must still find that report, and must use the parentRef exactly as written in the spec.
func TestBuildRouteStatusParentRefWithoutNamespace(t *testing.T) {
	specParentRef := gwv1.ParentReference{
		Group: new(gwv1.Group(gwv1.GroupName)),
		Kind:  new(gwv1.Kind("HTTPRoute")),
		Name:  "parent",
	}
	child := &gwv1.HTTPRoute{
		ObjectMeta: metav1.ObjectMeta{Name: "child", Namespace: "team1"},
		Spec: gwv1.HTTPRouteSpec{
			CommonRouteSpec: gwv1.CommonRouteSpec{ParentRefs: []gwv1.ParentReference{specParentRef}},
		},
	}

	rm := reports.NewReportMap()
	reportedParentRef := specParentRef
	reportedParentRef.Namespace = new(gwv1.Namespace("team1"))
	reports.NewReporter(&rm).Route(child).ParentRef(&reportedParentRef)

	status := rm.BuildRouteStatus(context.Background(), child, "kgateway.dev/kgateway")
	require.NotNil(t, status)
	require.Len(t, status.Parents, 1, "child route status should have an entry for its parentRef")
	require.Equal(t, specParentRef, status.Parents[0].ParentRef, "status parentRef must match the spec parentRef exactly")
}
