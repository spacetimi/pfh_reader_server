package faq

import "github.com/spacetimi/timi_shared_server/code/core/services/metadata_service/metadata_typedefs"

const MetadataKey = "FaqMetadata"

type MetadataFactory struct {   // Implements IMetadataFactory
}

type Metadata struct {          // Implements IMetadataItem
    FaqItems []Item
}

type Item struct {
    Question string
    Answer string
}

func (fmf MetadataFactory) Instantiate() metadata_typedefs.IMetadataItem {
	return &Metadata{}
}

func (fm Metadata) GetKey() string {
    return MetadataKey
}

func (fm Metadata) GetMetadataSpace() metadata_typedefs.MetadataSpace {
    return metadata_typedefs.METADATA_SPACE_APP
}
