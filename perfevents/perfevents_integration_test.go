// +build medium

/*
http://www.apache.org/licenses/LICENSE-2.0.txt


Copyright 2015 Intel Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package perfevents

import (
	"errors"
	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestPerfEventsCollector(t *testing.T) {
	Convey("GetMetricTypes functionality", t, func() {
		p := NewPerfeventsCollector()
		Convey("invalid init", func() {
			p.Init = func() error { return errors.New("error") }
			_, err := p.GetMetricTypes(plugin.Config{})
			So(err, ShouldNotBeNil)
		})
		Convey("set_supported_metrics", func() {
			cg := []string{"cgroup1", "cgroup2", "cgroup3"}
			events := []string{"event1", "event2", "event3"}
			a := get_supported_metrics(NsSubtype, cg, events)
			So(len(a), ShouldEqual, 9)
			So(a[len(a)-1].Namespace.Strings(), ShouldResemble, []string{NsVendor, NsClass, NsType, NsSubtype, "event3", "cgroup3"})
		})
		Convey("set_supported_metrics with special characters", func() {
			cg := []string{"machine_slice:machine-quemu\x2d1\x2dinstance\x2d001"}
			events := []string{"event1"}
			a := get_supported_metrics(NsSubtype, cg, events)
			So(a, ShouldNotBeNil)
			So(len(a), ShouldEqual, 1)
			So(a[0].Namespace.Strings(), ShouldResemble, []string{NsVendor, NsClass, NsType, NsSubtype, "event1", "machine_slice:machine-quemu\x2d1\x2dinstance\x2d001"})
		})
		Convey("flatten cgroup name", func() {
			cg := []string{"cg_root/cg_sub1/cg_sub2"}
			events := []string{"event1"}
			a := get_supported_metrics(NsSubtype, cg, events)
			So(a, ShouldNotBeNil)
			So(len(a), ShouldEqual, 1)
			So(a[0].Namespace.Strings(), ShouldResemble, []string{NsVendor, NsClass, NsType, NsSubtype, "event1", "cg_root:cg_sub1:cg_sub2"})
		})
	})
	Convey("CollectMetrics error cases", t, func() {
		p := NewPerfeventsCollector()
		Convey("empty list of requested metrics", func() {
			metricTypes := []plugin.Metric{}
			metrics, err := p.CollectMetrics(metricTypes)
			So(err, ShouldBeNil)
			So(metrics, ShouldBeEmpty)
		})
		Convey("namespace too short", func() {
			_, err := p.CollectMetrics(
				[]plugin.Metric{
					plugin.Metric{
						Namespace: plugin.NewNamespace("invalid"),
					},
				},
			)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "segments")
		})
		Convey("namespace wrong vendor", func() {
			_, err := p.CollectMetrics(
				[]plugin.Metric{
					plugin.Metric{
						Namespace: plugin.NewNamespace("invalid", NsClass, NsType, NsSubtype, "cycles", "A"),
					},
				},
			)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "1st")
		})
		Convey("namespace wrong class", func() {
			_, err := p.CollectMetrics(
				[]plugin.Metric{
					plugin.Metric{
						Namespace: plugin.NewNamespace(NsVendor, "invalid", NsType, NsSubtype, "cycles", "A"),
					},
				},
			)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "2nd")
		})
		Convey("namespace wrong type", func() {
			_, err := p.CollectMetrics(
				[]plugin.Metric{
					plugin.Metric{
						Namespace: plugin.NewNamespace(NsVendor, NsClass, "invalid", NsSubtype, "cycles", "A"),
					},
				},
			)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "3rd")
		})
		Convey("namespace wrong subtype", func() {
			_, err := p.CollectMetrics(
				[]plugin.Metric{
					plugin.Metric{
						Namespace: plugin.NewNamespace(NsVendor, NsClass, NsType, "invalid", "cycles", "A"),
					},
				},
			)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "4th")
		})
		Convey("namespace wrong event", func() {
			_, err := p.CollectMetrics(
				[]plugin.Metric{
					plugin.Metric{
						Namespace: plugin.NewNamespace(NsVendor, NsClass, NsType, NsSubtype, "invalid", "A"),
					},
				},
			)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "5th")
		})

	})
}
