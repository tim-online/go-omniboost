package main

// json := []byte(`{
// 	"name": "Kahn test 1",
// 	"url": "http://www.tim-online.nl/kahn.html",
// 	"description": "Kaaaaaaaaaaaahn",
// 	"variant": [
// 	{
// 		"sku": "kahn-test-1",
// 		"price": 40,
// 		"images": []
// 	}
// 	]
// }`)

// omni := NewClient()
// createRequest := &ProductCreateRequest{}
// json.Unmarshal(json, createRequest)
// omni.Products.Create(createRequest)

// Kijk: 

// import "golang.org/x/oauth2"

// pat := "mytoken"
// type TokenSource struct {
//     AccessToken string
// }

// func (t *TokenSource) Token() (*oauth2.Token, error) {
//     token := &oauth2.Token{
//         AccessToken: t.AccessToken,
//     }
//     return token, nil
// }

// tokenSource := &TokenSource{
//     AccessToken: pat,
// }
// oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
// client := godo.NewClient(oauthClient)
