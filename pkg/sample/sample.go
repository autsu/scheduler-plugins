package sample

import (
	"context"
	"math/rand/v2"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)

var _ framework.Plugin = &Plugin{}
var (
	_ framework.QueueSortPlugin  = &Plugin{}
	_ framework.PreFilterPlugin  = &Plugin{}
	_ framework.FilterPlugin     = &Plugin{}
	_ framework.PostFilterPlugin = &Plugin{}
	_ framework.PreScorePlugin   = &Plugin{}
	_ framework.ScorePlugin      = &Plugin{}
	_ framework.ReservePlugin    = &Plugin{}
	_ framework.PermitPlugin     = &Plugin{}
	_ framework.PreBindPlugin    = &Plugin{}
	_ framework.BindPlugin       = &Plugin{}
	_ framework.PostBindPlugin   = &Plugin{}
)

const (
	Name          = "Sample"
	annotationKey = "k8s.io/sample"
)

type Args struct {
	FavoriteColor  string `json:"favorite_color,omitempty"`
	FavoriteNumber int    `json:"favorite_number,omitempty"`
	ThanksTo       string `json:"thanks_to,omitempty"`
}

type Plugin struct {
	handle framework.Handle
	args   *Args
}

// NormalizeScore implements framework.ScoreExtensions.
// 对节点分数归一化
func (p *Plugin) NormalizeScore(
	ctx context.Context, state *framework.CycleState,
	pod *corev1.Pod, scores framework.NodeScoreList,
) *framework.Status {
	if len(scores) == 0 {
		return framework.NewStatus(framework.Error, "no nodes found")
	}
	klog.V(3).Infof("NormalizeScore: pod %s, scores = %+v", pod.Name, scores)
	max := scores[0].Score
	min := scores[0].Score
	for _, score := range scores[1:] {
		if score.Score > max {
			max = score.Score
		}
		if score.Score < min {
			min = score.Score
		}
	}
	if max == min {
		return framework.NewStatus(framework.Success, "OK")
	}
	// 归一化分数 = (原始分数 - 最小分数) * 100 / (最大分数 - 最小分数)
	for i := range scores {
		scores[i].Score = (scores[i].Score - min) * 100 / (max - min)
	}
	return framework.NewStatus(framework.Success, "OK")
}

// PostBind implements framework.PostBindPlugin.
func (p *Plugin) PostBind(
	ctx context.Context, state *framework.CycleState,
	pod *corev1.Pod, nodeName string,
) {
	klog.V(3).Infof("PostBind: pod %s to node %s", pod.Name, nodeName)
}

// Bind implements framework.BindPlugin.
func (p *Plugin) Bind(
	ctx context.Context, state *framework.CycleState,
	pod *corev1.Pod, nodeName string,
) *framework.Status {
	klog.V(3).Infof("Bind: pod %s to node %s", pod.Name, nodeName)

	binding := &corev1.Binding{
		ObjectMeta: metav1.ObjectMeta{
			Name: pod.Name,
		},
		Target: corev1.ObjectReference{
			APIVersion: "v1",
			Kind:       "Node",
			Name:       nodeName,
		},
	}

	if err := p.handle.ClientSet().CoreV1().Pods(pod.Namespace).Bind(ctx, binding, metav1.CreateOptions{}); err != nil {
		return framework.NewStatus(framework.Error, err.Error())
	}

	return framework.NewStatus(framework.Success, "OK")
}

// Permit implements framework.PermitPlugin.
// 延迟调度
func (p *Plugin) Permit(
	ctx context.Context, state *framework.CycleState,
	pod *corev1.Pod, nodeName string,
) (*framework.Status, time.Duration) {
	klog.V(3).Infof("Permit: pod %s to node %s, wait 10 seconds", pod.Name, nodeName)
	// 10s 后再调度到节点
	return framework.NewStatus(framework.Success, "OK"), time.Second * 10
}

// Reserve implements framework.ReservePlugin.
func (p *Plugin) Reserve(
	ctx context.Context, state *framework.CycleState,
	pod *corev1.Pod, nodeName string,
) *framework.Status {
	klog.V(3).Infof("Reserve: pod %s to node %s", pod.Name, nodeName)
	return framework.NewStatus(framework.Success, "OK")
}

// Unreserve implements framework.ReservePlugin.
func (p *Plugin) Unreserve(
	ctx context.Context, state *framework.CycleState,
	pod *corev1.Pod, nodeName string,
) {
	klog.V(3).Infof("Unreserve: pod %s from node %s", pod.Name, nodeName)
}

