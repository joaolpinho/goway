package product


const (
	WHERE_HEADER 	= "header"
	WHERE_PARAMS	= "params"
	WHERE_URL 	= "url"
)

type Product_v1 struct {
	Type 				string 				`json:"_type"`
	Id 				string 				`json:"id"`
	Code 				string 				`json:"code"`
	Name 				string 				`json:"name"`
	Version 			string 				`json:"version"`
	Routes	 			[]Routes_v1			`json:"routes"`
}

type Routes_v1  struct {
	ListenPath 			string 				`json:"listen_path"`
	Verb 				string 				`json:"verb"`
	ServiceName 			string 				`json:"service_name"`
	Code				string 				`json:"code"`
	Handlers 			[]string 			`json:"handlers"`
	Roles 				[]string 			`json:"roles"`
	InjectData			[]InjectData_v1			`json:"inject_data"`
	InjectGlobalData		bool				`json:"inject_global_data"`
	Asset				string				`json:"asset"`
	AssetId				string				`json:"asset_id"`
	Alias				string				`json:"alias"`
}

type InjectData_v1 struct {
	Where 				string 				`json:"where"`
	Code				string 				`json:"code"`
	Value				string 				`json:"value"`
}

type Client_v1 struct {
	ApiPath				string 				`json:"api_path"`
	Product				string				`json:"product"`
	ProductVersion			string				`json:"product_version"`
	Client				string 				`json:"client"`
	RemoveApiPath			bool				`json:"remove_api_path"`
	Version 			string				`json:"version"`
	GlobalInjectData		[]InjectData_v1			`json:"global_inject_data"`
	Routes	 			[]Routes_v1			`json:"routes"`
}