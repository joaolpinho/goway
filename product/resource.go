package product


const (
	REPOSITORY_LOCAL = "local"
)

type ProductResourceOptions struct {
	RepositoryType string
}


type ProductResource struct{
	RepositoryType string
}


func NewProductResource(options *ProductResourceOptions) *ProductResource{

	if(len(options.RepositoryType)==0){
		options.RepositoryType = REPOSITORY_LOCAL
	}

	return &ProductResource{
		RepositoryType: options.RepositoryType,
	}
}

func(r *ProductResource) GetRepository() IProductRepository{
	switch r.RepositoryType {
	case REPOSITORY_LOCAL:
		return NewLocalRepository()

	default:
		return NewLocalRepository()
	}

}