// Score implements framework.ScorePlugin.
// 给节点打分
func (p *Plugin) Score(
	ctx context.Context, state *framework.CycleState,
	pod *corev1.Pod, nodeName string,
) (int64, *framework.Status) {
	score := rand.Int64N(100)
	klog.V(3).Infof("Score: pod %s to node %s, node score = %d", pod.Name, nodeName, score)
	return score, framework.NewStatus(framework.Success, "OK")
}

// ScoreExtensions implements framework.ScorePlugin.
// 归一化
func (p *Plugin) ScoreExtensions() framework.ScoreExtensions {
	return p
}

var _ framework.ScoreExtensions = &Plugin{}

// PreScore implements framework.PreScorePlugin.
func (p *Plugin) PreScore(
	ctx context.Context, state *framework.CycleState,
	pod *corev1.Pod, nodes []*framework.NodeInfo,
) *framework.Status {
	klog.V(3).Infof("PreScore: pod %s", pod.Name)
	for _, node := range nodes {
		klog.V(3).Infof("PreScore: node %s", node.Node().Name)
	}
	return framework.NewStatus(framework.Success, "OK")
}

// Less implements framework.QueueSortPlugin.
func (p *Plugin) Less(x *framework.QueuedPodInfo, y *framework.QueuedPodInfo) bool {
	return x.Timestamp.Before(y.Timestamp)
}

// PostFilter implements framework.PostFilterPlugin.
func (p *Plugin) PostFilter(
	ctx context.Context, state *framework.CycleState,
	pod *corev1.Pod, filteredNodeStatusMap framework.NodeToStatusMap,
) (*framework.PostFilterResult, *framework.Status) {
	klog.V(3).Infof("PostFilter: pod %s", pod.Name)
	klog.V(3).Infof("PostFilter: filteredNodeStatusMap = %+v", filteredNodeStatusMap)
	return &framework.PostFilterResult{}, framework.NewStatus(framework.Success, "OK")
}

// Filter implements framework.FilterPlugin.
func (p *Plugin) Filter(
	ctx context.Context, state *framework.CycleState,
	pod *corev1.Pod, nodeInfo *framework.NodeInfo,
) *framework.Status {
	klog.V(3).Infof("Filter: pod %s, node %s", pod.Name, nodeInfo.Node().Name)
	if _, ok := nodeInfo.Node().Annotations[annotationKey]; !ok {
		return framework.NewStatus(framework.Unschedulable, "missing annotation "+annotationKey)
	}
	return framework.NewStatus(framework.Success, "OK")
}

// PreBind implements framework.PreBindPlugin.
func (*Plugin) PreBind(
	ctx context.Context, state *framework.CycleState,
	pod *corev1.Pod, nodeName string,
) *framework.Status {
	klog.V(3).Infof("PreBind: pod %s to node %s", pod.Name, nodeName)
	return framework.NewStatus(framework.Success, "OK")
}

// PreFilter implements framework.PreFilterPlugin.
func (p *Plugin) PreFilter(
	ctx context.Context, state *framework.CycleState,
	pod *corev1.Pod,
) (*framework.PreFilterResult, *framework.Status) {
	klog.V(3).Infof("PreFilter: pod %s", pod.Name)
	nodes, err := p.handle.SnapshotSharedLister().NodeInfos().List()
	if err != nil {
		return nil, framework.NewStatus(framework.Error, err.Error())
	}
	var workerNodes []*framework.NodeInfo
	for _, node := range nodes {
		var isControlPlane bool
		for _, taint := range node.Node().Spec.Taints {
			if taint.Key == "node-role.kubernetes.io/control-plane" {
				isControlPlane = true
				break
			}
		}
		if !isControlPlane {
			workerNodes = append(workerNodes, node)
		}
	}
	s := sets.New[string]()
	for _, node := range workerNodes {
		s.Insert(node.Node().Name)
	}
	return &framework.PreFilterResult{
		// 过滤的节点，传递给下一轮流程
		// 为 nil 代表所有节点都有资格
		NodeNames: s,
	}, framework.NewStatus(framework.Success, "OK")
}

// PreFilterExtensions implements framework.PreFilterPlugin.
func (p *Plugin) PreFilterExtensions() framework.PreFilterExtensions {
	klog.V(3).Infof("PreFilterExtensions")
	return nil
}

func (p *Plugin) Name() string {
	return Name
}

func New(ctx context.Context, configuration runtime.Object, f framework.Handle) (framework.Plugin, error) {
	return &Plugin{handle: f}, nil
}
