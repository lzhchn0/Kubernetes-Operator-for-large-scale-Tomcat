/*
SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and redis-operator contributors
SPDX-License-Identifier: Apache-2.0
*/

package operator

import (
	mynginxv1 "opencanon.com/api/v1"
	"opencanon.com/internal/generator"

	"embed"
	"flag"

	"github.com/pkg/errors"
	"github.com/sap/component-operator-runtime/pkg/component"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/sap/component-operator-runtime/pkg/operator"
)

// const Name = "opencanon.io/my-mar26"

const Name = "mymar"

//go:embed all:data
var data embed.FS

type Options struct {
	Name                  string
	DefaultServiceAccount string
	FlagPrefix            string
}

type Operator struct {
	options Options
}

var defaultOperator operator.Operator = New()

func GetName() string {
	return defaultOperator.GetName()
}

func InitScheme(scheme *runtime.Scheme) {
	defaultOperator.InitScheme(scheme)
}

func InitFlags(flagset *flag.FlagSet) {
	defaultOperator.InitFlags(flagset)
}

func ValidateFlags() error {
	return defaultOperator.ValidateFlags()
}

func GetUncacheableTypes() []client.Object {
	return defaultOperator.GetUncacheableTypes()
}

func Setup(mgr ctrl.Manager) error {
	return defaultOperator.Setup(mgr)
}

func New() *Operator {
	return NewWithOptions(Options{})
}

func NewWithOptions(options Options) *Operator {
	operator := &Operator{options: options}
	if operator.options.Name == "" {
		operator.options.Name = Name
	}
	return operator
}

func (o *Operator) GetName() string {
	return o.options.Name
}

func (o *Operator) InitScheme(scheme *runtime.Scheme) {
	utilruntime.Must(mynginxv1.AddToScheme(scheme))
}

func (o *Operator) InitFlags(flagset *flag.FlagSet) {

	flagset.StringVar(&o.options.DefaultServiceAccount,
		"default-service-account",
		o.options.DefaultServiceAccount,
		"Default service account name")
}

func (o *Operator) ValidateFlags() error {
	return nil
}

func (o *Operator) GetUncacheableTypes() []client.Object {
	return []client.Object{&mynginxv1.Monitoring{}}
}

func (o *Operator) Setup(mgr ctrl.Manager) error {
	resourceGenerator, err := generator.NewResourceGenerator(
		data,
		"data/charts/tomcat",
		mgr.GetClient(),
	)
	if err != nil {
		return errors.Wrap(err, "error initializing resource generator")
	}

	if err := component.NewReconciler[*mynginxv1.Monitoring](
		o.options.Name,
		resourceGenerator,
		component.ReconcilerOptions{
			DefaultServiceAccount: &o.options.DefaultServiceAccount,
		},
	).SetupWithManager(mgr); err != nil {
		return errors.Wrapf(err, "unable to create controller")
	}

	return nil
}
