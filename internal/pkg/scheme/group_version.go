// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package scheme

// GroupVersion 是 API 的唯一标识.
type GroupVersion struct {
	Group   string
	Version string
}

func (gv GroupVersion) WithResource(resource string) GroupVersionResource {
	return GroupVersionResource{Group: gv.Group, Version: gv.Version, Resource: resource}
}

// GroupResource 是组内资源的标识.
type GroupResource struct {
	Group    string
	Resource string
}

// GroupVersionResource 是指定版本的组内资源的标识.
type GroupVersionResource struct {
	Group    string
	Version  string
	Resource string
}

func (gvr GroupVersionResource) GroupResource() GroupResource {
	return GroupResource{Group: gvr.Group, Resource: gvr.Resource}
}
