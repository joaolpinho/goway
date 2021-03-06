package product


const (
	REPOSITORY_LOCAL = "local"
)

type ProductResourceOptions struct {
	Repository IProductRepository
}


type ProductResource struct{
	Repository IProductRepository
}


func NewProductResource(options *ProductResourceOptions) *ProductResource{

	if(options.Repository==nil){
		panic("options.Repository is required")
	}

	return &ProductResource{
		Repository: options.Repository,
	}
}

func(r *ProductResource) GetRepository() IProductRepository{
	return r.Repository
}



