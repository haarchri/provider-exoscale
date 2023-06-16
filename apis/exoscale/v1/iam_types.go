package v1

import (
	"reflect"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	// AccessKeyIDName is the environment variable name for the S3 access key ("username")
	AccessKeyIDName = "AWS_ACCESS_KEY_ID"
	// SecretAccessKeyName is the environment variable name for the S3 secret key ("password")
	SecretAccessKeyName = "AWS_SECRET_ACCESS_KEY"
)

// SOSSpec is the service type for Object Storage in exoscale
type SOSSpec struct {

	// +kubebuilder:validation:Required

	// Buckets is a list of buckets to which IAMKey has access to.
	Buckets []string `json:"buckets,omitempty"`
}

// ServicesSpec are the accessible exoscale services of the IAMKey.
type ServicesSpec struct {

	// +kubebuilder:validation:Required

	// SOSSpec is the Object Storage Service in exoscale.
	SOS SOSSpec `json:"sos,omitempty"`
}

// IAMKeyParameters are the configurable fields of IAMKey.
type IAMKeyParameters struct {
	// KeyName is the name of the Key as presented in the exoscale.com UI.
	// If empty, the value of `.metadata.annotations."crossplane.io/external-name"` is used.
	// There can be multiple keys that have the same key name in exoscale.com, but they will have different key IDs.
	KeyName string `json:"keyName,omitempty"`

	// +kubebuilder:validation:Required

	// Zone is the name of the zone where the IAM key is created.
	// The zone must be available in the S3 endpoint.
	// Cannot be changed after IAMKey is created.
	Zone string `json:"zone"`

	// +kubebuilder:validation:Required

	// Services is the exoscale service to which IAMKey gets access to.
	// Only object storage (sos) service is supported thus the IAMKey will be restricted to access only sos.
	Services ServicesSpec `json:"services,omitempty"`
}

// IAMKeySpec defines the desired state of an IAMKey.
type IAMKeySpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       IAMKeyParameters `json:"forProvider"`
}

// IAMKeyObservation contains the observed fields of an IAMKey.
type IAMKeyObservation struct {
	// KeyID is the observed unique ID as generated by exoscale.com.
	KeyID string `json:"keyID,omitempty"`

	// RoleID is the observed unique ID as generated by exoscale.com.
	RoleID string `json:"roleID,omitempty"`

	// KeyName is the observed key name as generated by exoscale.com.
	KeyName string `json:"keyName,omitempty"`

	// ServicesSpec is the exoscale service to which IAMKey gets access to.
	ServicesSpec `json:"services,omitempty"`
}

// IAMKeyStatus represents the observed state of a IAMKey.
type IAMKeyStatus struct {
	xpv1.ResourceStatus `json:",inline"`

	AtProvider IAMKeyObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="Synced",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="External Name",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="Key ID",type="string",JSONPath=".status.atProvider.keyID"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,exoscale}
// +kubebuilder:webhook:verbs=create;update,path=/validate-exoscale-crossplane-io-v1-iamkey,mutating=false,failurePolicy=fail,groups=exoscale.crossplane.io,resources=iamkeys,versions=v1,name=iamkeys.exoscale.crossplane.io,sideEffects=None,admissionReviewVersions=v1

// IAMKey is the API for creating IAM Object Storage Keys on exoscale.com.
type IAMKey struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IAMKeySpec   `json:"spec"`
	Status IAMKeyStatus `json:"status,omitempty"`
}

// GetKeyName returns the IAMKey key name in the following precedence:
//
//	.spec.forProvider.keyName
//	.metadata.annotations."crossplane.io/external-name"
//	.metadata.name
func (in *IAMKey) GetKeyName() string {
	if in.Spec.ForProvider.KeyName != "" {
		return in.Spec.ForProvider.KeyName
	}
	if name := meta.GetExternalName(in); name != "" {
		return name
	}
	return in.Name
}

// GetProviderConfigName returns the name of the ProviderConfig.
// Returns empty string if reference not given.
func (in *IAMKey) GetProviderConfigName() string {
	if ref := in.GetProviderConfigReference(); ref != nil {
		return ref.Name
	}
	return ""
}

// +kubebuilder:object:root=true

// IAMKeyList contains a list of IAMKey
type IAMKeyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IAMKey `json:"items"`
}

// IAMKey type metadata.
var (
	IAMKeyKind             = reflect.TypeOf(IAMKey{}).Name()
	IAMKeyGroupKind        = schema.GroupKind{Group: Group, Kind: IAMKeyKind}.String()
	IAMKeyGroupVersionKind = SchemeGroupVersion.WithKind(IAMKeyKind)
)

func init() {
	SchemeBuilder.Register(&IAMKey{}, &IAMKeyList{})
}
