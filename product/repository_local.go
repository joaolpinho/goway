package product

type LocalProductRepository struct {

}

func NewLocalRepository() *LocalProductRepository{
	return &LocalProductRepository{}
}

func(l *LocalProductRepository) GetAllProducts() []Product_v1 {

	return []Product_v1{

		Product_v1{
			Code: "customer",
			Name: "Customer Api",
			Version: "1",
			Routes: []Routes_v1{
				Routes_v1{
					ListenPath: "/api/facets",
					Verb: "GET",
					ServiceName: "authentication",
					Handlers: []string{"AUTHENTICATION", "METRICS"},
					Code:"auth_by_email",
					Roles: []string{},
					InjectData: []InjectData_v1{
						InjectData_v1{
							Where:"params",
							Code: "paramId",
							Value: "123456789",
						},
						InjectData_v1{
							Where:"header",
							Code: "headerId",
							Value: "9999999999",
						},
					},
					InjectGlobalData:true,
				},
				Routes_v1{
					ListenPath: "/api/facets",
					Verb: "POST",
					ServiceName: "authentication",
					Handlers: []string{"AUTHENTICATION", "METRICS"},
					Code:"auth_by_email",
					Roles: []string{},
					InjectData: []InjectData_v1{
						InjectData_v1{
							Where:"params",
							Code: "paramId",
							Value: "123456789",
						},
						InjectData_v1{
							Where:"header",
							Code: "headerId",
							Value: "9999999999",
						},
					},
					InjectGlobalData:true,
				},
			},
		},
		Product_v1{
			Code: "cockpit",
			Name: "cockpit Api",
			Version: "1",
			Routes: []Routes_v1{
				Routes_v1{
					ListenPath: "/auth/byemail",
					Verb: "GET",
					ServiceName: "authentication",
					Handlers: []string{"AUTHENTICATION", "METRICS"},
					Code:"auth_by_email",
					Roles: []string{},
				},
				Routes_v1{
					ListenPath: "/auth/renew",
					Verb: "GET",
					ServiceName: "authentication",
					Handlers:  []string{"AUTHENTICATION", "METRICS"},
					Code:"renew_token",
					Roles: []string{},
				},
			},
		},

	}
}

func(l *LocalProductRepository) GetAllClients() []Client_v1 {
	return []Client_v1{
		Client_v1{
			ApiPath:"12124578",
			Product:"customer",
			Client:"myorg",
			RemoveApiPath: true,
			Version: "1",
			GlobalInjectData:[]InjectData_v1{

			},
			Routes: []Routes_v1{
				Routes_v1{
					ListenPath: "/auth/byemail",
					Verb: "GET",
					ServiceName: "authentication-custom",
					Handlers: []string{"AUTHENTICATION", "METRICS"},
					Code:"auth_by_email",
					Roles: []string{},
				},
			},
		},
		Client_v1{
			ApiPath:"121245782",
			Product:"customer",
			Client:"myorg2",
			RemoveApiPath: true,
			Version: "1",
			GlobalInjectData:[]InjectData_v1{
				InjectData_v1{
					Where:"url",
					Code: "orgs",
					Value: "tlantic",
				},
				InjectData_v1{
					Where:"url",
					Code: "apps",
					Value: "customer",
				},
			},
			Routes: []Routes_v1{
				Routes_v1{
					ListenPath: "/auth/byemail",
					Verb: "GET",
					ServiceName: "authentication-custom",
					Handlers: []string{"AUTHENTICATION", "METRICS"},
					Code:"auth_by_email",
					Roles: []string{},
				},
			},
		},
	}
}
