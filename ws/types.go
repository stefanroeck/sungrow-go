package ws

type InverterParams struct {
	Protocol string
	Host     string
	Port     int
	User     string
	Password string
	Path     string
	Data     string
	Types    []string
}

type Keys map[string]Key

type KeyType string

var KeyTypes = struct {
	String KeyType
	Number KeyType
}{
	String: "string",
	Number: "number",
}

type Key struct {
	Name    string
	KeyType KeyType
}

type RequestConnect struct {
	Lang    string `json:"lang"`
	Token   string `json:"token"`
	Service string `json:"service"`
}
type ResponseConnect struct {
	ResultCode int    `json:"result_code"`
	ResultMsg  string `json:"result_msg"`
	ResultData struct {
		Service     string
		Token       string
		Uid         int
		TipsDisable int `json:"tips_disable"`
	} `json:"result_data"`
}

type RequestLogin struct {
	Lang     string `json:"lang"`
	Token    string `json:"token"`
	Service  string `json:"service"`
	Username string `json:"username"`
	Passwd   string `json:"passwd"`
}
type ResponseLogin struct {
	ResultCode int    `json:"result_code"`
	ResultMsg  string `json:"result_msg"`
	ResultData struct {
		Service     string
		Token       string
		Uid         int
		TipsDisable int `json:"tips_disable"`
	} `json:"result_data"`
}

type RequestReal struct {
	Lang       string `json:"lang"`
	Token      string `json:"token"`
	DevId      string `json:"dev_id"`
	Service    string `json:"service"`
	Time123456 int64  `json:"time123456"`
}
type ResponseReal struct {
	ResultCode int    `json:"result_code"`
	ResultMsg  string `json:"result_msg"`
	ResultData struct {
		Service string
		List    []struct {
			DataName  string `json:"data_name"`
			DataValue string `json:"data_value"`
			DataUnit  string `json:"data_unit"`
		}
		Count int
	} `json:"result_data"`
}
