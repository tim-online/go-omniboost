package omniboost

import "github.com/tim-online/go-omniboost/utils"

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
	Products []Product `json:"data"`
	Meta     *Meta     `json:"meta,omitempty"`
}

// ProductRoot represents a Product root
type productRoot struct {
	Product *Product `json:"product"`
	Meta    *Meta    `json:"meta,omitempty"`
}

func newProductRoot() *productRoot {
	return &productRoot{
		Product: &Product{},
		Meta:    newMeta(),
	}
}

type ProductCreateRequest struct {
	Name        string
	Url         *utils.URL
	Description string
	Variant     []Variant
}

type Product struct {
	ID          int        `jsonapi:"primary,products"`
	Name        string     `jsonapi:"attr,name"`
	Url         *utils.URL `jsonapi:"attr,url"`
	Description string     `jsonapi:"attr,description"`
	Variants    []Variant  `jsonapi:"relation,variants"`
}
