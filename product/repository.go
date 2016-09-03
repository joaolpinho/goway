package product


type IProductRepository interface {

	GetAllProducts() []Product_v1
	GetAllClients()	 []Client_v1

}
