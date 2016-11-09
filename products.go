package main

import "net/url"

const dropletBasePath = "v1/products"

type ProductsService struct {
	client *Client
}

// Create product
func (p *ProductsService) Create(createRequest *ProductCreateRequest) (*Product, *Response, error) {
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	path := dropletBasePath

	requestWithAttributes := struct {
		Attributes *ProductCreateRequest
	}{
		Attributes: createRequest,
	}

	req, err := p.client.NewRequest("POST", path, requestWithAttributes)
	if err != nil {
		return nil, nil, err
	}

	root := new(productRoot)
	resp, err := p.client.Do(req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Product, resp, err
}

// ProductRoot represents a Product root
type productRoot struct {
	Product *Product `json:"product"`
	// Links   *Links   `json:"links,omitempty"`
}

type ProductCreateRequest struct {
	Name        string
	url         *url.URL
	description string
	variant     []Variant
}

type Product struct {
	Name        string
	url         *url.URL
	description string
	variant     []Variant
}
