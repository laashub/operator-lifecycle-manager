package v1

import (
	"encoding/json"

	operatorsv1alpha1 "github.com/operator-framework/operator-lifecycle-manager/pkg/api/apis/operators/v1alpha1"
	"github.com/operator-framework/operator-registry/pkg/registry"
)

const (
	// The yaml attribute that specifies the related images of the ClusterServiceVersion
	relatedImages = "relatedImages"
)

// CreateCSVDescription creates a CSVDescription from a given CSV
func CreateCSVDescription(csv *operatorsv1alpha1.ClusterServiceVersion, csvJSON string) CSVDescription {
	desc := CSVDescription{
		DisplayName: csv.Spec.DisplayName,
		Version:     csv.Spec.Version,
		Provider: AppLink{
			Name: csv.Spec.Provider.Name,
			URL:  csv.Spec.Provider.URL,
		},
		Annotations:               csv.GetAnnotations(),
		LongDescription:           csv.Spec.Description,
		InstallModes:              csv.Spec.InstallModes,
		CustomResourceDefinitions: csv.Spec.CustomResourceDefinitions,
		APIServiceDefinitions:     csv.Spec.APIServiceDefinitions,
		NativeAPIs:                csv.Spec.NativeAPIs,
		MinKubeVersion:            csv.Spec.MinKubeVersion,
		RelatedImages:             GetImages(csvJSON),
	}

	icons := make([]Icon, len(csv.Spec.Icon))
	for i, icon := range csv.Spec.Icon {
		icons[i] = Icon{
			Base64Data: icon.Data,
			Mediatype:  icon.MediaType,
		}
	}

	if len(icons) > 0 {
		desc.Icon = icons
	}

	return desc
}

// GetImages returns a list of images listed in CSV (spec and deployments)
func GetImages(csvJSON string) []string {
	var images []string

	csv := &registry.ClusterServiceVersion{}
	err := json.Unmarshal([]byte(csvJSON), &csv)
	if err != nil {
		return images
	}

	imageSet, err := csv.GetOperatorImages()
	if err != nil {
		return images
	}

	relatedImgSet, err := csv.GetRelatedImages()
	if err != nil {
		return images
	}

	for k := range relatedImgSet {
		imageSet[k] = struct{}{}
	}

	for k := range imageSet {
		images = append(images, k)
	}

	return images
}
