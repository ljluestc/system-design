package fetcher

import (
    "crawler/internal/utils"
)

type FTPFetcher struct{}

func NewFTPFetcher() *FTPFetcher {
    return &FTPFetcher{}
}

func (f *FTPFetcher) Fetch(url string) ([]byte, error) {
    // Placeholder for FTP fetching
    return nil, utils.Errorf("FTP fetching not implemented")
}