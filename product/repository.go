package product


type IProductRepository interface {

	GetAllProducts() []Product_v1
	GetAllClients()	 []Client_v1
	CreateProduct(product *Product_v1) (bool, *Product_v1)
	CreateClient(client *Client_v1) (bool, *Client_v1)

}
