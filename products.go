package omniboost

import "net/url"

const productsBasePath = "v1/products"

type ProductsService struct {
	client *Client
}

// List all products
func (s *ProductsService) List() ([]Product, *Response, error) {
	path := productsBasePath
	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(productsRoot)
	resp, err := s.client.Do(req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Products, resp, err
}

// Create product
func (p *ProductsService) Create(createRequest *ProductCreateRequest) (*Product, *Response, error) {
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	path := productsBasePath

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

// ProductsRoot represents a Product root
type productsRoot struct {
	Products []Product `json:"products"`
	// Links   *Links   `json:"links,omitempty"`
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
	ID          int
	Name        string
	url         *url.URL
	description string
	variant     []Variant
}
