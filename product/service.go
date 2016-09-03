package product

func(r *ProductResource) GetAllProducts() []Product_v1{
	return r.GetRepository().GetAllProducts()
}

func(r *ProductResource) GetAllClients() []Client_v1{
	return r.GetRepository().GetAllClients()
}