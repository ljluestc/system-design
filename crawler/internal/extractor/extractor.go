package extractor

import (
    "crawler/pkg/models"
)

type Extractor struct{}

func New() *Extractor {
    return &Extractor{}
}

func (e *Extractor) Extract(content []byte) (models.Document, []string, error) {
    urls, text, err := ParseHTML(content)
    if err != nil {
        return models.Document{}, nil, err
    }
    doc := models.Document{
        Content: text,
    }
    return doc, urls, nil
}