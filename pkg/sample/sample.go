package sample

import (
	"context"
	"log"
	"log/slog"
	"os"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)

var (
	_ framework.Plugin           = &Plugin{}
	_ framework.PreFilterPlugin  = &Plugin{}
	_ framework.PreBindPlugin    = &Plugin{}
	_ framework.FilterPlugin     = &Plugin{}
	_ framework.PostFilterPlugin = &Plugin{}
)

var logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
	Level: slog.LevelDebug,
}))

const (
	Name = "Sample"
)

type Args struct {
	FavoriteColor  string `json:"favorite_color,omitempty"`
	FavoriteNumber int    `json:"favorite_number,omitempty"`
	ThanksTo       string `json:"thanks_to,omitempty"`
}

type Plugin struct {
	*kubernetes.Clientset
	handle framework.Handle
	args   *Args
}

// PostFilter implements framework.PostFilterPlugin.
func (p *Plugin) PostFilter(ctx context.Context, state *framework.CycleState, pod *v1.Pod, filteredNodeStatusMap framework.NodeToStatusMap) (*framework.PostFilterResult, *framework.Status) {
	log.Printf("PostFilter: pod %s", pod.Name)
	log.Printf("PostFilter: filteredNodeStatusMap = %+v", filteredNodeStatusMap)

	klog.V(3).Infof("PostFilter: pod %s", pod.Name)
	klog.V(3).Infof("PostFilter: filteredNodeStatusMap = %+v", filteredNodeStatusMap)
	return &framework.PostFilterResult{}, framework.NewStatus(framework.Success, "OK")
}

// Filter implements framework.FilterPlugin.
func (p *Plugin) Filter(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeInfo *framework.NodeInfo) *framework.Status {
	//curTime, err := state.Read(StateKeyCurTime)
	//if err != nil {
	//	return framework.NewStatus(framework.Error, err.Error())
	//}
	//log.Printf("Filter: state.Read(StateKeyCurTime) = %+v", curTime)
	log.Printf("Filter: nodeInfo.Name = %+v", nodeInfo.Node().Name)
	log.Printf("Filter: pod = %+v", pod.Name)

	// klog.V(3).Infof("Filter: state.Read(StateKeyCurTime) = %+v", curTime)
	klog.V(3).Infof("Filter: nodeInfo.Name = %+v", nodeInfo.Node().Name)
	klog.V(3).Infof("Filter: pod = %+v", pod.Name)
	return framework.NewStatus(framework.Success, "OK")
}

// PreBind implements framework.PreBindPlugin.
func (*Plugin) PreBind(ctx context.Context, state *framework.CycleState, p *v1.Pod, nodeName string) *framework.Status {
	log.Printf("PreBind: state = %+v", state)
	log.Printf("PreBind: pod = %+v", p.Name)
	log.Printf("PreBind: nodeName = %+v", nodeName)

	klog.V(3).Infof("PreBind: state = %+v", state)
	klog.V(3).Infof("PreBind: pod = %+v", p.Name)
	klog.V(3).Infof("PreBind: nodeName = %+v", nodeName)
	return framework.NewStatus(framework.Success, "OK")
}

// PreFilter implements framework.PreFilterPlugin.
func (p *Plugin) PreFilter(ctx context.Context, state *framework.CycleState, pod *v1.Pod) (*framework.PreFilterResult, *framework.Status) {
	log.Printf("PreFilter: state = %+v", state)
	log.Printf("PreFilter: pod = %+v", pod.Name)

	//if pod.Annotations == nil {
	//	return nil, framework.NewStatus(framework.Error, "pod annotations is nil")
	//}
	//if _, ok := pod.Annotations["xxx"]; !ok {
	//	return nil, framework.NewStatus(framework.Error, "xxx annotation is not set")
	//}
	//log.Info("PreFilter: pod %s", pod.Name)
	//klog.V(3).Infof("PreFiltering pod %s", pod.Name)
	//nodes, err := p.handle.SnapshotSharedLister().NodeInfos().List()
	//if err != nil {
	//	return nil, framework.NewStatus(framework.Error, err.Error())
	//}
	//if len(nodes) == 0 {
	//	return nil, framework.NewStatus(framework.Error, "no nodes found")
	//}
	//win := nodes[0]
	//log.Info("PreFilter: chose node %s", win.Node().Name)
	//klog.V(3).Infof("PreFiltering chose node %s", win.Node().Name)
	//state.Write(StateKeyCurTime, &StateDataCurTime{
	//	CurTime: time.Now().Format(time.DateTime),
	//})
	return &framework.PreFilterResult{
		// 过滤的节点，传递给下一轮流程
		// 为 nil 代表所有节点都有资格
		// NodeNames: sets.NewString(win.Node().Name),
	}, framework.NewStatus(framework.Success, "OK")
}

// PreFilterExtensions implements framework.PreFilterPlugin.
func (p *Plugin) PreFilterExtensions() framework.PreFilterExtensions {
	panic("unimplemented")
}

func (p *Plugin) Name() string {
	return Name
}

func New(ctx context.Context, configuration runtime.Object, f framework.Handle) (framework.Plugin, error) {
	return &Plugin{handle: f}, nil
}
